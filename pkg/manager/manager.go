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
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-ric/api/sb"
	"github.com/onosproject/onos-ric/pkg/southbound"
	"github.com/onosproject/onos-ric/pkg/southbound/monitor"
	"github.com/onosproject/onos-ric/pkg/store/device"
	"github.com/onosproject/onos-ric/pkg/store/telemetry"
	"github.com/onosproject/onos-ric/pkg/store/updates"
	topodevice "github.com/onosproject/onos-topo/api/device"
	"google.golang.org/grpc"
	"strings"
)

const ranSimulatorType = topodevice.Type("RanSimulator")

var log = logging.GetLogger("manager")
var mgr Manager

// NewManager initializes the RAN subsystem.
func NewManager(topoEndPoint string, opts []grpc.DialOption) (*Manager, error) {
	log.Info("Creating Manager")

	updatesStore, err := updates.NewDistributedStore()
	if err != nil {
		return nil, err
	}

	telemetryStore, err := telemetry.NewDistributedStore()
	if err != nil {
		return nil, err
	}

	deviceChangeStore, err := device.NewTopoStore(topoEndPoint, opts...)
	if err != nil {
		log.Info("Error in device change store")
		return nil, err
	}

	mgr = Manager{
		updatesStore:       updatesStore,
		telemetryStore:     telemetryStore,
		deviceChangesStore: deviceChangeStore,
		SbSessions:         make(map[sb.ECGI]*southbound.Session),
		controlUpdates:     make(chan sb.ControlUpdate),
		controlResponses:   make(chan sb.ControlResponse),
		telemetryUpdates:   make(chan sb.TelemetryMessage),
		topoMonitor: monitor.NewTopoMonitorBuilder().
			SetTopoChannel(make(chan *topodevice.ListResponse)).
			Build(),
	}

	go mgr.recvUpdates()

	go mgr.recvTelemetryUpdates()

	return &mgr, nil
}

// Manager single point of entry for the RAN system.
type Manager struct {
	updatesStore       updates.Store
	telemetryStore     telemetry.Store
	deviceChangesStore device.Store
	SbSessions         map[sb.ECGI]*southbound.Session
	controlUpdates     chan sb.ControlUpdate
	controlResponses   chan sb.ControlResponse
	telemetryUpdates   chan sb.TelemetryMessage
	topoMonitor        monitor.TopoMonitor
}

func (m *Manager) recvUpdates() {
	for update := range mgr.controlUpdates {
		m.storeControlUpdate(update)
	}
}

func (m *Manager) storeControlUpdate(update sb.ControlUpdate) {
	log.Infof("Got messageType %d", update.MessageType)
	switch x := update.S.(type) {
	case *sb.ControlUpdate_CellConfigReport:
		log.Infof("plmnid:%s, ecid:%s", x.CellConfigReport.Ecgi.PlmnId, x.CellConfigReport.Ecgi.Ecid)
		_ = m.updatesStore.Put(&update)
	case *sb.ControlUpdate_UEAdmissionRequest:
		log.Infof("plmnid:%s, ecid:%s, crnti:%s", x.UEAdmissionRequest.Ecgi.PlmnId, x.UEAdmissionRequest.Ecgi.Ecid, x.UEAdmissionRequest.Crnti)
		_ = m.updatesStore.Put(&update)
	case *sb.ControlUpdate_UEReleaseInd:
		log.Infof("delete ue - plmnid:%s, ecid:%s, crnti:%s", x.UEReleaseInd.Ecgi.PlmnId, x.UEReleaseInd.Ecgi.Ecid, x.UEReleaseInd.Crnti)
		// Delete UE's UEAdmissionRequest record from store
		key := updates.ID{
			PlmnID:      x.UEReleaseInd.Ecgi.PlmnId,
			Ecid:        x.UEReleaseInd.Ecgi.Ecid,
			Crnti:       x.UEReleaseInd.Crnti,
			MessageType: sb.MessageType_UE_ADMISSION_REQUEST,
		}.String()
		if err := m.updatesStore.DeleteWithKey(key); err != nil {
			log.Infof("Error deleting UEAdmissionRequest, key=%s", key)
		}
		// Delete UE's RadioMeasReportPerUE record from store
		key = updates.ID{
			PlmnID:      x.UEReleaseInd.Ecgi.PlmnId,
			Ecid:        x.UEReleaseInd.Ecgi.Ecid,
			Crnti:       x.UEReleaseInd.Crnti,
			MessageType: sb.MessageType_RADIO_MEAS_REPORT_PER_UE,
		}.String()
		if err := m.telemetryStore.DeleteWithKey(key); err != nil {
			log.Infof("Error deleting UEAdmissionRequest, key=%s", key)
		}
	default:
		log.Fatalf("ControlReport has unexpected type %T", x)
	}
}

