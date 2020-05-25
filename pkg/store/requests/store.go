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
	"github.com/onosproject/onos-lib-go/pkg/cluster"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-ric/api/sb/e2ap"
	"github.com/onosproject/onos-ric/api/store/requests"
	"github.com/onosproject/onos-ric/pkg/config"
	"github.com/onosproject/onos-ric/pkg/store/device"
	"github.com/onosproject/onos-ric/pkg/store/mastership"
	"io"
	"sync"
	"time"
)

var logger = logging.GetLogger("store", "requests")

const requestTimeout = 30 * time.Second

// TODO: Make backup count configurable
const syncBackupCount = 1;

// New creates a new indication
func New(request *e2ap.RicControlRequest) *Request {
	return &Request{
		RicControlRequest: request,
	}
}

// Request is a control request
type Request struct {
	*e2ap.RicControlRequest
	// Index is the request index
	Index Index
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
func NewDistributedStore(cluster cluster.Cluster, devices device.Store, masterships mastership.Store, config config.Config) (Store, error) {
	logger.Info("Creating distributed requests store")
	store := &store{
		cluster:        cluster,
		devices:        devices,
		masterships:    masterships,
		deviceRequests: make(map[device.Key]*deviceRequestsStore),
	}
	if err := store.open(); err != nil {
		return nil, err
	}
	return store, nil
}

var _ Store = &store{}

type store struct {
	cluster        cluster.Cluster
	devices        device.Store
	masterships    mastership.Store
	deviceRequests map[device.Key]*deviceRequestsStore
	mu             sync.RWMutex
}

func (s *store) open() error {
	getServer().registerHandler(s)
	ch := make(chan device.Event)
	err := s.devices.Watch(ch)
	if err != nil {
		return err
	}
	go func() {
		for event := range ch {
			if event.Type != device.EventUpdated {
				_, err := s.getDeviceStore(event.Device.ID.Key())
				if err != nil {
					logger.Errorf("Failed to initialize requests store for device %s: %s", event.Device.ID.Key(), err)
				}
			}
		}
	}()
	return nil
}

func (s *store) getDeviceStore(deviceKey device.Key) (*deviceRequestsStore, error) {
	s.mu.RLock()
	store, ok := s.deviceRequests[deviceKey]
	s.mu.RUnlock()
	if !ok {
		s.mu.Lock()
		defer s.mu.Unlock()
		store, ok = s.deviceRequests[deviceKey]
		if !ok {
			election, err := s.masterships.GetElection(mastership.Key(deviceKey))
			if err != nil {
				return nil, err
			}
			store, err := newDeviceIndicationsStore(deviceKey, s.cluster, election)
			if err != nil {
				return nil, err
			}
			s.deviceRequests[deviceKey] = store
		}
	}
	return store, nil
}

func getDeviceID(request *Request) device.ID {
	return device.ID(*request.Hdr.Ecgi)
}

func (s *store) Append(request *Request, opts ...AppendOption) error {
	store, err := s.getDeviceStore(getDeviceID(request).Key())
	if err != nil {
		return err
	}
	return store.Append(request, opts...)
}

func (s *store) Ack(request *Request) error {
	store, err := s.getDeviceStore(getDeviceID(request).Key())
	if err != nil {
		return err
	}
	return store.Ack(request)
}

func (s *store) Watch(deviceID device.ID, ch chan<- Event, opts ...WatchOption) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	wg := &sync.WaitGroup{}
	errCh := make(chan error)
	for _, store := range s.deviceRequests {
		wg.Add(1)
		go func(store Store) {
			err := store.Watch(deviceID, ch, opts...)
			if err != nil {
				errCh <- err
			}
			wg.Done()
		}(store)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		return err
	}
	return nil
}

func (s *store) append(ctx context.Context, request *requests.AppendRequest) (*requests.AppendResponse, error) {
	store, err := s.getDeviceStore(device.Key(request.DeviceID))
	if err != nil {
		return nil, err
	}
	return store.append(ctx, request)
}

func (s *store) ack(ctx context.Context, request *requests.AckRequest) (*requests.AckResponse, error) {
	store, err := s.getDeviceStore(device.Key(request.DeviceID))
	if err != nil {
		return nil, err
	}
	return store.ack(ctx, request)
}

func (s *store) backup(ctx context.Context, request *requests.BackupRequest) (*requests.BackupResponse, error) {
	store, err := s.getDeviceStore(device.Key(request.DeviceID))
	if err != nil {
		return nil, err
	}
	return store.backup(ctx, request)
}

func (s *store) Close() error {
	return nil
}
