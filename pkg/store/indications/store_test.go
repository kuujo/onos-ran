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
	"github.com/golang/mock/gomock"
	"github.com/onosproject/onos-lib-go/pkg/cluster"
	"github.com/onosproject/onos-ric/api/sb"
	"github.com/onosproject/onos-ric/api/sb/e2ap"
	"github.com/onosproject/onos-ric/api/sb/e2sm"
	"github.com/onosproject/onos-ric/api/store/indications"
	"github.com/onosproject/onos-ric/pkg/config"
	"github.com/onosproject/onos-ric/pkg/store/device"
	"github.com/onosproject/onos-ric/pkg/store/mastership"
	"github.com/onosproject/onos-ric/test/mocks/store/device"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"testing"
)

func TestStoreIndications(t *testing.T) {
	factory := cluster.NewTestFactory(func(id cluster.NodeID, server *grpc.Server) {
		indications.RegisterIndicationsServiceServer(server, newServer())
	})

	node1 := cluster.NodeID("node-1")
	node2 := cluster.NodeID("node-2")

	deviceStore := mock_device_store.NewMockStore(gomock.NewController(t))
	deviceStore.EXPECT().Watch(gomock.Any()).DoAndReturn(func(ch chan<- device.Event) error {
		go func() {
			ch <- device.Event{
				Type: device.EventNone,
				Device: device.Device{
					ID: device.ID{
						Ecid:   "test-ecid-1",
						PlmnId: "test-plmnid-1",
					},
				},
			}
		}()
		return nil
	}).AnyTimes()

	cluster1, err := factory.NewCluster(node1)
	assert.NoError(t, err)
	mastershipStore1, err := mastership.NewLocalStore("TestStoreIndications", node1)
	assert.NoError(t, err)
	store1, err := NewDistributedStore(cluster1, deviceStore, mastershipStore1, config.Config{})
	assert.NoError(t, err)
	assert.NotNil(t, store1)

	cluster2, err := factory.NewCluster(node2)
	assert.NoError(t, err)
	mastershipStore2, err := mastership.NewLocalStore("TestStoreIndications", node2)
	assert.NoError(t, err)
	store2, err := NewDistributedStore(cluster2, deviceStore, mastershipStore2, config.Config{})
	assert.NoError(t, err)
	assert.NotNil(t, store2)

	defer store1.Close()
	defer store2.Close()

	subscribeCh2 := make(chan Event)
	err = store2.Subscribe(subscribeCh2, WithReplay())
	assert.NoError(t, err)

	indication1 := &e2ap.RicIndication{
		Hdr: &e2sm.RicIndicationHeader{
			MessageType: sb.MessageType_UE_ADMISSION_REQUEST,
		},
		Msg: &e2sm.RicIndicationMessage{
			S: &e2sm.RicIndicationMessage_UEAdmissionRequest{
				UEAdmissionRequest: &sb.UEAdmissionRequest{
					Ecgi: &sb.ECGI{
						Ecid:   "test-ecid-1",
						PlmnId: "test-plmnid-1",
					},
					Crnti: "test-crnti-1",
				},
			},
		},
	}
	err = store1.Record(New(indication1))
	assert.Nil(t, err)

	event := <-subscribeCh2
	assert.Equal(t, "test-ecid-1", event.Indication.GetMsg().GetUEAdmissionRequest().Ecgi.Ecid)
	assert.Equal(t, "test-plmnid-1", event.Indication.GetMsg().GetUEAdmissionRequest().Ecgi.PlmnId)
}
