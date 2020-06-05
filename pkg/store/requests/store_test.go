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
	"github.com/golang/mock/gomock"
	"github.com/onosproject/onos-lib-go/pkg/cluster"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-ric/api/store/requests"
	"github.com/onosproject/onos-ric/pkg/config"
	"github.com/onosproject/onos-ric/pkg/store/device"
	"github.com/onosproject/onos-ric/pkg/store/mastership"
	"github.com/onosproject/onos-ric/test/mocks/store/device"
	"google.golang.org/grpc"
	"testing"
	"time"

	"github.com/onosproject/onos-ric/api/sb"
	"github.com/onosproject/onos-ric/api/sb/e2ap"
	"github.com/onosproject/onos-ric/api/sb/e2sm"
	"github.com/stretchr/testify/assert"
)

func TestStoreRequestNoBackup(t *testing.T) {
	logging.SetLevel(logging.DebugLevel)

	factory := cluster.NewTestFactory(func(id cluster.NodeID, server *grpc.Server) {
		requests.RegisterRequestsServiceServer(server, newServer())
	})

	node1 := cluster.NodeID("node-1")

	deviceECGI := sb.ECGI{
		PlmnId: "test",
		Ecid:   "test-1",
	}
	deviceID := device.ID(deviceECGI)

	deviceStore := mock_device_store.NewMockStore(gomock.NewController(t))
	deviceStore.EXPECT().Watch(gomock.Any()).DoAndReturn(func(ch chan<- device.Event) error {
		go func() {
			ch <- device.Event{
				Type: device.EventNone,
				Device: device.Device{
					ID: deviceID,
				},
			}
		}()
		return nil
	}).AnyTimes()

	cluster1, err := factory.NewCluster(node1)
	assert.NoError(t, err)
	mastershipStore1, err := mastership.NewLocalStore("TestStoreRequests", node1)
	assert.NoError(t, err)
	store1, err := NewDistributedStore(cluster1, deviceStore, mastershipStore1, config.Config{})
	assert.NoError(t, err)
	assert.NotNil(t, store1)

	time.Sleep(time.Second)

	defer store1.Close()

	watchCh := make(chan Event)
	err = store1.Watch(deviceID, watchCh, WithReplay())
	assert.NoError(t, err)

	err = store1.Append(New(&e2ap.RicControlRequest{
		Hdr: &e2sm.RicControlHeader{
			MessageType: sb.MessageType_HO_REQUEST,
			Ecgi:        &deviceECGI,
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
					Crntis: []string{"2"},
				},
			},
		},
	}))
	assert.NoError(t, err)

	err = store1.Append(New(&e2ap.RicControlRequest{
		Hdr: &e2sm.RicControlHeader{
			MessageType: sb.MessageType_HO_REQUEST,
			Ecgi:        &deviceECGI,
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
	}))
	assert.NoError(t, err)

	var request1 *Request
	select {
	case event := <-watchCh:
		assert.Equal(t, EventAppend, event.Type)
		assert.NoError(t, err)
		request1 = &event.Request
	case <-time.After(5 * time.Second):
		t.Log("Timed out waiting for EventAppend")
		t.FailNow()
	}

	var request2 *Request
	select {
	case event := <-watchCh:
		assert.Equal(t, EventAppend, event.Type)
		assert.NoError(t, err)
		request2 = &event.Request
	case <-time.After(5 * time.Second):
		t.Log("Timed out waiting for EventAppend")
		t.FailNow()
	}

	err = store1.Ack(request1)
	assert.NoError(t, err)

	select {
	case event := <-watchCh:
		assert.Equal(t, EventAck, event.Type)
	case <-time.After(5 * time.Second):
		t.Log("Timed out waiting for EventAppend")
		t.FailNow()
	}

	err = store1.Ack(request2)
	assert.NoError(t, err)

	select {
	case event := <-watchCh:
		assert.Equal(t, EventAck, event.Type)
	case <-time.After(5 * time.Second):
		t.Log("Timed out waiting for EventAppend")
		t.FailNow()
	}
}

