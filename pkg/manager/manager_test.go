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

package manager

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/onosproject/onos-lib-go/pkg/atomix"
	"github.com/onosproject/onos-ric/api/sb"
	"github.com/onosproject/onos-ric/api/sb/e2ap"
	"github.com/onosproject/onos-ric/api/sb/e2sm"
	"github.com/onosproject/onos-ric/pkg/config"
	"github.com/onosproject/onos-ric/pkg/store/device"
	"github.com/onosproject/onos-ric/pkg/store/indications"
	"github.com/onosproject/onos-ric/pkg/store/mastership"
	"github.com/onosproject/onos-ric/pkg/store/requests"
	"github.com/onosproject/onos-ric/test/mocks/store/device"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"sync"
	"testing"
	"time"
)

var oldMastershipStoreFactory func(configuration config.Config) (mastership.Store, error)
var oldIndicationsStoreFactory func(configuration config.Config) (indications.Store, error)
var oldRequestsStoreFactory func(configuration config.Config) (requests.Store, error)
var oldDeviceStoreFactory func(topoEndPoint string, opts ...grpc.DialOption) (store device.Store, err error)

func saveFactories() {
	oldMastershipStoreFactory = MastershipStoreFactory
	oldIndicationsStoreFactory = IndicationsStoreFactory
	oldRequestsStoreFactory = RequestsStoreFactory
	oldDeviceStoreFactory = DeviceStoreFactory
}

func restoreFactories() {
	MastershipStoreFactory = oldMastershipStoreFactory
	IndicationsStoreFactory = oldIndicationsStoreFactory
	RequestsStoreFactory = oldRequestsStoreFactory
	DeviceStoreFactory = oldDeviceStoreFactory
}

func makeNewManager(t *testing.T) *Manager {
	saveFactories()
	_, address := atomix.StartLocalNode()
	assert.NotNil(t, address)
	mockConfig := &config.Config{
		Mastership: config.MastershipConfig{},
		Atomix:     atomix.Config{},
	}
	mockTopoStore := mock_device_store.NewMockStore(gomock.NewController(t))
	mockTopoStore.EXPECT().Watch(gomock.Any()).AnyTimes()
	mockConfig.Atomix.Controller = string(address)
	config.WithConfig(mockConfig)
	MastershipStoreFactory = func(configuration config.Config) (mastership.Store, error) {
		return mastership.NewLocalStore("cluster1", "node1")
	}
	IndicationsStoreFactory = func(configuration config.Config) (indications.Store, error) {
		return indications.NewLocalStore()
	}
	RequestsStoreFactory = func(configuration config.Config) (requests.Store, error) {
		return requests.NewLocalStore()
	}
	DeviceStoreFactory = func(topoEndPoint string, opts ...grpc.DialOption) (store device.Store, err error) {
		return mockTopoStore, nil
	}

	newManager, err := NewManager("", nil)
	assert.NoError(t, err)
	assert.NotNil(t, newManager)
	newManager.Run()
	return newManager
}

func removeNewManager() {
	restoreFactories()
}

func IDToEcid(ID uint32) string {
	return fmt.Sprintf("ECID-%d", ID)
}

func IDToPlmnid(ID uint32) string {
	return fmt.Sprintf("PLMNID-%d", ID)
}

func IDToCrnti(ID uint32) string {
	return fmt.Sprintf("CRNTI-%d", ID)
}

func IDToECGI(ID uint32) *sb.ECGI {
	return &sb.ECGI{
		PlmnId: IDToPlmnid(ID),
		Ecid:   IDToEcid(ID),
	}
}

func generateRicIndicationRadioMeasReportPerUE(ID uint32) e2ap.RicIndication {
	radioReportServCells := []*sb.RadioRepPerServCell{
		{
			Ecgi:    IDToECGI(ID + 1000),
			CqiHist: []uint32{1, 2, 3},
		},
		{
			Ecgi:    IDToECGI(ID + 1000),
			CqiHist: []uint32{1, 2, 3},
		},
		{
			Ecgi:    IDToECGI(ID + 1000),
			CqiHist: []uint32{1, 2, 3},
		},
		{
			Ecgi:    IDToECGI(ID + 1000),
			CqiHist: []uint32{1, 2, 3},
		},
	}
	return e2ap.RicIndication{
		Hdr: &e2sm.RicIndicationHeader{MessageType: sb.MessageType_RADIO_MEAS_REPORT_PER_UE},
		Msg: &e2sm.RicIndicationMessage{S: &e2sm.RicIndicationMessage_RadioMeasReportPerUE{RadioMeasReportPerUE: &sb.RadioMeasReportPerUE{
			Ecgi:                 IDToECGI(ID),
			Crnti:                IDToCrnti(ID),
			RadioReportServCells: radioReportServCells,
		}}},
	}
}

