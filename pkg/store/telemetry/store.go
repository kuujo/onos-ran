// Copyright 2020-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package telemetry

import (
	"context"
	"fmt"
	"github.com/onosproject/onos-lib-go/pkg/atomix"
	"github.com/onosproject/onos-ric/api/store/telemetry"
	"github.com/onosproject/onos-ric/pkg/config"
	"github.com/onosproject/onos-ric/pkg/store/mastership"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/onosproject/onos-lib-go/pkg/logging"

	"github.com/atomix/go-client/pkg/client/map"
	"github.com/atomix/go-client/pkg/client/primitive"
	"github.com/atomix/go-client/pkg/client/util/net"
	"github.com/gogo/protobuf/proto"

	"github.com/onosproject/onos-ric/api/sb"
)

var log = logging.GetLogger("store", "telemetry")

const requestTimeout = 15 * time.Second
const retryInterval = time.Second
const primitiveName = "telemetry"

// NewDistributedStore creates a new distributed telemetry store
func NewDistributedStore(mastershipStore mastership.Store) (Store, error) {
	log.Info("Creating distributed telemetry store")
	ricConfig, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	database, err := atomix.GetDatabase(ricConfig.Atomix, ricConfig.Atomix.GetDatabase(atomix.DatabaseTypeConsensus))

	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	telemetry, err := database.GetMap(ctx, primitiveName)
	if err != nil {
		return nil, err
	}
	store := &distributedStore{
		dist:       telemetry,
		entries:    make(map[ID]entryStore),
		mastership: mastershipStore,
	}
	if err := store.open(); err != nil {
		return nil, err
	}
	return store, nil
}

// NewLocalStore returns a new local telemetry store
func NewLocalStore(mastershipStore mastership.Store) (Store, error) {
	_, address := atomix.StartLocalNode()
	return newLocalStore(address, mastershipStore)
}

// newLocalStore creates a new local telemetry store
func newLocalStore(address net.Address, mastershipStore mastership.Store) (Store, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	session, err := primitive.NewSession(ctx, primitive.Partition{ID: 1, Address: address})
	if err != nil {
		return nil, err
	}
	telemetryName := primitive.Name{
		Namespace: "local",
		Name:      primitiveName,
	}
	telemetry, err := _map.New(context.Background(), telemetryName, []*primitive.Session{session})
	if err != nil {
		return nil, err
	}

	return &distributedStore{
		dist:       telemetry,
		entries:    make(map[ID]entryStore),
		mastership: mastershipStore,
	}, nil
}

// Store is interface for telemetry store
type Store interface {
	io.Closer

	// Gets a telemetry message based on a given ID
	Get(ID) (*sb.TelemetryMessage, error)

	// GetRevision gets a telemetry message for the given revision
	GetRevision(Revision) (*sb.TelemetryMessage, error)

	// Puts a telemetry message to the store
	Put(*sb.TelemetryMessage) error

	// Removes a telemetry message from the store
	Delete(ID) error

	// DeleteRevision deletes a telemetry message revision from the store
	DeleteRevision(revision Revision) error

	// List all of the last up to date telemetry messages
	List(ch chan<- sb.TelemetryMessage) error

	// Watch watches telemetry updates
	Watch(ch chan<- sb.TelemetryMessage, opts ...WatchOption) error

	// Clear deletes all entries from the store
	Clear() error
}

// WatchOption is a telemetry store watch option
type WatchOption interface {
	apply(options *watchOptions)
}

// watchOptions is a struct of telemetry store watch options
type watchOptions struct {
	replay bool
}

// WithReplay returns a watch option that replays existing telemetry
func WithReplay() WatchOption {
	return &watchReplayOption{
		replay: true,
	}
}

// watchReplayOption is an option for configuring whether to replay telemetry on watch calls
type watchReplayOption struct {
	replay bool
}

func (o *watchReplayOption) apply(options *watchOptions) {
	options.replay = o.replay
}

// distributedStore is an Atomix based store
type distributedStore struct {
	dist       _map.Map
	entries    map[ID]entryStore
	mastership mastership.Store
	mu         sync.RWMutex
}

// open opens the store
func (s *distributedStore) open() error {
	ch := make(chan *_map.Event)
	if err := s.dist.Watch(context.Background(), ch, _map.WithReplay()); err != nil {
		return err
	}
	go func() {
		for event := range ch {
			id := parseID(event.Entry.Key)
			s.mu.RLock()
			store, ok := s.entries[id]
			s.mu.RUnlock()
			if ok {
				if entry, err := decodeEntry(event.Entry); err == nil {
					tombstone := event.Type == _map.EventRemoved
					go store.update(entry, tombstone)
				}
			}
		}
	}()
	return nil
}

