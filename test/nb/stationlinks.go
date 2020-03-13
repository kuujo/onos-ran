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
	"context"
	"github.com/onosproject/onos-ric/api/nb"
	"github.com/stretchr/testify/assert"
	"io"
	"regexp"
	"testing"
	"time"
)

func readStationLinks(t *testing.T) map[string]*nb.StationLinkInfo {
	ids := make(map[string]*nb.StationLinkInfo)

	// Make a client to connect to the onos-ric northbound API
	client := makeNBClientOrFail(t)

	// Create a list stations request
	request := &nb.StationLinkListRequest{
		Subscribe: false,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	stationLinks, stationLinksErr := client.ListStationLinks(ctx, request)
	assert.NoError(t, stationLinksErr)
	assert.NotNil(t, stationLinks)

	for {
		stationInfo, err := stationLinks.Recv()
		if err == io.EOF {
			break
		}
		assert.NoError(t, err)
		ids[stationInfo.GetEcgi().Ecid] = stationInfo
	}

	return ids
}

// TestNBStationLinksAPI tests the NB stations API
func (s *TestSuite) TestNBStationLinksAPI(t *testing.T) {

	const (
		expectedStationLinkCount = 9
		expectedPLMNID           = defaultPlmnid
	)

	// Wait for simulator to respond
	waitForSimulatorOrFail(t)

	// Save the station infos into a map indexed by the station's ECID
	links := readStationLinks(t)

	// Make sure the data returned are correct
	re := regexp.MustCompile(cellNamePattern)
	assert.Equal(t, expectedStationLinkCount, len(links))

	for id, link := range links {
		switch id {
		case cell0Name:
		case cell1Name:
		case cell2Name:
		case cell3Name:
		case cell4Name:
		case cell5Name:
		case cell6Name:
		case cell7Name:
		case cell8Name:
			assert.Equal(t, id, link.GetEcgi().Ecid)
			assert.Equal(t, expectedPLMNID, link.Ecgi.Plmnid)
			for _, neighbor := range link.NeighborECGI {
				assert.Equal(t, expectedPLMNID, neighbor.Plmnid)
				assert.True(t, re.MatchString(neighbor.Ecid))
			}
		default:
			assert.Failf(t, "Unexpected station ID %s", id)
		}
	}
}
