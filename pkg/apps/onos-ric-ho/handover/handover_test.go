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
	"sync"
	"testing"
	"time"

	"github.com/onosproject/onos-ric/api/nb"
	"github.com/stretchr/testify/assert"
)

// Define sample UELinkInfo messages
var sampleMsg1 *nb.UELinkInfo // no HO
var sampleMsg2 *nb.UELinkInfo
var targetSampleMsg2 string
var sampleMsg3 *nb.UELinkInfo
var targetSampleMsg3 string
var sampleMsg4 *nb.UELinkInfo
var targetSampleMsg4 string

var testPLMNID = "315010"

func GenSampleUELinkInfoMsgs() {
	var cqV1 []*nb.ChannelQuality
	cqV1S := &nb.ChannelQuality{
		TargetEcgi: &nb.ECGI{
			Plmnid: testPLMNID,
			Ecid:   "0001",
		},
		CqiHist: 15,
	}
	cqV1 = append(cqV1, cqV1S)
	cqV1N1 := &nb.ChannelQuality{
		TargetEcgi: &nb.ECGI{
			Plmnid: testPLMNID,
			Ecid:   "0002",
		},
		CqiHist: 12,
	}
	cqV1 = append(cqV1, cqV1N1)
	cqV1N2 := &nb.ChannelQuality{
		TargetEcgi: &nb.ECGI{
			Plmnid: testPLMNID,
			Ecid:   "0003",
		},
		CqiHist: 11,
	}
	cqV1 = append(cqV1, cqV1N2)
	cqV1N3 := &nb.ChannelQuality{
		TargetEcgi: &nb.ECGI{
			Plmnid: testPLMNID,
			Ecid:   "0004",
		},
		CqiHist: 1,
	}
	cqV1 = append(cqV1, cqV1N3)
	sampleMsg1 = &nb.UELinkInfo{
		Crnti: "1",
		Ecgi: &nb.ECGI{
			Plmnid: testPLMNID,
			Ecid:   "0001",
		},
		ChannelQualities: cqV1,
	}

	targetSampleMsg2 = "0002"
	var cqV2 []*nb.ChannelQuality
	cqV2S := &nb.ChannelQuality{
		TargetEcgi: &nb.ECGI{
			Plmnid: testPLMNID,
			Ecid:   "0001",
		},
		CqiHist: 10,
	}
	cqV2 = append(cqV2, cqV2S)
	cqV2N1 := &nb.ChannelQuality{
		TargetEcgi: &nb.ECGI{
			Plmnid: testPLMNID,
			Ecid:   "0002",
		},
		CqiHist: 12,
	}
	cqV2 = append(cqV2, cqV2N1)
	cqV2N2 := &nb.ChannelQuality{
		TargetEcgi: &nb.ECGI{
			Plmnid: testPLMNID,
			Ecid:   "0003",
		},
		CqiHist: 11,
	}
	cqV2 = append(cqV2, cqV2N2)
	cqV2N3 := &nb.ChannelQuality{
		TargetEcgi: &nb.ECGI{
			Plmnid: testPLMNID,
			Ecid:   "0004",
		},
		CqiHist: 1,
	}
	cqV2 = append(cqV2, cqV2N3)
	sampleMsg2 = &nb.UELinkInfo{
		Crnti: "1",
		Ecgi: &nb.ECGI{
			Plmnid: testPLMNID,
			Ecid:   "0001",
		},
		ChannelQualities: cqV2,
	}

	targetSampleMsg3 = "0003"
	var cqV3 []*nb.ChannelQuality
	cqV3S := &nb.ChannelQuality{
		TargetEcgi: &nb.ECGI{
			Plmnid: testPLMNID,
			Ecid:   "0001",
		},
		CqiHist: 10,
	}
	cqV3 = append(cqV3, cqV3S)
	cqV3N1 := &nb.ChannelQuality{
		TargetEcgi: &nb.ECGI{
			Plmnid: testPLMNID,
			Ecid:   "0002",
		},
		CqiHist: 9,
	}
	cqV3 = append(cqV3, cqV3N1)
	cqV3N2 := &nb.ChannelQuality{
		TargetEcgi: &nb.ECGI{
			Plmnid: testPLMNID,
			Ecid:   "0003",
		},
		CqiHist: 11,
	}
	cqV3 = append(cqV3, cqV3N2)
	cqV3N3 := &nb.ChannelQuality{
		TargetEcgi: &nb.ECGI{
			Plmnid: testPLMNID,
			Ecid:   "0004",
		},
		CqiHist: 1,
	}
	cqV3 = append(cqV3, cqV3N3)
	sampleMsg3 = &nb.UELinkInfo{
		Crnti: "1",
		Ecgi: &nb.ECGI{
			Plmnid: testPLMNID,
			Ecid:   "0001",
		},
		ChannelQualities: cqV3,
	}

	targetSampleMsg4 = "0004"
	var cqV4 []*nb.ChannelQuality
	cqV4S := &nb.ChannelQuality{
		TargetEcgi: &nb.ECGI{
			Plmnid: testPLMNID,
			Ecid:   "0001",
		},
		CqiHist: 10,
	}
	cqV4 = append(cqV4, cqV4S)
	cqV4N1 := &nb.ChannelQuality{
		TargetEcgi: &nb.ECGI{
			Plmnid: testPLMNID,
			Ecid:   "0002",
		},
		CqiHist: 9,
	}
	cqV4 = append(cqV4, cqV4N1)
	cqV4N2 := &nb.ChannelQuality{
		TargetEcgi: &nb.ECGI{
			Plmnid: testPLMNID,
			Ecid:   "0003",
		},
		CqiHist: 11,
	}
	cqV4 = append(cqV4, cqV4N2)
	cqV4N3 := &nb.ChannelQuality{
		TargetEcgi: &nb.ECGI{
			Plmnid: testPLMNID,
			Ecid:   "0004",
		},
		CqiHist: 15,
	}
	cqV4 = append(cqV4, cqV4N3)
	sampleMsg4 = &nb.UELinkInfo{
		Crnti: "2",
		Ecgi: &nb.ECGI{
			Plmnid: testPLMNID,
			Ecid:   "0001",
		},
		ChannelQualities: cqV4,
	}
}

// Case 1. Serving station is the best one
func TestHODecisionMakerCase1(t *testing.T) {
	GenSampleUELinkInfoMsgs()
	hoReqV1 := HODecisionMaker(sampleMsg1)
	assert.Nil(t, hoReqV1)
}

// Case 2-1. 1st neighbor station is the best one
func TestHODecisionMakerCase2Dot1(t *testing.T) {
	GenSampleUELinkInfoMsgs()
	hoReqV2 := HODecisionMaker(sampleMsg2)
	assert.Equal(t, (hoReqV2).GetDstStation().GetEcid(), targetSampleMsg2)
}

// Case 2-2. 2nd neighbor station is the best one
func TestHODecisionMakerCase2Dot2(t *testing.T) {
	GenSampleUELinkInfoMsgs()
	hoReqV2 := HODecisionMaker(sampleMsg3)
	assert.Equal(t, hoReqV2.GetDstStation().GetEcid(), targetSampleMsg3)
}

// Case 2-3. 3rd neighbor station is the best one
func TestHODecisionMakerCase2Dot3(t *testing.T) {
	GenSampleUELinkInfoMsgs()
	hoReqV2 := HODecisionMaker(sampleMsg4)
	assert.Equal(t, hoReqV2.GetDstStation().GetEcid(), targetSampleMsg4)
}

