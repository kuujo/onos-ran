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
	"fmt"
	"strings"

	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-ric/api/sb"
	"github.com/onosproject/onos-ric/api/sb/e2ap"
	"github.com/onosproject/onos-ric/pkg/config"
	"github.com/onosproject/onos-ric/pkg/southbound"
	"github.com/onosproject/onos-ric/pkg/southbound/monitor"
	"github.com/onosproject/onos-ric/pkg/store/device"
	"github.com/onosproject/onos-ric/pkg/store/indications"
	"github.com/onosproject/onos-ric/pkg/store/mastership"
	topodevice "github.com/onosproject/onos-topo/api/device"
	"google.golang.org/grpc"
)

const e2NodeType = topodevice.Type("E2Node")

var log = logging.GetLogger("manager")
var mgr Manager

// MastershipStoreFactory creates the mastership store
var MastershipStoreFactory = func(configuration config.Config) (mastership.Store, error) {
	return mastership.NewDistributedStore(configuration)
}

// IndicationsStoreFactory creates the indications store
var IndicationsStoreFactory = func(configuration config.Config) (indications.Store, error) {
	return indications.NewDistributedStore(configuration)
}

// DeviceStoreFactory creates the device store
var DeviceStoreFactory = func(topoEndPoint string, opts ...grpc.DialOption) (device.Store, error) {
	return device.NewTopoStore(topoEndPoint, opts...)
}

// NewManager initializes the RAN subsystem.
func NewManager(topoEndPoint string, enableMetrics bool, opts []grpc.DialOption) (*Manager, error) {
	log.Info("Creating Manager")

	configuration, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	indicationsStore, err := IndicationsStoreFactory(configuration)
	if err != nil {
		return nil, err
	}

	deviceChangeStore, err := DeviceStoreFactory(topoEndPoint, opts...)
	if err != nil {
		log.Info("Error in device change store")
		return nil, err
	}

	return InitializeManager(indicationsStore, deviceChangeStore, enableMetrics), nil
}

// InitializeManager initializes the manager structure with the given data
func InitializeManager(indicationsStore indications.Store, deviceChangesStore device.Store, enableMetrics bool) *Manager {
	mgr = Manager{
		indicationsStore:   indicationsStore,
		deviceChangesStore: deviceChangesStore,
		SbSessions:         make(map[sb.ECGI]southbound.E2),
		enableMetrics:      enableMetrics,
		topoMonitor: monitor.NewTopoMonitorBuilder().
			SetTopoChannel(make(chan *topodevice.ListResponse)).
			Build(),
		e2Chan: make(chan southbound.E2),
	}
	return &mgr
}

// Manager single point of entry for the RAN system.
type Manager struct {
	indicationsStore   indications.Store
	deviceChangesStore device.Store
	SbSessions         map[sb.ECGI]southbound.E2
	topoMonitor        monitor.TopoMonitor
	enableMetrics      bool
	e2Chan             chan southbound.E2 // TBD - Make this into a dispatcher
}

// StoreRicControlResponse - write the RicControlResponse to store
func (m *Manager) StoreRicControlResponse(update e2ap.RicIndication) {
	msgType := update.GetHdr().GetMessageType()
	switch msgType {
	case sb.MessageType_CELL_CONFIG_REPORT:
		_ = m.indicationsStore.Put(indications.GetID(&update), &update)
	default:
		log.Fatalf("RicControlResponse has unexpected type %d", msgType)
	}
}

