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
	"fmt"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-ric/api/nb"
	"sync"
	"time"
)

var log = logging.GetLogger("ho", "southbound")

// A3EventMap has all a3Event structure for each UE
var A3EventMap map[string]*a3Event

// A3EventMapMutex is a mutex to lock A3EventMap
var A3EventMapMutex sync.RWMutex

// a3Event is a structure including variables and its channels related with A3 handover event
type a3Event struct {
	chanUELinkInfo    chan *nb.UELinkInfo
	lastUELinkInfoMsg *nb.UELinkInfo
	targetCell        *nb.ECGI
	startTime         time.Time
	hystCQI           int
	a3OffsetCQI       int
	timeToTrigger     int
}

// HODecisionMaker decide whether the UE in UELinkInfo should do handover or not
// if HODecisionMakerWithHOParams matures, this function will be removed.
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

// HODecisionMakerWithHOParams makes a handover decision according to CQI and basic handover parameters
func HODecisionMakerWithHOParams(ueInfo *nb.UELinkInfo, hoReqChan chan *nb.HandOverRequest, hystCQI int, a3OffsetCQI int, TTTMs int) {
	targetCellID, _, targetCQIDelta := getTargetCellInfo(ueInfo, true, hystCQI, a3OffsetCQI)
	tStart := time.Now()

	// HO is unnecessary
	if targetCellID.String() == ueInfo.Ecgi.String() {
		// HO is unnecessary since sCell is the best cell
		return
	} else if targetCQIDelta <= 0 {
		// HO is unnecessary since CQIDelta value is less than or equal to 0
		return
	}

	// According to UELinkInfo message, HO is necessary
	var a3 *a3Event
	A3EventMapMutex.Lock()
	if _, ok := A3EventMap[getUEID(ueInfo.Crnti, ueInfo.Ecgi)]; !ok {
		// Start A3Event
		a3 = &a3Event{
			chanUELinkInfo:    make(chan *nb.UELinkInfo),
			lastUELinkInfoMsg: ueInfo,
			startTime:         tStart,
			targetCell:        targetCellID,
			hystCQI:           hystCQI,
			a3OffsetCQI:       a3OffsetCQI,
			timeToTrigger:     TTTMs,
		}
		A3EventMap[getUEID(ueInfo.Crnti, ueInfo.Ecgi)] = a3
		go startA3Event(a3, hoReqChan)
	} else {
		// A3Event was started and new UELinkInfo message arrives for the UE
		a3 = A3EventMap[getUEID(ueInfo.Crnti, ueInfo.Ecgi)]
		a3.chanUELinkInfo <- ueInfo
	}
	A3EventMapMutex.Unlock()
}

// startA3Event starts A3 handover event; normally it will be initiated as a goroutine
func startA3Event(event *a3Event, hoReqChan chan *nb.HandOverRequest) {
	defer delA3Event(event)
	for {
		remainingTime := (time.Duration(event.timeToTrigger) * time.Millisecond).Nanoseconds() - time.Since(event.startTime).Nanoseconds()

		// if remaining time is equal to or less than 0, trigger HO
		if remainingTime <= 0 {
			sendHOReq(event.lastUELinkInfoMsg.Crnti, event.lastUELinkInfoMsg.Ecgi, event.targetCell, hoReqChan)
			return
		}

		select {
		case ueInfo, ok := <-event.chanUELinkInfo:
			if !ok {
				log.Error("UELinkInfo channel is broken in A3Event due to an unexpected error")
				return
			}

			targetCellID, _, targetCQIDelta := getTargetCellInfo(ueInfo, false, event.hystCQI, event.a3OffsetCQI)

			// Discard HO event: HO is unnecessary since sCell becomes the best cell or CQIDelta becomes negative
			if targetCellID.String() == event.lastUELinkInfoMsg.Ecgi.String() {
				// In the A3Event, HO is unnecessary because sCell becomes the best cell
				return
			} else if targetCQIDelta <= 0 {
				// In the A3Event, HO is unnecessary because CQIDelta value is less than or equal to 0
				return
			}

			// targetCell has been changed: re-start A3Event with the new target cell
			if targetCellID.String() != event.targetCell.String() {
				// In the A3Event, target cell has been changed - restart A3 Event and reset timer
				newA3 := &a3Event{
					chanUELinkInfo:    make(chan *nb.UELinkInfo),
					lastUELinkInfoMsg: ueInfo,
					startTime:         time.Now(),
					targetCell:        targetCellID,
					hystCQI:           event.hystCQI,
					a3OffsetCQI:       event.a3OffsetCQI,
					timeToTrigger:     event.timeToTrigger,
				}
				A3EventMapMutex.Lock()
				A3EventMap[getUEID(ueInfo.Crnti, ueInfo.Ecgi)] = newA3
				A3EventMapMutex.Unlock()
				go startA3Event(newA3, hoReqChan)
				return
			}

			// still HO is required and target cell is the same as before
			if targetCellID.String() == event.targetCell.String() {
				// In the A3Event, new UELinkInfo message arrives, but target cell is the same as before
				continue
			}

		case <-time.After(time.Nanosecond * time.Duration(remainingTime)):
			sendHOReq(event.lastUELinkInfoMsg.Crnti, event.lastUELinkInfoMsg.Ecgi, event.targetCell, hoReqChan)
			return
		}
	}
}

