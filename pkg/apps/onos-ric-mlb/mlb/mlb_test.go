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
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/onosproject/onos-ric/api/nb"
)

var testPLMNID = "315010"
var testMaxUEs = uint32(5)
var testThreshold = float64(0.5)
var testSampleStas1 map[string]*nb.StationInfo
var testSampleStaLinks1 map[string]*nb.StationLinkInfo
var testSampleUEs1 map[string]*nb.UEInfo
var testSampleStas2 map[string]*nb.StationInfo
var testSampleStaLinks2 map[string]*nb.StationLinkInfo
var testSampleUEs2 map[string]*nb.UEInfo

func genSampleRNIBs() {
	// Sample 1
	testSampleStas1 = make(map[string]*nb.StationInfo)
	testSampleStaLinks1 = make(map[string]*nb.StationLinkInfo)
	testSampleUEs1 = make(map[string]*nb.UEInfo)

	// StationInfo
	for i := 1; i < 10; i++ {
		tmpEcid := fmt.Sprintf("000%d", i)
		tmpSta := &nb.StationInfo{
			Ecgi: &nb.ECGI{
				Plmnid: testPLMNID,
				Ecid:   tmpEcid,
			},
			MaxNumConnectedUes: testMaxUEs,
		}

		testSampleStas1[tmpSta.Ecgi.String()] = tmpSta
	}
	// StationLinkInfo
	// Sta1: 2, 4
	t1StaLink1 := &nb.StationLinkInfo{
		Ecgi: &nb.ECGI{
			Plmnid: testPLMNID,
			Ecid:   "0001",
		},
		NeighborECGI: []*nb.ECGI{
			{
				Plmnid: testPLMNID,
				Ecid:   "0002",
			},
			{
				Plmnid: testPLMNID,
				Ecid:   "0004",
			},
		},
	}

	// Sta2: 1, 3, 5
	t1StaLink2 := &nb.StationLinkInfo{
		Ecgi: &nb.ECGI{
			Plmnid: testPLMNID,
			Ecid:   "0002",
		},
		NeighborECGI: []*nb.ECGI{
			{
				Plmnid: testPLMNID,
				Ecid:   "0001",
			},
			{
				Plmnid: testPLMNID,
				Ecid:   "0003",
			},
			{
				Plmnid: testPLMNID,
				Ecid:   "0005",
			},
		},
	}

	// Sta3: 2, 6
	t1StaLink3 := &nb.StationLinkInfo{
		Ecgi: &nb.ECGI{
			Plmnid: testPLMNID,
			Ecid:   "0003",
		},
		NeighborECGI: []*nb.ECGI{
			{
				Plmnid: testPLMNID,
				Ecid:   "0002",
			},
			{
				Plmnid: testPLMNID,
				Ecid:   "0006",
			},
		},
	}

	// Sta4: 1, 5, 7
	t1StaLink4 := &nb.StationLinkInfo{
		Ecgi: &nb.ECGI{
			Plmnid: testPLMNID,
			Ecid:   "0004",
		},
		NeighborECGI: []*nb.ECGI{
			{
				Plmnid: testPLMNID,
				Ecid:   "0001",
			},
			{
				Plmnid: testPLMNID,
				Ecid:   "0005",
			},
			{
				Plmnid: testPLMNID,
				Ecid:   "0007",
			},
		},
	}

	// Sta5: 2, 4, 6, 8
	t1StaLink5 := &nb.StationLinkInfo{
		Ecgi: &nb.ECGI{
			Plmnid: testPLMNID,
			Ecid:   "0005",
		},
		NeighborECGI: []*nb.ECGI{
			{
				Plmnid: testPLMNID,
				Ecid:   "0002",
			},
			{
				Plmnid: testPLMNID,
				Ecid:   "0004",
			},
			{
				Plmnid: testPLMNID,
				Ecid:   "0006",
			},
			{
				Plmnid: testPLMNID,
				Ecid:   "0008",
			},
		},
	}

	// Sta6: 3, 5, 9
	t1StaLink6 := &nb.StationLinkInfo{
		Ecgi: &nb.ECGI{
			Plmnid: testPLMNID,
			Ecid:   "0006",
		},
		NeighborECGI: []*nb.ECGI{
			{
				Plmnid: testPLMNID,
				Ecid:   "0003",
			},
			{
				Plmnid: testPLMNID,
				Ecid:   "0005",
			},
			{
				Plmnid: testPLMNID,
				Ecid:   "0009",
			},
		},
	}

	// Sta7: 4, 8
	t1StaLink7 := &nb.StationLinkInfo{
		Ecgi: &nb.ECGI{
			Plmnid: testPLMNID,
			Ecid:   "0007",
		},
		NeighborECGI: []*nb.ECGI{
			{
				Plmnid: testPLMNID,
				Ecid:   "0004",
			},
			{
				Plmnid: testPLMNID,
				Ecid:   "0008",
			},
		},
	}

	// Sta8: 5, 7, 9
	t1StaLink8 := &nb.StationLinkInfo{
		Ecgi: &nb.ECGI{
			Plmnid: testPLMNID,
			Ecid:   "0008",
		},
		NeighborECGI: []*nb.ECGI{
			{
				Plmnid: testPLMNID,
				Ecid:   "0005",
			},
			{
				Plmnid: testPLMNID,
				Ecid:   "0007",
			},
			{
				Plmnid: testPLMNID,
				Ecid:   "0009",
			},
		},
	}

	// Sta9: 6, 8
	t1StaLink9 := &nb.StationLinkInfo{
		Ecgi: &nb.ECGI{
			Plmnid: testPLMNID,
			Ecid:   "0009",
		},
		NeighborECGI: []*nb.ECGI{
			{
				Plmnid: testPLMNID,
				Ecid:   "0006",
			},
			{
				Plmnid: testPLMNID,
				Ecid:   "0008",
			},
		},
	}
	testSampleStaLinks1[t1StaLink1.Ecgi.String()] = t1StaLink1
	testSampleStaLinks1[t1StaLink2.Ecgi.String()] = t1StaLink2
	testSampleStaLinks1[t1StaLink3.Ecgi.String()] = t1StaLink3
	testSampleStaLinks1[t1StaLink4.Ecgi.String()] = t1StaLink4
	testSampleStaLinks1[t1StaLink5.Ecgi.String()] = t1StaLink5
	testSampleStaLinks1[t1StaLink6.Ecgi.String()] = t1StaLink6
	testSampleStaLinks1[t1StaLink7.Ecgi.String()] = t1StaLink7
	testSampleStaLinks1[t1StaLink8.Ecgi.String()] = t1StaLink8
	testSampleStaLinks1[t1StaLink9.Ecgi.String()] = t1StaLink9

	// UELinkInfo
	// UE1 - conn w/ STA8
	t1UeInfo1 := &nb.UEInfo{
		Crnti: "00001",
		Ecgi: &nb.ECGI{
			Plmnid: testPLMNID,
			Ecid:   "0008",
		},
		Imsi: "315010999900001",
	}
	t1UeInfo2 := &nb.UEInfo{
		Crnti: "00002",
		Ecgi: &nb.ECGI{
			Plmnid: testPLMNID,
			Ecid:   "0008",
		},
		Imsi: "315010999900002",
	}
	t1UeInfo3 := &nb.UEInfo{
		Crnti: "00003",
		Ecgi: &nb.ECGI{
			Plmnid: testPLMNID,
			Ecid:   "0008",
		},
		Imsi: "315010999900003",
	}

	testSampleUEs1[t1UeInfo1.Imsi] = t1UeInfo1
	testSampleUEs1[t1UeInfo2.Imsi] = t1UeInfo2
	testSampleUEs1[t1UeInfo3.Imsi] = t1UeInfo3

	// Sample 2
	testSampleStas2 = make(map[string]*nb.StationInfo)
	testSampleStaLinks2 = make(map[string]*nb.StationLinkInfo)
	testSampleUEs2 = make(map[string]*nb.UEInfo)

	// StationInfo
	for i := 1; i < 10; i++ {
		tmpEcid := fmt.Sprintf("000%d", i)
		tmpSta := &nb.StationInfo{
			Ecgi: &nb.ECGI{
				Plmnid: testPLMNID,
				Ecid:   tmpEcid,
			},
			MaxNumConnectedUes: testMaxUEs,
		}

		testSampleStas2[tmpSta.Ecgi.String()] = tmpSta
	}
	// StationLinkInfo
	// Sta1: 2, 4
	t2StaLink1 := &nb.StationLinkInfo{
		Ecgi: &nb.ECGI{
			Plmnid: testPLMNID,
			Ecid:   "0001",
		},
		NeighborECGI: []*nb.ECGI{
			{
				Plmnid: testPLMNID,
				Ecid:   "0002",
			},
			{
				Plmnid: testPLMNID,
				Ecid:   "0004",
			},
		},
	}

	// Sta2: 1, 3, 5
	t2StaLink2 := &nb.StationLinkInfo{
		Ecgi: &nb.ECGI{
			Plmnid: testPLMNID,
			Ecid:   "0002",
		},
		NeighborECGI: []*nb.ECGI{
			{
				Plmnid: testPLMNID,
				Ecid:   "0001",
			},
			{
				Plmnid: testPLMNID,
				Ecid:   "0003",
			},
			{
				Plmnid: testPLMNID,
				Ecid:   "0005",
			},
		},
	}

	// Sta3: 2, 6
	t2StaLink3 := &nb.StationLinkInfo{
		Ecgi: &nb.ECGI{
			Plmnid: testPLMNID,
			Ecid:   "0003",
		},
		NeighborECGI: []*nb.ECGI{
			{
				Plmnid: testPLMNID,
				Ecid:   "0002",
			},
			{
				Plmnid: testPLMNID,
				Ecid:   "0006",
			},
		},
	}

	// Sta4: 1, 5, 7
	t2StaLink4 := &nb.StationLinkInfo{
		Ecgi: &nb.ECGI{
			Plmnid: testPLMNID,
			Ecid:   "0004",
		},
		NeighborECGI: []*nb.ECGI{
			{
				Plmnid: testPLMNID,
				Ecid:   "0001",
			},
			{
				Plmnid: testPLMNID,
				Ecid:   "0005",
			},
			{
				Plmnid: testPLMNID,
				Ecid:   "0007",
			},
		},
	}

	// Sta5: 2, 4, 6, 8
	t2StaLink5 := &nb.StationLinkInfo{
		Ecgi: &nb.ECGI{
			Plmnid: testPLMNID,
			Ecid:   "0005",
		},
		NeighborECGI: []*nb.ECGI{
			{
				Plmnid: testPLMNID,
				Ecid:   "0002",
			},
			{
				Plmnid: testPLMNID,
				Ecid:   "0004",
			},
			{
				Plmnid: testPLMNID,
				Ecid:   "0006",
			},
			{
				Plmnid: testPLMNID,
				Ecid:   "0008",
			},
		},
	}

	// Sta6: 3, 5, 9
	t2StaLink6 := &nb.StationLinkInfo{
		Ecgi: &nb.ECGI{
			Plmnid: testPLMNID,
			Ecid:   "0006",
		},
		NeighborECGI: []*nb.ECGI{
			{
				Plmnid: testPLMNID,
				Ecid:   "0003",
			},
			{
				Plmnid: testPLMNID,
				Ecid:   "0005",
			},
			{
				Plmnid: testPLMNID,
				Ecid:   "0009",
			},
		},
	}

	// Sta7: 4, 8
	t2StaLink7 := &nb.StationLinkInfo{
		Ecgi: &nb.ECGI{
			Plmnid: testPLMNID,
			Ecid:   "0007",
		},
		NeighborECGI: []*nb.ECGI{
			{
				Plmnid: testPLMNID,
				Ecid:   "0004",
			},
			{
				Plmnid: testPLMNID,
				Ecid:   "0008",
			},
		},
	}

	// Sta8: 5, 7, 9
	t2StaLink8 := &nb.StationLinkInfo{
		Ecgi: &nb.ECGI{
			Plmnid: testPLMNID,
			Ecid:   "0008",
		},
		NeighborECGI: []*nb.ECGI{
			{
				Plmnid: testPLMNID,
				Ecid:   "0005",
			},
			{
				Plmnid: testPLMNID,
				Ecid:   "0007",
			},
			{
				Plmnid: testPLMNID,
				Ecid:   "0009",
			},
		},
	}

	// Sta9: 6, 8
	t2StaLink9 := &nb.StationLinkInfo{
		Ecgi: &nb.ECGI{
			Plmnid: testPLMNID,
			Ecid:   "0009",
		},
		NeighborECGI: []*nb.ECGI{
			{
				Plmnid: testPLMNID,
				Ecid:   "0006",
			},
			{
				Plmnid: testPLMNID,
				Ecid:   "0008",
			},
		},
	}
	testSampleStaLinks2[t2StaLink1.Ecgi.String()] = t2StaLink1
	testSampleStaLinks2[t2StaLink2.Ecgi.String()] = t2StaLink2
	testSampleStaLinks2[t2StaLink3.Ecgi.String()] = t2StaLink3
	testSampleStaLinks2[t2StaLink4.Ecgi.String()] = t2StaLink4
	testSampleStaLinks2[t2StaLink5.Ecgi.String()] = t2StaLink5
	testSampleStaLinks2[t2StaLink6.Ecgi.String()] = t2StaLink6
	testSampleStaLinks2[t2StaLink7.Ecgi.String()] = t2StaLink7
	testSampleStaLinks2[t2StaLink8.Ecgi.String()] = t2StaLink8
	testSampleStaLinks2[t2StaLink9.Ecgi.String()] = t2StaLink9

	// UELinkInfo
	t2UeInfo1 := &nb.UEInfo{
		Crnti: "00001",
		Ecgi: &nb.ECGI{
			Plmnid: testPLMNID,
			Ecid:   "0008",
		},
		Imsi: "315010999900001",
	}
	t2UeInfo2 := &nb.UEInfo{
		Crnti: "00002",
		Ecgi: &nb.ECGI{
			Plmnid: testPLMNID,
			Ecid:   "0008",
		},
		Imsi: "315010999900002",
	}
	t2UeInfo3 := &nb.UEInfo{
		Crnti: "00003",
		Ecgi: &nb.ECGI{
			Plmnid: testPLMNID,
			Ecid:   "0008",
		},
		Imsi: "315010999900003",
	}
	t2UeInfo4 := &nb.UEInfo{
		Crnti: "00004",
		Ecgi: &nb.ECGI{
			Plmnid: testPLMNID,
			Ecid:   "0005",
		},
		Imsi: "315010999900004",
	}
	t2UeInfo5 := &nb.UEInfo{
		Crnti: "00005",
		Ecgi: &nb.ECGI{
			Plmnid: testPLMNID,
			Ecid:   "0005",
		},
		Imsi: "315010999900005",
	}
	t2UeInfo6 := &nb.UEInfo{
		Crnti: "00006",
		Ecgi: &nb.ECGI{
			Plmnid: testPLMNID,
			Ecid:   "0005",
		},
		Imsi: "315010999900006",
	}

	testSampleUEs2[t2UeInfo1.Imsi] = t2UeInfo1
	testSampleUEs2[t2UeInfo2.Imsi] = t2UeInfo2
	testSampleUEs2[t2UeInfo3.Imsi] = t2UeInfo3
	testSampleUEs2[t2UeInfo4.Imsi] = t2UeInfo4
	testSampleUEs2[t2UeInfo5.Imsi] = t2UeInfo5
	testSampleUEs2[t2UeInfo6.Imsi] = t2UeInfo6
}

