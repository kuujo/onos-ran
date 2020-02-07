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

// Store is ran store
type Store interface {
	// Gets a control update message based on a given ID
	Get(ID) (sb.ControlUpdate, error)

	// Puts a control update message to the store
	Put(sb.ControlUpdate) error

	// List all of the last up to date control update messages
	List() []sb.ControlUpdate
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

// RanStore is responsible for tracking the RAN data.
type updatesStore struct {
	controlUpdates map[ID]sb.ControlUpdate
	mu             sync.RWMutex
}

// NewStore creates a new Control updates store
func NewStore() (Store, error) {
	log.Info("Creating Control Updates Store")
	return &updatesStore{
		controlUpdates: make(map[ID]sb.ControlUpdate),
	}, nil
}
