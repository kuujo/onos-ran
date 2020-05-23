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
	"errors"
	"github.com/google/uuid"
	"github.com/onosproject/onos-lib-go/pkg/cluster"
	"github.com/onosproject/onos-ric/api/store/indications"
	"github.com/onosproject/onos-ric/pkg/store/device"
	"github.com/onosproject/onos-ric/pkg/store/mastership"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
)

// newDeviceIndicationsStore creates a new indications store for a single device
func newDeviceIndicationsStore(deviceKey device.Key, cluster cluster.Cluster, election mastership.Election) (*deviceIndicationsStore, error) {
	store := &deviceIndicationsStore{
		deviceKey:   deviceKey,
		cluster:     cluster,
		election:    election,
		cache:       make(map[int]*Indication),
		subscribers: make([]chan<- Event, 0),
		streams:     make(map[uuid.UUID]indications.IndicationsService_SubscribeServer),
		recordCh:    make(chan Indication),
	}
	if err := store.open(); err != nil {
		return nil, err
	}
	return store, nil
}

// deviceIndicationsStore is a store of indications for a single device
type deviceIndicationsStore struct {
	deviceKey   device.Key
	cluster     cluster.Cluster
	election    mastership.Election
	cache       map[int]*Indication
	subscribers []chan<- Event
	state       *mastership.State
	streams     map[uuid.UUID]indications.IndicationsService_SubscribeServer
	recordCh    chan Indication
	mu          sync.RWMutex
}

func (s *deviceIndicationsStore) open() error {
	ch := make(chan mastership.State)
	if err := s.election.Watch(ch); err != nil {
		return err
	}
	go func() {
		for state := range ch {
			s.mu.Lock()
			s.state = &state
			s.mu.Unlock()
		}
	}()
	go func() {
		for indication := range s.recordCh {
			s.mu.RLock()
			if s.state != nil {
				for _, stream := range s.streams {
					err := stream.Send(&indications.SubscribeResponse{
						Device:     string(s.deviceKey),
						Term:       uint64(s.state.Term),
						Indication: indication.RicIndication,
					})
					if err != nil {
						log.Errorf("Failed to send indication: %s", err)
					}
				}
			}
			s.mu.RUnlock()
		}
	}()
	return nil
}

func (s *deviceIndicationsStore) Record(indication *Indication) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.state == nil || s.state.Master != s.cluster.Node().ID {
		return errors.New("not the master")
	}
	s.recordCh <- *indication
	return nil
}

func (s *deviceIndicationsStore) List(ch chan<- Indication) error {
	go func() {
		s.mu.RLock()
		defer s.mu.RUnlock()
		for _, indication := range s.cache {
			ch <- *indication
		}
	}()
	return nil
}

func (s *deviceIndicationsStore) Subscribe(ch chan<- Event, opts ...SubscribeOption) error {
	options := &subscribeOptions{}
	for _, opt := range opts {
		opt.applySubscribe(options)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.subscribers = append(s.subscribers, ch)

	if options.replay {
		go func() {
			s.mu.RLock()
			defer s.mu.RUnlock()
			for _, indication := range s.cache {
				ch <- Event{
					Type:       EventNone,
					Indication: *indication,
				}
			}
		}()
	}
	return nil
}

func (s *deviceIndicationsStore) subscribe(request *indications.SubscribeRequest, stream indications.IndicationsService_SubscribeServer) error {
	s.mu.RLock()
	state := s.state
	s.mu.RUnlock()

	if state == nil {
		return status.Error(codes.NotFound, "not the master")
	} else if state.Master != s.cluster.Node().ID || state.Term != mastership.Term(request.Term) {
		return status.Error(codes.NotFound, "not the master for the given term")
	}

	id := uuid.New()

	s.mu.Lock()
	s.streams[id] = stream
	s.mu.Unlock()

	<-stream.Context().Done()

	s.mu.Lock()
	delete(s.streams, id)
	s.mu.Unlock()

	return nil
}

func (s *deviceIndicationsStore) Close() error {
	return nil
}