// sendHOReq is a function to send HO request message to C1 interface through a channel
func sendHOReq(crnti string, sCellID *nb.ECGI, nCellID *nb.ECGI, hoReqChan chan *nb.HandOverRequest) {
	hoReq := &nb.HandOverRequest{
		Crnti:      crnti,
		SrcStation: sCellID,
		DstStation: nCellID,
	}

	hoReqChan <- hoReq
}

// getTargetCellInfo returns the target cell information of the handover such as the target cell's ID, its CQI, and its CQI delta
func getTargetCellInfo(ueInfo *nb.UELinkInfo, enterEvent bool, hystCQI int, a3OffsetCQI int) (*nb.ECGI, uint32, int) {
	sCellCQI := getSCellCQI(ueInfo)
	bestECGI := ueInfo.Ecgi
	bestCQIDelta := 0
	bestCQI := sCellCQI

	for _, chanQuality := range ueInfo.ChannelQualities {
		tmpCQIDelta := calcCQIDelta(sCellCQI, chanQuality.CqiHist, enterEvent, hystCQI, a3OffsetCQI)
		if tmpCQIDelta > bestCQIDelta {
			bestECGI = chanQuality.TargetEcgi
			bestCQIDelta = tmpCQIDelta
		}
	}

	return bestECGI, bestCQI, bestCQIDelta
}

// calcCQIDelta calculates a CQI delta (difference) with CQIs, A3-offset, and hysteresis values
func calcCQIDelta(sCellCQI uint32, nCellCQI uint32, enterEvent bool, hystCQI int, a3OffsetCQI int) int {
	if enterEvent {
		// equation when entering A3 event
		return int(nCellCQI) - (int(sCellCQI) + a3OffsetCQI + hystCQI)
	}
	// equation when leaving A3 event
	return int(nCellCQI) - (int(sCellCQI) + a3OffsetCQI - hystCQI)
}

// getSCellCQI find out serving cell's CQI value from UELinkInfo message
func getSCellCQI(ueInfo *nb.UELinkInfo) uint32 {
	for _, chanQuality := range ueInfo.ChannelQualities {
		if chanQuality.TargetEcgi.String() == ueInfo.Ecgi.String() {
			return chanQuality.CqiHist
		}
	}

	log.Error("UELinkInfo Message is wrong: ChannelQuality in UELinkInfo has no CQI Hist of SCell")
	return 0
}

// getUEID makes a UE's ID to a string type
func getUEID(crnti string, ecgi *nb.ECGI) string {
	return fmt.Sprintf("%s:%s", crnti, ecgi.String())
}

// InitA3EventMap initializes A3EventMap; normally, it is called when the gRPC connection for UELinkInfo is established
func InitA3EventMap() {
	A3EventMapMutex.Lock()
	for _, v := range A3EventMap {
		close(v.chanUELinkInfo)
	}
	A3EventMap = nil
	A3EventMap = make(map[string]*a3Event)
	A3EventMapMutex.Unlock()
}

// delA3Event gets rid of a3Event object from A3EventMap, gracefully
func delA3Event(a3 *a3Event) {
	A3EventMapMutex.Lock()
	close(a3.chanUELinkInfo)
	delete(A3EventMap, getUEID(a3.lastUELinkInfoMsg.Crnti, a3.lastUELinkInfoMsg.Ecgi))
	A3EventMapMutex.Unlock()
}
