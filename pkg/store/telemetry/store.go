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
	"io"
	"time"

	_map "github.com/atomix/go-client/pkg/client/map"
	"github.com/atomix/go-client/pkg/client/primitive"
	"github.com/atomix/go-client/pkg/client/util/net"
	"github.com/gogo/protobuf/proto"
	"github.com/onosproject/onos-ric/pkg/store/utils"

	"github.com/onosproject/onos-ric/api/sb"
	log "k8s.io/klog"
)

const timeout = 15 * time.Second
const primitiveName = "telemetry"

// NewDistributedStore creates a new distributed telemetry store
func NewDistributedStore() (Store, error) {
	log.Info("Creating distributed telemetry store")
	database, err := utils.GetAtomixDatabase()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	telemetry, err := database.GetMap(ctx, primitiveName)
	if err != nil {
		return nil, err
	}
	return &atomixStore{
		telemetry: telemetry,
	}, nil
}

// NewLocalStore returns a new local telemetry store
func NewLocalStore() (Store, error) {
	_, address := utils.StartLocalNode()
	return newLocalStore(address)
}

// newLocalStore creates a new local telemetry store
func newLocalStore(address net.Address) (Store, error) {
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
	return &atomixStore{
		telemetry: telemetry,
	}, nil
}

// Store is interface for telemetry store
type Store interface {
	io.Closer

	// Gets a telemetry message based on a given ID
	Get(ID) (*sb.TelemetryMessage, error)

	// Puts a telemetry message to the store
	Put(*sb.TelemetryMessage) error

	// Removes a telemetry message from the store
	Delete(*sb.TelemetryMessage) error

	// List all of the last up to date telemetry messages
	List(ch chan<- sb.TelemetryMessage) error

	// Watch watches telemetry updates
	Watch(ch chan<- sb.TelemetryMessage, opts ...WatchOption) error

	// Delete a telemetry message based on a given ID

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

// atomixStore is an Atomix based store
type atomixStore struct {
	telemetry _map.Map
}

func (s *atomixStore) Get(id ID) (*sb.TelemetryMessage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	entry, err := s.telemetry.Get(ctx, id.String())
	if err != nil {
		return nil, err
	} else if entry.Value == nil {
		return nil, nil
	}
	telemetry := &sb.TelemetryMessage{}
	if err := proto.Unmarshal(entry.Value, telemetry); err != nil {
		return nil, err
	}
	return telemetry, nil
}

func (s *atomixStore) Put(telemetry *sb.TelemetryMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	bytes, err := proto.Marshal(telemetry)
	if err != nil {
		return err
	}
	_, err = s.telemetry.Put(ctx, getKey(telemetry).String(), bytes)
	if err != nil {
		return err
	}
	return nil
}

func (s *atomixStore) Delete(telemetry *sb.TelemetryMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err := s.telemetry.Remove(ctx, getKey(telemetry).String())
	if err != nil {
		return err
	}
	return nil
}

func (s *atomixStore) List(ch chan<- sb.TelemetryMessage) error {
	entryCh := make(chan *_map.Entry)
	if err := s.telemetry.Entries(context.Background(), entryCh); err != nil {
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

func (s *atomixStore) Watch(ch chan<- sb.TelemetryMessage, opts ...WatchOption) error {
	options := &watchOptions{}
	for _, opt := range opts {
		opt.apply(options)
	}
	watchOptions := []_map.WatchOption{}
	if options.replay {
		watchOptions = append(watchOptions, _map.WithReplay())
	}

	watchCh := make(chan *_map.Event)
	if err := s.telemetry.Watch(context.Background(), watchCh, watchOptions...); err != nil {
		return err
	}
	go func() {
		defer close(ch)
		for event := range watchCh {
			if event.Type != _map.EventRemoved {
				telemetry := &sb.TelemetryMessage{}
				if err := proto.Unmarshal(event.Entry.Value, telemetry); err == nil {
					ch <- *telemetry
				}
			}
		}
	}()
	return nil
}

func (s *atomixStore) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return s.telemetry.Close(ctx)
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

func getKey(telemetry *sb.TelemetryMessage) ID {
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
	id := ID{
		PlmnID:      ecgi.PlmnId,
		Ecid:        ecgi.Ecid,
		Crnti:       crnti,
		MessageType: telemetry.GetMessageType(),
	}
	return id
}
