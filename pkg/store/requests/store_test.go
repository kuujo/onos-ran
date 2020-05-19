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

package requests

import (
	"github.com/onosproject/onos-ric/pkg/store/device"
	"testing"
	"time"

	"github.com/onosproject/onos-ric/api/sb"
	"github.com/onosproject/onos-ric/api/sb/e2ap"
	"github.com/onosproject/onos-ric/api/sb/e2sm"
	"github.com/stretchr/testify/assert"
)

func TestStoreRequests(t *testing.T) {
	testStore, err := NewLocalStore()
	assert.NoError(t, err)
	assert.NotNil(t, testStore)

	defer testStore.Close()

	deviceID := device.ID{
		PlmnId: "test",
		Ecid:   "test-1",
	}

	watchCh := make(chan Event)
	err = testStore.Watch(deviceID, watchCh, WithReplay())
	assert.NoError(t, err)

	err = testStore.Append(&Request{
		DeviceID: deviceID,
		RicControlRequest: &e2ap.RicControlRequest{
			Hdr: &e2sm.RicControlHeader{
				MessageType: sb.MessageType_HO_REQUEST,
			},
			Msg: &e2sm.RicControlMessage{
				S: &e2sm.RicControlMessage_HORequest{
					HORequest: &sb.HORequest{
						Crnti: "1",
						EcgiS: &sb.ECGI{
							PlmnId: "test",
							Ecid:   "test-1",
						},
						EcgiT: &sb.ECGI{
							PlmnId: "test",
							Ecid:   "test-2",
						},
						Crntis: []string{"1"},
					},
				},
			},
		},
	})
	assert.NoError(t, err)

	select {
	case event := <-watchCh:
		assert.Equal(t, EventAppend, event.Type)
		assert.Equal(t, deviceID, event.Request.DeviceID)
		err = testStore.Ack(&event.Request)
		assert.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Log("Timed out waiting for EventAppend")
		t.FailNow()
	}

	select {
	case event := <-watchCh:
		assert.Equal(t, EventRemove, event.Type)
	case <-time.After(5 * time.Second):
		t.Log("Timed out waiting for EventRemove")
		t.FailNow()
	}
}