// Test case 1: only single station is the overloaded station.
func TestMLBDecisionMaker1(t *testing.T) {
	genSampleRNIBs()

	txPwrSta1 := nb.StationPowerOffset_PA_DB_0
	txPwrSta2 := nb.StationPowerOffset_PA_DB_0
	txPwrSta3 := nb.StationPowerOffset_PA_DB_0
	txPwrSta4 := nb.StationPowerOffset_PA_DB_0
	txPwrSta5 := nb.StationPowerOffset_PA_DB_0
	txPwrSta6 := nb.StationPowerOffset_PA_DB_0
	txPwrSta7 := nb.StationPowerOffset_PA_DB_0
	txPwrSta8 := nb.StationPowerOffset_PA_DB_0
	txPwrSta9 := nb.StationPowerOffset_PA_DB_0

	testResult, _ := MLBDecisionMaker(testSampleStas1, testSampleStaLinks1, testSampleUEs1, &testThreshold)

	assert.Equal(t, len(*testResult), 4)

	for _, rpReq := range *testResult {
		switch rpReq.GetEcgi().GetEcid() {
		case "0001":
			txPwrSta1 = rpReq.GetOffset()
		case "0002":
			txPwrSta2 = rpReq.GetOffset()
		case "0003":
			txPwrSta3 = rpReq.GetOffset()
		case "0004":
			txPwrSta4 = rpReq.GetOffset()
		case "0005":
			txPwrSta5 = rpReq.GetOffset()
		case "0006":
			txPwrSta6 = rpReq.GetOffset()
		case "0007":
			txPwrSta7 = rpReq.GetOffset()
		case "0008":
			txPwrSta8 = rpReq.GetOffset()
		case "0009":
			txPwrSta9 = rpReq.GetOffset()
		}
	}

	assert.Equal(t, txPwrSta1, nb.StationPowerOffset_PA_DB_0)
	assert.Equal(t, txPwrSta2, nb.StationPowerOffset_PA_DB_0)
	assert.Equal(t, txPwrSta3, nb.StationPowerOffset_PA_DB_0)
	assert.Equal(t, txPwrSta4, nb.StationPowerOffset_PA_DB_0)
	assert.Equal(t, txPwrSta5, nb.StationPowerOffset_PA_DB_3)
	assert.Equal(t, txPwrSta6, nb.StationPowerOffset_PA_DB_0)
	assert.Equal(t, txPwrSta7, nb.StationPowerOffset_PA_DB_3)
	assert.Equal(t, txPwrSta8, nb.StationPowerOffset_PA_DB_MINUS3)
	assert.Equal(t, txPwrSta9, nb.StationPowerOffset_PA_DB_3)
}