func generateRicIndicationsUEAdmissionRequest(ID uint32) e2ap.RicIndication {
	return e2ap.RicIndication{
		Hdr: &e2sm.RicIndicationHeader{MessageType: sb.MessageType_UE_ADMISSION_REQUEST},
		Msg: &e2sm.RicIndicationMessage{S: &e2sm.RicIndicationMessage_UEAdmissionRequest{UEAdmissionRequest: &sb.UEAdmissionRequest{
			Crnti: IDToCrnti(ID),
			Ecgi:  IDToECGI(ID),
		}}},
	}
}

func generateRicIndicationsUEReleaseRequest(ID uint32) e2ap.RicIndication {
	return e2ap.RicIndication{
		Hdr: &e2sm.RicIndicationHeader{MessageType: sb.MessageType_UE_RELEASE_IND},
		Msg: &e2sm.RicIndicationMessage{S: &e2sm.RicIndicationMessage_UEReleaseInd{UEReleaseInd: &sb.UEReleaseInd{
			Crnti:        IDToCrnti(ID),
			Ecgi:         IDToECGI(ID),
			ReleaseCause: sb.ReleaseCause_RELEASE_INACTIVITY,
		}}},
	}
}

func generateRicControlRequest(ID uint32) e2ap.RicControlRequest {
	return e2ap.RicControlRequest{
		Hdr: &e2sm.RicControlHeader{
			MessageType: sb.MessageType_CELL_CONFIG_REQUEST,
			Crnti:       []string{IDToCrnti(ID)},
			Ecgi:        IDToECGI(ID),
		},
		Msg: &e2sm.RicControlMessage{S: &e2sm.RicControlMessage_CellConfigRequest{}},
	}
}

func generateRicIndicationsRicControlResponse(ID uint32) e2ap.RicIndication {
	return e2ap.RicIndication{
		Hdr: &e2sm.RicIndicationHeader{MessageType: sb.MessageType_CELL_CONFIG_REPORT},
		Msg: &e2sm.RicIndicationMessage{S: &e2sm.RicIndicationMessage_CellConfigReport{CellConfigReport: &sb.CellConfigReport{
			Ecgi: IDToECGI(ID),
		}}},
	}
}

func checkEvent(t *testing.T, expected e2ap.RicIndication, actual indications.Event, eventType indications.EventType) {
	assert.Equal(t, eventType, actual.Type)
	actualReport := actual.Indication.Msg.GetRadioMeasReportPerUE()
	expectedReport := expected.Msg.GetRadioMeasReportPerUE()
	assert.Equal(t, expectedReport, actualReport)
}

func Test_Create(t *testing.T) {
	newManager := makeNewManager(t)

	testManager := GetManager()
	assert.Equal(t, newManager, testManager)
	newManager.Close()
	removeNewManager()
}

func Test_TelemetrySubscribe(t *testing.T) {
	newManager := makeNewManager(t)

	telemetryMessage := generateRicIndicationRadioMeasReportPerUE(1)

	ch := make(chan indications.Event)

	err := newManager.SubscribeIndications(ch)
	assert.NoError(t, err)

	var wg sync.WaitGroup
	wg.Add(1)
	var e1 indications.Event
	var e2 indications.Event
	go func() {
		e1 = <-ch
		e2 = <-ch

		wg.Done()
	}()

	err = newManager.StoreTelemetry(telemetryMessage)
	assert.NoError(t, err, "error storing telemetry %v", err)
	err = newManager.DeleteTelemetry(IDToPlmnid(1), IDToEcid(1), IDToCrnti(1))
	assert.NoError(t, err, "error deleting telemetry %v", err)

	wg.Wait()

	// First event should be an insert
	checkEvent(t, telemetryMessage, e1, indications.EventInsert)

	// Second event should be delete
	assert.Equal(t, indications.EventDelete, e2.Type)
	removeNewManager()
}

func Test_ListIndications(t *testing.T) {
	newManager := makeNewManager(t)

	telemetryMessage := generateRicIndicationRadioMeasReportPerUE(1)
	err := newManager.StoreTelemetry(telemetryMessage)
	assert.NoError(t, err, "error storing telemetry %v", err)

	ch := make(chan e2ap.RicIndication)

	err = newManager.ListIndications(ch)
	assert.NoError(t, err)

	var wg sync.WaitGroup
	wg.Add(1)
	var ri e2ap.RicIndication

	go func() {
		ri = <-ch
		wg.Done()
	}()

	wg.Wait()

	assert.Equal(t, telemetryMessage.Msg.GetRadioMeasReportPerUE(), ri.Msg.GetRadioMeasReportPerUE())
	removeNewManager()
}

