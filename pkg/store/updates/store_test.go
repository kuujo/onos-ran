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

	"github.com/onosproject/onos-ric/pkg/store/mastership"
	timestore "github.com/onosproject/onos-ric/pkg/store/time"

	"github.com/onosproject/onos-ric/api/sb"
	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	mastershipStore, err := mastership.NewLocalStore("test", "1")
	assert.NoError(t, err)
	timeStore, err := timestore.NewStore(mastershipStore)
	assert.NoError(t, err)
	testStore, err := NewLocalStore(timeStore)
	assert.NoError(t, err)
	assert.NotNil(t, testStore)

	defer testStore.Close()
	defer timeStore.Close()
	defer testStore.Close()

	watchCh2 := make(chan Event)
	err = testStore.Watch(watchCh2, WithReplay())
	assert.NoError(t, err)

	controlUpdate3 := &sb.ControlUpdate{
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
	id3 := ID(sb.NewID(controlUpdate3.GetMessageType(), ecigTest3.GetPlmnId(), ecigTest3.GetEcid(), "test-crnti-3"))
	value3, err := testStore.Get(id3)
	assert.Nil(t, err)
	assert.Equal(t, value3.MessageType.String(), sb.MessageType_UE_ADMISSION_STATUS.String())

	event := <-watchCh2
	assert.Equal(t, "test-ecid-3", event.Message.GetUEAdmissionStatus().Ecgi.Ecid)
	assert.Equal(t, "test-plmnid-3", event.Message.GetUEAdmissionStatus().Ecgi.PlmnId)

	err = testStore.Delete(id3)
	assert.Nil(t, err)

}