// Test case 2:
func TestMLBDecisionMaker2(t *testing.T) {
	genSampleRNIBs()

	txPwrSta1 := nb.StationPowerOffset_PA_DB_0
	txPwrSta2 := nb.StationPowerOffset_PA_DB_0
	txPwrSta3 := nb.StationPowerOffset_PA_DB_0
	txPwrSta4 := nb.StationPowerOffset_PA_DB_0
	txPwrSta5 := nb.StationPowerOffset_PA_DB_0
	txPwrSta6 := nb.StationPowerOffset_PA_DB_0
	txPwrSta7 := nb.StationPowerOffset_PA_DB_0
	txPwrSta8 := nb.StationPowerOffset_PA_DB_0
	txPwrSta9 := nb.StationPowerOffset_PA_DB_0

	testResult, _ := MLBDecisionMaker(testSampleStas2, testSampleStaLinks2, testSampleUEs2, &testThreshold)

	assert.Equal(t, len(*testResult), 7)

	for _, rpReq := range *testResult {
		switch rpReq.GetEcgi().GetEcid() {
		case "0001":
			txPwrSta1 = rpReq.GetOffset()
		case "0002":
			txPwrSta2 = rpReq.GetOffset()
		case "0003":
			txPwrSta3 = rpReq.GetOffset()
		case "0004":
			txPwrSta4 = rpReq.GetOffset()
		case "0005":
			txPwrSta5 = rpReq.GetOffset()
		case "0006":
			txPwrSta6 = rpReq.GetOffset()
		case "0007":
			txPwrSta7 = rpReq.GetOffset()
		case "0008":
			txPwrSta8 = rpReq.GetOffset()
		case "0009":
			txPwrSta9 = rpReq.GetOffset()
		}
	}

	assert.Equal(t, txPwrSta1, nb.StationPowerOffset_PA_DB_0)
	assert.Equal(t, txPwrSta2, nb.StationPowerOffset_PA_DB_3)
	assert.Equal(t, txPwrSta3, nb.StationPowerOffset_PA_DB_0)
	assert.Equal(t, txPwrSta4, nb.StationPowerOffset_PA_DB_3)
	assert.Equal(t, txPwrSta5, nb.StationPowerOffset_PA_DB_MINUS3)
	assert.Equal(t, txPwrSta6, nb.StationPowerOffset_PA_DB_3)
	assert.Equal(t, txPwrSta7, nb.StationPowerOffset_PA_DB_3)
	assert.Equal(t, txPwrSta8, nb.StationPowerOffset_PA_DB_MINUS3)
	assert.Equal(t, txPwrSta9, nb.StationPowerOffset_PA_DB_3)
}