// getKeyStore uses a double checked lock to get or create a key store
func (s *distributedStore) getEntryStore(id ID) entryStore {
	s.mu.RLock()
	store, ok := s.entries[id]
	s.mu.RUnlock()
	if !ok {
		s.mu.Lock()
		store, ok = s.entries[id]
		if !ok {
			store = newEntryStore(s.dist, s.mastership)
			s.entries[id] = store
		}
		s.mu.Unlock()
	}
	return store
}

func (s *distributedStore) Get(id ID) (*sb.TelemetryMessage, error) {
	return s.getEntryStore(id).get(Revision{ID: id})
}

func (s *distributedStore) GetRevision(revision Revision) (*sb.TelemetryMessage, error) {
	return s.getEntryStore(revision.ID).get(revision)
}

func (s *distributedStore) Put(telemetry *sb.TelemetryMessage) error {
	return s.getEntryStore(newID(telemetry)).put(telemetry)
}

func (s *distributedStore) Delete(id ID) error {
	return s.getEntryStore(id).delete(Revision{ID: id})
}

func (s *distributedStore) DeleteRevision(revision Revision) error {
	return s.getEntryStore(revision.ID).delete(revision)
}

func (s *distributedStore) List(ch chan<- sb.TelemetryMessage) error {
	entryCh := make(chan *_map.Entry)
	if err := s.dist.Entries(context.Background(), entryCh); err != nil {
		return err
	}
	go func() {
		defer close(ch)
		for entry := range entryCh {
			telemetry := &sb.TelemetryMessage{}
			if err := proto.Unmarshal(entry.Value, telemetry); err == nil {
				ch <- *telemetry
			}
		}
	}()
	return nil
}

func (s *distributedStore) Watch(ch chan<- sb.TelemetryMessage, opts ...WatchOption) error {
	options := &watchOptions{}
	for _, opt := range opts {
		opt.apply(options)
	}
	watchOptions := []_map.WatchOption{}
	if options.replay {
		watchOptions = append(watchOptions, _map.WithReplay())
	}

	watchCh := make(chan *_map.Event)
	if err := s.dist.Watch(context.Background(), watchCh, watchOptions...); err != nil {
		return err
	}
	go func() {
		defer close(ch)
		for event := range watchCh {
			if event.Type != _map.EventRemoved {
				if telemetry, err := decodeEntry(event.Entry); err == nil {
					ch <- *telemetry.Message
				}
			}
		}
	}()
	return nil
}

func (s *distributedStore) Clear() error {
	return s.dist.Clear(context.Background())
}

func (s *distributedStore) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	return s.dist.Close(ctx)
}

type waiterContext struct {
	term      Term
	timestamp Timestamp
	ch        chan *telemetry.TelemetryEntry
}

const keySep = ":"

// parseID parses the telemetry ID from a string
func parseID(key string) ID {
	parts := strings.Split(key, keySep)
	i, _ := strconv.Atoi(parts[0])
	return ID{
		MessageType: sb.MessageType(i),
		PlmnID:      parts[1],
		Ecid:        parts[2],
		Crnti:       parts[3],
	}
}

// ID is a telemetry store identifier
type ID struct {
	MessageType sb.MessageType
	PlmnID      string
	Ecid        string
	Crnti       string
}

func (i ID) String() string {
	return fmt.Sprintf("%d%s%s%s%s%s%s", i.MessageType, keySep, i.PlmnID, keySep, i.Ecid, keySep, i.Crnti)
}

// Term is a telemetry store term
type Term uint64

// Timestamp is a telemetry store timestamp
type Timestamp uint64

// Revision is a telemetry message revision
type Revision struct {
	// ID is the revision identifier
	ID
	// Term is the term in which the revision was created
	Term Term
	// Timestamp is the timestamp at which the revision was created
	Timestamp Timestamp
}

func newID(telemetry *sb.TelemetryMessage) ID {
	var ecgi sb.ECGI
	var crnti string
	switch telemetry.MessageType {
	case sb.MessageType_RADIO_MEAS_REPORT_PER_UE:
		ecgi = *telemetry.GetRadioMeasReportPerUE().GetEcgi()
		crnti = telemetry.GetRadioMeasReportPerUE().GetCrnti()
	case sb.MessageType_RADIO_MEAS_REPORT_PER_CELL:
		ecgi = *telemetry.GetRadioMeasReportPerCell().GetEcgi()
	case sb.MessageType_SCHED_MEAS_REPORT_PER_UE:
	case sb.MessageType_SCHED_MEAS_REPORT_PER_CELL:
		ecgi = *telemetry.GetSchedMeasReportPerCell().GetEcgi()
	}
	return ID{
		PlmnID:      ecgi.PlmnId,
		Ecid:        ecgi.Ecid,
		Crnti:       crnti,
		MessageType: telemetry.GetMessageType(),
	}
}