/* --- new test code --- */
func TestSendHOReq(t *testing.T) {
	tmpCrnti := "1"
	tmpSCell := &nb.ECGI{
		Plmnid: testPLMNID,
		Ecid:   "0001",
	}
	tmpNCell := &nb.ECGI{
		Plmnid: testPLMNID,
		Ecid:   "0002",
	}
	tmpHOChan := make(chan *nb.HandOverRequest)
	go sendHOReq(tmpCrnti, tmpSCell, tmpNCell, tmpHOChan)
	recvHOReq := <-tmpHOChan
	assert.Equal(t, tmpCrnti, recvHOReq.Crnti)
	assert.Equal(t, tmpSCell, recvHOReq.SrcStation)
	assert.Equal(t, tmpNCell, recvHOReq.DstStation)
}

func TestGetTargetCellInfo_Hyst0_A3Offset0_StartEvent(t *testing.T) {
	hystCqi := 0
	a3Offset := 0
	GenSampleUELinkInfoMsgs()
	// for Sample1 message
	s1Ecgi, s1BestCqi, s1BestCqiDelta := getTargetCellInfo(sampleMsg1, true, hystCqi, a3Offset)
	assert.Equal(t, sampleMsg1.Ecgi, s1Ecgi)
	assert.Equal(t, sampleMsg1.ChannelQualities[0].CqiHist, s1BestCqi)
	assert.Equal(t, 0, s1BestCqiDelta)

	// for Sample2 message
	s2Ecgi, s2BestCqi, s2BestCqiDelta := getTargetCellInfo(sampleMsg2, true, hystCqi, a3Offset)
	assert.Equal(t, sampleMsg2.ChannelQualities[1].TargetEcgi, s2Ecgi)
	assert.Equal(t, sampleMsg2.ChannelQualities[1].CqiHist, s2BestCqi)
	assert.Equal(t, int(sampleMsg2.ChannelQualities[1].CqiHist)-int(getSCellCQI(sampleMsg2))-a3Offset-hystCqi, s2BestCqiDelta)

	// for Sample3 message
	s3Ecgi, s3BestCqi, s3BestCqiDelta := getTargetCellInfo(sampleMsg3, true, hystCqi, a3Offset)
	assert.Equal(t, sampleMsg3.ChannelQualities[2].TargetEcgi, s3Ecgi)
	assert.Equal(t, sampleMsg3.ChannelQualities[2].CqiHist, s3BestCqi)
	assert.Equal(t, int(sampleMsg3.ChannelQualities[2].CqiHist)-int(getSCellCQI(sampleMsg3))-a3Offset-hystCqi, s3BestCqiDelta)

	// for Sample4 message
	s4Ecgi, s4BestCqi, s4BestCqiDelta := getTargetCellInfo(sampleMsg4, true, hystCqi, a3Offset)
	assert.Equal(t, sampleMsg4.ChannelQualities[3].TargetEcgi, s4Ecgi)
	assert.Equal(t, sampleMsg4.ChannelQualities[3].CqiHist, s4BestCqi)
	assert.Equal(t, int(sampleMsg4.ChannelQualities[3].CqiHist)-int(getSCellCQI(sampleMsg4))-a3Offset-hystCqi, s4BestCqiDelta)
}

func TestGetTargetCellInfo_Hyst0_A3Offset0_EndEvent(t *testing.T) {
	hystCqi := 0
	a3Offset := 0
	GenSampleUELinkInfoMsgs()
	// for Sample1 message
	s1Ecgi, s1BestCqi, s1BestCqiDelta := getTargetCellInfo(sampleMsg1, false, hystCqi, a3Offset)
	assert.Equal(t, sampleMsg1.Ecgi, s1Ecgi)
	assert.Equal(t, sampleMsg1.ChannelQualities[0].CqiHist, s1BestCqi)
	assert.Equal(t, 0, s1BestCqiDelta)

	// for Sample2 message
	s2Ecgi, s2BestCqi, s2BestCqiDelta := getTargetCellInfo(sampleMsg2, false, hystCqi, a3Offset)
	assert.Equal(t, sampleMsg2.ChannelQualities[1].TargetEcgi, s2Ecgi)
	assert.Equal(t, sampleMsg2.ChannelQualities[1].CqiHist, s2BestCqi)
	assert.Equal(t, int(sampleMsg2.ChannelQualities[1].CqiHist)-int(getSCellCQI(sampleMsg2))-a3Offset+hystCqi, s2BestCqiDelta)

	// for Sample3 message
	s3Ecgi, s3BestCqi, s3BestCqiDelta := getTargetCellInfo(sampleMsg3, false, hystCqi, a3Offset)
	assert.Equal(t, sampleMsg3.ChannelQualities[2].TargetEcgi, s3Ecgi)
	assert.Equal(t, sampleMsg3.ChannelQualities[2].CqiHist, s3BestCqi)
	assert.Equal(t, int(sampleMsg3.ChannelQualities[2].CqiHist)-int(getSCellCQI(sampleMsg3))-a3Offset+hystCqi, s3BestCqiDelta)

	// for Sample4 message
	s4Ecgi, s4BestCqi, s4BestCqiDelta := getTargetCellInfo(sampleMsg4, false, hystCqi, a3Offset)
	assert.Equal(t, sampleMsg4.ChannelQualities[3].TargetEcgi, s4Ecgi)
	assert.Equal(t, sampleMsg4.ChannelQualities[3].CqiHist, s4BestCqi)
	assert.Equal(t, int(sampleMsg4.ChannelQualities[3].CqiHist)-int(getSCellCQI(sampleMsg4))-a3Offset+hystCqi, s4BestCqiDelta)
}

func TestGetTargetCellInfo_Hyst0_A3Offset1_StartEvent(t *testing.T) {
	hystCqi := 0
	a3Offset := 1
	GenSampleUELinkInfoMsgs()
	// for Sample1 message
	s1Ecgi, s1BestCqi, s1BestCqiDelta := getTargetCellInfo(sampleMsg1, true, hystCqi, a3Offset)
	assert.Equal(t, sampleMsg1.Ecgi, s1Ecgi)
	assert.Equal(t, sampleMsg1.ChannelQualities[0].CqiHist, s1BestCqi)
	assert.Equal(t, 0, s1BestCqiDelta)

	// for Sample2 message
	s2Ecgi, s2BestCqi, s2BestCqiDelta := getTargetCellInfo(sampleMsg2, true, hystCqi, a3Offset)
	assert.Equal(t, sampleMsg2.ChannelQualities[1].TargetEcgi, s2Ecgi)
	assert.Equal(t, sampleMsg2.ChannelQualities[1].CqiHist, s2BestCqi)
	assert.Equal(t, int(sampleMsg2.ChannelQualities[1].CqiHist)-int(getSCellCQI(sampleMsg2))-a3Offset-hystCqi, s2BestCqiDelta)

	// for Sample3 message
	s3Ecgi, s3BestCqi, s3BestCqiDelta := getTargetCellInfo(sampleMsg3, true, hystCqi, a3Offset)
	assert.Equal(t, sampleMsg3.ChannelQualities[0].TargetEcgi, s3Ecgi)
	assert.Equal(t, sampleMsg3.ChannelQualities[0].CqiHist, s3BestCqi)
	assert.Equal(t, 0, s3BestCqiDelta)

	// for Sample4 message
	s4Ecgi, s4BestCqi, s4BestCqiDelta := getTargetCellInfo(sampleMsg4, true, hystCqi, a3Offset)
	assert.Equal(t, sampleMsg4.ChannelQualities[3].TargetEcgi, s4Ecgi)
	assert.Equal(t, sampleMsg4.ChannelQualities[3].CqiHist, s4BestCqi)
	assert.Equal(t, int(sampleMsg4.ChannelQualities[3].CqiHist)-int(getSCellCQI(sampleMsg4))-a3Offset-hystCqi, s4BestCqiDelta)
}