func getStaUEJointMap(i int) map[string]*StaUeJointLink {
	genSampleRNIBs()
	staUeJointMap := make(map[string]*StaUeJointLink)

	tStart := time.Now()
	if i == 1 {
		for _, v := range testSampleStas1 {
			tmpStaUeJointLink := &StaUeJointLink{
				TimeStamp:   tStart,
				PlmnID:      v.GetEcgi().GetPlmnid(),
				Ecid:        v.GetEcgi().GetEcid(),
				MaxNumUes:   v.GetMaxNumConnectedUes(),
				NumUes:      0,
				Pa:          0,
				ElapsedTime: 0,
			}
			staUeJointMap[v.Ecgi.String()] = tmpStaUeJointLink
		}
	} else {
		for _, v := range testSampleStas2 {
			tmpStaUeJointLink := &StaUeJointLink{
				TimeStamp:   tStart,
				PlmnID:      v.GetEcgi().GetPlmnid(),
				Ecid:        v.GetEcgi().GetEcid(),
				MaxNumUes:   v.GetMaxNumConnectedUes(),
				NumUes:      0,
				Pa:          0,
				ElapsedTime: 0,
			}
			staUeJointMap[v.Ecgi.String()] = tmpStaUeJointLink
		}
	}
	return staUeJointMap
}

