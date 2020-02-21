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
	"testing"

	"github.com/onosproject/onos-ran/api/sb"
	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	testStore, err := NewStore()

	watchCh1 := make(chan sb.TelemetryMessage)
	err = testStore.Watch(watchCh1)
	assert.NoError(t, err)

	assert.Nil(t, err)
	assert.NotNil(t, testStore)
	telemetry1 := sb.TelemetryMessage{
		MessageType: sb.MessageType_RADIO_MEAS_REPORT_PER_CELL,
		S: &sb.TelemetryMessage_RadioMeasReportPerCell{
			RadioMeasReportPerCell: &sb.RadioMeasReportPerCell{
				Ecgi: &sb.ECGI{
					Ecid:   "test-ecid",
					PlmnId: "test-plmnid",
				},
			},
		},
	}
	err = testStore.Put(telemetry1)
	assert.Nil(t, err)
	ecigTest1 := sb.ECGI{
		Ecid:   "test-ecid",
		PlmnId: "test-plmnid",
	}
	id1 := ID{
		Ecid:        ecigTest1.GetEcid(),
		PlmnID:      ecigTest1.GetPlmnId(),
		MessageType: telemetry1.GetMessageType(),
	}
	value1, err := testStore.Get(id1)
	assert.Nil(t, err)
	assert.Equal(t, value1.MessageType.String(), sb.MessageType_RADIO_MEAS_REPORT_PER_CELL.String())

	event := <-watchCh1
	assert.Equal(t, "test-ecid", event.GetRadioMeasReportPerCell().Ecgi.Ecid)
	assert.Equal(t, "test-plmnid", event.GetRadioMeasReportPerCell().Ecgi.PlmnId)

	watchCh2 := make(chan sb.TelemetryMessage)
	err = testStore.Watch(watchCh2, WithReplay())

	event = <-watchCh2
	assert.Equal(t, "test-ecid", event.GetRadioMeasReportPerCell().Ecgi.Ecid)
	assert.Equal(t, "test-plmnid", event.GetRadioMeasReportPerCell().Ecgi.PlmnId)

	telemetry2 := sb.TelemetryMessage{
		MessageType: sb.MessageType_RADIO_MEAS_REPORT_PER_CELL,
		S: &sb.TelemetryMessage_RadioMeasReportPerCell{
			RadioMeasReportPerCell: &sb.RadioMeasReportPerCell{
				Ecgi: &sb.ECGI{
					Ecid:   "test-ecid-2",
					PlmnId: "test-plmnid-2",
				},
			},
		},
	}

	err = testStore.Put(telemetry2)
	assert.Nil(t, err)
	ecigTest2 := sb.ECGI{
		Ecid:   "test-ecid-2",
		PlmnId: "test-plmnid-2",
	}
	id2 := ID{
		Ecid:        ecigTest2.GetEcid(),
		PlmnID:      ecigTest2.GetPlmnId(),
		MessageType: telemetry2.GetMessageType(),
	}
	value2, err := testStore.Get(id2)
	assert.Nil(t, err)
	assert.Equal(t, value2.MessageType.String(), sb.MessageType_RADIO_MEAS_REPORT_PER_CELL.String())

	event = <-watchCh1
	assert.Equal(t, "test-ecid-2", event.GetRadioMeasReportPerCell().Ecgi.Ecid)
	assert.Equal(t, "test-plmnid-2", event.GetRadioMeasReportPerCell().Ecgi.PlmnId)

	event = <-watchCh2
	assert.Equal(t, "test-ecid-2", event.GetRadioMeasReportPerCell().Ecgi.Ecid)
	assert.Equal(t, "test-plmnid-2", event.GetRadioMeasReportPerCell().Ecgi.PlmnId)

	telemetry3 := sb.TelemetryMessage{
		MessageType: sb.MessageType_RADIO_MEAS_REPORT_PER_UE,
		S: &sb.TelemetryMessage_RadioMeasReportPerUE{
			RadioMeasReportPerUE: &sb.RadioMeasReportPerUE{
				Ecgi: &sb.ECGI{
					Ecid:   "test-ecid-3",
					PlmnId: "test-plmnid-3",
				},
				Crnti: "test-crnti-3",
			},
		},
	}
	err = testStore.Put(telemetry3)
	assert.Nil(t, err)

	ecigTest3 := sb.ECGI{
		Ecid:   "test-ecid-3",
		PlmnId: "test-plmnid-3",
	}
	id3 := ID{
		Ecid:        ecigTest3.GetEcid(),
		PlmnID:      ecigTest3.GetPlmnId(),
		Crnti:       "test-crnti-3",
		MessageType: telemetry3.GetMessageType(),
	}
	value3, err := testStore.Get(id3)
	assert.Nil(t, err)
	assert.Equal(t, value3.MessageType.String(), sb.MessageType_RADIO_MEAS_REPORT_PER_UE.String())

	event = <-watchCh1
	assert.Equal(t, "test-ecid-3", event.GetRadioMeasReportPerUE().Ecgi.Ecid)
	assert.Equal(t, "test-plmnid-3", event.GetRadioMeasReportPerUE().Ecgi.PlmnId)

	event = <-watchCh2
	assert.Equal(t, "test-ecid-3", event.GetRadioMeasReportPerUE().Ecgi.Ecid)
	assert.Equal(t, "test-plmnid-3", event.GetRadioMeasReportPerUE().Ecgi.PlmnId)

}
