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

package cluster

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetNodeIDNoEnv(t *testing.T) {
	nodeID := GetNodeID()
	assert.Equal(t, "onos-ric", string(nodeID))
}

func TestGetNodeIdEnv(t *testing.T) {
	myNodeID := "my-node-id"
	assert.NoError(t, os.Setenv(nodeIDEnv, myNodeID))
	nodeID := GetNodeID()
	assert.Equal(t, myNodeID, string(nodeID))
}
