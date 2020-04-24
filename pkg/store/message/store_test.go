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

package message

import (
	"github.com/onosproject/onos-ric/api/store/message"
	"github.com/onosproject/onos-ric/pkg/store/mastership"
	timestore "github.com/onosproject/onos-ric/pkg/store/time"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_PutGetListDelete(t *testing.T) {
	mastershipStore, err := mastership.NewLocalStore("test1", "1")
	assert.NoError(t, err)
	timeStore, err := timestore.NewStore(mastershipStore)
	assert.NoError(t, err)
	testStore, err := NewLocalStore("foo", timeStore)
	assert.NoError(t, err)
	assert.NotNil(t, testStore)

	defer testStore.Close()
	defer timeStore.Close()
	defer testStore.Close()

	// Add two things to the store
	pk1 := NewPartitionKey("thing", "1")
	id1 := NewID("thing", "1")
	entry1 := message.MessageEntry{Id: id1.String(), PartitionKey: pk1.String()}

	pk2 := NewPartitionKey("thing", "2")
	id2 := NewID("thing", "2")
	entry2 := message.MessageEntry{Id: id2.String(), PartitionKey: pk2.String()}

	err = testStore.Put(pk1, id1, &entry1)
	assert.NoError(t, err)

	err = testStore.Put(pk2, id2, &entry2)
	assert.NoError(t, err)

	// Get thing1
	v, err := testStore.Get(pk1, id1)
	assert.NoError(t, err)
	assert.Equal(t, id1.String(), v.Id)

	v, err = testStore.Get(pk1, id1, WithRevision(Revision{1, 1}))
	assert.NoError(t, err)
	assert.Equal(t, id1.String(), v.Id)

	err = testStore.Delete(pk2, id2)
	assert.NoError(t, err)

	lch := make(chan message.MessageEntry)
	err = testStore.List(lch)
	assert.NoError(t, err)

	le := <-lch
	assert.Equal(t, pk1.String(), le.PartitionKey)
	assert.Equal(t, id1.String(), le.Id)

	err = testStore.Clear()
	assert.NoError(t, err)

	err = testStore.Close()
	assert.NoError(t, err)
}

func Test_Watch(t *testing.T) {
	mastershipStore, err := mastership.NewLocalStore("test2", "1")
	assert.NoError(t, err)
	timeStore, err := timestore.NewStore(mastershipStore)
	assert.NoError(t, err)
	testStore, err := NewLocalStore("foo", timeStore)
	assert.NoError(t, err)
	assert.NotNil(t, testStore)

	defer testStore.Close()
	defer timeStore.Close()
	defer testStore.Close()

	wch := make(chan Event)
	err = testStore.Watch(wch, WithReplay())
	assert.NoError(t, err)

	pk1 := NewPartitionKey("thing", "1")
	id1 := NewID("thing", "1")
	entry1 := message.MessageEntry{Id: id1.String(), PartitionKey: pk1.String()}

	err = testStore.Put(pk1, id1, &entry1)
	assert.NoError(t, err)

	we := <-wch
	assert.Equal(t, EventInsert, we.Type)
	assert.Equal(t, pk1.String(), we.Message.PartitionKey)
	assert.Equal(t, id1.String(), we.Message.Id)
}
