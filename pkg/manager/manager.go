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
	"github.com/onosproject/onos-lib-go/pkg/atomix"
	"github.com/onosproject/onos-lib-go/pkg/cluster"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-ric/api/sb"
	"github.com/onosproject/onos-ric/api/sb/e2ap"
	"github.com/onosproject/onos-ric/pkg/config"
	"github.com/onosproject/onos-ric/pkg/southbound/e2"
	"github.com/onosproject/onos-ric/pkg/store/device"
	"github.com/onosproject/onos-ric/pkg/store/indications"
	"github.com/onosproject/onos-ric/pkg/store/mastership"
	"github.com/onosproject/onos-ric/pkg/store/requests"
	"google.golang.org/grpc"
)

var log = logging.GetLogger("manager")
var mgr Manager

// ClusterFactory creates the cluster
var ClusterFactory = func(configuration config.Config) (cluster.Cluster, error) {
	client, err := atomix.GetClient(configuration.Atomix)
	if err != nil {
		return nil, err
	}
	return cluster.New(client)
}

// MastershipStoreFactory creates the mastership store
var MastershipStoreFactory = func(cluster cluster.Cluster, configuration config.Config) (mastership.Store, error) {
	return mastership.NewDistributedStore(cluster, configuration)
}

// IndicationsStoreFactory creates the indications store
var IndicationsStoreFactory = func(cluster cluster.Cluster, deviceStore device.Store, mastershipStore mastership.Store, configuration config.Config) (indications.Store, error) {
	return indications.NewDistributedStore(cluster, deviceStore, mastershipStore, configuration)
}

// RequestsStoreFactory creates the requests store
var RequestsStoreFactory = func(cluster cluster.Cluster, deviceStore device.Store, mastershipStore mastership.Store, configuration config.Config) (requests.Store, error) {
	return requests.NewDistributedStore(cluster, deviceStore, mastershipStore, configuration)
}

// DeviceStoreFactory creates the device store
var DeviceStoreFactory = func(topoEndPoint string, opts ...grpc.DialOption) (device.Store, error) {
	return device.NewTopoStore(topoEndPoint, opts...)
}

// NewManager initializes the RAN subsystem.
func NewManager(topoEndPoint string, opts []grpc.DialOption) (*Manager, error) {
	log.Info("Creating Manager")

	configuration, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	cluster, err := ClusterFactory(configuration)
	if err != nil {
		return nil, err
	}

	deviceStore, err := DeviceStoreFactory(topoEndPoint, opts...)
	if err != nil {
		log.Info("Error in device change store")
		return nil, err
	}

	mastershipStore, err := MastershipStoreFactory(cluster, configuration)
	if err != nil {
		return nil, err
	}

	indicationsStore, err := IndicationsStoreFactory(cluster, deviceStore, mastershipStore, configuration)
	if err != nil {
		return nil, err
	}

	requestsStore, err := RequestsStoreFactory(cluster, deviceStore, mastershipStore, configuration)
	if err != nil {
		return nil, err
	}

	sessionMgr, err := e2.NewSessionManager(deviceStore, mastershipStore, requestsStore, indicationsStore)
	if err != nil {
		return nil, err
	}

	return InitializeManager(indicationsStore, requestsStore, sessionMgr), nil
}

// InitializeManager initializes the manager structure with the given data
func InitializeManager(indicationsStore indications.Store, requestsStore requests.Store, sessionMgr *e2.SessionManager) *Manager {
	mgr = Manager{
		indicationsStore: indicationsStore,
		requestsStore:    requestsStore,
		sessions:         sessionMgr,
	}
	return &mgr
}

// Manager single point of entry for the RAN system.
type Manager struct {
	indicationsStore indications.Store
	requestsStore    requests.Store
	sessions         *e2.SessionManager
}

// StoreRicControlRequest stores a RicControlRequest to the store
func (m *Manager) StoreRicControlRequest(request *e2ap.RicControlRequest) error {
	log.Debugf("Handling request %+v", request)
	return m.requestsStore.Append(requests.New(request))
}

// StoreRicControlResponse - write the RicControlResponse to store
func (m *Manager) StoreRicControlResponse(update e2ap.RicIndication) error {
	msgType := update.GetHdr().GetMessageType()
	switch msgType {
	case sb.MessageType_CELL_CONFIG_REPORT:
		return m.indicationsStore.Record(indications.New(&update))
	default:
		return fmt.Errorf("RicControlResponse has unexpected type %d", msgType)
	}
}

// StoreControlUpdate - put the control update in the atomix store
func (m *Manager) StoreControlUpdate(update e2ap.RicIndication) error {
	log.Infof("Got messageType %s", update.GetHdr().MessageType)
	switch update.GetHdr().GetMessageType() {
	case sb.MessageType_UE_ADMISSION_REQUEST:
		msg := update.GetMsg().GetUEAdmissionRequest()
		log.Infof("plmnid:%s, ecid:%s, crnti:%s, imsi:%d", msg.GetEcgi().GetPlmnId(), msg.GetEcgi().GetEcid(), msg.GetCrnti(), msg.GetImsi())
		return m.indicationsStore.Record(indications.New(&update))
	default:
		return fmt.Errorf("ControlReport has unexpected type %T", update.GetHdr().GetMessageType())
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

// ListIndications lists control updates
func (m *Manager) ListIndications(ch chan<- e2ap.RicIndication) error {
	indicationCh := make(chan indications.Indication)
	err := m.indicationsStore.List(indicationCh)
	if err != nil {
		return err
	}
	go func() {
		defer close(ch)
		for indication := range indicationCh {
			ch <- *indication.RicIndication
		}
	}()
	return nil
}

// SubscribeIndications subscribes the given channel to control updates
func (m *Manager) SubscribeIndications(ch chan<- indications.Event, opts ...indications.SubscribeOption) error {
	indicationCh := make(chan indications.Event)
	err := m.indicationsStore.Subscribe(indicationCh, opts...)
	if err != nil {
		return err
	}
	go func() {
		for indication := range indicationCh {
			log.Debugf("Handling response %+v", indication.Indication.RicIndication)
			ch <- indication
		}
	}()
	return nil
}

// Run starts a synchronizer based on the devices and the northbound services.
func (m *Manager) Run() {
	log.Info("Starting Manager")
}

//Close kills the channels and manager related objects
func (m *Manager) Close() {
	_ = m.sessions.Close()
	log.Info("Closing Manager")
}

// GetManager returns the initialized and running instance of manager.
// Should be called only after NewManager and Run are done.
func GetManager() *Manager {
	return &mgr
}

// StoreTelemetry - put the telemetry update in the atomix store
// Only handles MessageType_RADIO_MEAS_REPORT_PER_UE at the moment
func (m *Manager) StoreTelemetry(update e2ap.RicIndication) error {
	err := m.indicationsStore.Record(indications.New(&update))
	if err != nil {
		return err
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
		return nil
	default:
		return fmt.Errorf("telemetry update has unexpected type %T", update.GetHdr().GetMessageType())
	}
}