func TestGetTargetCellInfo_Hyst0_A3Offset1_EndEvent(t *testing.T) {
	hystCqi := 0
	a3Offset := 1
	GenSampleUELinkInfoMsgs()
	// for Sample1 message
	s1Ecgi, s1BestCqi, s1BestCqiDelta := getTargetCellInfo(sampleMsg1, false, hystCqi, a3Offset)
	assert.Equal(t, sampleMsg1.Ecgi, s1Ecgi)
	assert.Equal(t, sampleMsg1.ChannelQualities[0].CqiHist, s1BestCqi)
	assert.Equal(t, 0, s1BestCqiDelta)

	// for Sample2 message
	s2Ecgi, s2BestCqi, s2BestCqiDelta := getTargetCellInfo(sampleMsg2, false, hystCqi, a3Offset)
	assert.Equal(t, sampleMsg2.ChannelQualities[1].TargetEcgi, s2Ecgi)
	assert.Equal(t, sampleMsg2.ChannelQualities[1].CqiHist, s2BestCqi)
	assert.Equal(t, int(sampleMsg2.ChannelQualities[1].CqiHist)-int(getSCellCQI(sampleMsg2))-a3Offset-hystCqi, s2BestCqiDelta)

	// for Sample3 message
	s3Ecgi, s3BestCqi, s3BestCqiDelta := getTargetCellInfo(sampleMsg3, false, hystCqi, a3Offset)
	assert.Equal(t, sampleMsg3.ChannelQualities[0].TargetEcgi, s3Ecgi)
	assert.Equal(t, sampleMsg3.ChannelQualities[0].CqiHist, s3BestCqi)
	assert.Equal(t, 0, s3BestCqiDelta)

	// for Sample4 message
	s4Ecgi, s4BestCqi, s4BestCqiDelta := getTargetCellInfo(sampleMsg4, false, hystCqi, a3Offset)
	assert.Equal(t, sampleMsg4.ChannelQualities[3].TargetEcgi, s4Ecgi)
	assert.Equal(t, sampleMsg4.ChannelQualities[3].CqiHist, s4BestCqi)
	assert.Equal(t, int(sampleMsg4.ChannelQualities[3].CqiHist)-int(getSCellCQI(sampleMsg4))-a3Offset-hystCqi, s4BestCqiDelta)
}

func TestGetTargetCellInfo_Hyst1_A3Offset0_StartEvent(t *testing.T) {
	hystCqi := 1
	a3Offset := 0
	GenSampleUELinkInfoMsgs()
	// for Sample1 message
	s1Ecgi, s1BestCqi, s1BestCqiDelta := getTargetCellInfo(sampleMsg1, true, hystCqi, a3Offset)
	assert.Equal(t, sampleMsg1.Ecgi, s1Ecgi)
	assert.Equal(t, sampleMsg1.ChannelQualities[0].CqiHist, s1BestCqi)
	assert.Equal(t, 0, s1BestCqiDelta)

	// for Sample2 message
	s2Ecgi, s2BestCqi, s2BestCqiDelta := getTargetCellInfo(sampleMsg2, true, hystCqi, a3Offset)
	assert.Equal(t, sampleMsg2.ChannelQualities[1].TargetEcgi, s2Ecgi)
	assert.Equal(t, sampleMsg2.ChannelQualities[1].CqiHist, s2BestCqi)
	assert.Equal(t, int(sampleMsg2.ChannelQualities[1].CqiHist)-int(getSCellCQI(sampleMsg2))-a3Offset-hystCqi, s2BestCqiDelta)

	// for Sample3 message
	s3Ecgi, s3BestCqi, s3BestCqiDelta := getTargetCellInfo(sampleMsg3, true, hystCqi, a3Offset)
	assert.Equal(t, sampleMsg3.ChannelQualities[0].TargetEcgi, s3Ecgi)
	assert.Equal(t, sampleMsg3.ChannelQualities[0].CqiHist, s3BestCqi)
	assert.Equal(t, 0, s3BestCqiDelta)

	// for Sample4 message
	s4Ecgi, s4BestCqi, s4BestCqiDelta := getTargetCellInfo(sampleMsg4, true, hystCqi, a3Offset)
	assert.Equal(t, sampleMsg4.ChannelQualities[3].TargetEcgi, s4Ecgi)
	assert.Equal(t, sampleMsg4.ChannelQualities[3].CqiHist, s4BestCqi)
	assert.Equal(t, int(sampleMsg4.ChannelQualities[3].CqiHist)-int(getSCellCQI(sampleMsg4))-a3Offset-hystCqi, s4BestCqiDelta)
}

func TestGetTargetCellInfo_Hyst1_A3Offset0_EndEvent(t *testing.T) {
	hystCqi := 1
	a3Offset := 0
	GenSampleUELinkInfoMsgs()
	// for Sample1 message
	s1Ecgi, s1BestCqi, s1BestCqiDelta := getTargetCellInfo(sampleMsg1, false, hystCqi, a3Offset)
	assert.Equal(t, sampleMsg1.Ecgi, s1Ecgi)
	assert.Equal(t, sampleMsg1.ChannelQualities[0].CqiHist, s1BestCqi)
	assert.Equal(t, 1, s1BestCqiDelta)

	// for Sample2 message
	s2Ecgi, s2BestCqi, s2BestCqiDelta := getTargetCellInfo(sampleMsg2, false, hystCqi, a3Offset)
	assert.Equal(t, sampleMsg2.ChannelQualities[1].TargetEcgi, s2Ecgi)
	assert.Equal(t, sampleMsg2.ChannelQualities[1].CqiHist, s2BestCqi)
	assert.Equal(t, int(sampleMsg2.ChannelQualities[1].CqiHist)-int(getSCellCQI(sampleMsg2))-a3Offset+hystCqi, s2BestCqiDelta)

	// for Sample3 message
	s3Ecgi, s3BestCqi, s3BestCqiDelta := getTargetCellInfo(sampleMsg3, false, hystCqi, a3Offset)
	assert.Equal(t, sampleMsg3.ChannelQualities[2].TargetEcgi, s3Ecgi)
	assert.Equal(t, sampleMsg3.ChannelQualities[2].CqiHist, s3BestCqi)
	assert.Equal(t, int(sampleMsg3.ChannelQualities[2].CqiHist)-int(getSCellCQI(sampleMsg3))-a3Offset+hystCqi, s3BestCqiDelta)

	// for Sample4 message
	s4Ecgi, s4BestCqi, s4BestCqiDelta := getTargetCellInfo(sampleMsg4, false, hystCqi, a3Offset)
	assert.Equal(t, sampleMsg4.ChannelQualities[3].TargetEcgi, s4Ecgi)
	assert.Equal(t, sampleMsg4.ChannelQualities[3].CqiHist, s4BestCqi)
	assert.Equal(t, int(sampleMsg4.ChannelQualities[3].CqiHist)-int(getSCellCQI(sampleMsg4))-a3Offset+hystCqi, s4BestCqiDelta)
}

