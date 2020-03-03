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

	"github.com/onosproject/onos-lib-go/pkg/logging"

	"github.com/onosproject/onos-ric/api/nb"
)

var log = logging.GetLogger("ho", "handover")

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
func HODecisionMaker(ueinfo []*nb.UELinkInfo) []*nb.HandOverRequest {

	var resultHoReqs []*nb.HandOverRequest

	for _, l := range ueinfo {
		servStationID := l.GetEcgi()
		numNeighborCells := len(l.GetChannelQualities())
		bestStationID := l.GetChannelQualities()[0].GetTargetEcgi()
		bestCQI := l.GetChannelQualities()[0].GetCqiHist()

		for i := 1; i < numNeighborCells; i++ {
			tmpCQI := l.GetChannelQualities()[i].GetCqiHist()
			if bestCQI < tmpCQI {
				bestStationID = l.GetChannelQualities()[i].GetTargetEcgi()
				bestCQI = tmpCQI
			}
		}

		if servStationID.GetEcid() == bestStationID.GetEcid() && servStationID.GetPlmnid() == bestStationID.GetPlmnid() {
			log.Infof("No need to trigger HO - UE: %s (p:%s,e:%s)", l.GetCrnti(), l.GetEcgi().GetPlmnid(), l.GetEcgi().GetEcid())
			continue
		}

		hoReq := &nb.HandOverRequest{
			Crnti: l.GetCrnti(),
			SrcStation: &nb.ECGI{
				Plmnid: servStationID.GetPlmnid(),
				Ecid:   servStationID.GetEcid(),
			},
			DstStation: &nb.ECGI{
				Plmnid: bestStationID.GetPlmnid(),
				Ecid:   bestStationID.GetEcid(),
			},
		}
		resultHoReqs = append(resultHoReqs, hoReq)
	}
	return resultHoReqs
}
