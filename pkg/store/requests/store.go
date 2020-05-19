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

package requests

import (
	"context"
	"errors"
	"fmt"
	"github.com/atomix/go-client/pkg/client/log"
	"github.com/atomix/go-client/pkg/client/primitive"
	"github.com/atomix/go-client/pkg/client/util/net"
	"github.com/gogo/protobuf/proto"
	"github.com/onosproject/onos-lib-go/pkg/atomix"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-ric/api/sb/e2ap"
	"github.com/onosproject/onos-ric/pkg/config"
	"github.com/onosproject/onos-ric/pkg/store/device"
	"io"
	"sync"
	"time"
)

var logger = logging.GetLogger("store", "requests")

const requestTimeout = 15 * time.Second

const primitiveName = "indications"
const databaseType = atomix.DatabaseTypeCache

// Index is a message index
type Index int64

// New creates a new indication
func New(deviceID device.ID, request *e2ap.RicControlRequest) *Request {
	return &Request{
		DeviceID:          deviceID,
		RicControlRequest: request,
	}
}

// Request is a control request
type Request struct {
	*e2ap.RicControlRequest
	// Index is the request index
	Index Index
	// DeviceID is the request device identifier
	DeviceID device.ID
}

// Event is a store event
type Event struct {
	// Type is the event type
	Type EventType
	// Request is the request
	Request Request
}

// EventType is a store event type
type EventType string

const (
	// EventNone indicates an event that was not triggered but replayed
	EventNone EventType = ""
	// EventAppend indicates the message was appended to the store
	EventAppend EventType = "append"
	// EventRemove indicates the message was removed from the store
	EventRemove EventType = "update"
)

// Store is interface for store
type Store interface {
	io.Closer

	// Append appends a new request
	Append(*Request, ...AppendOption) error

	// Ack acknowledges a request
	Ack(*Request) error

	// Watch watches the store for changes
	Watch(device.ID, chan<- Event, ...WatchOption) error
}

// NewDistributedStore creates a new distributed indications store
func NewDistributedStore(config config.Config) (Store, error) {
	logger.Info("Creating distributed message store")
	database, err := atomix.GetDatabase(config.Atomix, config.Atomix.GetDatabase(databaseType))
	if err != nil {
		return nil, err
	}
	return &store{
		factory: func(id device.ID) (log.Log, error) {
			ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
			defer cancel()
			name := fmt.Sprintf("%s-%s-%s", primitiveName, id.PlmnId, id.Ecid)
			return database.GetLog(ctx, name)
		},
		logs: make(map[device.ID]log.Log),
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
	return &store{
		factory: func(id device.ID) (log.Log, error) {
			logName := primitive.Name{
				Namespace: "local",
				Name:      fmt.Sprintf("%s-%s-%s", primitiveName, id.PlmnId, id.Ecid),
			}
			return log.New(context.Background(), logName, []*primitive.Session{session})
		},
		logs: make(map[device.ID]log.Log),
	}, nil
}

var _ Store = &store{}

type store struct {
	factory func(device.ID) (log.Log, error)
	logs    map[device.ID]log.Log
	mu      sync.RWMutex
}

func (s *store) getLog(deviceID device.ID) (log.Log, error) {
	s.mu.RLock()
	log, ok := s.logs[deviceID]
	s.mu.RUnlock()
	if ok {
		return log, nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	log, ok = s.logs[deviceID]
	if ok {
		return log, nil
	}

	log, err := s.factory(deviceID)
	if err != nil {
		return nil, err
	}
	s.logs[deviceID] = log
	return log, nil
}

func (s *store) Append(request *Request, opts ...AppendOption) error {
	deviceLog, err := s.getLog(request.DeviceID)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	bytes, err := encode(request.RicControlRequest)
	if err != nil {
		return err
	}
	entry, err := deviceLog.Append(ctx, bytes)
	if err != nil {
		return err
	}
	request.Index = Index(entry.Index)
	return nil
}

func (s *store) Ack(request *Request) error {
	deviceLog, err := s.getLog(request.DeviceID)
	if err != nil {
		return err
	}
	if request.Index == 0 {
		return errors.New("cannot acknowledge request at unknown index")
	}
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	_, err = deviceLog.Remove(ctx, int64(request.Index))
	return err
}

func (s *store) Watch(deviceID device.ID, ch chan<- Event, opts ...WatchOption) error {
	watchOpts := &watchOptions{}
	for _, opt := range opts {
		opt.applyWatch(watchOpts)
	}

	options := []log.WatchOption{}
	if watchOpts.replay {
		options = append(options, log.WithReplay())
	}

	deviceLog, err := s.getLog(deviceID)
	if err != nil {
		return err
	}

	eventCh := make(chan *log.Event)
	err = deviceLog.Watch(context.Background(), eventCh, options...)
	if err != nil {
		return err
	}
	go func() {
		for event := range eventCh {
			request, err := decode(event.Entry.Value)
			if err == nil {
				var eventType EventType
				switch event.Type {
				case log.EventNone:
					eventType = EventNone
				case log.EventAppended:
					eventType = EventAppend
				case log.EventRemoved:
					eventType = EventRemove
				}
				ch <- Event{
					Type: eventType,
					Request: Request{
						RicControlRequest: request,
						DeviceID:          deviceID,
						Index:             Index(event.Entry.Index),
					},
				}
			}
		}
	}()
	return nil
}

func (s *store) Close() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var returnErr error
	for _, log := range s.logs {
		ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
		err := log.Close(ctx)
		cancel()
		if err != nil {
			returnErr = err
		}
	}
	return returnErr
}

func decode(bytes []byte) (*e2ap.RicControlRequest, error) {
	m := &e2ap.RicControlRequest{}
	if err := proto.Unmarshal(bytes, m); err != nil {
		return nil, err
	}
	return m, nil
}

func encode(m *e2ap.RicControlRequest) ([]byte, error) {
	return proto.Marshal(m)
}