func TestCalcCQIDelta(t *testing.T) {
	sCellCQI := uint32(7)
	nCellCQI := uint32(13)
	// case 1: hyst 0; a3-offset 0
	c1Start := calcCQIDelta(sCellCQI, nCellCQI, true, 0, 0)
	c1End := calcCQIDelta(sCellCQI, nCellCQI, false, 0, 0)
	assert.Equal(t, c1Start, 6)
	assert.Equal(t, c1End, 6)
	// case 2: hyst 1; a3-offset 0
	c2Start := calcCQIDelta(sCellCQI, nCellCQI, true, 1, 0)
	c2End := calcCQIDelta(sCellCQI, nCellCQI, false, 1, 0)
	assert.Equal(t, c2Start, 5)
	assert.Equal(t, c2End, 7)
	// case 3: hyst 3; a3-offset 0
	c3Start := calcCQIDelta(sCellCQI, nCellCQI, true, 3, 0)
	c3End := calcCQIDelta(sCellCQI, nCellCQI, false, 3, 0)
	assert.Equal(t, c3Start, 3)
	assert.Equal(t, c3End, 9)
	// case 4: hyst 0; a3-offset 1
	c4Start := calcCQIDelta(sCellCQI, nCellCQI, true, 0, 1)
	c4End := calcCQIDelta(sCellCQI, nCellCQI, false, 0, 1)
	assert.Equal(t, c4Start, 5)
	assert.Equal(t, c4End, 5)
	// case 5: hyst 0; a3-offset 3
	c5Start := calcCQIDelta(sCellCQI, nCellCQI, true, 0, 3)
	c5End := calcCQIDelta(sCellCQI, nCellCQI, false, 0, 3)
	assert.Equal(t, c5Start, 3)
	assert.Equal(t, c5End, 3)
	// case 6: hyst 1; a3-offset 1
	c6Start := calcCQIDelta(sCellCQI, nCellCQI, true, 1, 1)
	c6End := calcCQIDelta(sCellCQI, nCellCQI, false, 1, 1)
	assert.Equal(t, c6Start, 4)
	assert.Equal(t, c6End, 6)
	// case 7: hyst 3; a3-offset 3
	c7Start := calcCQIDelta(sCellCQI, nCellCQI, true, 3, 3)
	c7End := calcCQIDelta(sCellCQI, nCellCQI, false, 3, 3)
	assert.Equal(t, c7Start, 0)
	assert.Equal(t, c7End, 6)
}

func TestGetScellCQI(t *testing.T) {
	GenSampleUELinkInfoMsgs()
	// for Sample1Msg
	cqi1 := getSCellCQI(sampleMsg1)
	assert.Equal(t, cqi1, uint32(15))
	// for Sample2Msg
	cqi2 := getSCellCQI(sampleMsg2)
	assert.Equal(t, cqi2, uint32(10))
	// for Sample3Msg
	cqi3 := getSCellCQI(sampleMsg3)
	assert.Equal(t, cqi3, uint32(10))
	// for Sample4Msg
	cqi4 := getSCellCQI(sampleMsg4)
	assert.Equal(t, cqi4, uint32(10))
}

func TestGetScellCQI_WrongInfo(t *testing.T) {
	GenSampleUELinkInfoMsgs()
	sampleMsg1.ChannelQualities = nil
	cqi := getSCellCQI(sampleMsg1)
	assert.Equal(t, uint32(0), cqi)
}

func TestGetUEID(t *testing.T) {
	GenSampleUELinkInfoMsgs()
	// for Sample1Msg
	ue1 := getUEID(sampleMsg1.Crnti, sampleMsg1.Ecgi)
	assert.Equal(t, ue1, fmt.Sprintf("%s:%s", sampleMsg1.Crnti, sampleMsg1.Ecgi.String()))
	// for Sample2Msg
	ue2 := getUEID(sampleMsg2.Crnti, sampleMsg2.Ecgi)
	assert.Equal(t, ue2, fmt.Sprintf("%s:%s", sampleMsg2.Crnti, sampleMsg2.Ecgi.String()))
	// for Sample3Msg
	ue3 := getUEID(sampleMsg3.Crnti, sampleMsg3.Ecgi)
	assert.Equal(t, ue3, fmt.Sprintf("%s:%s", sampleMsg3.Crnti, sampleMsg3.Ecgi.String()))
	// for Sample4Msg
	ue4 := getUEID(sampleMsg4.Crnti, sampleMsg4.Ecgi)
	assert.Equal(t, ue4, fmt.Sprintf("%s:%s", sampleMsg4.Crnti, sampleMsg4.Ecgi.String()))
}

func TestInitA3EventMap(t *testing.T) {
	GenSampleUELinkInfoMsgs()
	InitA3EventMap()
	assert.NotNil(t, A3EventMap)
	assert.Equal(t, 0, len(A3EventMap))
	a3 := &a3Event{
		chanUELinkInfo:    make(chan *nb.UELinkInfo),
		lastUELinkInfoMsg: sampleMsg1,
		startTime:         time.Now(),
		targetCell:        sampleMsg1.Ecgi,
		hystCQI:           0,
		a3OffsetCQI:       0,
		timeToTrigger:     0,
	}
	A3EventMapMutex.Lock()
	A3EventMap[getUEID(sampleMsg1.Crnti, sampleMsg1.Ecgi)] = a3
	A3EventMapMutex.Unlock()
	assert.NotNil(t, A3EventMap)
	assert.Equal(t, 1, len(A3EventMap))
	InitA3EventMap()
	assert.NotNil(t, A3EventMap)
	assert.Equal(t, 0, len(A3EventMap))
}

func TestDelA3Event(t *testing.T) {
	InitA3EventMap()
	GenSampleUELinkInfoMsgs()
	// make one A3Event
	targetCellIDSample1, _, _ := getTargetCellInfo(sampleMsg1, true, 0, 0)
	a3EventSample1 := &a3Event{
		chanUELinkInfo:    make(chan *nb.UELinkInfo),
		lastUELinkInfoMsg: sampleMsg1,
		startTime:         time.Now(),
		targetCell:        targetCellIDSample1,
		hystCQI:           0,
		a3OffsetCQI:       0,
		timeToTrigger:     0,
	}
	A3EventMapMutex.Lock()
	A3EventMap[getUEID(sampleMsg1.Crnti, sampleMsg1.Ecgi)] = a3EventSample1
	A3EventMapMutex.Unlock()

	assert.Equal(t, len(A3EventMap), 1)
	delA3Event(a3EventSample1)
	assert.Equal(t, len(A3EventMap), 0)
}

func TestHODecisionMakerWithHOParams1_H0A0T0(t *testing.T) {
	var mu sync.RWMutex
	hyst := 0
	a3offset := 0
	ttt := 0
	tmpHOChan := make(chan *nb.HandOverRequest, 1)
	numHOs := 0
	expectedNumHOs := 0

	go func() {
		for range tmpHOChan {
			mu.Lock()
			numHOs++
			mu.Unlock()
		}
	}()

	InitA3EventMap()
	GenSampleUELinkInfoMsgs()
	mu.Lock()
	HODecisionMakerWithHOParams(sampleMsg1, tmpHOChan, hyst, a3offset, ttt)
	mu.Unlock()
	time.Sleep(time.Duration(ttt)*time.Millisecond*10 + 1*time.Second)
	mu.Lock()
	assert.Equal(t, expectedNumHOs, numHOs)
	mu.Unlock()
}

