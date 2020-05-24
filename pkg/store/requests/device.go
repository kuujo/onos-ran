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
	"github.com/onosproject/onos-lib-go/pkg/cluster"
	"github.com/onosproject/onos-ric/api/store/requests"
	"github.com/onosproject/onos-ric/pkg/store/device"
	"github.com/onosproject/onos-ric/pkg/store/mastership"
	"sync"
)

// newDeviceIndicationsStore creates a new indications store for a single device
func newDeviceIndicationsStore(deviceKey device.Key, cluster cluster.Cluster, election mastership.Election) (*deviceRequestsStore, error) {
	store := &deviceRequestsStore{
		deviceKey: deviceKey,
		cluster:   cluster,
		election:  election,
		log:       newLog(),
	}
	if err := store.open(); err != nil {
		return nil, err
	}
	return store, nil
}

// deviceRequestsStore is a store of requests for a single device
type deviceRequestsStore struct {
	deviceKey device.Key
	cluster   cluster.Cluster
	election  mastership.Election
	state     *mastership.State
	handler   storeHandler
	log       Log
	mu        sync.RWMutex
}

func (s *deviceRequestsStore) open() error {
	ch := make(chan mastership.State)
	if err := s.election.Watch(ch); err != nil {
		return err
	}
	go s.processElectionChanges(ch)
	return nil
}

func (s *deviceRequestsStore) processElectionChanges(ch <-chan mastership.State) {
	for state := range ch {
		s.mu.Lock()
		s.state = &state
		if state.Master == s.cluster.Node().ID {
			if s.handler != nil {
				s.handler.Close()
			}
			handler, err := newMasterStore(s.deviceKey, s.cluster, state, s.log)
			if err != nil {
				logger.Errorf("Failed to initialize master store: %s", err)
			} else {
				s.handler = handler
			}
		} else if state.Master != s.cluster.Node().ID {
			if s.handler != nil {
				s.handler.Close()
			}
			handler, err := newBackupStore(s.deviceKey, s.cluster, state, s.log)
			if err != nil {
				logger.Errorf("Failed to initialize backup store: %s", err)
			} else {
				s.handler = handler
			}
		}
		s.mu.Unlock()
	}
}

func (s *deviceRequestsStore) Append(request *Request, opts ...AppendOption) error {
	s.mu.RLock()
	handler := s.handler
	s.mu.RUnlock()
	if handler == nil {
		return errors.New("no master found")
	}
	return handler.Append(request, opts...)
}

func (s *deviceRequestsStore) Ack(request *Request) error {
	s.mu.RLock()
	handler := s.handler
	s.mu.RUnlock()
	if handler == nil {
		return errors.New("no master found")
	}
	return handler.Ack(request)
}

func (s *deviceRequestsStore) Watch(deviceID device.ID, ch chan<- Event, opts ...WatchOption) error {
	options := &watchOptions{}
	for _, opt := range opts {
		opt.applyWatch(options)
	}
	// TODO
	return nil
}

func (s *deviceRequestsStore) append(ctx context.Context, request *requests.AppendRequest) (*requests.AppendResponse, error) {
	s.mu.RLock()
	handler := s.handler
	s.mu.RUnlock()
	if handler == nil {
		return nil, errors.New("no master found")
	}
	return handler.append(ctx, request)
}

func (s *deviceRequestsStore) ack(ctx context.Context, request *requests.AckRequest) (*requests.AckResponse, error) {
	s.mu.RLock()
	handler := s.handler
	s.mu.RUnlock()
	if handler == nil {
		return nil, errors.New("no master found")
	}
	return handler.ack(ctx, request)
}

func (s *deviceRequestsStore) backup(ctx context.Context, request *requests.BackupRequest) (*requests.BackupResponse, error) {
	s.mu.RLock()
	handler := s.handler
	s.mu.RUnlock()
	if handler == nil {
		return nil, errors.New("no master found")
	}
	return handler.backup(ctx, request)
}

func (s *deviceRequestsStore) Close() error {
	return nil
}

var _ storeHandler = &deviceRequestsStore{}