func TestCountUEs1(t *testing.T) {
	genSampleRNIBs()
	staUeJointMap := getStaUEJointMap(1)
	countUEs(staUeJointMap, testSampleUEs1)
	for _, v := range staUeJointMap {
		switch v.Ecid {
		case "0001":
			assert.Equal(t, int32(0), v.NumUes)
		case "0002":
			assert.Equal(t, int32(0), v.NumUes)
		case "0003":
			assert.Equal(t, int32(0), v.NumUes)
		case "0004":
			assert.Equal(t, int32(0), v.NumUes)
		case "0005":
			assert.Equal(t, int32(0), v.NumUes)
		case "0006":
			assert.Equal(t, int32(0), v.NumUes)
		case "0007":
			assert.Equal(t, int32(0), v.NumUes)
		case "0008":
			assert.Equal(t, int32(3), v.NumUes)
		case "0009":
			assert.Equal(t, int32(0), v.NumUes)
		}
	}
}

func TestCountUEs2(t *testing.T) {
	genSampleRNIBs()
	staUeJointMap := getStaUEJointMap(2)
	countUEs(staUeJointMap, testSampleUEs2)
	for _, v := range staUeJointMap {
		switch v.Ecid {
		case "0001":
			assert.Equal(t, int32(0), v.NumUes)
		case "0002":
			assert.Equal(t, int32(0), v.NumUes)
		case "0003":
			assert.Equal(t, int32(0), v.NumUes)
		case "0004":
			assert.Equal(t, int32(0), v.NumUes)
		case "0005":
			assert.Equal(t, int32(3), v.NumUes)
		case "0006":
			assert.Equal(t, int32(0), v.NumUes)
		case "0007":
			assert.Equal(t, int32(0), v.NumUes)
		case "0008":
			assert.Equal(t, int32(3), v.NumUes)
		case "0009":
			assert.Equal(t, int32(0), v.NumUes)
		}
	}
}

