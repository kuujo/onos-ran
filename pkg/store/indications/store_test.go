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

package indications

import (
	"testing"
	"time"

	"github.com/onosproject/onos-ric/api/sb"
	"github.com/onosproject/onos-ric/api/sb/e2ap"
	"github.com/onosproject/onos-ric/api/sb/e2sm"
	"github.com/stretchr/testify/assert"
)

func TestStoreUpdates(t *testing.T) {
	testStore, err := NewLocalStore()
	assert.NoError(t, err)
	assert.NotNil(t, testStore)

	defer testStore.Close()

	watchCh2 := make(chan Event)
	err = testStore.Watch(watchCh2, WithReplay())
	assert.NoError(t, err)

	controlUpdate3 := &e2ap.RicIndication{
		Hdr: &e2sm.RicIndicationHeader{
			MessageType: sb.MessageType_UE_ADMISSION_REQUEST,
		},
		Msg: &e2sm.RicIndicationMessage{
			S: &e2sm.RicIndicationMessage_UEAdmissionRequest{
				UEAdmissionRequest: &sb.UEAdmissionRequest{
					Ecgi: &sb.ECGI{
						Ecid:   "test-ecid-3",
						PlmnId: "test-plmnid-3",
					},
					Crnti: "test-crnti-3",
				},
			},
		},
	}
	err = testStore.Put(GetID(controlUpdate3), controlUpdate3)
	assert.Nil(t, err)

	ecigTest3 := sb.ECGI{
		Ecid:   "test-ecid-3",
		PlmnId: "test-plmnid-3",
	}
	id3 := NewID(controlUpdate3.GetHdr().GetMessageType(), ecigTest3.GetPlmnId(), ecigTest3.GetEcid(), "test-crnti-3")
	value3, err := testStore.Get(id3)
	assert.Nil(t, err)
	assert.Equal(t, value3.GetHdr().MessageType.String(), sb.MessageType_UE_ADMISSION_REQUEST.String())

	id4 := NewID(controlUpdate3.GetHdr().GetMessageType(), ecigTest3.GetPlmnId(), ecigTest3.GetEcid(), "test-crnti-none")
	value4, err := testStore.Get(id4)
	assert.Nil(t, err)
	assert.Nil(t, value4)

	event := <-watchCh2
	assert.Equal(t, "test-ecid-3", event.Message.GetMsg().GetUEAdmissionRequest().Ecgi.Ecid)
	assert.Equal(t, "test-plmnid-3", event.Message.GetMsg().GetUEAdmissionRequest().Ecgi.PlmnId)

	err = testStore.Delete(id4)
	assert.Nil(t, err)

	err = testStore.Delete(id3)
	assert.Nil(t, err)

	err = testStore.Clear()
	assert.Nil(t, err)
}

