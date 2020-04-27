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
	"sync"
	"time"

	"github.com/onosproject/onos-lib-go/pkg/logging"

	"github.com/onosproject/onos-ric/api/nb"
)

var log = logging.GetLogger("mlb")

// StaLinkMapMutex is the mutex to lock and unlock StaLinkMap
var StaLinkMapMutex sync.RWMutex

// UEMapMutex is the mutex to lock and unlock UEMap
var UEMapMutex sync.RWMutex

// RNIBCellMapMutex is the mutex to lock and unlock RNIBCellMap
var RNIBCellMapMutex sync.RWMutex

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
func MLBDecisionMaker(staMap map[string]*nb.StationInfo, staLinkMap map[string]*nb.StationLinkInfo, ueMap map[string]*nb.UEInfo, threshold *float64) (*[]nb.RadioPowerRequest, map[string]*StaUeJointLink) {
	var mlbReqs []nb.RadioPowerRequest
	staUeJointMap := make(map[string]*StaUeJointLink)

	// Decide whose tx power should be reduced
	// init staUEJointLinkList
	t := time.Now()
	for _, v := range staMap {
		tmpStaUeJointLink := &StaUeJointLink{
			TimeStamp:   t,
			PlmnID:      v.GetEcgi().GetPlmnid(),
			Ecid:        v.GetEcgi().GetEcid(),
			MaxNumUes:   v.GetMaxNumConnectedUes(),
			NumUes:      0,
			Pa:          0,
			ElapsedTime: 0,
		}
		staUeJointMap[v.Ecgi.String()] = tmpStaUeJointLink
	}

	countUEs(staUeJointMap, ueMap)

	overloadedStas := getOverloadedStationList(staUeJointMap, threshold)
	stasToBeExpanded := getStasToBeExpanded(staUeJointMap, overloadedStas, staLinkMap)

	for _, os := range *overloadedStas {
		tmpMlbReq := &nb.RadioPowerRequest{
			Ecgi: &nb.ECGI{
				Plmnid: os.PlmnID,
				Ecid:   os.Ecid,
			},
			Offset: nb.StationPowerOffset_PA_DB_MINUS3,
		}
		mlbReqs = append(mlbReqs, *tmpMlbReq)
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
	}

	var numTotalUes int32
	for _, e := range staUeJointMap {
		// for debug -- should be removed
		log.Infof("STA(p:%s,e:%s) - numUEs:%d (threshold:%d * %f)\n", e.PlmnID, e.Ecid, e.NumUes, e.MaxNumUes, *threshold)
		numTotalUes += e.NumUes
	}
	log.Infof("Total num of reported UEs: %d", numTotalUes)

	return &mlbReqs, staUeJointMap
}

func countUEs(staJointMap map[string]*StaUeJointLink, ueMap map[string]*nb.UEInfo) {
	for _, v := range ueMap {
		if _, ok := staJointMap[v.GetEcgi().String()]; ok {
			staJointMap[v.GetEcgi().String()].NumUes++
		} else {
			log.Errorf("UE %s is connected to the unregistered eNB %s (no CellConfig message for %s)", v.GetCrnti(), v.GetEcgi().String(), v.GetEcgi().String())
		}
	}
}

// getOverloadedStationList gets the list of overloaded stations.
func getOverloadedStationList(staUeJointLinks map[string]*StaUeJointLink, threshold *float64) *[]StaUeJointLink {
	var resultOverloadedStations []StaUeJointLink
	for _, v := range staUeJointLinks {
		if float64(v.NumUes) > (*threshold)*float64(v.MaxNumUes) {
			v.Pa = -3
			resultOverloadedStations = append(resultOverloadedStations, *v)
		}
	}
	return &resultOverloadedStations
}

// getStasToBeExpanded gets the list of stations which have the coverage to be Expanded.
func getStasToBeExpanded(staUeJointLinks map[string]*StaUeJointLink, overloadedStas *[]StaUeJointLink, staLinkMap map[string]*nb.StationLinkInfo) *[]StaUeJointLink {
	var resultStasToBeExpanded []StaUeJointLink

	for _, os := range *overloadedStas {
		nEcgis := getNeighborStaEcgi(os.PlmnID, os.Ecid, staLinkMap)
		for _, n := range nEcgis {
			if staUeJointLinks[n.String()].Pa == 0 {
				staUeJointLinks[n.String()].Pa = 3
				resultStasToBeExpanded = append(resultStasToBeExpanded, *staUeJointLinks[n.String()])
			}
		}
	}

	return &resultStasToBeExpanded
}

// getNeighborStaEcgi gets neighbor ECGI list for the STA having given plmnid and ecid.
func getNeighborStaEcgi(plmnid string, ecid string, staLinkMap map[string]*nb.StationLinkInfo) []*nb.ECGI {
	for _, v := range staLinkMap {
		if v.GetEcgi().GetPlmnid() == plmnid && v.GetEcgi().GetEcid() == ecid {
			return v.GetNeighborECGI()
		}
	}

	return nil
}
