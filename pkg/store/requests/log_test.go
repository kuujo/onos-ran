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
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestWriteReadBatch(t *testing.T) {
	log := newLog()
	assert.Equal(t, Index(0), log.Writer().Index())

	entry := log.Writer().Write("foo")
	assert.Equal(t, Index(1), entry.Index)
	assert.Equal(t, "foo", entry.Value)
	assert.Equal(t, Index(1), log.Writer().Index())

	entry = log.Writer().Write("bar")
	assert.Equal(t, Index(2), entry.Index)
	assert.Equal(t, "bar", entry.Value)
	assert.Equal(t, Index(2), log.Writer().Index())

	reader := log.OpenReader(0)
	batch := reader.ReadBatch()
	assert.Equal(t, Index(0), batch.PrevIndex)
	assert.Len(t, batch.Entries, 2)
	assert.Equal(t, Index(1), batch.Entries[0].Index)
	assert.Equal(t, "foo", batch.Entries[0].Value)
	assert.Equal(t, Index(2), batch.Entries[1].Index)
	assert.Equal(t, "bar", batch.Entries[1].Value)

	go func() {
		time.Sleep(time.Second)
		entry = log.Writer().Write("baz")
		assert.Equal(t, Index(3), entry.Index)
		assert.Equal(t, "baz", entry.Value)
		assert.Equal(t, Index(3), log.Writer().Index())
	}()

	batch = reader.ReadBatch()
	assert.Equal(t, Index(2), batch.PrevIndex)
	assert.Len(t, batch.Entries, 1)
	assert.Equal(t, Index(3), batch.Entries[0].Index)
	assert.Equal(t, "baz", batch.Entries[0].Value)
}

func TestWriteAwaitBatch(t *testing.T) {
	log := newLog()
	assert.Equal(t, Index(0), log.Writer().Index())

	entry := log.Writer().Write("foo")
	assert.Equal(t, Index(1), entry.Index)
	assert.Equal(t, "foo", entry.Value)
	assert.Equal(t, Index(1), log.Writer().Index())

	entry = log.Writer().Write("bar")
	assert.Equal(t, Index(2), entry.Index)
	assert.Equal(t, "bar", entry.Value)
	assert.Equal(t, Index(2), log.Writer().Index())

	reader := log.OpenReader(0)
	batch := <-reader.AwaitBatch()
	assert.Equal(t, Index(0), batch.PrevIndex)
	assert.Len(t, batch.Entries, 2)
	assert.Equal(t, Index(1), batch.Entries[0].Index)
	assert.Equal(t, "foo", batch.Entries[0].Value)
	assert.Equal(t, Index(2), batch.Entries[1].Index)
	assert.Equal(t, "bar", batch.Entries[1].Value)

	go func() {
		time.Sleep(time.Second)
		entry = log.Writer().Write("baz")
		assert.Equal(t, Index(3), entry.Index)
		assert.Equal(t, "baz", entry.Value)
		assert.Equal(t, Index(3), log.Writer().Index())
	}()

	batch = <-reader.AwaitBatch()
	assert.Equal(t, Index(2), batch.PrevIndex)
	assert.Len(t, batch.Entries, 1)
	assert.Equal(t, Index(3), batch.Entries[0].Index)
	assert.Equal(t, "baz", batch.Entries[0].Value)
}
