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
	"github.com/onosproject/onos-lib-go/pkg/cluster"
	"github.com/onosproject/onos-ric/api/sb/e2ap"
	"github.com/onosproject/onos-ric/api/store/requests"
	"github.com/onosproject/onos-ric/pkg/config"
	"github.com/onosproject/onos-ric/pkg/store/device"
	"github.com/onosproject/onos-ric/pkg/store/mastership"
	"sync"
)

// newDeviceRequestsStore creates a new requests store for a single device
func newDeviceRequestsStore(deviceKey device.Key, cluster cluster.Cluster, election mastership.Election, config config.RequestsStoreConfig) (*deviceRequestsStore, error) {
	store := &deviceRequestsStore{
		config:    config,
		deviceKey: deviceKey,
		cluster:   cluster,
		election:  election,
		state:     &deviceStoreState{},
		log:       newLog(),
	}
	logger.Debugf("Initializing device store %s", deviceKey)
	if err := store.open(); err != nil {
		logger.Errorf("Failed to initialize device store %s: %s", deviceKey, err)
		return nil, err
	}
	return store, nil
}

// deviceRequestsStore is a store of requests for a single device
type deviceRequestsStore struct {
	config    config.RequestsStoreConfig
	deviceKey device.Key
	cluster   cluster.Cluster
	election  mastership.Election
	state     *deviceStoreState
	handler   storeHandler
	log       Log
	mu        sync.RWMutex
}

func (s *deviceRequestsStore) open() error {
	ch := make(chan mastership.State)
	if err := s.election.Watch(ch); err != nil {
		return err
	}
	state, err := s.election.GetState()
	if err != nil {
		return err
	} else if state == nil {
		return errors.New("failed to join mastership election")
	}
	logger.Debugf("Initial election state for %s: %v", s.deviceKey, *state)
	if err := s.processElectionChange(*state); err != nil {
		return err
	}
	go s.processElectionChanges(ch)
	return nil
}

func (s *deviceRequestsStore) processElectionChanges(ch <-chan mastership.State) {
	for state := range ch {
		logger.Debugf("Election state changed for %s: %v", s.deviceKey, state)
		s.mu.Lock()
		err := s.processElectionChange(state)
		if err != nil {
			logger.Errorf("Failed to process election state change for %s: %s", s.deviceKey, err)
		}
		s.mu.Unlock()
	}
}

func (s *deviceRequestsStore) processElectionChange(state mastership.State) error {
	s.state.setMastership(&state)
	if state.Master == s.cluster.Node().ID {
		logger.Debugf("Transitioning to master role for %s", s.deviceKey)
		handler, err := newMasterStore(s.deviceKey, s.cluster, s.state, s.log, s.config)
		if err != nil {
			return fmt.Errorf("failed to initialize master store: %s", err)
		}
		s.handler = handler
	} else {
		var index int
		for i := 1; i < len(state.Replicas); i++ {
			if state.Replicas[i] == cluster.ReplicaID(s.cluster.Node().ID) {
				index = i
			}
		}
		if index != 0 && index <= s.config.Backups+s.config.AsyncBackups {
			logger.Debugf("Transitioning to backup role for %s", s.deviceKey)
			handler, err := newBackupStore(s.deviceKey, s.cluster, s.state, s.log, s.config)
			if err != nil {
				return fmt.Errorf("failed to initialize backup store: %s", err)
			}
			s.handler = handler
		} else {
			logger.Debugf("Transitioning to standby role for %s", s.deviceKey)
			handler, err := newStandbyStore(s.deviceKey, s.cluster, s.state, s.log, s.config)
			if err != nil {
				return fmt.Errorf("failed to initialize standby store: %s", err)
			}
			s.handler = handler
		}
	}
	return nil
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

	var reader Reader
	nextIndex := s.log.Writer().Index() + 1
	if options.replay {
		reader = s.log.OpenReader(0)
	} else {
		reader = s.log.OpenReader(nextIndex)
	}

	go func() {
		for {
			batch := reader.ReadUntil(s.state.getCommitIndex())
			for _, entry := range batch.Entries {
				var eventType EventType
				if entry.Index < nextIndex {
					eventType = EventNone
				} else {
					eventType = EventAppend
				}
				request := New(entry.Value.(*e2ap.RicControlRequest))
				request.Index = entry.Index
				ch <- Event{
					Type:    eventType,
					Request: *request,
				}
			}
		}
	}()
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
