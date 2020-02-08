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

package hoapphandover

import (
	"github.com/onosproject/onos-ran/api/nb"
)

// HODecisionMaker decide whether the UE in UELinkInfo should do handover or not
func HODecisionMaker(ueinfo *nb.UELinkInfo) nb.HandOverRequest {
	servStationID := ueinfo.GetEcgi()

	numNeighborCells := len(ueinfo.GetChannelQualities())
	bestStationID := ueinfo.GetChannelQualities()[0].GetTargetEcgi()
	bestCQI := ueinfo.GetChannelQualities()[0].GetCqiHist()

	// start loop
	for i := 1; i < numNeighborCells; i++ {
		tmpCQI := ueinfo.GetChannelQualities()[i].GetCqiHist()
		if bestCQI < tmpCQI {
			bestStationID = ueinfo.GetChannelQualities()[i].GetTargetEcgi()
			bestCQI = tmpCQI
		}
	}

	// no need to trigger handover, since best stations is the serving station
	if servStationID.GetEcid() == bestStationID.GetEcid() && servStationID.GetPlmnid() == bestStationID.GetPlmnid() {
		return nb.HandOverRequest{}
	}

	return nb.HandOverRequest{
		Crnti: ueinfo.GetCrnti(),
		SrcStation: &nb.ECGI{
			Plmnid: servStationID.GetPlmnid(),
			Ecid:   servStationID.GetEcid(),
		},
		DstStation: &nb.ECGI{
			Plmnid: bestStationID.GetPlmnid(),
			Ecid:   bestStationID.GetEcid(),
		},
	}

}