func TestHODecisionMakerWithHOParams1_H1A0T0(t *testing.T) {
	var mu sync.RWMutex
	hyst := 1
	a3offset := 0
	ttt := 0
	tmpHOChan := make(chan *nb.HandOverRequest)
	numHOs := 0
	expectedNumHOs := 0

	go func() {
		for range tmpHOChan {
			mu.Lock()
			numHOs++
			mu.Unlock()
		}
	}()

	InitA3EventMap()
	GenSampleUELinkInfoMsgs()
	mu.Lock()
	HODecisionMakerWithHOParams(sampleMsg1, tmpHOChan, hyst, a3offset, ttt)
	mu.Unlock()
	time.Sleep(time.Duration(ttt)*time.Millisecond*10 + 1*time.Second)
	mu.Lock()
	assert.Equal(t, expectedNumHOs, numHOs)
	mu.Unlock()
}
func TestHODecisionMakerWithHOParams1_H3A0T0(t *testing.T) {
	var mu sync.RWMutex
	hyst := 3
	a3offset := 0
	ttt := 0
	tmpHOChan := make(chan *nb.HandOverRequest)
	numHOs := 0
	expectedNumHOs := 0

	go func() {
		for range tmpHOChan {
			mu.Lock()
			numHOs++
			mu.Unlock()
		}
	}()

	InitA3EventMap()
	GenSampleUELinkInfoMsgs()
	mu.Lock()
	HODecisionMakerWithHOParams(sampleMsg1, tmpHOChan, hyst, a3offset, ttt)
	time.Sleep(time.Duration(ttt)*time.Millisecond*10 + 1*time.Second)
	mu.Unlock()
	time.Sleep(time.Duration(ttt)*time.Millisecond*10 + 1*time.Second)
	mu.Lock()
	assert.Equal(t, expectedNumHOs, numHOs)
	mu.Unlock()
}

func TestHODecisionMakerWithHOParams1_H0A1T0(t *testing.T) {
	var mu sync.RWMutex
	hyst := 0
	a3offset := 1
	ttt := 0
	tmpHOChan := make(chan *nb.HandOverRequest)
	numHOs := 0
	expectedNumHOs := 0

	go func() {
		for range tmpHOChan {
			mu.Lock()
			numHOs++
			mu.Unlock()
		}
	}()

	InitA3EventMap()
	GenSampleUELinkInfoMsgs()
	mu.Lock()
	HODecisionMakerWithHOParams(sampleMsg1, tmpHOChan, hyst, a3offset, ttt)
	mu.Unlock()
	time.Sleep(time.Duration(ttt)*time.Millisecond*10 + 1*time.Second)
	mu.Lock()
	assert.Equal(t, expectedNumHOs, numHOs)
	mu.Unlock()
}

func TestHODecisionMakerWithHOParams1_H0A3T0(t *testing.T) {
	var mu sync.RWMutex
	hyst := 0
	a3offset := 3
	ttt := 0
	tmpHOChan := make(chan *nb.HandOverRequest)
	numHOs := 0
	expectedNumHOs := 0

	go func() {
		for range tmpHOChan {
			mu.Lock()
			numHOs++
			mu.Unlock()
		}
	}()

	InitA3EventMap()
	GenSampleUELinkInfoMsgs()
	mu.Lock()
	HODecisionMakerWithHOParams(sampleMsg1, tmpHOChan, hyst, a3offset, ttt)
	mu.Unlock()
	time.Sleep(time.Duration(ttt)*time.Millisecond*10 + 1*time.Second)
	mu.Lock()
	assert.Equal(t, expectedNumHOs, numHOs)
	mu.Unlock()
}

func TestHODecisionMakerWithHOParams1_H1A1T0(t *testing.T) {
	var mu sync.RWMutex
	hyst := 1
	a3offset := 1
	ttt := 0
	tmpHOChan := make(chan *nb.HandOverRequest)
	numHOs := 0
	expectedNumHOs := 0

	go func() {
		for range tmpHOChan {
			mu.Lock()
			numHOs++
			mu.Unlock()
		}
	}()

	InitA3EventMap()
	GenSampleUELinkInfoMsgs()
	mu.Lock()
	HODecisionMakerWithHOParams(sampleMsg1, tmpHOChan, hyst, a3offset, ttt)
	mu.Unlock()
	time.Sleep(time.Duration(ttt)*time.Millisecond*10 + 1*time.Second)
	mu.Lock()
	assert.Equal(t, expectedNumHOs, numHOs)
	mu.Unlock()
}

func TestHODecisionMakerWithHOParams1_H3A3T0(t *testing.T) {
	var mu sync.RWMutex
	hyst := 3
	a3offset := 3
	ttt := 0
	tmpHOChan := make(chan *nb.HandOverRequest)
	numHOs := 0
	expectedNumHOs := 0

	go func() {
		for range tmpHOChan {
			mu.Lock()
			numHOs++
			mu.Unlock()
		}
	}()

	InitA3EventMap()
	GenSampleUELinkInfoMsgs()
	mu.Lock()
	HODecisionMakerWithHOParams(sampleMsg1, tmpHOChan, hyst, a3offset, ttt)
	mu.Unlock()
	time.Sleep(time.Duration(ttt)*time.Millisecond*10 + 1*time.Second)
	mu.Lock()
	assert.Equal(t, expectedNumHOs, numHOs)
	mu.Unlock()
}

func TestHODecisionMakerWithHOParams2_H0A0T0(t *testing.T) {
	var mu sync.RWMutex
	hyst := 0
	a3offset := 0
	ttt := 0
	tmpHOChan := make(chan *nb.HandOverRequest, 1)
	numHOs := 0
	expectedNumHOs := 1
	var hoEvents *nb.HandOverRequest

	go func() {
		for ho := range tmpHOChan {
			mu.Lock()
			numHOs++
			hoEvents = ho
			mu.Unlock()
		}
	}()

	InitA3EventMap()
	GenSampleUELinkInfoMsgs()
	mu.Lock()
	HODecisionMakerWithHOParams(sampleMsg2, tmpHOChan, hyst, a3offset, ttt)
	mu.Unlock()
	time.Sleep(time.Duration(ttt)*time.Millisecond*10 + 1*time.Second)
	mu.Lock()
	assert.Equal(t, expectedNumHOs, numHOs)
	assert.Equal(t, "0002", hoEvents.DstStation.Ecid)
	mu.Unlock()
}

func TestHODecisionMakerWithHOParams2_H1A0T0(t *testing.T) {
	var mu sync.RWMutex
	hyst := 1
	a3offset := 0
	ttt := 0
	tmpHOChan := make(chan *nb.HandOverRequest)
	numHOs := 0
	expectedNumHOs := 1
	var hoEvents *nb.HandOverRequest

	go func() {

		for ho := range tmpHOChan {
			mu.Lock()
			numHOs++
			hoEvents = ho
			mu.Unlock()
		}
	}()

	InitA3EventMap()
	GenSampleUELinkInfoMsgs()
	mu.Lock()
	HODecisionMakerWithHOParams(sampleMsg2, tmpHOChan, hyst, a3offset, ttt)
	mu.Unlock()
	time.Sleep(time.Duration(ttt)*time.Millisecond*10 + 1*time.Second)
	mu.Lock()
	assert.Equal(t, expectedNumHOs, numHOs)
	assert.Equal(t, "0002", hoEvents.DstStation.Ecid)
	mu.Unlock()
}

