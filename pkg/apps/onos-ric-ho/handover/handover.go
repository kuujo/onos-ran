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
	"time"

	"github.com/onosproject/onos-ric/api/nb"
)

// HOEvent represents a single HO event
type HOEvent struct {
	TimeStamp   time.Time
	CRNTI       string
	SrcPlmnID   string
	SrcEcid     string
	DstPlmnID   string
	DstEcid     string
	ElapsedTime int64
}

// HODecisionMaker decide whether the UE in UELinkInfo should do handover or not
func HODecisionMaker(ueInfo *nb.UELinkInfo) *nb.HandOverRequest {

	servStationID := ueInfo.GetEcgi()
	numNeighborCells := len(ueInfo.GetChannelQualities())
	bestStationID := ueInfo.GetChannelQualities()[0].GetTargetEcgi()
	bestCQI := ueInfo.GetChannelQualities()[0].GetCqiHist()

	for i := 1; i < numNeighborCells; i++ {
		tmpCQI := ueInfo.GetChannelQualities()[i].GetCqiHist()
		if bestCQI < tmpCQI {
			bestStationID = ueInfo.GetChannelQualities()[i].GetTargetEcgi()
			bestCQI = tmpCQI
		}
	}

	if servStationID.GetEcid() == bestStationID.GetEcid() && servStationID.GetPlmnid() == bestStationID.GetPlmnid() {
		return nil
	}

	hoReq := &nb.HandOverRequest{
		Crnti: ueInfo.GetCrnti(),
		SrcStation: &nb.ECGI{
			Plmnid: servStationID.GetPlmnid(),
			Ecid:   servStationID.GetEcid(),
		},
		DstStation: &nb.ECGI{
			Plmnid: bestStationID.GetPlmnid(),
			Ecid:   bestStationID.GetEcid(),
		},
	}
	return hoReq
}
