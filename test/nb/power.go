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
	"testing"
	"time"
)

func setStationsPower(t *testing.T, offset nb.StationPowerOffset, attempts int) {
	stations := readStations(t)
	client := makeNBClientOrFail(t)

	for _, station := range stations {
		for i := 1; i <= attempts; i++ {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			request := &nb.RadioPowerRequest{
				Ecgi:   &nb.ECGI{Plmnid: defaultPlmnid, Ecid: station.Ecgi.Ecid},
				Offset: offset,
			}
			response, err := client.SetRadioPower(ctx, request)
			cancel()
			assert.NoError(t, err)
			assert.NotNil(t, response)
		}
	}
	// Allow the links to settle with the new settings
	time.Sleep(15 * time.Second)
}

func linksAreAllPresent(baseLinks map[string]*nb.UELinkInfo, currentLinks map[string]*nb.UELinkInfo) bool {
	found := 0
	for key := range baseLinks {
		currentLink := currentLinks[key]
		if currentLink.Ecgi.Ecid == baseLinks[key].Ecgi.Ecid {
			found++
		}
	}
	return found == len(baseLinks)
}

// TestNBPowerAPI tests the NB stations API
func (s *TestSuite) TestNBPowerAPI(t *testing.T) {

	// Wait for simulator to respond
	waitForSimulatorOrFail(t)

	// get the uelinks before any changes
	uelinksBeforePowerDown := readLinks(t)

	//  turn the power down to 0 for all stations
	setStationsPower(t, nb.StationPowerOffset_PA_DB_MINUS6, 2)
	setStationsPower(t, nb.StationPowerOffset_PA_DB_1, 2)

	// get the uelinks after the power down change
	uelinksAfterPowerDown := readLinks(t)

	//  make sure that links are still there
	assert.True(t, linksAreAllPresent(uelinksBeforePowerDown, uelinksAfterPowerDown),
		"After power down, one or more links have higher CQI")

	// turn the power back up for all stations
	setStationsPower(t, nb.StationPowerOffset_PA_DB_3, 5)

	// get the uelinks after the power up change
	uelinksAfterPowerUp := readLinks(t)

	// make sure the links are still there
	assert.True(t, linksAreAllPresent(uelinksAfterPowerDown, uelinksAfterPowerUp),
		"After power up, one or more links have lower CQI")
}