func TestHODecisionMakerWithHOParams2_H3A0T0(t *testing.T) {
	var mu sync.RWMutex
	hyst := 3
	a3offset := 0
	ttt := 0
	tmpHOChan := make(chan *nb.HandOverRequest)
	numHOs := 0
	expectedNumHOs := 0

	go func() {
		for range tmpHOChan {
			mu.Lock()
			numHOs++
			mu.Unlock()
		}
	}()

	InitA3EventMap()
	GenSampleUELinkInfoMsgs()
	mu.Lock()
	HODecisionMakerWithHOParams(sampleMsg2, tmpHOChan, hyst, a3offset, ttt)
	mu.Unlock()
	time.Sleep(time.Duration(ttt)*time.Millisecond*10 + 1*time.Second)
	mu.Lock()
	assert.Equal(t, expectedNumHOs, numHOs)
	mu.Unlock()
}

func TestHODecisionMakerWithHOParams2_H0A1T0(t *testing.T) {
	var mu sync.RWMutex
	hyst := 0
	a3offset := 1
	ttt := 0
	tmpHOChan := make(chan *nb.HandOverRequest)
	numHOs := 0
	expectedNumHOs := 1
	var hoEvents *nb.HandOverRequest

	go func() {
		for ho := range tmpHOChan {
			mu.Lock()
			numHOs++
			hoEvents = ho
			mu.Unlock()
		}
	}()

	InitA3EventMap()
	GenSampleUELinkInfoMsgs()
	mu.Lock()
	HODecisionMakerWithHOParams(sampleMsg2, tmpHOChan, hyst, a3offset, ttt)
	mu.Unlock()
	time.Sleep(time.Duration(ttt)*time.Millisecond*10 + 1*time.Second)
	mu.Lock()
	assert.Equal(t, expectedNumHOs, numHOs)
	assert.Equal(t, "0002", hoEvents.DstStation.Ecid)
	mu.Unlock()
}

func TestHODecisionMakerWithHOParams2_H0A3T0(t *testing.T) {
	var mu sync.RWMutex
	hyst := 0
	a3offset := 3
	ttt := 0
	tmpHOChan := make(chan *nb.HandOverRequest)
	numHOs := 0
	expectedNumHOs := 0

	go func() {
		for range tmpHOChan {
			mu.Lock()
			numHOs++
			mu.Unlock()
		}
	}()

	InitA3EventMap()
	GenSampleUELinkInfoMsgs()
	mu.Lock()
	HODecisionMakerWithHOParams(sampleMsg2, tmpHOChan, hyst, a3offset, ttt)
	mu.Unlock()
	time.Sleep(time.Duration(ttt)*time.Millisecond*10 + 1*time.Second)
	mu.Lock()
	assert.Equal(t, expectedNumHOs, numHOs)
	mu.Unlock()
}

func TestHODecisionMakerWithHOParams2_H1A1T0(t *testing.T) {
	var mu sync.RWMutex
	hyst := 1
	a3offset := 1
	ttt := 0
	tmpHOChan := make(chan *nb.HandOverRequest)
	numHOs := 0
	expectedNumHOs := 0

	go func() {
		for range tmpHOChan {
			mu.Lock()
			numHOs++
			mu.Unlock()
		}
	}()

	InitA3EventMap()
	GenSampleUELinkInfoMsgs()
	mu.Lock()
	HODecisionMakerWithHOParams(sampleMsg2, tmpHOChan, hyst, a3offset, ttt)
	mu.Unlock()
	time.Sleep(time.Duration(ttt)*time.Millisecond*10 + 1*time.Second)
	mu.Lock()
	assert.Equal(t, expectedNumHOs, numHOs)
	mu.Unlock()
}

func TestHODecisionMakerWithHOParams2_H3A3T0(t *testing.T) {
	var mu sync.RWMutex
	hyst := 3
	a3offset := 3
	ttt := 0
	tmpHOChan := make(chan *nb.HandOverRequest)
	numHOs := 0
	expectedNumHOs := 0

	go func() {
		for range tmpHOChan {
			mu.Lock()
			numHOs++
			mu.Unlock()
		}
	}()

	InitA3EventMap()
	GenSampleUELinkInfoMsgs()
	mu.Lock()
	HODecisionMakerWithHOParams(sampleMsg2, tmpHOChan, hyst, a3offset, ttt)
	mu.Unlock()
	time.Sleep(time.Duration(ttt)*time.Millisecond*10 + 1*time.Second)
	mu.Lock()
	assert.Equal(t, expectedNumHOs, numHOs)
	mu.Unlock()
}

func TestHODecisionMakerWithHOParams3_H0A0T0(t *testing.T) {
	var mu sync.RWMutex
	hyst := 0
	a3offset := 0
	ttt := 0
	tmpHOChan := make(chan *nb.HandOverRequest, 1)
	numHOs := 0
	expectedNumHOs := 1
	var hoEvents *nb.HandOverRequest

	go func() {
		for ho := range tmpHOChan {
			mu.Lock()
			numHOs++
			hoEvents = ho
			mu.Unlock()
		}
	}()

	InitA3EventMap()
	GenSampleUELinkInfoMsgs()
	mu.Lock()
	HODecisionMakerWithHOParams(sampleMsg3, tmpHOChan, hyst, a3offset, ttt)
	mu.Unlock()
	time.Sleep(time.Duration(ttt)*time.Millisecond*10 + 1*time.Second)
	mu.Lock()
	assert.Equal(t, expectedNumHOs, numHOs)
	assert.Equal(t, "0003", hoEvents.DstStation.Ecid)
	mu.Unlock()
}

func TestHODecisionMakerWithHOParams3_H1A0T0(t *testing.T) {
	var mu sync.RWMutex
	hyst := 1
	a3offset := 0
	ttt := 0
	tmpHOChan := make(chan *nb.HandOverRequest)
	numHOs := 0
	expectedNumHOs := 0

	go func() {
		for range tmpHOChan {
			mu.Lock()
			numHOs++
			mu.Unlock()
		}
	}()

	InitA3EventMap()
	GenSampleUELinkInfoMsgs()
	mu.Lock()
	HODecisionMakerWithHOParams(sampleMsg3, tmpHOChan, hyst, a3offset, ttt)
	mu.Unlock()
	time.Sleep(time.Duration(ttt)*time.Millisecond*10 + 1*time.Second)
	mu.Lock()
	assert.Equal(t, expectedNumHOs, numHOs)
	mu.Unlock()
}

func TestHODecisionMakerWithHOParams3_H3A0T0(t *testing.T) {
	var mu sync.RWMutex
	hyst := 3
	a3offset := 0
	ttt := 0
	tmpHOChan := make(chan *nb.HandOverRequest)
	numHOs := 0
	expectedNumHOs := 0

	go func() {
		for range tmpHOChan {
			mu.Lock()
			numHOs++
			mu.Unlock()
		}
	}()

	InitA3EventMap()
	GenSampleUELinkInfoMsgs()
	mu.Lock()
	HODecisionMakerWithHOParams(sampleMsg3, tmpHOChan, hyst, a3offset, ttt)
	mu.Unlock()
	time.Sleep(time.Duration(ttt)*time.Millisecond*10 + 1*time.Second)
	mu.Lock()
	assert.Equal(t, expectedNumHOs, numHOs)
	mu.Unlock()
}

