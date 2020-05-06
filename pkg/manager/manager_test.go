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
	"github.com/onosproject/onos-ric/pkg/store/time"
	store "github.com/onosproject/onos-ric/test/mocks/store/device"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"sync"
	"testing"
	time2 "time"
)

func makeNewManager(t *testing.T) *Manager {
	_, address := atomix.StartLocalNode()
	assert.NotNil(t, address)
	mockConfig := &config.Config{
		Mastership: config.MastershipConfig{},
		Atomix:     atomix.Config{},
	}
	mockTopoStore := store.NewMockStore(gomock.NewController(t))
	mockConfig.Atomix.Controller = string(address)
	config.WithConfig(mockConfig)
	MastershipStoreFactory = func(configuration config.Config) (mastership.Store, error) {
		return mastership.NewLocalStore("cluster1", "node1")
	}
	IndicationsStoreFactory = func(configuration config.Config, timeStore time.Store) (indications.Store, error) {
		return indications.NewLocalStore(timeStore)
	}
	DeviceStoreFactory = func(topoEndPoint string, opts ...grpc.DialOption) (store device.Store, err error) {
		return mockTopoStore, nil
	}

	newManager, err := NewManager("", false, nil)
	assert.NoError(t, err)
	assert.NotNil(t, newManager)
	return newManager
}

func Test_Create(t *testing.T) {
	newManager := makeNewManager(t)

	testManager := GetManager()
	assert.Equal(t, newManager, testManager)
	newManager.Close()
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

func checkEvent(t *testing.T, expected e2ap.RicIndication, actual indications.Event, eventType indications.EventType) {
	assert.Equal(t, eventType, actual.Type)
	actualReport := actual.Message.Msg.GetRadioMeasReportPerUE()
	expectedReport := expected.Msg.GetRadioMeasReportPerUE()
	assert.Equal(t, expectedReport, actualReport)
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

	newManager.StoreTelemetry(telemetryMessage)
	// This sleep is to work around a bug in the store that causes List() to miss new
	// Indications if called too soon after Put()
	time2.Sleep(time2.Second * 2)
	err = newManager.DeleteTelemetry(IDToPlmnid(1), IDToEcid(1), IDToCrnti(1))
	assert.NoError(t, err, "error deleting telemetry %v", err)

	wg.Wait()

	// First event should be an insert
	checkEvent(t, telemetryMessage, e1, indications.EventInsert)

	// Second event should be delete
	assert.Equal(t, indications.EventDelete, e2.Type)
}

func Test_ListIndications(t *testing.T) {
	newManager := makeNewManager(t)

	telemetryMessage := generateRicIndicationRadioMeasReportPerUE(1)
	newManager.StoreTelemetry(telemetryMessage)

	// This sleep is to work around a bug in the store that causes List() to miss new
	// Indications if called too soon after Put()
	time2.Sleep(time2.Second * 2)

	ch := make(chan e2ap.RicIndication)

	err := newManager.ListIndications(ch)
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
}

func Test_GetIndications(t *testing.T) {
	newManager := makeNewManager(t)

	telemetryMessage := generateRicIndicationRadioMeasReportPerUE(1)
	newManager.StoreTelemetry(telemetryMessage)

	// This sleep is to work around a bug in the store that causes Get() to miss new
	// Indications if called too soon after Put()
	time2.Sleep(time2.Second * 2)

	inds, err := newManager.GetIndications()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(inds))
	assert.Equal(t, telemetryMessage.Msg.GetRadioMeasReportPerUE(), inds[0].Msg.GetRadioMeasReportPerUE())
}

func Test_StoreControlUpdateUEAdmission(t *testing.T) {
	newManager := makeNewManager(t)

	UEAdmissionMessage := generateRicIndicationsUEAdmissionRequest(55)
	newManager.StoreControlUpdate(UEAdmissionMessage)

	// This sleep is to work around a bug in the store that causes List() to miss new
	// Indications if called too soon after Put()
	time2.Sleep(time2.Second * 2)

	inds, err := newManager.GetIndications()
	assert.NoError(t, err)
	assert.NotNil(t, inds)
	assert.Equal(t, 1, len(inds))
	assert.Equal(t, UEAdmissionMessage.Msg.GetUEAdmissionRequest(), inds[0].Msg.GetUEAdmissionRequest())
}

func Test_StoreControlUpdateUERelease(t *testing.T) {
	newManager := makeNewManager(t)

	// Add a UE
	UEAdmissionMessage := generateRicIndicationsUEAdmissionRequest(55)
	newManager.StoreControlUpdate(UEAdmissionMessage)

	// This sleep is to work around a bug in the store that causes List() to miss new
	// Indications if called too soon after Put()
	time2.Sleep(time2.Second * 2)

	// Check that the UE made it into the store
	indsBefore, err := newManager.GetIndications()
	assert.NoError(t, err)
	assert.NotNil(t, indsBefore)
	assert.Equal(t, 1, len(indsBefore))
	assert.Equal(t, UEAdmissionMessage.Msg.GetUEAdmissionRequest(), indsBefore[0].Msg.GetUEAdmissionRequest())

	// Now release the UE
	UEReleaseMessage := generateRicIndicationsUEReleaseRequest(55)
	newManager.StoreControlUpdate(UEReleaseMessage)
	time2.Sleep(time2.Second * 2)

	indsAfter, err := newManager.GetIndications()
	assert.NoError(t, err)
	assert.NotNil(t, indsAfter)
	assert.Empty(t, indsAfter)
}

// GetUEAdmissionByID(ecgi *sb.ECGI, crnti string)
func Test_GetUEAdmissionByID(t *testing.T) {
	newManager := makeNewManager(t)
	ID := uint32(778899)

	// Add a UE
	UEAdmissionMessage := generateRicIndicationsUEAdmissionRequest(ID)
	newManager.StoreControlUpdate(UEAdmissionMessage)

	// This sleep is to work around a bug in the store that causes List() to miss new
	// Indications if called too soon after Put()
	time2.Sleep(time2.Second * 2)

	// Check that the UE made it into the store
	crnti := IDToCrnti(ID)
	ecgi := IDToECGI(ID)
	ue, err := newManager.GetUEAdmissionByID(ecgi, crnti)
	assert.NoError(t, err)
	assert.Equal(t, UEAdmissionMessage.Msg.GetUEAdmissionRequest(), ue.Msg.GetUEAdmissionRequest())
}
