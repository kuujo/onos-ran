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
	"fmt"
	"github.com/onosproject/onos-ran/api/nb"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
	"time"
)

func readStations(t *testing.T) map[string]*nb.StationInfo {
	ids := make(map[string]*nb.StationInfo)

	// Make a client to connect to the onos-ran northbound API
	client := makeNBClientOrFail(t)

	// Create a list stations request
	request := &nb.StationListRequest{
		Subscribe: false,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	stations, stationsErr := client.ListStations(ctx, request)
	assert.NoError(t, stationsErr)
	assert.NotNil(t, stations)

	for {
		stationInfo, err := stations.Recv()
		if err == io.EOF {
			break
		}
		assert.NoError(t, err)
		ids[stationInfo.GetEcgi().Ecid] = stationInfo
	}

	return ids
}

// TestNBStationsAPI tests the NB stations API
func (s *TestSuite) TestNBStationsAPI(t *testing.T) {

	const (
		expectedStationInfoCount          = 9
		expectedPLMNID                    = "001001"
		expectedMaxNumConnectedUes uint32 = 5
	)

	// Wait for simulator to respond
	assert.NoError(t, waitForSimulator())

	// Save the station infos into a map indexed by the station's ECID
	ids := readStations(t)

	// Make sure the data returned are correct
	assert.Equal(t, expectedStationInfoCount, len(ids))
	for stationIndex := 1; stationIndex <= expectedStationInfoCount; stationIndex++ {
		id := fmt.Sprintf("000000%d", stationIndex)
		station, stationFound := ids[id]
		assert.True(t, stationFound)
		assert.Equal(t, id, station.GetEcgi().Ecid)
		assert.Equal(t, expectedPLMNID, station.Ecgi.Plmnid)
		assert.Equal(t, expectedMaxNumConnectedUes, station.MaxNumConnectedUes)
	}
}