func TestWrongCase_CountUEs(t *testing.T) {
	genSampleRNIBs()
	staUeJointMap := getStaUEJointMap(1)
	for _, v := range testSampleUEs1 {
		v.Ecgi.Ecid = "1111"
	}

	countUEs(staUeJointMap, testSampleUEs1)
	for _, v := range staUeJointMap {
		switch v.Ecid {
		case "0001":
			assert.Equal(t, int32(0), v.NumUes)
		case "0002":
			assert.Equal(t, int32(0), v.NumUes)
		case "0003":
			assert.Equal(t, int32(0), v.NumUes)
		case "0004":
			assert.Equal(t, int32(0), v.NumUes)
		case "0005":
			assert.Equal(t, int32(0), v.NumUes)
		case "0006":
			assert.Equal(t, int32(0), v.NumUes)
		case "0007":
			assert.Equal(t, int32(0), v.NumUes)
		case "0008":
			assert.Equal(t, int32(0), v.NumUes)
		case "0009":
			assert.Equal(t, int32(0), v.NumUes)
		}
	}
}

func TestGetOverloadedStationList1(t *testing.T) {
	genSampleRNIBs()
	staUeJointMap := getStaUEJointMap(1)
	countUEs(staUeJointMap, testSampleUEs1)
	overloadedStas := getOverloadedStationList(staUeJointMap, &testThreshold)
	assert.Equal(t, 1, len(*overloadedStas))
	for _, v := range *overloadedStas {
		switch v.Ecid {
		case "0008":
		default:
			t.Errorf("%s STA should not be an overloaded STA", v.Ecid)
		}
	}
}

