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

package dispatcher

import (
	"github.com/onosproject/onos-ric/pkg/store/cluster"
	"github.com/onosproject/onos-ric/pkg/store/mastership"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type testKeyResolver struct {
}

func (r testKeyResolver) Resolve(request Request) (mastership.Key, error) {
	return mastership.Key(request.ID), nil
}

func TestMastershipFilter(t *testing.T) {
	node1 := cluster.NodeID("node1")
	node2 := cluster.NodeID("node2")

	store1, err := mastership.NewLocalStore("TestMastershipFilter", node1)
	assert.NoError(t, err)

	store2, err := mastership.NewLocalStore("TestMastershipFilter", node2)
	assert.NoError(t, err)

	filter1 := &MastershipFilter{
		Store:    store1,
		Resolver: testKeyResolver{},
	}

	filter2 := &MastershipFilter{
		Store:    store2,
		Resolver: testKeyResolver{},
	}

	key1 := mastership.Key("key1")
	key2 := mastership.Key("key2")

	election1, err := store1.GetElection(key1)
	assert.NoError(t, err)
	master, err := election1.IsMaster()
	assert.NoError(t, err)
	assert.True(t, master)

	election1, err = store2.GetElection(key1)
	assert.NoError(t, err)
	master, err = election1.IsMaster()
	assert.NoError(t, err)
	assert.False(t, master)

	election2, err := store2.GetElection(key2)
	assert.NoError(t, err)
	master, err = election2.IsMaster()
	assert.NoError(t, err)
	assert.True(t, master)

	election2, err = store1.GetElection(key2)
	assert.NoError(t, err)
	master, err = election2.IsMaster()
	assert.NoError(t, err)
	assert.False(t, master)

	assert.True(t, filter1.Accept(Request{string(key1)}))
	assert.False(t, filter2.Accept(Request{string(key1)}))
	assert.True(t, filter2.Accept(Request{string(key2)}))
	assert.False(t, filter1.Accept(Request{string(key2)}))

	ch := make(chan mastership.State)
	election1, err = store2.GetElection(key1)
	assert.NoError(t, err)
	err = election1.Watch(ch)
	assert.NoError(t, err)

	err = store1.Close()
	assert.NoError(t, err)

	select {
	case <-ch:
	case <-time.After(5 * time.Second):
		t.FailNow()
	}

	assert.True(t, filter2.Accept(Request{string(key1)}))
	assert.True(t, filter2.Accept(Request{string(key2)}))
}
