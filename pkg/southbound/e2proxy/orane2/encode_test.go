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

package orane2

import (
	"gotest.tools/assert"
	"testing"
)

func Test_xerEncodeRicInd(t *testing.T) {
	ri, err := buildRicIndication()
	assert.NilError(t, err)

	//err = dumpXer(10, ri)
	//assert.NilError(t, err)

	bytes, err := EncodeRicIndXer(10, ri)
	assert.NilError(t, err)
	t.Logf("Xer encoded \n%s", string(bytes))
}

func Test_xerEncodeRicIndDump(t *testing.T) {
	ri, err := buildRicIndication()
	assert.NilError(t, err)

	err = dumpXer(ri)
	assert.NilError(t, err)
}

func Test_perEncodeRicInd(t *testing.T) {
	ri, err := buildRicIndication()
	assert.NilError(t, err)

	bytes, err := EncRicIndToBuffer(ri)
	assert.NilError(t, err)
	assert.Equal(t, 125, len(bytes))
	t.Logf("Per encoded \n%v", bytes)
}
