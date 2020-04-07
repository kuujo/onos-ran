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

package time

import (
	"github.com/onosproject/onos-ric/pkg/store/mastership"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStore(t *testing.T) {
	mastershipStore, err := mastership.NewLocalStore("test", "1")
	assert.NoError(t, err)
	timeStore, err := NewStore(mastershipStore)
	assert.NoError(t, err)
	assert.NotNil(t, timeStore)

	defer timeStore.Close()
	defer mastershipStore.Close()

	clock, err := timeStore.GetLogicalClock("foo")
	assert.NoError(t, err)

	timestamp, err := clock.GetTimestamp()
	assert.NoError(t, err)
	assert.Equal(t, Epoch(1), timestamp.Epoch)
	assert.Equal(t, LogicalTime(1), timestamp.LogicalTime)

	timestamp, err = clock.GetTimestamp()
	assert.NoError(t, err)
	assert.Equal(t, Epoch(1), timestamp.Epoch)
	assert.Equal(t, LogicalTime(2), timestamp.LogicalTime)
}
