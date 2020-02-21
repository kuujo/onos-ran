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

	telemetryStore, err := telemetry.NewDistributedStore()
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
		make(chan sb.TelemetryMessage),
	}

	go mgr.recvUpdates()

	go mgr.recvTelemetryUpdates()

	return &mgr, nil
}

// Manager single point of entry for the RAN system.
type Manager struct {
	updatesStore     updates.Store
	telemetryStore   telemetry.Store
	SB               *southbound.Sessions
	controlUpdates   chan sb.ControlUpdate
	controlResponses chan sb.ControlResponse
	telemetryUpdates chan sb.TelemetryMessage
}

func (m *Manager) recvUpdates() {
	for update := range mgr.controlUpdates {
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

// SubscribeControlUpdates subscribes the given channel to control updates
func (m *Manager) SubscribeControlUpdates(ch chan<- sb.ControlUpdate) error {
	return m.updatesStore.Watch(ch, updates.WithReplay())
}

// GetTelemetry gets telemeter messages
func (m *Manager) GetTelemetry(ch chan<- sb.TelemetryMessage) error {
	return m.telemetryStore.List(ch)
}

// SubscribeTelemetry subscribes the given channel to telemetry events
func (m *Manager) SubscribeTelemetry(ch chan<- sb.TelemetryMessage) error {
	return m.telemetryStore.Watch(ch, telemetry.WithReplay())
}

// Run starts a synchronizer based on the devices and the northbound services.
func (m *Manager) Run() {
	log.Info("Starting Manager")
	m.SB.Run(m.controlUpdates, m.controlResponses, m.telemetryUpdates)
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

func (m *Manager) recvTelemetryUpdates() {
	for update := range mgr.telemetryUpdates {
		_ = m.telemetryStore.Put(&update)
		switch x := update.S.(type) {
		case *sb.TelemetryMessage_RadioMeasReportPerUE:
			log.Infof("RadioMeasReport plmnid:%s ecid:%s crnti:%s cqis:%d(ecid:%s),%d(ecid:%s),%d(ecid:%s)",
				x.RadioMeasReportPerUE.Ecgi.PlmnId,
				x.RadioMeasReportPerUE.Ecgi.Ecid,
				x.RadioMeasReportPerUE.Crnti,
				x.RadioMeasReportPerUE.RadioReportServCells[0].CqiHist[0],
				x.RadioMeasReportPerUE.RadioReportServCells[0].GetEcgi().GetEcid(),
				x.RadioMeasReportPerUE.RadioReportServCells[1].CqiHist[0],
				x.RadioMeasReportPerUE.RadioReportServCells[1].GetEcgi().GetEcid(),
				x.RadioMeasReportPerUE.RadioReportServCells[2].CqiHist[0],
				x.RadioMeasReportPerUE.RadioReportServCells[2].GetEcgi().GetEcid(),
			)
		default:
			log.Fatalf("Telemetry update has unexpected type %T", x)
			log.Infof("Got telemetry messageType %d", update.MessageType)
		}
	}
}

// DeleteTelemetry deletes telemetry when a handover happens
func (m *Manager) DeleteTelemetry(plmnid string, ecid string, crnti string) {
	id := telemetry.ID{
		PlmnID:      plmnid,
		Ecid:        ecid,
		Crnti:       crnti,
		MessageType: sb.MessageType_RADIO_MEAS_REPORT_PER_UE,
	}
	err := m.telemetryStore.Delete(id)
	if err != nil {
		log.Error(err)
	}
}
