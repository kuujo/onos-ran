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

package indications

import (
	"context"
	"fmt"
	"github.com/atomix/go-client/pkg/client/map"
	"github.com/atomix/go-client/pkg/client/primitive"
	"github.com/atomix/go-client/pkg/client/util/net"
	"github.com/gogo/protobuf/proto"
	"github.com/onosproject/onos-lib-go/pkg/atomix"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-ric/api/sb"
	"github.com/onosproject/onos-ric/api/sb/e2ap"
	"github.com/onosproject/onos-ric/pkg/config"
	"io"
	"time"
)

var log = logging.GetLogger("store", "indications")

const requestTimeout = 15 * time.Second

const primitiveName = "indications"
const databaseType = atomix.DatabaseTypeCache

// Revision is a message revision number
type Revision uint64

// ID is a message identifier
type ID string

// NewID creates a new telemetry store ID
func NewID(messageType sb.MessageType, plmnidn, ecid, crnti string) ID {
	return ID(fmt.Sprintf("%s:%s:%s:%s", messageType, plmnidn, ecid, crnti))
}

// GetID gets the telemetry store ID for the given message
func GetID(message *e2ap.RicIndication) ID {
	var ecgi sb.ECGI
	var crnti string
	msgType := message.GetHdr().GetMessageType()
	switch msgType {
	case sb.MessageType_RADIO_MEAS_REPORT_PER_UE:
		ecgi = *message.GetMsg().GetRadioMeasReportPerUE().GetEcgi()
		crnti = message.GetMsg().GetRadioMeasReportPerUE().GetCrnti()
	case sb.MessageType_RADIO_MEAS_REPORT_PER_CELL:
		ecgi = *message.GetMsg().GetRadioMeasReportPerCell().GetEcgi()
	case sb.MessageType_CELL_CONFIG_REPORT:
		ecgi = *message.GetMsg().GetCellConfigReport().GetEcgi()
	case sb.MessageType_UE_ADMISSION_REQUEST:
		ecgi = *message.GetMsg().GetUEAdmissionRequest().GetEcgi()
		crnti = message.GetMsg().GetUEAdmissionRequest().GetCrnti()
	}
	return NewID(msgType, ecgi.PlmnId, ecgi.Ecid, crnti)
}

// Event is a store event
type Event struct {
	// Type is the event type
	Type EventType
	// ID is the message identifier
	ID ID
	// Message is the event message
	Message e2ap.RicIndication
}

// EventType is a store event type
type EventType string

const (
	// EventNone indicates an event that was not triggered but replayed
	EventNone EventType = ""
	// EventInsert indicates the message was inserted into the store
	EventInsert EventType = "insert"
	// EventUpdate indicates the message was updated in the store
	EventUpdate EventType = "update"
	// EventDelete indicates the message was deleted from the store
	EventDelete EventType = "remove"
)

// Store is interface for store
type Store interface {
	io.Closer

	// Gets a message based on a given ID
	Get(ID, ...GetOption) (*e2ap.RicIndication, error)

	// Puts a message to the store
	Put(ID, *e2ap.RicIndication, ...PutOption) error

	// Removes a message from the store
	Delete(ID, ...DeleteOption) error

	// List all of the last up to date messages
	List(ch chan<- e2ap.RicIndication) error

	// Watch watches the store for changes
	Watch(ch chan<- Event, opts ...WatchOption) error

	// Clear deletes all messages from the store
	Clear() error
}

// NewDistributedStore creates a new distributed indications store
func NewDistributedStore(config config.Config) (Store, error) {
	log.Info("Creating distributed message store")
	database, err := atomix.GetDatabase(config.Atomix, config.Atomix.GetDatabase(databaseType))
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	messages, err := database.GetMap(ctx, primitiveName)
	if err != nil {
		return nil, err
	}
	return &store{
		dist: messages,
	}, nil
}

// NewLocalStore returns a new local indications store
func NewLocalStore() (Store, error) {
	_, address := atomix.StartLocalNode()
	return newLocalStore(address)
}

// newLocalStore creates a new local indications store
func newLocalStore(address net.Address) (Store, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	session, err := primitive.NewSession(ctx, primitive.Partition{ID: 1, Address: address})
	if err != nil {
		return nil, err
	}
	primitiveName := primitive.Name{
		Namespace: "local",
		Name:      primitiveName,
	}
	messages, err := _map.New(context.Background(), primitiveName, []*primitive.Session{session})
	if err != nil {
		return nil, err
	}
	return &store{
		dist: messages,
	}, nil
}

var _ Store = &store{}

type store struct {
	dist _map.Map
}

func (s *store) Get(id ID, opts ...GetOption) (*e2ap.RicIndication, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	entry, err := s.dist.Get(ctx, string(id))
	if err != nil {
		return nil, err
	} else if entry == nil {
		return nil, nil
	}
	return decode(entry.Value)
}

func (s *store) Put(id ID, message *e2ap.RicIndication, opts ...PutOption) error {
	putOpts := &putOptions{}
	for _, opt := range opts {
		opt.applyPut(putOpts)
	}
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	bytes, err := encode(message)
	if err != nil {
		return err
	}
	options := []_map.PutOption{}
	if putOpts.revision > 0 {
		options = append(options, _map.IfVersion(int64(putOpts.revision)))
	}
	_, err = s.dist.Put(ctx, string(id), bytes, options...)
	return err
}

func (s *store) Delete(id ID, opts ...DeleteOption) error {
	deleteOpts := &deleteOptions{}
	for _, opt := range opts {
		opt.applyDelete(deleteOpts)
	}
	options := []_map.RemoveOption{}
	if deleteOpts.revision > 0 {
		options = append(options, _map.IfVersion(int64(deleteOpts.revision)))
	}
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	_, err := s.dist.Remove(ctx, string(id), options...)
	return err
}

func (s *store) List(ch chan<- e2ap.RicIndication) error {
	entryCh := make(chan *_map.Entry)
	if err := s.dist.Entries(context.Background(), entryCh); err != nil {
		return err
	}
	go func() {
		defer close(ch)
		for entry := range entryCh {
			message, err := decode(entry.Value)
			if err == nil {
				ch <- *message
			}
		}
	}()
	return nil
}

func (s *store) Watch(ch chan<- Event, opts ...WatchOption) error {
	watchOpts := &watchOptions{}
	for _, opt := range opts {
		opt.applyWatch(watchOpts)
	}

	options := []_map.WatchOption{}
	if watchOpts.replay {
		options = append(options, _map.WithReplay())
	}

	eventCh := make(chan *_map.Event)
	err := s.dist.Watch(context.Background(), eventCh, options...)
	if err != nil {
		return err
	}
	go func() {
		for event := range eventCh {
			message, err := decode(event.Entry.Value)
			if err == nil {
				var eventType EventType
				switch event.Type {
				case _map.EventNone:
					eventType = EventNone
				case _map.EventInserted:
					eventType = EventInsert
				case _map.EventUpdated:
					eventType = EventUpdate
				case _map.EventRemoved:
					eventType = EventDelete
				}
				ch <- Event{
					Type:    eventType,
					Message: *message,
				}
			}
		}
	}()
	return nil
}

func (s *store) Clear() error {
	return s.dist.Clear(context.Background())
}

func (s *store) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	return s.dist.Close(ctx)
}

func decode(bytes []byte) (*e2ap.RicIndication, error) {
	m := &e2ap.RicIndication{}
	if err := proto.Unmarshal(bytes, m); err != nil {
		return nil, err
	}
	return m, nil
}

func encode(m *e2ap.RicIndication) ([]byte, error) {
	return proto.Marshal(m)
}