func TestStoreRequestsSyncBackup(t *testing.T) {
	logging.SetLevel(logging.DebugLevel)

	factory := cluster.NewTestFactory(func(id cluster.NodeID, server *grpc.Server) {
		requests.RegisterRequestsServiceServer(server, newServer())
	})

	node1 := cluster.NodeID("node-1")
	node2 := cluster.NodeID("node-2")

	deviceECGI := sb.ECGI{
		PlmnId: "test",
		Ecid:   "test-1",
	}
	deviceID := device.ID(deviceECGI)

	deviceStore := mock_device_store.NewMockStore(gomock.NewController(t))
	deviceStore.EXPECT().Watch(gomock.Any()).DoAndReturn(func(ch chan<- device.Event) error {
		go func() {
			ch <- device.Event{
				Type: device.EventNone,
				Device: device.Device{
					ID: deviceID,
				},
			}
		}()
		return nil
	}).AnyTimes()

	one := 1
	cluster1, err := factory.NewCluster(node1)
	assert.NoError(t, err)
	mastershipStore1, err := mastership.NewLocalStore("TestStoreRequests", node1)
	assert.NoError(t, err)
	store1, err := NewDistributedStore(cluster1, deviceStore, mastershipStore1, config.Config{Stores: config.StoresConfig{Requests: config.RequestsStoreConfig{SyncBackups: &one}}})
	assert.NoError(t, err)
	assert.NotNil(t, store1)

	cluster2, err := factory.NewCluster(node2)
	assert.NoError(t, err)
	mastershipStore2, err := mastership.NewLocalStore("TestStoreRequests", node2)
	assert.NoError(t, err)
	store2, err := NewDistributedStore(cluster2, deviceStore, mastershipStore2, config.Config{Stores: config.StoresConfig{Requests: config.RequestsStoreConfig{SyncBackups: &one}}})
	assert.NoError(t, err)
	assert.NotNil(t, store2)

	time.Sleep(time.Second)

	defer store1.Close()
	defer store2.Close()

	watchCh := make(chan Event)
	err = store1.Watch(deviceID, watchCh, WithReplay())
	assert.NoError(t, err)

	err = store1.Append(New(&e2ap.RicControlRequest{
		Hdr: &e2sm.RicControlHeader{
			MessageType: sb.MessageType_HO_REQUEST,
			Ecgi:        &deviceECGI,
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
	}))
	assert.NoError(t, err)

	select {
	case event := <-watchCh:
		assert.Equal(t, EventAppend, event.Type)
		err = store1.Ack(&event.Request)
		assert.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Log("Timed out waiting for EventAppend")
		t.FailNow()
	}

	select {
	case event := <-watchCh:
		assert.Equal(t, EventAck, event.Type)
	case <-time.After(5 * time.Second):
		t.Log("Timed out waiting for EventAck")
		t.FailNow()
	}

	err = store1.Append(New(&e2ap.RicControlRequest{
		Hdr: &e2sm.RicControlHeader{
			MessageType: sb.MessageType_HO_REQUEST,
			Ecgi:        &deviceECGI,
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
	}))
	assert.NoError(t, err)

	err = store1.Append(New(&e2ap.RicControlRequest{
		Hdr: &e2sm.RicControlHeader{
			MessageType: sb.MessageType_HO_REQUEST,
			Ecgi:        &deviceECGI,
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
	}))
	assert.NoError(t, err)

	select {
	case event := <-watchCh:
		assert.Equal(t, EventAppend, event.Type)
		err = store1.Ack(&event.Request)
		assert.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Log("Timed out waiting for EventAppend")
		t.FailNow()
	}

	select {
	case event := <-watchCh:
		assert.Equal(t, EventAppend, event.Type)
		err = store1.Ack(&event.Request)
		assert.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Log("Timed out waiting for EventAppend")
		t.FailNow()
	}

	select {
	case event := <-watchCh:
		assert.Equal(t, EventAck, event.Type)
	case <-time.After(5 * time.Second):
		t.Log("Timed out waiting for EventAck")
		t.FailNow()
	}

	select {
	case event := <-watchCh:
		assert.Equal(t, EventAck, event.Type)
	case <-time.After(5 * time.Second):
		t.Log("Timed out waiting for EventAck")
		t.FailNow()
	}
}

func TestStoreRequestsAsyncBackup(t *testing.T) {
	logging.SetLevel(logging.DebugLevel)

	factory := cluster.NewTestFactory(func(id cluster.NodeID, server *grpc.Server) {
		requests.RegisterRequestsServiceServer(server, newServer())
	})

	node1 := cluster.NodeID("node-1")
	node2 := cluster.NodeID("node-2")

	deviceECGI := sb.ECGI{
		PlmnId: "test",
		Ecid:   "test-1",
	}
	deviceID := device.ID(deviceECGI)

	deviceStore := mock_device_store.NewMockStore(gomock.NewController(t))
	deviceStore.EXPECT().Watch(gomock.Any()).DoAndReturn(func(ch chan<- device.Event) error {
		go func() {
			ch <- device.Event{
				Type: device.EventNone,
				Device: device.Device{
					ID: deviceID,
				},
			}
		}()
		return nil
	}).AnyTimes()

	zero := 0
	cluster1, err := factory.NewCluster(node1)
	assert.NoError(t, err)
	mastershipStore1, err := mastership.NewLocalStore("TestStoreRequests", node1)
	assert.NoError(t, err)
	store1, err := NewDistributedStore(cluster1, deviceStore, mastershipStore1, config.Config{Stores: config.StoresConfig{Requests: config.RequestsStoreConfig{SyncBackups: &zero}}})
	assert.NoError(t, err)
	assert.NotNil(t, store1)

	cluster2, err := factory.NewCluster(node2)
	assert.NoError(t, err)
	mastershipStore2, err := mastership.NewLocalStore("TestStoreRequests", node2)
	assert.NoError(t, err)
	store2, err := NewDistributedStore(cluster2, deviceStore, mastershipStore2, config.Config{Stores: config.StoresConfig{Requests: config.RequestsStoreConfig{SyncBackups: &zero}}})
	assert.NoError(t, err)
	assert.NotNil(t, store2)

	time.Sleep(time.Second)

	defer store1.Close()
	defer store2.Close()

	watchCh := make(chan Event)
	err = store1.Watch(deviceID, watchCh, WithReplay())
	assert.NoError(t, err)

	err = store1.Append(New(&e2ap.RicControlRequest{
		Hdr: &e2sm.RicControlHeader{
			MessageType: sb.MessageType_HO_REQUEST,
			Ecgi:        &deviceECGI,
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
	}))
	assert.NoError(t, err)

	select {
	case event := <-watchCh:
		assert.Equal(t, EventAppend, event.Type)
		err = store1.Ack(&event.Request)
		assert.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Log("Timed out waiting for EventAppend")
		t.FailNow()
	}

	select {
	case event := <-watchCh:
		assert.Equal(t, EventAck, event.Type)
	case <-time.After(5 * time.Second):
		t.Log("Timed out waiting for EventAck")
		t.FailNow()
	}

	err = store1.Append(New(&e2ap.RicControlRequest{
		Hdr: &e2sm.RicControlHeader{
			MessageType: sb.MessageType_HO_REQUEST,
			Ecgi:        &deviceECGI,
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
	}))
	assert.NoError(t, err)

	err = store1.Append(New(&e2ap.RicControlRequest{
		Hdr: &e2sm.RicControlHeader{
			MessageType: sb.MessageType_HO_REQUEST,
			Ecgi:        &deviceECGI,
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
	}))
	assert.NoError(t, err)

	select {
	case event := <-watchCh:
		assert.Equal(t, EventAppend, event.Type)
		err = store1.Ack(&event.Request)
		assert.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Log("Timed out waiting for EventAppend")
		t.FailNow()
	}

	select {
	case event := <-watchCh:
		assert.Equal(t, EventAppend, event.Type)
		err = store1.Ack(&event.Request)
		assert.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Log("Timed out waiting for EventAppend")
		t.FailNow()
	}

	select {
	case event := <-watchCh:
		assert.Equal(t, EventAck, event.Type)
	case <-time.After(5 * time.Second):
		t.Log("Timed out waiting for EventAck")
		t.FailNow()
	}

	select {
	case event := <-watchCh:
		assert.Equal(t, EventAck, event.Type)
	case <-time.After(5 * time.Second):
		t.Log("Timed out waiting for EventAck")
		t.FailNow()
	}
}

func TestStoreRequestsSyncAsyncBackup(t *testing.T) {
	logging.SetLevel(logging.DebugLevel)

	factory := cluster.NewTestFactory(func(id cluster.NodeID, server *grpc.Server) {
		requests.RegisterRequestsServiceServer(server, newServer())
	})

	node1 := cluster.NodeID("node-1")
	node2 := cluster.NodeID("node-2")
	node3 := cluster.NodeID("node-3")

	deviceECGI := sb.ECGI{
		PlmnId: "test",
		Ecid:   "test-1",
	}
	deviceID := device.ID(deviceECGI)

	deviceStore := mock_device_store.NewMockStore(gomock.NewController(t))
	deviceStore.EXPECT().Watch(gomock.Any()).DoAndReturn(func(ch chan<- device.Event) error {
		go func() {
			ch <- device.Event{
				Type: device.EventNone,
				Device: device.Device{
					ID: deviceID,
				},
			}
		}()
		return nil
	}).AnyTimes()

	one := 1
	cluster1, err := factory.NewCluster(node1)
	assert.NoError(t, err)
	mastershipStore1, err := mastership.NewLocalStore("TestStoreRequests", node1)
	assert.NoError(t, err)
	store1, err := NewDistributedStore(cluster1, deviceStore, mastershipStore1, config.Config{Stores: config.StoresConfig{Requests: config.RequestsStoreConfig{SyncBackups: &one}}})
	assert.NoError(t, err)
	assert.NotNil(t, store1)

	cluster2, err := factory.NewCluster(node2)
	assert.NoError(t, err)
	mastershipStore2, err := mastership.NewLocalStore("TestStoreRequests", node2)
	assert.NoError(t, err)
	store2, err := NewDistributedStore(cluster2, deviceStore, mastershipStore2, config.Config{Stores: config.StoresConfig{Requests: config.RequestsStoreConfig{SyncBackups: &one}}})
	assert.NoError(t, err)
	assert.NotNil(t, store2)

	cluster3, err := factory.NewCluster(node3)
	assert.NoError(t, err)
	mastershipStore3, err := mastership.NewLocalStore("TestStoreRequests", node3)
	assert.NoError(t, err)
	store3, err := NewDistributedStore(cluster3, deviceStore, mastershipStore3, config.Config{Stores: config.StoresConfig{Requests: config.RequestsStoreConfig{SyncBackups: &one}}})
	assert.NoError(t, err)
	assert.NotNil(t, store3)

	time.Sleep(time.Second)

	defer store1.Close()
	defer store2.Close()
	defer store3.Close()

	watchCh := make(chan Event)
	err = store1.Watch(deviceID, watchCh, WithReplay())
	assert.NoError(t, err)

	err = store1.Append(New(&e2ap.RicControlRequest{
		Hdr: &e2sm.RicControlHeader{
			MessageType: sb.MessageType_HO_REQUEST,
			Ecgi:        &deviceECGI,
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
	}))
	assert.NoError(t, err)

	select {
	case event := <-watchCh:
		assert.Equal(t, EventAppend, event.Type)
		err = store1.Ack(&event.Request)
		assert.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Log("Timed out waiting for EventAppend")
		t.FailNow()
	}

	select {
	case event := <-watchCh:
		assert.Equal(t, EventAck, event.Type)
	case <-time.After(5 * time.Second):
		t.Log("Timed out waiting for EventAck")
		t.FailNow()
	}

	err = store1.Append(New(&e2ap.RicControlRequest{
		Hdr: &e2sm.RicControlHeader{
			MessageType: sb.MessageType_HO_REQUEST,
			Ecgi:        &deviceECGI,
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
	}))
	assert.NoError(t, err)

	err = store1.Append(New(&e2ap.RicControlRequest{
		Hdr: &e2sm.RicControlHeader{
			MessageType: sb.MessageType_HO_REQUEST,
			Ecgi:        &deviceECGI,
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
	}))
	assert.NoError(t, err)

	select {
	case event := <-watchCh:
		assert.Equal(t, EventAppend, event.Type)
		err = store1.Ack(&event.Request)
		assert.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Log("Timed out waiting for EventAppend")
		t.FailNow()
	}

	select {
	case event := <-watchCh:
		assert.Equal(t, EventAppend, event.Type)
		err = store1.Ack(&event.Request)
		assert.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Log("Timed out waiting for EventAppend")
		t.FailNow()
	}

	select {
	case event := <-watchCh:
		assert.Equal(t, EventAck, event.Type)
	case <-time.After(5 * time.Second):
		t.Log("Timed out waiting for EventAck")
		t.FailNow()
	}

	select {
	case event := <-watchCh:
		assert.Equal(t, EventAck, event.Type)
	case <-time.After(5 * time.Second):
		t.Log("Timed out waiting for EventAck")
		t.FailNow()
	}
}