func TestGetOverloadedStationList2(t *testing.T) {
	genSampleRNIBs()
	staUeJointMap := getStaUEJointMap(2)
	countUEs(staUeJointMap, testSampleUEs2)
	overloadedStas := getOverloadedStationList(staUeJointMap, &testThreshold)
	assert.Equal(t, 2, len(*overloadedStas))
	for _, v := range *overloadedStas {
		switch v.Ecid {
		case "0005":
		case "0008":
		default:
			t.Errorf("%s STA should not be an overloaded STA", v.Ecid)
		}
	}
}

func TestGetStasToBeExpanded1(t *testing.T) {
	genSampleRNIBs()
	staUeJointMap := getStaUEJointMap(1)
	countUEs(staUeJointMap, testSampleUEs1)
	overloadedStas := getOverloadedStationList(staUeJointMap, &testThreshold)
	assert.Equal(t, 1, len(*overloadedStas))
	stasToBeExpanded := getStasToBeExpanded(staUeJointMap, overloadedStas, testSampleStaLinks1)
	assert.Equal(t, 3, len(*stasToBeExpanded))
	for _, v := range *stasToBeExpanded {
		switch v.Ecid {
		case "0005":
		case "0007":
		case "0009":
		default:
			t.Errorf("%s STA should not be an overloaded STA", v.Ecid)
		}
	}
}

func TestGetStasToBeExpanded2(t *testing.T) {
	genSampleRNIBs()
	staUeJointMap := getStaUEJointMap(2)
	countUEs(staUeJointMap, testSampleUEs2)
	overloadedStas := getOverloadedStationList(staUeJointMap, &testThreshold)
	assert.Equal(t, 2, len(*overloadedStas))
	stasToBeExpanded := getStasToBeExpanded(staUeJointMap, overloadedStas, testSampleStaLinks2)
	assert.Equal(t, 5, len(*stasToBeExpanded))
	for _, v := range *stasToBeExpanded {
		switch v.Ecid {
		case "0002":
		case "0004":
		case "0006":
		case "0007":
		case "0009":
		default:
			t.Errorf("%s STA should not be an overloaded STA", v.Ecid)
		}
	}
}

func TestGetNeighborStaEcgi_ID0001(t *testing.T) {
	ecid := "0001"
	genSampleRNIBs()
	nCells := getNeighborStaEcgi(testPLMNID, ecid, testSampleStaLinks1)
	for _, v := range nCells {
		switch v.Ecid {
		case "0002":
		case "0004":
		default:
			t.Errorf("ECID:%s should not be a neighbor cell of ECID:%s", v.Ecid, ecid)
		}
	}
}