func TestStoreTelemetry(t *testing.T) {
	testStore, err := NewLocalStore()
	assert.NoError(t, err)
	assert.NotNil(t, testStore)

	defer testStore.Close()

	watchCh1 := make(chan Event)
	err = testStore.Watch(watchCh1)
	assert.NoError(t, err)

	telemetry1 := &e2ap.RicIndication{
		Hdr: &e2sm.RicIndicationHeader{
			MessageType: sb.MessageType_RADIO_MEAS_REPORT_PER_CELL,
		},
		Msg: &e2sm.RicIndicationMessage{
			S: &e2sm.RicIndicationMessage_RadioMeasReportPerCell{
				RadioMeasReportPerCell: &sb.RadioMeasReportPerCell{
					Ecgi: &sb.ECGI{
						Ecid:   "test-ecid",
						PlmnId: "test-plmnid",
					},
				},
			},
		},
	}
	err = testStore.Put(GetID(telemetry1), telemetry1)
	assert.Nil(t, err)
	ecigTest1 := sb.ECGI{
		Ecid:   "test-ecid",
		PlmnId: "test-plmnid",
	}
	id1 := NewID(telemetry1.GetHdr().GetMessageType(), ecigTest1.GetPlmnId(), ecigTest1.GetEcid(), "")
	value1, err := testStore.Get(id1)
	assert.Nil(t, err)
	assert.Equal(t, value1.GetHdr().GetMessageType().String(), sb.MessageType_RADIO_MEAS_REPORT_PER_CELL.String())

	event := <-watchCh1
	assert.Equal(t, "test-ecid", event.Message.GetMsg().GetRadioMeasReportPerCell().Ecgi.Ecid)
	assert.Equal(t, "test-plmnid", event.Message.GetMsg().GetRadioMeasReportPerCell().Ecgi.PlmnId)

	watchCh2 := make(chan Event)
	err = testStore.Watch(watchCh2, WithReplay())
	assert.NoError(t, err)

	event = <-watchCh2
	assert.Equal(t, "test-ecid", event.Message.GetMsg().GetRadioMeasReportPerCell().Ecgi.Ecid)
	assert.Equal(t, "test-plmnid", event.Message.GetMsg().GetRadioMeasReportPerCell().Ecgi.PlmnId)

	telemetry2 := &e2ap.RicIndication{
		Hdr: &e2sm.RicIndicationHeader{
			MessageType: sb.MessageType_RADIO_MEAS_REPORT_PER_CELL,
		},
		Msg: &e2sm.RicIndicationMessage{
			S: &e2sm.RicIndicationMessage_RadioMeasReportPerCell{
				RadioMeasReportPerCell: &sb.RadioMeasReportPerCell{
					Ecgi: &sb.ECGI{
						Ecid:   "test-ecid-2",
						PlmnId: "test-plmnid-2",
					},
				},
			},
		},
	}

	err = testStore.Put(GetID(telemetry2), telemetry2)
	assert.Nil(t, err)
	ecigTest2 := sb.ECGI{
		Ecid:   "test-ecid-2",
		PlmnId: "test-plmnid-2",
	}
	id2 := NewID(telemetry2.GetHdr().GetMessageType(), ecigTest2.GetPlmnId(), ecigTest2.GetEcid(), "")
	value2, err := testStore.Get(id2)
	assert.Nil(t, err)
	assert.Equal(t, value2.GetHdr().MessageType.String(), sb.MessageType_RADIO_MEAS_REPORT_PER_CELL.String())

	event = <-watchCh1
	assert.Equal(t, "test-ecid-2", event.Message.GetMsg().GetRadioMeasReportPerCell().Ecgi.Ecid)
	assert.Equal(t, "test-plmnid-2", event.Message.GetMsg().GetRadioMeasReportPerCell().Ecgi.PlmnId)

	event = <-watchCh2
	assert.Equal(t, "test-ecid-2", event.Message.GetMsg().GetRadioMeasReportPerCell().Ecgi.Ecid)
	assert.Equal(t, "test-plmnid-2", event.Message.GetMsg().GetRadioMeasReportPerCell().Ecgi.PlmnId)

	telemetry3 := &e2ap.RicIndication{
		Hdr: &e2sm.RicIndicationHeader{
			MessageType: sb.MessageType_RADIO_MEAS_REPORT_PER_UE,
		},
		Msg: &e2sm.RicIndicationMessage{
			S: &e2sm.RicIndicationMessage_RadioMeasReportPerUE{
				RadioMeasReportPerUE: &sb.RadioMeasReportPerUE{
					Ecgi: &sb.ECGI{
						Ecid:   "test-ecid-3",
						PlmnId: "test-plmnid-3",
					},
					Crnti: "test-crnti-3",
				},
			},
		},
	}
	err = testStore.Put(GetID(telemetry3), telemetry3)
	assert.Nil(t, err)

	ecigTest3 := sb.ECGI{
		Ecid:   "test-ecid-3",
		PlmnId: "test-plmnid-3",
	}
	id3 := NewID(telemetry3.GetHdr().GetMessageType(), ecigTest3.GetPlmnId(), ecigTest3.GetEcid(), "test-crnti-3")
	value3, err := testStore.Get(id3)
	assert.Nil(t, err)
	assert.Equal(t, value3.GetHdr().MessageType.String(), sb.MessageType_RADIO_MEAS_REPORT_PER_UE.String())

	event = <-watchCh1
	assert.Equal(t, "test-ecid-3", event.Message.GetMsg().GetRadioMeasReportPerUE().Ecgi.Ecid)
	assert.Equal(t, "test-plmnid-3", event.Message.GetMsg().GetRadioMeasReportPerUE().Ecgi.PlmnId)

	event = <-watchCh2
	assert.Equal(t, "test-ecid-3", event.Message.GetMsg().GetRadioMeasReportPerUE().Ecgi.Ecid)
	assert.Equal(t, "test-plmnid-3", event.Message.GetMsg().GetRadioMeasReportPerUE().Ecgi.PlmnId)

	listCh := make(chan e2ap.RicIndication)
	err = testStore.List(listCh)
	assert.NoError(t, err)

	select {
	case <-listCh:
	case <-time.After(5 * time.Second):
		t.Fail()
	}
	select {
	case <-listCh:
	case <-time.After(5 * time.Second):
		t.Fail()
	}
	select {
	case <-listCh:
	case <-time.After(5 * time.Second):
		t.Fail()
	}

	err = testStore.Delete(id1)
	assert.Nil(t, err)

	err = testStore.Delete(id2)
	assert.Nil(t, err)

	err = testStore.Delete(id3)
	assert.Nil(t, err)
}
