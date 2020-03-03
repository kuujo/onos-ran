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

package mlbapploadbalance

import (
	"time"

	"github.com/onosproject/onos-lib-go/pkg/logging"

	"github.com/onosproject/onos-ric/api/nb"
)

var log = logging.GetLogger("mlb")

// StaUeJointLink is the joint list of StationInfo and UELinkInfo.
type StaUeJointLink struct {
	TimeStamp   time.Time
	PlmnID      string
	Ecid        string
	MaxNumUes   uint32
	NumUes      int32
	Pa          int32
	ElapsedTime int64
}

// MLBDecisionMaker decides stations to adjust transmission power.
func MLBDecisionMaker(stas []*nb.StationInfo, staLinks []nb.StationLinkInfo, ueLinks []*nb.UELinkInfo, threshold *float64) (*[]nb.RadioPowerRequest, []StaUeJointLink) {
	var mlbReqs []nb.RadioPowerRequest
	var mlbEvents []StaUeJointLink

	// 1. Decide whose tx power should be reduced
	// init staUeJointLinkList
	var staUeJointLinkList []StaUeJointLink
	for _, s := range stas {
		tmpStaUeJointLink := &StaUeJointLink{
			TimeStamp:   time.Now(),
			PlmnID:      s.GetEcgi().GetPlmnid(),
			Ecid:        s.GetEcgi().GetEcid(),
			MaxNumUes:   s.GetMaxNumConnectedUes(),
			NumUes:      0,
			Pa:          0,
			ElapsedTime: 0,
		}
		staUeJointLinkList = append(staUeJointLinkList, *tmpStaUeJointLink)
	}

	// fill ueLinks in each staUeJointLinkList
	setUeLinks(&staUeJointLinkList, ueLinks)

	overloadedStas := getOverloadedStationList(&staUeJointLinkList, threshold)
	stasToBeExpanded := getStasToBeExpanded(&staUeJointLinkList, overloadedStas, &staLinks)

	// 2. Decide how much tx power should be reduced? (static, or dynamic according to CQI values?)
	// - For static, just + or - 3 dB
	for _, os := range *overloadedStas {
		tmpMlbReq := &nb.RadioPowerRequest{
			Ecgi: &nb.ECGI{
				Plmnid: os.PlmnID,
				Ecid:   os.Ecid,
			},
			Offset: nb.StationPowerOffset_PA_DB_MINUS3,
		}
		mlbReqs = append(mlbReqs, *tmpMlbReq)
		mlbEvents = append(mlbEvents, os)
	}

	for _, ss := range *stasToBeExpanded {
		tmpMlbReq := &nb.RadioPowerRequest{
			Ecgi: &nb.ECGI{
				Plmnid: ss.PlmnID,
				Ecid:   ss.Ecid,
			},
			Offset: nb.StationPowerOffset_PA_DB_3,
		}
		mlbReqs = append(mlbReqs, *tmpMlbReq)
		mlbEvents = append(mlbEvents, ss)
	}

	var numTotalUes int32
	for _, e := range staUeJointLinkList {
		// for debug -- should be removed
		log.Infof("STA(p:%s,e:%s) - numUEs:%d (threshold:%d * %f)\n", e.PlmnID, e.Ecid, e.NumUes, e.MaxNumUes, *threshold)
		numTotalUes += e.NumUes
	}
	log.Infof("Total num of reported UEs: %d", numTotalUes)

	// To-Do: For dynamic, sort UE's CQI values and pick UEs should be handed over:
	// if max CQI < 10; - 1 dB, otherwise, -3 dB

	// 3. Return values
	return &mlbReqs, mlbEvents
}

// setUeLinks sets UELink info into StaUeJointLink struct.
func setUeLinks(staJointList *[]StaUeJointLink, ueLinks []*nb.UELinkInfo) {
	for _, l := range ueLinks {
		tmpSta := getStaUeJointLink(l.GetEcgi().GetPlmnid(), l.GetEcgi().GetEcid(), staJointList)
		tmpSta.NumUes++
	}
}

// getOverloadedStationList gets the list of overloaded stations.
func getOverloadedStationList(staUeLinkList *[]StaUeJointLink, threshold *float64) *[]StaUeJointLink {
	var resultOverloadedStations []StaUeJointLink

	for i := 0; i < len(*staUeLinkList); i++ {
		if float64((*staUeLinkList)[i].NumUes) > (*threshold)*float64((*staUeLinkList)[i].MaxNumUes) {
			(*staUeLinkList)[i].Pa = -3
			resultOverloadedStations = append(resultOverloadedStations, (*staUeLinkList)[i])
		}
	}

	return &resultOverloadedStations
}

// getStasToBeExpanded gets the list of stations which have the coverage to be Expanded.
func getStasToBeExpanded(staUeLinkList *[]StaUeJointLink, overloadedStas *[]StaUeJointLink, staLinks *[]nb.StationLinkInfo) *[]StaUeJointLink {
	var resultStasToBeExpanded []StaUeJointLink

	for _, os := range *overloadedStas {
		nEcgis := getNeighborStaEcgi(os.PlmnID, os.Ecid, staLinks)
		for _, n := range nEcgis {
			tmpStaUeJointLink := getStaUeJointLink(n.GetPlmnid(), n.GetEcid(), staUeLinkList)
			if (*tmpStaUeJointLink).Pa == 0 {
				(*tmpStaUeJointLink).Pa = 3
				resultStasToBeExpanded = append(resultStasToBeExpanded, *tmpStaUeJointLink)
			}
		}
	}

	return &resultStasToBeExpanded
}

// getNeighborStaEcgi gets neighbor ECGI list for the STA having given plmnid and ecid.
func getNeighborStaEcgi(plmnid string, ecid string, staLinks *[]nb.StationLinkInfo) []*nb.ECGI {
	for _, s := range *staLinks {
		if s.GetEcgi().GetPlmnid() == plmnid && s.GetEcgi().GetEcid() == ecid {
			return s.GetNeighborECGI()
		}
	}
	return nil
}

// getStaUeJointLink gets the StaUeJointLink having given plmnid and ecid.
func getStaUeJointLink(plmnid string, ecid string, staUeLinkList *[]StaUeJointLink) *StaUeJointLink {
	for i := 0; i < len(*staUeLinkList); i++ {
		if (*staUeLinkList)[i].PlmnID == plmnid && (*staUeLinkList)[i].Ecid == ecid {
			return &(*staUeLinkList)[i]
		}
	}
	return nil
}