func TestGetNeighborStaEcgi_ID0002(t *testing.T) {
	ecid := "0002"
	genSampleRNIBs()
	nCells := getNeighborStaEcgi(testPLMNID, ecid, testSampleStaLinks1)
	for _, v := range nCells {
		switch v.Ecid {
		case "0001":
		case "0003":
		case "0005":
		default:
			t.Errorf("ECID:%s should not be a neighbor cell of ECID:%s", v.Ecid, ecid)
		}
	}
}

func TestGetNeighborStaEcgi_ID0003(t *testing.T) {
	ecid := "0003"
	genSampleRNIBs()
	nCells := getNeighborStaEcgi(testPLMNID, ecid, testSampleStaLinks1)
	for _, v := range nCells {
		switch v.Ecid {
		case "0002":
		case "0006":
		default:
			t.Errorf("ECID:%s should not be a neighbor cell of ECID:%s", v.Ecid, ecid)
		}
	}
}

func TestGetNeighborStaEcgi_ID0004(t *testing.T) {
	ecid := "0004"
	genSampleRNIBs()
	nCells := getNeighborStaEcgi(testPLMNID, ecid, testSampleStaLinks1)
	for _, v := range nCells {
		switch v.Ecid {
		case "0001":
		case "0005":
		case "0007":
		default:
			t.Errorf("ECID:%s should not be a neighbor cell of ECID:%s", v.Ecid, ecid)
		}
	}
}

func TestGetNeighborStaEcgi_ID0005(t *testing.T) {
	ecid := "0005"
	genSampleRNIBs()
	nCells := getNeighborStaEcgi(testPLMNID, ecid, testSampleStaLinks1)
	for _, v := range nCells {
		switch v.Ecid {
		case "0002":
		case "0004":
		case "0006":
		case "0008":
		default:
			t.Errorf("ECID:%s should not be a neighbor cell of ECID:%s", v.Ecid, ecid)
		}
	}
}

func TestGetNeighborStaEcgi_ID0006(t *testing.T) {
	ecid := "0006"
	genSampleRNIBs()
	nCells := getNeighborStaEcgi(testPLMNID, ecid, testSampleStaLinks1)
	for _, v := range nCells {
		switch v.Ecid {
		case "0003":
		case "0005":
		case "0009":
		default:
			t.Errorf("ECID:%s should not be a neighbor cell of ECID:%s", v.Ecid, ecid)
		}
	}
}

func TestGetNeighborStaEcgi_ID0007(t *testing.T) {
	ecid := "0007"
	genSampleRNIBs()
	nCells := getNeighborStaEcgi(testPLMNID, ecid, testSampleStaLinks1)
	for _, v := range nCells {
		switch v.Ecid {
		case "0004":
		case "0008":
		default:
			t.Errorf("ECID:%s should not be a neighbor cell of ECID:%s", v.Ecid, ecid)
		}
	}
}

func TestGetNeighborStaEcgi_ID0008(t *testing.T) {
	ecid := "0008"
	genSampleRNIBs()
	nCells := getNeighborStaEcgi(testPLMNID, ecid, testSampleStaLinks1)
	for _, v := range nCells {
		switch v.Ecid {
		case "0005":
		case "0007":
		case "0009":
		default:
			t.Errorf("ECID:%s should not be a neighbor cell of ECID:%s", v.Ecid, ecid)
		}
	}
}

func TestGetNeighborStaEcgi_ID0009(t *testing.T) {
	ecid := "0009"
	genSampleRNIBs()
	nCells := getNeighborStaEcgi(testPLMNID, ecid, testSampleStaLinks1)
	for _, v := range nCells {
		switch v.Ecid {
		case "0006":
		case "0008":
		default:
			t.Errorf("ECID:%s should not be a neighbor cell of ECID:%s", v.Ecid, ecid)
		}
	}
}

func TestGetNeighborStaEcgi_WrongID(t *testing.T) {
	ecid := "0010"
	genSampleRNIBs()
	nCells := getNeighborStaEcgi(testPLMNID, ecid, testSampleStaLinks1)
	assert.Nil(t, nCells)
}
