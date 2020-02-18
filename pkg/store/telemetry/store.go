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

// Store is interface for telemetry store
type Store interface {
	// Gets a telemetry message based on a given ID
	Get(ID) (sb.TelemetryMessage, error)

	// Puts a telemetry message to the store
	Put(sb.TelemetryMessage) error

	// List all of the last up to date telemetry messages
	List() []sb.TelemetryMessage

	// Delete a telemetry message based on a given ID
	Delete(ID) error
}

// Get gets a telemetry message based on a given ID
func (s *telemetryStore) Get(id ID) (sb.TelemetryMessage, error) {
	s.mu.RLock()
	if telemetry, ok := s.telemetry[id]; ok {
		s.mu.RUnlock()
		return telemetry, nil
	}
	s.mu.RUnlock()
	return sb.TelemetryMessage{}, fmt.Errorf("not found")
}

func getKey(telemetry sb.TelemetryMessage) ID {
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

// Put puts a telemetry  message in the store
func (s *telemetryStore) Put(telemetry sb.TelemetryMessage) error {
	id := getKey(telemetry)
	s.mu.Lock()
	s.telemetry[id] = telemetry
	s.mu.Unlock()
	return nil
}

// List gets all of the telemetry messages in the store
func (s *telemetryStore) List() []sb.TelemetryMessage {
	var telemetryMessages []sb.TelemetryMessage
	s.mu.RLock()
	for _, value := range s.telemetry {
		telemetryMessages = append(telemetryMessages, value)
	}
	s.mu.RUnlock()
	return telemetryMessages
}

func (s *telemetryStore) Delete(id ID) error {
	s.mu.Lock()
	delete(s.telemetry, id)
	s.mu.Unlock()
	return nil
}

// telemetryStore is responsible for tracking the RAN telemetry data
type telemetryStore struct {
	telemetry map[ID]sb.TelemetryMessage
	mu        sync.RWMutex
}

// NewStore creates a new RAN store controller.
func NewStore() (Store, error) {
	log.Info("Creating Telemetry Store")
	return &telemetryStore{
		telemetry: make(map[ID]sb.TelemetryMessage),
	}, nil
}
