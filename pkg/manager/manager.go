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
func NewManager(topoEndPoint string, enableMetrics bool, opts []grpc.DialOption) (*Manager, error) {
	log.Info("Creating Manager")

	updatesStore, err := updates.NewDistributedStore()
	if err != nil {
		return nil, err
	}
	if err = updatesStore.Clear(); err != nil {
		log.Error("Error clearing Updates store %s", err.Error())
	}

	telemetryStore, err := telemetry.NewDistributedStore()
	if err != nil {
		return nil, err
	}

	// Should always clear out the stores on startup because it will be out of sync with ran-simulator
	if err = telemetryStore.Clear(); err != nil {
		log.Error("Error clearing Telemetry store %s", err.Error())
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
		enableMetrics:      enableMetrics,
		topoMonitor: monitor.NewTopoMonitorBuilder().
			SetTopoChannel(make(chan *topodevice.ListResponse)).
			Build(),
		dispatcher:       newDispatcher(),
		telemetryChannel: make(chan Event),
	}
	return &mgr, nil
}

// Manager single point of entry for the RAN system.
type Manager struct {
	updatesStore       updates.Store
	telemetryStore     telemetry.Store
	deviceChangesStore device.Store
	SbSessions         map[sb.ECGI]*southbound.Session
	topoMonitor        monitor.TopoMonitor
	enableMetrics      bool
	dispatcher         *Dispatcher
	telemetryChannel   chan Event
}

// StoreControlUpdate - put the control update in the atomix store
func (m *Manager) StoreControlUpdate(update sb.ControlUpdate) {
	log.Infof("Got messageType %s", update.MessageType)
	switch x := update.S.(type) {
	case *sb.ControlUpdate_CellConfigReport:
		log.Infof("plmnid:%s, ecid:%s", x.CellConfigReport.Ecgi.PlmnId, x.CellConfigReport.Ecgi.Ecid)
		_ = m.updatesStore.Put(&update)
	case *sb.ControlUpdate_UEAdmissionRequest:
		log.Infof("plmnid:%s, ecid:%s, crnti:%s, imsi:%d", x.UEAdmissionRequest.Ecgi.PlmnId, x.UEAdmissionRequest.Ecgi.Ecid, x.UEAdmissionRequest.Crnti, x.UEAdmissionRequest.Imsi)
		_ = m.updatesStore.Put(&update)
	case *sb.ControlUpdate_UEReleaseInd:
		log.Infof("delete ue - plmnid:%s, ecid:%s, crnti:%s", x.UEReleaseInd.Ecgi.PlmnId, x.UEReleaseInd.Ecgi.Ecid, x.UEReleaseInd.Crnti)

		err := m.DeleteTelemetry(x.UEReleaseInd.Ecgi.PlmnId, x.UEReleaseInd.Ecgi.Ecid, x.UEReleaseInd.Crnti)
		if err != nil {
			log.Errorf("%s", err)
		}
		err = m.DeleteUEAdmissionRequest(x.UEReleaseInd.Ecgi.PlmnId, x.UEReleaseInd.Ecgi.Ecid, x.UEReleaseInd.Crnti)
		if err != nil {
			log.Errorf("%s", err)
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

// GetUEAdmissionByID retrieve a single value from the updates store
func (m *Manager) GetUEAdmissionByID(ecgi *sb.ECGI, crnti string) (*sb.ControlUpdate, error) {
	id := updates.ID{
		PlmnID:      ecgi.PlmnId,
		Ecid:        ecgi.Ecid,
		Crnti:       crnti,
		MessageType: sb.MessageType_UE_ADMISSION_REQUEST,
	}
	return m.updatesStore.Get(id)
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
func (m *Manager) SubscribeTelemetry(ch chan<- sb.TelemetryMessage, withReplay bool) error {
	if withReplay {
		return m.telemetryStore.Watch(ch, telemetry.WithReplay())
	}
	return m.telemetryStore.Watch(ch)
}

// RegisterTelemetryListener :
func (m *Manager) RegisterTelemetryListener(name string) (chan Event, error) {
	return m.dispatcher.registerTelemetryListener(name)
}

// UnregisterTelemetryListener :
func (m *Manager) UnregisterTelemetryListener(name string) {
	m.dispatcher.unregisterTelemetryListener(name)
}

// Run starts a synchronizer based on the devices and the northbound services.
func (m *Manager) Run() {
	log.Info("Starting Manager")

	m.topoMonitor.TopoEventHandler(m.topoEventHandler)

	err := mgr.deviceChangesStore.Watch(m.topoMonitor.TopoChannel())
	if err != nil {
		log.Errorf("Error listening to topo service: %s", err.Error())
	}

	go m.dispatcher.listenTelemetryEvents(m.telemetryChannel)
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
				session.ControlUpdateHandlerFunc = m.StoreControlUpdate
				session.TelemetryUpdateHandlerFunc = m.StoreTelemetry
				session.EnableMetrics = m.enableMetrics
				m.SbSessions[ecgi] = session
				session.Run(device.GetDevice().GetTLS(), device.GetDevice().GetCredentials())
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

// StoreTelemetry - put the telemetry update in the atomix store
// Only handles MessageType_RADIO_MEAS_REPORT_PER_UE at the moment
func (m *Manager) StoreTelemetry(update sb.TelemetryMessage) {
	m.telemetryChannel <- Event{
		Type:   update.GetMessageType().String(),
		Object: update,
	}

	err := m.telemetryStore.Put(&update)
	if err != nil {
		log.Fatalf("Could not put message %v in telemetry store %s", update, err.Error())
	}

	switch update.MessageType {
	case sb.MessageType_RADIO_MEAS_REPORT_PER_UE:
		x, ok := update.S.(*sb.TelemetryMessage_RadioMeasReportPerUE)
		if !ok {
			log.Fatalf("Telemetry update has unexpected type %T", x)
		}
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
		log.Fatalf("Telemetry update has unexpected type %T", update.MessageType)
	}
}

// DeleteTelemetry deletes telemetry when a handover happens
func (m *Manager) DeleteTelemetry(plmnid string, ecid string, crnti string) error {
	id := telemetry.ID{
		PlmnID:      plmnid,
		Ecid:        ecid,
		Crnti:       crnti,
		MessageType: sb.MessageType_RADIO_MEAS_REPORT_PER_UE,
	}.String()
	if err := m.telemetryStore.DeleteWithKey(id); err != nil {
		log.Infof("Error deleting Telemetry, key=%s", id)
		return err
	}
	return nil
}

// DeleteUEAdmissionRequest deletes UpdateControls
func (m *Manager) DeleteUEAdmissionRequest(plmnid string, ecid string, crnti string) error {
	id := updates.ID{
		PlmnID:      plmnid,
		Ecid:        ecid,
		Crnti:       crnti,
		MessageType: sb.MessageType_UE_ADMISSION_REQUEST,
	}.String()
	if err := m.updatesStore.DeleteWithKey(id); err != nil {
		log.Infof("Error deleting UEAdmissionRequest, key=%s", id)
		return err
	}
	return nil
}

// ecgiFromTopoID topo device is formatted like "001001-0001786" PlmnId:Ecid
func ecgiFromTopoID(id topodevice.ID) sb.ECGI {
	parts := strings.Split(string(id), "-")
	return sb.ECGI{Ecid: parts[1], PlmnId: parts[0]}
}
