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

package updates

import (
	"context"
	"fmt"
	"github.com/onosproject/onos-lib-go/pkg/atomix"
	"io"
	"time"

	"github.com/onosproject/onos-lib-go/pkg/logging"

	_map "github.com/atomix/go-client/pkg/client/map"
	"github.com/atomix/go-client/pkg/client/primitive"
	"github.com/atomix/go-client/pkg/client/util/net"
	"github.com/gogo/protobuf/proto"

	"github.com/onosproject/onos-ric/api/sb"
)

var log = logging.GetLogger("store", "updates")

const timeout = 15 * time.Second
const primitiveName = "control-updates"

// NewDistributedStore creates a new distributed control updates store
func NewDistributedStore() (Store, error) {
	log.Info("Creating distributed control updates store")
	database, err := atomix.GetAtomixDatabase()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	updates, err := database.GetMap(ctx, primitiveName)
	if err != nil {
		return nil, err
	}
	return &atomixStore{
		updates: updates,
	}, nil
}

// NewLocalStore returns a new local control updates store
func NewLocalStore() (Store, error) {
	_, address := atomix.StartLocalNode()
	return newLocalStore(address)
}

// newLocalStore creates a new local control updates store
func newLocalStore(address net.Address) (Store, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	session, err := primitive.NewSession(ctx, primitive.Partition{ID: 1, Address: address})
	if err != nil {
		return nil, err
	}
	updatesName := primitive.Name{
		Namespace: "local",
		Name:      primitiveName,
	}
	updates, err := _map.New(context.Background(), updatesName, []*primitive.Session{session})
	if err != nil {
		return nil, err
	}
	return &atomixStore{
		updates: updates,
	}, nil
}

// Store is ran store
type Store interface {
	io.Closer

	// Gets a control update message based on a given ID
	Get(ID) (*sb.ControlUpdate, error)

	// Puts a control update message to the store
	Put(*sb.ControlUpdate) error

	// Deletes a control update message from the store
	Delete(*sb.ControlUpdate) error

	// Deletes a control update message from store using key
	DeleteWithKey(key string) error

	// List all of the last up to date control update messages
	List(ch chan<- sb.ControlUpdate) error

	// Watch watches the store for control update messages
	Watch(ch chan<- sb.ControlUpdate, opts ...WatchOption) error

	// Clear deletes all entries from the store
	Clear() error
}

// WatchOption is a RAN store watch option
type WatchOption interface {
	apply(options *watchOptions)
}

// watchOptions is a struct of RAN store watch options
type watchOptions struct {
	replay bool
}

// WithReplay returns a watch option that replays existing updates
func WithReplay() WatchOption {
	return &watchReplayOption{
		replay: true,
	}
}

// watchReplayOption is an option for configuring whether to replay updates on watch calls
type watchReplayOption struct {
	replay bool
}

func (o *watchReplayOption) apply(options *watchOptions) {
	options.replay = o.replay
}

// atomixStore is an Atomix based store
type atomixStore struct {
	updates _map.Map
}

func (s *atomixStore) Get(id ID) (*sb.ControlUpdate, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	entry, err := s.updates.Get(ctx, id.String())
	if err != nil {
		return nil, err
	} else if entry == nil || entry.Value == nil {
		return nil, nil
	}
	update := &sb.ControlUpdate{}
	if err := proto.Unmarshal(entry.Value, update); err != nil {
		return nil, err
	}
	return update, nil
}

func (s *atomixStore) Put(update *sb.ControlUpdate) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	bytes, err := proto.Marshal(update)
	if err != nil {
		return err
	}
	_, err = s.updates.Put(ctx, getKey(update).String(), bytes)
	if err != nil {
		return err
	}
	return nil
}

func (s *atomixStore) Delete(update *sb.ControlUpdate) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err := s.updates.Remove(ctx, getKey(update).String())
	if err != nil {
		return err
	}
	return nil
}

func (s *atomixStore) DeleteWithKey(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err := s.updates.Remove(ctx, key)
	if err != nil {
		return err
	}
	return nil
}

func (s *atomixStore) List(ch chan<- sb.ControlUpdate) error {
	entryCh := make(chan *_map.Entry)
	if err := s.updates.Entries(context.Background(), entryCh); err != nil {
		return err
	}
	go func() {
		defer close(ch)
		for entry := range entryCh {
			update := &sb.ControlUpdate{}
			if err := proto.Unmarshal(entry.Value, update); err == nil {
				ch <- *update
			}
		}
	}()
	return nil
}

func (s *atomixStore) Watch(ch chan<- sb.ControlUpdate, opts ...WatchOption) error {
	options := &watchOptions{}
	for _, opt := range opts {
		opt.apply(options)
	}
	watchOptions := []_map.WatchOption{}
	if options.replay {
		watchOptions = append(watchOptions, _map.WithReplay())
	}

	watchCh := make(chan *_map.Event)
	if err := s.updates.Watch(context.Background(), watchCh, watchOptions...); err != nil {
		return err
	}
	go func() {
		defer close(ch)
		for event := range watchCh {
			if event.Type != _map.EventRemoved {
				update := &sb.ControlUpdate{}
				if err := proto.Unmarshal(event.Entry.Value, update); err == nil {
					ch <- *update
				}
			}
		}
	}()
	return nil
}

func (s *atomixStore) Clear() error {
	return s.updates.Clear(context.Background())
}

func (s *atomixStore) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return s.updates.Close(ctx)
}

func getKey(update *sb.ControlUpdate) ID {
	var ecgi sb.ECGI
	var crnti string
	switch update.MessageType {
	case sb.MessageType_CELL_CONFIG_REPORT:
		ecgi = *update.GetCellConfigReport().GetEcgi()
	case sb.MessageType_RRM_CONFIG_STATUS:
		ecgi = *update.GetRRMConfigStatus().GetEcgi()
	case sb.MessageType_UE_ADMISSION_REQUEST:
		ecgi = *update.GetUEAdmissionRequest().GetEcgi()
		crnti = update.GetUEAdmissionRequest().GetCrnti()
	case sb.MessageType_UE_ADMISSION_STATUS:
		ecgi = *update.GetUEAdmissionStatus().GetEcgi()
		crnti = update.GetUEAdmissionStatus().GetCrnti()
	case sb.MessageType_UE_CONTEXT_UPDATE:
		ecgi = *update.GetUEContextUpdate().GetEcgi()
		crnti = update.GetUEContextUpdate().GetCrnti()
	case sb.MessageType_BEARER_ADMISSION_REQUEST:
		ecgi = *update.GetBearerAdmissionRequest().GetEcgi()
		crnti = update.GetBearerAdmissionStatus().GetCrnti()
	case sb.MessageType_BEARER_ADMISSION_STATUS:
		ecgi = *update.GetBearerAdmissionStatus().GetEcgi()
		crnti = update.GetBearerAdmissionStatus().GetCrnti()
	}
	id := ID{
		PlmnID:      ecgi.PlmnId,
		Ecid:        ecgi.Ecid,
		Crnti:       crnti,
		MessageType: update.GetMessageType(),
	}
	return id
}

// ID store id
type ID struct {
	MessageType sb.MessageType
	PlmnID      string
	Ecid        string
	Crnti       string
}

func (i ID) String() string {
	return fmt.Sprintf("%d:%s:%s:%s", i.MessageType, i.PlmnID, i.Ecid, i.Crnti)
}
