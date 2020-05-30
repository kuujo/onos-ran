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

package mastership

import (
	"github.com/onosproject/onos-lib-go/pkg/atomix"
	"github.com/onosproject/onos-lib-go/pkg/cluster"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMastershipElection(t *testing.T) {
	_, address := atomix.StartLocalNode()

	store1, err := newLocalElection(1, "a", address)
	assert.NoError(t, err)

	store2, err := newLocalElection(1, "b", address)
	assert.NoError(t, err)

	store2Ch := make(chan State)
	err = store2.Watch(store2Ch)
	assert.NoError(t, err)

	store3, err := newLocalElection(1, "c", address)
	assert.NoError(t, err)

	store3Ch := make(chan State)
	err = store3.Watch(store3Ch)
	assert.NoError(t, err)

	master, err := store1.IsMaster()
	assert.NoError(t, err)
	assert.True(t, master)

	master, err = store2.IsMaster()
	assert.NoError(t, err)
	assert.False(t, master)

	master, err = store3.IsMaster()
	assert.NoError(t, err)
	assert.False(t, master)

	err = store1.Close()
	assert.NoError(t, err)

	mastership := <-store2Ch
	assert.Equal(t, cluster.NodeID("a"), mastership.Master)
	assert.Len(t, mastership.Replicas, 3)

	mastership = <-store2Ch
	assert.Equal(t, cluster.NodeID("b"), mastership.Master)
	assert.Len(t, mastership.Replicas, 2)

	master, err = store2.IsMaster()
	assert.NoError(t, err)
	assert.True(t, master)

	mastership = <-store3Ch
	assert.Equal(t, cluster.NodeID("b"), mastership.Master)

	master, err = store3.IsMaster()
	assert.NoError(t, err)
	assert.False(t, master)

	err = store2.Close()
	assert.NoError(t, err)

	mastership = <-store3Ch
	assert.Equal(t, cluster.NodeID("c"), mastership.Master)

	master, err = store3.IsMaster()
	assert.NoError(t, err)
	assert.True(t, master)

	_ = store3.Close()
}