// StoreControlUpdate - put the control update in the atomix store
func (m *Manager) StoreControlUpdate(update e2ap.RicIndication) {
	log.Infof("Got messageType %s", update.GetHdr().MessageType)
	switch update.GetHdr().GetMessageType() {
	case sb.MessageType_UE_ADMISSION_REQUEST:
		msg := update.GetMsg().GetUEAdmissionRequest()
		log.Infof("plmnid:%s, ecid:%s, crnti:%s, imsi:%d", msg.GetEcgi().GetPlmnId(), msg.GetEcgi().GetEcid(), msg.GetCrnti(), msg.GetImsi())
		_ = m.indicationsStore.Put(indications.GetID(&update), &update)
	case sb.MessageType_UE_RELEASE_IND:
		msg := update.GetMsg().GetUEReleaseInd()
		log.Infof("delete ue - plmnid:%s, ecid:%s, crnti:%s", msg.GetEcgi().GetPlmnId(), msg.GetEcgi().GetEcid(), msg.GetCrnti())

		err := m.DeleteTelemetry(msg.GetEcgi().GetPlmnId(), msg.GetEcgi().GetEcid(), msg.GetCrnti())
		if err != nil {
			log.Errorf("%s", err)
		}
		err = m.DeleteUEAdmissionRequest(msg.GetEcgi().GetPlmnId(), msg.GetEcgi().GetEcid(), msg.GetCrnti())
		if err != nil {
			log.Errorf("%s", err)
		}
	default:
		log.Fatalf("ControlReport has unexpected type %T", update.GetHdr().GetMessageType())
	}
}

// GetIndications gets update indications
func (m *Manager) GetIndications() ([]e2ap.RicIndication, error) {
	ch := make(chan e2ap.RicIndication)
	if err := m.ListIndications(ch); err != nil {
		return nil, err
	}
	messages := make([]e2ap.RicIndication, 0)
	for update := range ch {
		messages = append(messages, update)
	}
	return messages, nil
}

// GetUEAdmissionByID retrieve a single value from the updates store
func (m *Manager) GetUEAdmissionByID(ecgi *sb.ECGI, crnti string) (*e2ap.RicIndication, error) {
	return m.indicationsStore.Get(indications.NewID(sb.MessageType_UE_ADMISSION_REQUEST, ecgi.PlmnId, ecgi.Ecid, crnti))
}

// ListIndications lists control updates
func (m *Manager) ListIndications(ch chan<- e2ap.RicIndication) error {
	return m.indicationsStore.List(ch)
}

// SubscribeIndications subscribes the given channel to control updates
func (m *Manager) SubscribeIndications(ch chan<- indications.Event, opts ...indications.WatchOption) error {
	return m.indicationsStore.Watch(ch, opts...)
}

// Run starts a synchronizer based on the devices and the northbound services.
func (m *Manager) Run() {
	log.Info("Starting Manager")

	config.NewE2Config().Run(m.e2Chan)

	m.topoMonitor.TopoEventHandler(m.topoEventHandler)

	err := mgr.deviceChangesStore.Watch(m.topoMonitor.TopoChannel())
	if err != nil {
		log.Errorf("Error listening to topo service: %s", err.Error())
	}
}