func TestHODecisionMakerWithHOParams3_H0A1T0(t *testing.T) {
	var mu sync.RWMutex
	hyst := 0
	a3offset := 1
	ttt := 0
	tmpHOChan := make(chan *nb.HandOverRequest)
	numHOs := 0
	expectedNumHOs := 0

	go func() {
		for range tmpHOChan {
			mu.Lock()
			numHOs++
			mu.Unlock()
		}
	}()

	InitA3EventMap()
	GenSampleUELinkInfoMsgs()
	mu.Lock()
	HODecisionMakerWithHOParams(sampleMsg3, tmpHOChan, hyst, a3offset, ttt)
	mu.Unlock()
	time.Sleep(time.Duration(ttt)*time.Millisecond*10 + 1*time.Second)
	mu.Lock()
	assert.Equal(t, expectedNumHOs, numHOs)
	mu.Unlock()
}

func TestHODecisionMakerWithHOParams3_H0A3T0(t *testing.T) {
	var mu sync.RWMutex
	hyst := 0
	a3offset := 3
	ttt := 0
	tmpHOChan := make(chan *nb.HandOverRequest)
	numHOs := 0
	expectedNumHOs := 0

	go func() {
		for range tmpHOChan {
			mu.Lock()
			numHOs++
			mu.Unlock()
		}
	}()

	InitA3EventMap()
	GenSampleUELinkInfoMsgs()
	mu.Lock()
	HODecisionMakerWithHOParams(sampleMsg3, tmpHOChan, hyst, a3offset, ttt)
	mu.Unlock()
	time.Sleep(time.Duration(ttt)*time.Millisecond*10 + 1*time.Second)
	mu.Lock()
	assert.Equal(t, expectedNumHOs, numHOs)
	mu.Unlock()
}

func TestHODecisionMakerWithHOParams3_H1A1T0(t *testing.T) {
	var mu sync.RWMutex
	hyst := 1
	a3offset := 1
	ttt := 0
	tmpHOChan := make(chan *nb.HandOverRequest)
	numHOs := 0
	expectedNumHOs := 0

	go func() {
		for range tmpHOChan {
			mu.Lock()
			numHOs++
			mu.Unlock()
		}
	}()

	InitA3EventMap()
	GenSampleUELinkInfoMsgs()
	mu.Lock()
	HODecisionMakerWithHOParams(sampleMsg3, tmpHOChan, hyst, a3offset, ttt)
	mu.Unlock()
	time.Sleep(time.Duration(ttt)*time.Millisecond*10 + 1*time.Second)
	mu.Lock()
	assert.Equal(t, expectedNumHOs, numHOs)
	mu.Unlock()
}

func TestHODecisionMakerWithHOParams3_H3A3T0(t *testing.T) {
	var mu sync.RWMutex
	hyst := 3
	a3offset := 3
	ttt := 0
	tmpHOChan := make(chan *nb.HandOverRequest)
	numHOs := 0
	expectedNumHOs := 0

	go func() {
		for range tmpHOChan {
			mu.Lock()
			numHOs++
			mu.Unlock()
		}
	}()

	InitA3EventMap()
	GenSampleUELinkInfoMsgs()
	mu.Lock()
	HODecisionMakerWithHOParams(sampleMsg3, tmpHOChan, hyst, a3offset, ttt)
	mu.Unlock()
	time.Sleep(time.Duration(ttt)*time.Millisecond*10 + 1*time.Second)
	mu.Lock()
	assert.Equal(t, expectedNumHOs, numHOs)
	mu.Unlock()
}

func TestHODecisionMakerWithHOParams3_to_3_H0A0T1000(t *testing.T) {
	var mu sync.RWMutex
	hyst := 0
	a3offset := 0
	ttt := 1000
	tmpHOChan := make(chan *nb.HandOverRequest)
	numHOs := 0
	expectedNumHOs := 1

	go func() {
		for ho := range tmpHOChan {
			mu.Lock()
			numHOs++
			mu.Unlock()
			tmpHOChan <- ho

		}
	}()

	InitA3EventMap()
	GenSampleUELinkInfoMsgs()
	mu.Lock()
	HODecisionMakerWithHOParams(sampleMsg3, tmpHOChan, hyst, a3offset, ttt)
	time.Sleep(100 * time.Millisecond)
	HODecisionMakerWithHOParams(sampleMsg2, tmpHOChan, hyst, a3offset, ttt)
	mu.Unlock()
	time.Sleep(time.Duration(ttt)*time.Millisecond*10 + 1*time.Second)
	mu.Lock()
	assert.Equal(t, expectedNumHOs, numHOs)
	assert.Equal(t, "0002", (<-tmpHOChan).DstStation.Ecid)
	mu.Unlock()
}

func TestHODecisionMakerWithHOParams3_to_3_H0A0T1000_Scalable_500Events(t *testing.T) {
	var mu sync.RWMutex
	numEvents := 500
	hyst := 0
	a3offset := 0
	ttt := 1000
	tmpHOChan := make(chan *nb.HandOverRequest)
	numHOs := 0

	go func() {
		for ho := range tmpHOChan {
			mu.Lock()
			numHOs++
			mu.Unlock()
			tmpHOChan <- ho

		}
	}()

	defer func() {
		// Panic handler
		time.Sleep(time.Duration(ttt)*time.Millisecond*10 + 1*time.Second)
		if r := recover(); r != nil {
			t.Errorf("The have a panic - please check detailed log")
		}
	}()

	InitA3EventMap()
	GenSampleUELinkInfoMsgs()
	mu.Lock()
	for i := 0; i < numEvents; i++ {
		HODecisionMakerWithHOParams(sampleMsg1, tmpHOChan, hyst, a3offset, ttt)
		HODecisionMakerWithHOParams(sampleMsg2, tmpHOChan, hyst, a3offset, ttt)
		HODecisionMakerWithHOParams(sampleMsg3, tmpHOChan, hyst, a3offset, ttt)
		HODecisionMakerWithHOParams(sampleMsg4, tmpHOChan, hyst, a3offset, ttt)
	}
	mu.Unlock()
}

func TestStartA3Event_S1_TTT0(t *testing.T) {
	timeToTrigger := 0
	GenSampleUELinkInfoMsgs()
	InitA3EventMap()
	hoChan := make(chan *nb.HandOverRequest)
	tStart := time.Now()
	sampleA3Event := &a3Event{
		chanUELinkInfo:    make(chan *nb.UELinkInfo),
		lastUELinkInfoMsg: sampleMsg1,
		targetCell:        sampleMsg1.Ecgi,
		startTime:         tStart,
		hystCQI:           0,
		a3OffsetCQI:       0,
		timeToTrigger:     timeToTrigger,
	}
	go startA3Event(sampleA3Event, hoChan)
	ho := <-hoChan
	eTime := time.Since(tStart)

	if eTime.Milliseconds() < int64(timeToTrigger) {
		t.Errorf("TimeToTrigger is not working - elapsed time: %dms, TTT: %dms", eTime.Milliseconds(), timeToTrigger)
	}
	assert.Equal(t, ho.DstStation, sampleA3Event.targetCell)
}

func TestStartA3Event_S1_TTT1000(t *testing.T) {
	timeToTrigger := 1000
	GenSampleUELinkInfoMsgs()
	InitA3EventMap()
	hoChan := make(chan *nb.HandOverRequest)
	tStart := time.Now()
	sampleA3Event := &a3Event{
		chanUELinkInfo:    make(chan *nb.UELinkInfo),
		lastUELinkInfoMsg: sampleMsg1,
		targetCell:        sampleMsg1.Ecgi,
		startTime:         tStart,
		hystCQI:           0,
		a3OffsetCQI:       0,
		timeToTrigger:     timeToTrigger,
	}
	go startA3Event(sampleA3Event, hoChan)
	ho := <-hoChan
	eTime := time.Since(tStart)
	if eTime.Milliseconds() < int64(timeToTrigger) {
		t.Errorf("TimeToTrigger is not working - elapsed time: %dms, TTT: %dms", eTime.Milliseconds(), timeToTrigger)
	}
	assert.Equal(t, ho.DstStation, sampleA3Event.targetCell)
}

