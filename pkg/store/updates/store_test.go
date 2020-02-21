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

package updates

import (
	"testing"

	"github.com/onosproject/onos-ran/api/sb"
	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	testStore, err := NewStore()

	watchCh1 := make(chan sb.ControlUpdate)
	err = testStore.Watch(watchCh1)
	assert.NoError(t, err)

	assert.Nil(t, err)
	assert.NotNil(t, testStore)
	controlUpdate1 := sb.ControlUpdate{
		MessageType: sb.MessageType_CELL_CONFIG_REPORT,
		S: &sb.ControlUpdate_CellConfigReport{
			CellConfigReport: &sb.CellConfigReport{
				Ecgi: &sb.ECGI{
					Ecid:   "test-ecid",
					PlmnId: "test-plmnid",
				},
			},
		},
	}
	err = testStore.Put(controlUpdate1)
	assert.Nil(t, err)
	ecigTest1 := sb.ECGI{
		Ecid:   "test-ecid",
		PlmnId: "test-plmnid",
	}
	id1 := ID{
		Ecid:        ecigTest1.GetEcid(),
		PlmnID:      ecigTest1.GetPlmnId(),
		MessageType: controlUpdate1.GetMessageType(),
	}
	value1, err := testStore.Get(id1)
	assert.Nil(t, err)
	assert.Equal(t, value1.MessageType.String(), sb.MessageType_CELL_CONFIG_REPORT.String())

	event := <-watchCh1
	assert.Equal(t, "test-ecid", event.GetCellConfigReport().Ecgi.Ecid)
	assert.Equal(t, "test-plmnid", event.GetCellConfigReport().Ecgi.PlmnId)

	watchCh2 := make(chan sb.ControlUpdate)
	err = testStore.Watch(watchCh2, WithReplay())

	event = <-watchCh2
	assert.Equal(t, "test-ecid", event.GetCellConfigReport().Ecgi.Ecid)
	assert.Equal(t, "test-plmnid", event.GetCellConfigReport().Ecgi.PlmnId)

	controlUpdate2 := sb.ControlUpdate{
		MessageType: sb.MessageType_CELL_CONFIG_REPORT,
		S: &sb.ControlUpdate_CellConfigReport{
			CellConfigReport: &sb.CellConfigReport{
				Ecgi: &sb.ECGI{
					Ecid:   "test-ecid-2",
					PlmnId: "test-plmnid-2",
				},
			},
		},
	}
	err = testStore.Put(controlUpdate2)
	assert.Nil(t, err)

	ecigTest2 := sb.ECGI{
		Ecid:   "test-ecid-2",
		PlmnId: "test-plmnid-2",
	}
	id2 := ID{
		Ecid:        ecigTest2.GetEcid(),
		PlmnID:      ecigTest2.GetPlmnId(),
		MessageType: controlUpdate2.GetMessageType(),
	}
	value2, err := testStore.Get(id2)
	assert.Nil(t, err)
	assert.Equal(t, value2.MessageType.String(), sb.MessageType_CELL_CONFIG_REPORT.String())

	event = <-watchCh1
	assert.Equal(t, "test-ecid-2", event.GetCellConfigReport().Ecgi.Ecid)
	assert.Equal(t, "test-plmnid-2", event.GetCellConfigReport().Ecgi.PlmnId)

	event = <-watchCh2
	assert.Equal(t, "test-ecid-2", event.GetCellConfigReport().Ecgi.Ecid)
	assert.Equal(t, "test-plmnid-2", event.GetCellConfigReport().Ecgi.PlmnId)

	controlUpdate3 := sb.ControlUpdate{
		MessageType: sb.MessageType_UE_ADMISSION_STATUS,
		S: &sb.ControlUpdate_UEAdmissionStatus{
			UEAdmissionStatus: &sb.UEAdmissionStatus{
				Ecgi: &sb.ECGI{
					Ecid:   "test-ecid-3",
					PlmnId: "test-plmnid-3",
				},
				Crnti: "test-crnti-3",
			},
		},
	}
	err = testStore.Put(controlUpdate3)
	assert.Nil(t, err)

	ecigTest3 := sb.ECGI{
		Ecid:   "test-ecid-3",
		PlmnId: "test-plmnid-3",
	}
	id3 := ID{
		Ecid:        ecigTest3.GetEcid(),
		PlmnID:      ecigTest3.GetPlmnId(),
		Crnti:       "test-crnti-3",
		MessageType: controlUpdate3.GetMessageType(),
	}
	value3, err := testStore.Get(id3)
	assert.Nil(t, err)
	assert.Equal(t, value3.MessageType.String(), sb.MessageType_UE_ADMISSION_STATUS.String())

	event = <-watchCh1
	assert.Equal(t, "test-ecid-3", event.GetUEAdmissionStatus().Ecgi.Ecid)
	assert.Equal(t, "test-plmnid-3", event.GetUEAdmissionStatus().Ecgi.PlmnId)

	event = <-watchCh2
	assert.Equal(t, "test-ecid-3", event.GetUEAdmissionStatus().Ecgi.Ecid)
	assert.Equal(t, "test-plmnid-3", event.GetUEAdmissionStatus().Ecgi.PlmnId)

}