func (m *Manager) topoEventHandler(topoChannel chan *topodevice.ListResponse) {
	log.Infof("Watching topo channel")
	for device := range topoChannel { // Block here waiting for topo events
		if device.GetDevice().GetType() != e2NodeType {
			continue
		}
		if device.Type == topodevice.ListResponse_NONE || device.Type == topodevice.ListResponse_ADDED {
			ecgi, err := ecgiFromTopoID(device.GetDevice().GetID())
			if err != nil {
				log.Warnf("Topo device event rejected %s", err)
				continue
			}
			log.Infof("Topo Device added %s", ecgi.String())
			deviceEndpoint := sb.Endpoint(device.GetDevice().GetAddress())
			session, err := southbound.NewSession()
			if err != nil {
				log.Errorf("Unable to create new session %s", err.Error())
				continue
			}
			if session != nil {
				m.SbSessions[ecgi] = session
				session.Run(ecgi, deviceEndpoint, device.GetDevice().GetTLS(),
					device.GetDevice().GetCredentials(), m.StoreRicControlResponse,
					m.StoreControlUpdate, m.StoreTelemetry, m.enableMetrics, m.e2Chan)
			} else {
				log.Errorf("Error creating new session for %v", ecgi)
				continue
			}
		} else if device.Type == topodevice.ListResponse_REMOVED {
			ecgi, err := ecgiFromTopoID(device.GetDevice().GetID())
			if err != nil {
				log.Warnf("Topo device event rejected %s", err)
				continue
			}
			log.Infof("Topo Device removing %s", ecgi.String())
			session, ok := m.SbSessions[ecgi]
			if !ok {
				log.Warnf("Unable to find session for %s", ecgi.String())
				continue
			}
			session.Close()
			delete(m.SbSessions, ecgi)
			// Delete all entries in the indicationStore for this cell
			ricIndications := make(chan e2ap.RicIndication)
			if err := m.ListIndications(ricIndications); err != nil {
				log.Warnf("Unable to list indications %s", err)
				continue
			}
			for ricIndication := range ricIndications {
				switch ricIndication.GetHdr().GetMessageType() {
				case sb.MessageType_CELL_CONFIG_REPORT:
					cellConfigReport := ricIndication.GetMsg().GetCellConfigReport()
					if cellConfigReport.GetEcgi().String() == ecgi.String() {
						ricIndication := ricIndication // pin
						if err := m.indicationsStore.Delete(indications.GetID(&ricIndication)); err != nil {
							log.Warnf("Unable to delete CELL_CONFIG_REPORT indication %s", err)
							continue
						}
					}
				}
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
func (m *Manager) StoreTelemetry(update e2ap.RicIndication) {
	err := m.indicationsStore.Put(indications.GetID(&update), &update)
	if err != nil {
		log.Fatalf("Could not put message %v in telemetry store %s", update, err.Error())
	}

	switch update.GetHdr().MessageType {
	case sb.MessageType_RADIO_MEAS_REPORT_PER_UE:
		msg := update.GetMsg()
		log.Infof("RadioMeasReport plmnid:%s ecid:%s crnti:%s cqis:%d(ecid:%s),%d(ecid:%s),%d(ecid:%s),%d(ecid:%s)",
			msg.GetRadioMeasReportPerUE().GetEcgi().GetPlmnId(),
			msg.GetRadioMeasReportPerUE().GetEcgi().GetEcid(),
			msg.GetRadioMeasReportPerUE().GetCrnti(),
			msg.GetRadioMeasReportPerUE().RadioReportServCells[0].CqiHist[0],
			msg.GetRadioMeasReportPerUE().RadioReportServCells[0].GetEcgi().GetEcid(),
			msg.GetRadioMeasReportPerUE().RadioReportServCells[1].CqiHist[0],
			msg.GetRadioMeasReportPerUE().RadioReportServCells[1].GetEcgi().GetEcid(),
			msg.GetRadioMeasReportPerUE().RadioReportServCells[2].CqiHist[0],
			msg.GetRadioMeasReportPerUE().RadioReportServCells[2].GetEcgi().GetEcid(),
			msg.GetRadioMeasReportPerUE().RadioReportServCells[3].CqiHist[0],
			msg.GetRadioMeasReportPerUE().RadioReportServCells[3].GetEcgi().GetEcid(),
		)
	default:
		log.Fatalf("Telemetry update has unexpected type %T", update.GetHdr().GetMessageType())
	}
}

// DeleteTelemetry deletes telemetry when a handover happens
func (m *Manager) DeleteTelemetry(plmnid string, ecid string, crnti string) error {
	id := indications.NewID(sb.MessageType_RADIO_MEAS_REPORT_PER_UE, plmnid, ecid, crnti)
	if err := m.indicationsStore.Delete(id); err != nil {
		log.Infof("Error deleting Telemetry, key=%s", id)
		return err
	}
	return nil
}

// DeleteUEAdmissionRequest deletes UpdateControls
func (m *Manager) DeleteUEAdmissionRequest(plmnid string, ecid string, crnti string) error {
	id := indications.NewID(sb.MessageType_UE_ADMISSION_REQUEST, plmnid, ecid, crnti)
	if err := m.indicationsStore.Delete(id); err != nil {
		log.Infof("Error deleting UEAdmissionRequest, key=%s", id)
		return err
	}
	return nil
}

// ecgiFromTopoID topo device is formatted like "315010-0001786" PlmnId-Ecid
func ecgiFromTopoID(id topodevice.ID) (sb.ECGI, error) {
	if !strings.Contains(string(id), "-") {
		return sb.ECGI{}, fmt.Errorf("unexpected format for E2Node ID %s", id)
	}
	parts := strings.Split(string(id), "-")
	if len(parts) != 2 {
		return sb.ECGI{}, fmt.Errorf("unexpected format for E2Node ID %s", id)
	}
	return sb.ECGI{Ecid: parts[1], PlmnId: parts[0]}, nil
}
