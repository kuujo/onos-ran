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
	"fmt"
	"sync"

	"github.com/onosproject/onos-ran/api/sb"
	log "k8s.io/klog"
)

// ID store id
type ID struct {
	MessageType sb.MessageType
	PlmnID      string
	Ecid        string
	Crnti       string
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

// Store is ran store
type Store interface {
	// Gets a control update message based on a given ID
	Get(ID) (sb.ControlUpdate, error)

	// Puts a control update message to the store
	Put(sb.ControlUpdate) error

	// List all of the last up to date control update messages
	List() []sb.ControlUpdate

	// Watch watches the store for control update messages
	Watch(ch chan<- sb.ControlUpdate, opts ...WatchOption) error
}

// processEvents processes update events in order
func (s *updatesStore) processEvents() {
	for event := range s.events {
		s.mu.RLock()
		for _, watcher := range s.watchers {
			watcher <- event
		}
		s.mu.RUnlock()
	}
}

// enqueueEvent enqueues the given event to be propagated to watchers
func (s *updatesStore) enqueueEvent(update sb.ControlUpdate) {
	s.events <- update
}

// Get gets a control update message based on a given ID
func (s *updatesStore) Get(id ID) (sb.ControlUpdate, error) {
	s.mu.RLock()
	if controlUpdate, ok := s.controlUpdates[id]; ok {
		s.mu.RUnlock()
		return controlUpdate, nil
	}
	s.mu.RUnlock()
	return sb.ControlUpdate{}, fmt.Errorf("not found")
}

func getKey(update sb.ControlUpdate) ID {
	var ecgi sb.ECGI
	var crnti string
	switch update.MessageType {
	case sb.MessageType_CELL_CONFIG_REPORT:
		ecgi = *update.GetCellConfigReport().GetEcgi()
	case sb.MessageType_RRM_CONFIG_STATUS:
		ecgi = *update.GetRRMConfigStatus().GetEcgi()
	case sb.MessageType_UE_ADMISSION_REQUEST:
		ecgi = *update.GetUEAdmissionRequest().GetEcgi()
		crnti = update.GetBearerAdmissionRequest().GetCrnti()
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

// Put puts a control update message in the store
func (s *updatesStore) Put(update sb.ControlUpdate) error {
	id := getKey(update)
	s.mu.Lock()
	s.controlUpdates[id] = update
	s.mu.Unlock()
	s.enqueueEvent(update)
	return nil
}

// List gets all of the control update messages in the store
func (s *updatesStore) List() []sb.ControlUpdate {
	var controlUpdates []sb.ControlUpdate
	s.mu.RLock()
	for _, value := range s.controlUpdates {
		controlUpdates = append(controlUpdates, value)
	}
	s.mu.RUnlock()
	return controlUpdates
}

// Watch watches the store for control update messages
func (s *updatesStore) Watch(ch chan<- sb.ControlUpdate, opts ...WatchOption) error {
	options := &watchOptions{}
	for _, opt := range opts {
		opt.apply(options)
	}

	done := make(chan struct{})
	go func() {
		s.mu.Lock()
		if options.replay {
			for _, update := range s.controlUpdates {
				ch <- update
			}
		}
		s.watchers = append(s.watchers, ch)
		s.mu.Unlock()
		close(done)
	}()
	<-done
	return nil
}

// RanStore is responsible for tracking the RAN data.
type updatesStore struct {
	controlUpdates map[ID]sb.ControlUpdate
	events         chan sb.ControlUpdate
	watchers       []chan<- sb.ControlUpdate
	mu             sync.RWMutex
}

// NewStore creates a new Control updates store
func NewStore() (Store, error) {
	log.Info("Creating Control Updates Store")
	store := &updatesStore{
		controlUpdates: make(map[ID]sb.ControlUpdate),
		events:         make(chan sb.ControlUpdate),
		watchers:       make([]chan<- sb.ControlUpdate, 0),
	}
	go store.processEvents()
	return store, nil
}