func Test_GetIndications(t *testing.T) {
	newManager := makeNewManager(t)

	telemetryMessage := generateRicIndicationRadioMeasReportPerUE(1)
	err := newManager.StoreTelemetry(telemetryMessage)
	assert.NoError(t, err, "error storing telemetry %v", err)

	inds, err := newManager.GetIndications()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(inds))
	assert.Equal(t, telemetryMessage.Msg.GetRadioMeasReportPerUE(), inds[0].Msg.GetRadioMeasReportPerUE())
	removeNewManager()
}

func Test_StoreControlUpdateUEAdmission(t *testing.T) {
	newManager := makeNewManager(t)

	UEAdmissionMessage := generateRicIndicationsUEAdmissionRequest(55)
	err := newManager.StoreControlUpdate(UEAdmissionMessage)
	assert.NoError(t, err)

	inds, err := newManager.GetIndications()
	assert.NoError(t, err)
	assert.NotNil(t, inds)
	assert.Equal(t, 1, len(inds))
	assert.Equal(t, UEAdmissionMessage.Msg.GetUEAdmissionRequest(), inds[0].Msg.GetUEAdmissionRequest())
	removeNewManager()
}

func Test_StoreControlUpdateUERelease(t *testing.T) {
	newManager := makeNewManager(t)

	// Add a UE
	UEAdmissionMessage := generateRicIndicationsUEAdmissionRequest(55)
	err := newManager.StoreControlUpdate(UEAdmissionMessage)
	assert.NoError(t, err)

	// Check that the UE made it into the store
	indsBefore, err := newManager.GetIndications()
	assert.NoError(t, err)
	assert.NotNil(t, indsBefore)
	assert.Equal(t, 1, len(indsBefore))
	assert.Equal(t, UEAdmissionMessage.Msg.GetUEAdmissionRequest(), indsBefore[0].Msg.GetUEAdmissionRequest())

	// Now release the UE
	UEReleaseMessage := generateRicIndicationsUEReleaseRequest(55)
	err = newManager.StoreControlUpdate(UEReleaseMessage)
	assert.NoError(t, err)

	time.Sleep(2 * time.Second)

	indsAfter, err := newManager.GetIndications()
	assert.NoError(t, err)
	assert.NotNil(t, indsAfter)
	assert.Empty(t, indsAfter)
	removeNewManager()
}

func Test_GetUEAdmissionByID(t *testing.T) {
	newManager := makeNewManager(t)
	ID := uint32(778899)

	// Add a UE
	UEAdmissionMessage := generateRicIndicationsUEAdmissionRequest(ID)
	err := newManager.StoreControlUpdate(UEAdmissionMessage)
	assert.NoError(t, err)

	// Check that the UE made it into the store
	crnti := IDToCrnti(ID)
	ecgi := IDToECGI(ID)
	ue, err := newManager.GetUEAdmissionByID(ecgi, crnti)
	assert.NoError(t, err)
	assert.Equal(t, UEAdmissionMessage.Msg.GetUEAdmissionRequest(), ue.Msg.GetUEAdmissionRequest())
	removeNewManager()
}

func Test_RicControlMessages(t *testing.T) {
	newManager := makeNewManager(t)
	ID := uint32(556677)
	deviceID := device.ID{PlmnId: IDToPlmnid(ID), Ecid: IDToEcid(ID)}

	requestMessage := generateRicControlRequest(ID)
	err := newManager.StoreRicControlRequest(deviceID, &requestMessage)
	assert.NoError(t, err)

	time.Sleep(2 * time.Second)

	indsAfterControlRequest, err := newManager.GetIndications()
	assert.NoError(t, err)
	assert.NotNil(t, indsAfterControlRequest)
	assert.Empty(t, indsAfterControlRequest)

	response := generateRicIndicationsRicControlResponse(ID)
	err = newManager.StoreRicControlResponse(response)
	assert.NoError(t, err)

	time.Sleep(2 * time.Second)

	indsAfterControlResponse, err := newManager.GetIndications()
	assert.NoError(t, err)
	assert.NotNil(t, indsAfterControlResponse)
	assert.Equal(t, 1, len(indsAfterControlResponse))
	assert.Equal(t, IDToEcid(ID), indsAfterControlResponse[0].Msg.GetCellConfigReport().Ecgi.Ecid)
	removeNewManager()
}

// Test_StoresBadConfig : Tests default store implementations with bad configurations
func Test_StoresBadConfig(t *testing.T) {
	var cfg config.Config

	mship, err := MastershipStoreFactory(cfg)
	assert.Error(t, err)
	assert.Nil(t, mship)

	inds, err := IndicationsStoreFactory(cfg)
	assert.Error(t, err)
	assert.Nil(t, inds)

	reqs, err := RequestsStoreFactory(cfg)
	assert.Error(t, err)
	assert.Nil(t, reqs)

	devs, err := DeviceStoreFactory("abc")
	assert.Error(t, err)
	assert.Nil(t, devs)
}
