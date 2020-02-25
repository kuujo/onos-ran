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

package nb

import (
	"github.com/onosproject/onos-ric/api/nb"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
	"time"
)

// TestNBUELinksAPI tests the NB stations API
func (s *TestSuite) TestNBUELinksAPI(t *testing.T) {
	const expectedPLMNID = "001001"

	// Wait for simulator to respond
	assert.NoError(t, waitForSimulator())

	var ids map[string]*nb.UELinkInfo

	for attempt := 1; attempt <= 10; attempt++ {
		ids = readLinks(t)

		if len(ids) != 0 {
			break
		}
		time.Sleep(5 * time.Second)
	}

	if len(ids) == 0 {
		assert.FailNow(t, "No links were received")
	}

	// Make sure the data returned are correct
	re := regexp.MustCompile("00000[0-9][0-9]")

	for _, link := range ids {
		assert.Equal(t, expectedPLMNID, link.Ecgi.Plmnid)
		assert.True(t, re.MatchString(link.Ecgi.Ecid))

		for _, cq := range link.ChannelQualities {
			assert.Equal(t, expectedPLMNID, cq.TargetEcgi.Plmnid)
			assert.True(t, re.MatchString(cq.TargetEcgi.Ecid))
			assert.True(t, cq.CqiHist < 1000)
		}
	}
}