func TestStartA3Event_S2_TTT0(t *testing.T) {
	timeToTrigger := 0
	GenSampleUELinkInfoMsgs()
	InitA3EventMap()
	hoChan := make(chan *nb.HandOverRequest)
	tStart := time.Now()
	sampleA3Event := &a3Event{
		chanUELinkInfo:    make(chan *nb.UELinkInfo),
		lastUELinkInfoMsg: sampleMsg2,
		targetCell:        sampleMsg2.ChannelQualities[1].TargetEcgi,
		startTime:         tStart,
		hystCQI:           0,
		a3OffsetCQI:       0,
		timeToTrigger:     timeToTrigger,
	}
	go startA3Event(sampleA3Event, hoChan)
	ho := <-hoChan
	eTime := time.Since(tStart)

	if eTime.Milliseconds() < int64(timeToTrigger) {
		t.Errorf("TimeToTrigger is not working - elapsed time: %dms, TTT: %dms", eTime.Milliseconds(), timeToTrigger)
	}
	assert.Equal(t, ho.DstStation, sampleA3Event.targetCell)
}

func TestStartA3Event_S2_TTT1000(t *testing.T) {
	timeToTrigger := 1000
	GenSampleUELinkInfoMsgs()
	InitA3EventMap()
	hoChan := make(chan *nb.HandOverRequest)
	tStart := time.Now()
	sampleA3Event := &a3Event{
		chanUELinkInfo:    make(chan *nb.UELinkInfo),
		lastUELinkInfoMsg: sampleMsg2,
		targetCell:        sampleMsg2.ChannelQualities[1].TargetEcgi,
		startTime:         tStart,
		hystCQI:           0,
		a3OffsetCQI:       0,
		timeToTrigger:     timeToTrigger,
	}
	go startA3Event(sampleA3Event, hoChan)
	ho := <-hoChan
	eTime := time.Since(tStart)

	if eTime.Milliseconds() < int64(timeToTrigger) {
		t.Errorf("TimeToTrigger is not working - elapsed time: %dms, TTT: %dms", eTime.Milliseconds(), timeToTrigger)
	}
	assert.Equal(t, ho.DstStation, sampleA3Event.targetCell)
}

func TestStartA3Event_S2_3_TTT1000(t *testing.T) {
	timeToTrigger := 1000
	GenSampleUELinkInfoMsgs()
	InitA3EventMap()
	hoChan := make(chan *nb.HandOverRequest)
	tStart := time.Now()
	sampleA3Event := &a3Event{
		chanUELinkInfo:    make(chan *nb.UELinkInfo),
		lastUELinkInfoMsg: sampleMsg2,
		targetCell:        sampleMsg2.ChannelQualities[1].TargetEcgi,
		startTime:         tStart,
		hystCQI:           0,
		a3OffsetCQI:       0,
		timeToTrigger:     timeToTrigger,
	}
	go startA3Event(sampleA3Event, hoChan)
	time.Sleep(100 * time.Millisecond)
	sampleA3Event.chanUELinkInfo <- sampleMsg3

	time.Sleep(time.Duration(timeToTrigger)*time.Millisecond*10 + 1*time.Second)
	select {
	case ho := <-hoChan:
		assert.Equal(t, "0003", ho.DstStation.Ecid)
	case <-time.After(10 * time.Second):
		t.Error("Timeout: HO should be happened but not happened for 20 second")
	}
}

func TestStartA3Event_S2_2_TTT1000(t *testing.T) {
	timeToTrigger := 1000
	GenSampleUELinkInfoMsgs()
	InitA3EventMap()
	hoChan := make(chan *nb.HandOverRequest)
	tStart := time.Now()
	sampleA3Event := &a3Event{
		chanUELinkInfo:    make(chan *nb.UELinkInfo),
		lastUELinkInfoMsg: sampleMsg2,
		targetCell:        sampleMsg2.ChannelQualities[1].TargetEcgi,
		startTime:         tStart,
		hystCQI:           0,
		a3OffsetCQI:       0,
		timeToTrigger:     timeToTrigger,
	}
	go startA3Event(sampleA3Event, hoChan)
	time.Sleep(100 * time.Millisecond)
	sampleA3Event.chanUELinkInfo <- sampleMsg2

	time.Sleep(time.Duration(timeToTrigger)*time.Millisecond*10 + 1*time.Second)
	select {
	case ho := <-hoChan:
		assert.Equal(t, "0002", ho.DstStation.Ecid)
	case <-time.After(10 * time.Second):
		t.Error("Timeout: HO should be happened but not happened for 20 second")
	}
}

func TestStartA3Event_S2_1_TTT1000(t *testing.T) {
	timeToTrigger := 1000
	GenSampleUELinkInfoMsgs()
	InitA3EventMap()
	hoChan := make(chan *nb.HandOverRequest)
	tStart := time.Now()
	sampleA3Event := &a3Event{
		chanUELinkInfo:    make(chan *nb.UELinkInfo),
		lastUELinkInfoMsg: sampleMsg2,
		targetCell:        sampleMsg2.ChannelQualities[1].TargetEcgi,
		startTime:         tStart,
		hystCQI:           0,
		a3OffsetCQI:       0,
		timeToTrigger:     timeToTrigger,
	}
	go startA3Event(sampleA3Event, hoChan)
	time.Sleep(100 * time.Millisecond)
	sampleA3Event.chanUELinkInfo <- sampleMsg1

	time.Sleep(time.Duration(timeToTrigger)*time.Millisecond*10 + 1*time.Second)
	select {
	case ho := <-hoChan:
		t.Errorf("HO should not be happened, but happened %s to %s", ho.SrcStation.Ecid, ho.DstStation.Ecid)
	case <-time.After(10 * time.Second):
		assert.Equal(t, 0, len(hoChan))
	}
}

func TestStartA3Event_S2_3_TTT1000_ChanForceClose(t *testing.T) {
	timeToTrigger := 1000
	GenSampleUELinkInfoMsgs()
	InitA3EventMap()
	hoChan := make(chan *nb.HandOverRequest)
	tStart := time.Now()
	sampleA3Event := &a3Event{
		chanUELinkInfo:    make(chan *nb.UELinkInfo),
		lastUELinkInfoMsg: sampleMsg2,
		targetCell:        sampleMsg2.ChannelQualities[1].TargetEcgi,
		startTime:         tStart,
		hystCQI:           0,
		a3OffsetCQI:       0,
		timeToTrigger:     timeToTrigger,
	}

	defer func() {
		// Panic handler
		time.Sleep(time.Duration(timeToTrigger)*time.Millisecond*10 + 1*time.Second)
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		} else {
			return
		}
	}()

	go startA3Event(sampleA3Event, hoChan)
	time.Sleep(100 * time.Millisecond)
	A3EventMapMutex.Lock()
	close(sampleA3Event.chanUELinkInfo)
	A3EventMapMutex.Unlock()
	time.Sleep(100 * time.Millisecond)
	A3EventMapMutex.Lock()
	sampleA3Event.chanUELinkInfo <- sampleMsg3
	A3EventMapMutex.Unlock()
}
