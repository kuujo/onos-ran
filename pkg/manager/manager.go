// Copyright 2019-present Open Networking Foundation.
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

// Package manager is is the main coordinator for the ONOS RAN subsystem.
package manager

import (
	"github.com/onosproject/onos-ran/api/sb"
	"github.com/onosproject/onos-ran/pkg/southbound"
	"github.com/onosproject/onos-ran/pkg/store/telemetry"
	"github.com/onosproject/onos-ran/pkg/store/updates"

	log "k8s.io/klog"
)

var mgr Manager

// NewManager initializes the RAN subsystem.
func NewManager() (*Manager, error) {
	log.Info("Creating Manager")

	updatesStore, err := updates.NewStore()
	if err != nil {
		return nil, err
	}

	telemetryStore, err := telemetry.NewStore()
	if err != nil {
		return nil, err
	}

	sbSession, err := southbound.NewSessions()
	if err != nil {
		return nil, err
	}

	mgr = Manager{
		updatesStore,
		telemetryStore,
		sbSession,
		make(chan sb.ControlUpdate),
		make(chan sb.ControlResponse),
	}

	go mgr.recvUpdates()

	return &mgr, nil
}

// Manager single point of entry for the RAN system.
type Manager struct {
	updatesStore   updates.Store
	telemetryStore telemetry.Store
	SB             *southbound.Sessions
	updates        chan sb.ControlUpdate
	responses      chan sb.ControlResponse
}

func (m *Manager) recvUpdates() {
	for update := range mgr.updates {
		_ = m.updatesStore.Put(update)
		log.Infof("Got messageType %d", update.MessageType)
		switch x := update.S.(type) {
		case *sb.ControlUpdate_CellConfigReport:
			log.Infof("plmnid:%s, ecid:%s", x.CellConfigReport.Ecgi.PlmnId, x.CellConfigReport.Ecgi.Ecid)
		case *sb.ControlUpdate_UEAdmissionRequest:
			log.Infof("plmnid:%s, ecid:%s, crnti:%s", x.UEAdmissionRequest.Ecgi.PlmnId, x.UEAdmissionRequest.Ecgi.Ecid, x.UEAdmissionRequest.Crnti)
		default:
			log.Fatalf("ControlReport has unexpected type %T", x)
		}
	}
}

// GetControlUpdates gets a control update based on a given ID
func (m *Manager) GetControlUpdates() ([]sb.ControlUpdate, error) {
	return m.updatesStore.List(), nil
}

// GetTelemetry gets telemeter messages
func (m *Manager) GetTelemetry() ([]sb.TelemetryMessage, error) {
	return m.telemetryStore.List(), nil
}

// Run starts a synchronizer based on the devices and the northbound services.
func (m *Manager) Run() {
	log.Info("Starting Manager")
	m.SB.Run(m.updates, m.responses)
}

//Close kills the channels and manager related objects
func (m *Manager) Close() {
	log.Info("Closing Manager")
}

// GetManager returns the initialized and running instance of manager.
// Should be called only after NewManager and Run are done.
func GetManager() *Manager {
	return &mgr
}