// GetControlUpdates gets control updates
func (m *Manager) GetControlUpdates() ([]sb.ControlUpdate, error) {
	ch := make(chan sb.ControlUpdate)
	if err := m.ListControlUpdates(ch); err != nil {
		return nil, err
	}
	messages := make([]sb.ControlUpdate, 0)
	for update := range ch {
		messages = append(messages, update)
	}
	return messages, nil
}

// ListControlUpdates lists control updates
func (m *Manager) ListControlUpdates(ch chan<- sb.ControlUpdate) error {
	return m.updatesStore.List(ch)
}

// SubscribeControlUpdates subscribes the given channel to control updates
func (m *Manager) SubscribeControlUpdates(ch chan<- sb.ControlUpdate) error {
	return m.updatesStore.Watch(ch, updates.WithReplay())
}

// GetTelemetry gets telemeter messages
func (m *Manager) GetTelemetry() ([]sb.TelemetryMessage, error) {
	ch := make(chan sb.TelemetryMessage)
	if err := m.ListTelemetry(ch); err != nil {
		return nil, err
	}
	messages := make([]sb.TelemetryMessage, 0)
	for telemetry := range ch {
		messages = append(messages, telemetry)
	}
	return messages, nil
}

// ListTelemetry lists telemeter messages
func (m *Manager) ListTelemetry(ch chan<- sb.TelemetryMessage) error {
	return m.telemetryStore.List(ch)
}

// SubscribeTelemetry subscribes the given channel to telemetry events
func (m *Manager) SubscribeTelemetry(ch chan<- sb.TelemetryMessage) error {
	return m.telemetryStore.Watch(ch, telemetry.WithReplay())
}

// Run starts a synchronizer based on the devices and the northbound services.
func (m *Manager) Run() {
	log.Info("Starting Manager")

	m.topoMonitor.TopoEventHandler(m.topoEventHandler)

	err := mgr.deviceChangesStore.Watch(m.topoMonitor.TopoChannel())
	if err != nil {
		log.Errorf("Error listening to topo service: %s", err.Error())
	}
}

func (m *Manager) topoEventHandler(topoChannel chan *topodevice.ListResponse) {
	log.Infof("Watching topo channel")
	for device := range topoChannel {
		log.Infof("Device received %s", device.GetDevice().GetID())
		if device.GetDevice().GetType() != ranSimulatorType {
			continue
		}
		if device.Type == topodevice.ListResponse_NONE || device.Type == topodevice.ListResponse_ADDED {
			ecgi := ecgiFromTopoID(device.GetDevice().GetID())
			deviceEndpoint := sb.Endpoint(device.GetDevice().GetAddress())
			session, err := southbound.NewSession(ecgi, deviceEndpoint)
			if err != nil {
				log.Fatalf("Unable to create new session %s", err.Error())
			}
			if session != nil {
				m.SbSessions[ecgi] = session
				session.Run(device.GetDevice().GetTLS(), device.GetDevice().GetCredentials(),
					m.controlUpdates, m.controlResponses, m.telemetryUpdates)
			} else {
				log.Fatalf("Error creating new session for %v", ecgi)
			}
		} else {
			log.Warnf("Topo device event not yet handled %s %v", device.String())
		}
	}
}

//Close kills the channels and manager related objects
func (m *Manager) Close() {
	m.topoMonitor.Close()
	log.Info("Closing Manager")
}

// GetManager returns the initialized and running instance of manager.
// Should be called only after NewManager and Run are done.
func GetManager() *Manager {
	return &mgr
}

func (m *Manager) recvTelemetryUpdates() {
	for update := range mgr.telemetryUpdates {
		m.storeTelemetry(update)
	}
}

func (m *Manager) storeTelemetry(update sb.TelemetryMessage) {
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

// DeleteTelemetry deletes telemetry when a handover happens
func (m *Manager) DeleteTelemetry(plmnid string, ecid string, crnti string) error {
	id := telemetry.ID{
		PlmnID:      plmnid,
		Ecid:        ecid,
		Crnti:       crnti,
		MessageType: sb.MessageType_RADIO_MEAS_REPORT_PER_UE,
	}
	telemetryMsg, err := m.telemetryStore.Get(id)
	if err != nil {
		return err
	}
	err = m.telemetryStore.Delete(telemetryMsg)
	if err != nil {
		return err
	}

	return nil

}

// ecgiFromTopoID topo device is formatted like "001001-0001786" PlmnId:Ecid
func ecgiFromTopoID(id topodevice.ID) sb.ECGI {
	parts := strings.Split(string(id), "-")
	return sb.ECGI{Ecid: parts[1], PlmnId: parts[0]}
}
