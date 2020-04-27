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

	"github.com/stretchr/testify/assert"

	"github.com/onosproject/onos-ric/api/nb"
)

// Test case 1: only single station is the overloaded station.
func TestMLBDecisionMaker1(t *testing.T) {

	// Define common parameter
	plmnid := "315010"
	maxUes := uint32(5)
	threshold := 0.5

	// Define test cases
	tStas := make(map[string]*nb.StationInfo)
	tStaLinks := make(map[string]*nb.StationLinkInfo)
	tUes := make(map[string]*nb.UEInfo)

	// StationInfo
	for i := 1; i < 10; i++ {
		tmpEcid := fmt.Sprintf("000%d", i)
		tmpSta := &nb.StationInfo{
			Ecgi: &nb.ECGI{
				Plmnid: plmnid,
				Ecid:   tmpEcid,
			},
			MaxNumConnectedUes: maxUes,
		}

		tStas[tmpSta.Ecgi.String()] = tmpSta
	}
	// StationLinkInfo
	// Sta1: 2, 4
	tStaLink1 := &nb.StationLinkInfo{
		Ecgi: &nb.ECGI{
			Plmnid: plmnid,
			Ecid:   "0001",
		},
		NeighborECGI: []*nb.ECGI{
			{
				Plmnid: plmnid,
				Ecid:   "0002",
			},
			{
				Plmnid: plmnid,
				Ecid:   "0004",
			},
		},
	}

	// Sta2: 1, 3, 5
	tStaLink2 := &nb.StationLinkInfo{
		Ecgi: &nb.ECGI{
			Plmnid: plmnid,
			Ecid:   "0002",
		},
		NeighborECGI: []*nb.ECGI{
			{
				Plmnid: plmnid,
				Ecid:   "0001",
			},
			{
				Plmnid: plmnid,
				Ecid:   "0003",
			},
			{
				Plmnid: plmnid,
				Ecid:   "0005",
			},
		},
	}

	// Sta3: 2, 6
	tStaLink3 := &nb.StationLinkInfo{
		Ecgi: &nb.ECGI{
			Plmnid: plmnid,
			Ecid:   "0003",
		},
		NeighborECGI: []*nb.ECGI{
			{
				Plmnid: plmnid,
				Ecid:   "0002",
			},
			{
				Plmnid: plmnid,
				Ecid:   "0006",
			},
		},
	}

	// Sta4: 1, 5, 7
	tStaLink4 := &nb.StationLinkInfo{
		Ecgi: &nb.ECGI{
			Plmnid: plmnid,
			Ecid:   "0004",
		},
		NeighborECGI: []*nb.ECGI{
			{
				Plmnid: plmnid,
				Ecid:   "0001",
			},
			{
				Plmnid: plmnid,
				Ecid:   "0005",
			},
			{
				Plmnid: plmnid,
				Ecid:   "0007",
			},
		},
	}

	// Sta5: 2, 4, 6, 8
	tStaLink5 := &nb.StationLinkInfo{
		Ecgi: &nb.ECGI{
			Plmnid: plmnid,
			Ecid:   "0005",
		},
		NeighborECGI: []*nb.ECGI{
			{
				Plmnid: plmnid,
				Ecid:   "0002",
			},
			{
				Plmnid: plmnid,
				Ecid:   "0004",
			},
			{
				Plmnid: plmnid,
				Ecid:   "0006",
			},
			{
				Plmnid: plmnid,
				Ecid:   "0008",
			},
		},
	}

	// Sta6: 3, 5, 9
	tStaLink6 := &nb.StationLinkInfo{
		Ecgi: &nb.ECGI{
			Plmnid: plmnid,
			Ecid:   "0006",
		},
		NeighborECGI: []*nb.ECGI{
			{
				Plmnid: plmnid,
				Ecid:   "0003",
			},
			{
				Plmnid: plmnid,
				Ecid:   "0005",
			},
			{
				Plmnid: plmnid,
				Ecid:   "0009",
			},
		},
	}

	// Sta7: 4, 8
	tStaLink7 := &nb.StationLinkInfo{
		Ecgi: &nb.ECGI{
			Plmnid: plmnid,
			Ecid:   "0007",
		},
		NeighborECGI: []*nb.ECGI{
			{
				Plmnid: plmnid,
				Ecid:   "0004",
			},
			{
				Plmnid: plmnid,
				Ecid:   "0008",
			},
		},
	}

	// Sta8: 5, 7, 9
	tStaLink8 := &nb.StationLinkInfo{
		Ecgi: &nb.ECGI{
			Plmnid: plmnid,
			Ecid:   "0008",
		},
		NeighborECGI: []*nb.ECGI{
			{
				Plmnid: plmnid,
				Ecid:   "0005",
			},
			{
				Plmnid: plmnid,
				Ecid:   "0007",
			},
			{
				Plmnid: plmnid,
				Ecid:   "0009",
			},
		},
	}

	// Sta9: 6, 8
	tStaLink9 := &nb.StationLinkInfo{
		Ecgi: &nb.ECGI{
			Plmnid: plmnid,
			Ecid:   "0009",
		},
		NeighborECGI: []*nb.ECGI{
			{
				Plmnid: plmnid,
				Ecid:   "0006",
			},
			{
				Plmnid: plmnid,
				Ecid:   "0008",
			},
		},
	}
	tStaLinks[tStaLink1.Ecgi.String()] = tStaLink1
	tStaLinks[tStaLink2.Ecgi.String()] = tStaLink2
	tStaLinks[tStaLink3.Ecgi.String()] = tStaLink3
	tStaLinks[tStaLink4.Ecgi.String()] = tStaLink4
	tStaLinks[tStaLink5.Ecgi.String()] = tStaLink5
	tStaLinks[tStaLink6.Ecgi.String()] = tStaLink6
	tStaLinks[tStaLink7.Ecgi.String()] = tStaLink7
	tStaLinks[tStaLink8.Ecgi.String()] = tStaLink8
	tStaLinks[tStaLink9.Ecgi.String()] = tStaLink9

	// UELinkInfo
	// UE1 - conn w/ STA8
	tUeInfo1 := &nb.UEInfo{
		Crnti: "00001",
		Ecgi: &nb.ECGI{
			Plmnid: plmnid,
			Ecid:   "0008",
		},
		Imsi: "315010999900001",
	}
	tUeInfo2 := &nb.UEInfo{
		Crnti: "00002",
		Ecgi: &nb.ECGI{
			Plmnid: plmnid,
			Ecid:   "0008",
		},
		Imsi: "315010999900002",
	}
	tUeInfo3 := &nb.UEInfo{
		Crnti: "00003",
		Ecgi: &nb.ECGI{
			Plmnid: plmnid,
			Ecid:   "0008",
		},
		Imsi: "315010999900003",
	}

	tUes[tUeInfo1.Imsi] = tUeInfo1
	tUes[tUeInfo2.Imsi] = tUeInfo2
	tUes[tUeInfo3.Imsi] = tUeInfo3
	txPwrSta1 := nb.StationPowerOffset_PA_DB_0
	txPwrSta2 := nb.StationPowerOffset_PA_DB_0
	txPwrSta3 := nb.StationPowerOffset_PA_DB_0
	txPwrSta4 := nb.StationPowerOffset_PA_DB_0
	txPwrSta5 := nb.StationPowerOffset_PA_DB_0
	txPwrSta6 := nb.StationPowerOffset_PA_DB_0
	txPwrSta7 := nb.StationPowerOffset_PA_DB_0
	txPwrSta8 := nb.StationPowerOffset_PA_DB_0
	txPwrSta9 := nb.StationPowerOffset_PA_DB_0

	testResult, _ := MLBDecisionMaker(tStas, tStaLinks, tUes, &threshold)

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
	// Define common parameter
	plmnid := "315010"
	maxUes := uint32(5)
	threshold := 0.5

	// Define test cases
	tStas := make(map[string]*nb.StationInfo)
	tStaLinks := make(map[string]*nb.StationLinkInfo)
	tUes := make(map[string]*nb.UEInfo)

	// StationInfo
	for i := 1; i < 10; i++ {
		tmpEcid := fmt.Sprintf("000%d", i)
		tmpSta := &nb.StationInfo{
			Ecgi: &nb.ECGI{
				Plmnid: plmnid,
				Ecid:   tmpEcid,
			},
			MaxNumConnectedUes: maxUes,
		}

		tStas[tmpSta.Ecgi.String()] = tmpSta
	}
	// StationLinkInfo
	// Sta1: 2, 4
	tStaLink1 := &nb.StationLinkInfo{
		Ecgi: &nb.ECGI{
			Plmnid: plmnid,
			Ecid:   "0001",
		},
		NeighborECGI: []*nb.ECGI{
			{
				Plmnid: plmnid,
				Ecid:   "0002",
			},
			{
				Plmnid: plmnid,
				Ecid:   "0004",
			},
		},
	}

	// Sta2: 1, 3, 5
	tStaLink2 := &nb.StationLinkInfo{
		Ecgi: &nb.ECGI{
			Plmnid: plmnid,
			Ecid:   "0002",
		},
		NeighborECGI: []*nb.ECGI{
			{
				Plmnid: plmnid,
				Ecid:   "0001",
			},
			{
				Plmnid: plmnid,
				Ecid:   "0003",
			},
			{
				Plmnid: plmnid,
				Ecid:   "0005",
			},
		},
	}

	// Sta3: 2, 6
	tStaLink3 := &nb.StationLinkInfo{
		Ecgi: &nb.ECGI{
			Plmnid: plmnid,
			Ecid:   "0003",
		},
		NeighborECGI: []*nb.ECGI{
			{
				Plmnid: plmnid,
				Ecid:   "0002",
			},
			{
				Plmnid: plmnid,
				Ecid:   "0006",
			},
		},
	}

	// Sta4: 1, 5, 7
	tStaLink4 := &nb.StationLinkInfo{
		Ecgi: &nb.ECGI{
			Plmnid: plmnid,
			Ecid:   "0004",
		},
		NeighborECGI: []*nb.ECGI{
			{
				Plmnid: plmnid,
				Ecid:   "0001",
			},
			{
				Plmnid: plmnid,
				Ecid:   "0005",
			},
			{
				Plmnid: plmnid,
				Ecid:   "0007",
			},
		},
	}

	// Sta5: 2, 4, 6, 8
	tStaLink5 := &nb.StationLinkInfo{
		Ecgi: &nb.ECGI{
			Plmnid: plmnid,
			Ecid:   "0005",
		},
		NeighborECGI: []*nb.ECGI{
			{
				Plmnid: plmnid,
				Ecid:   "0002",
			},
			{
				Plmnid: plmnid,
				Ecid:   "0004",
			},
			{
				Plmnid: plmnid,
				Ecid:   "0006",
			},
			{
				Plmnid: plmnid,
				Ecid:   "0008",
			},
		},
	}

	// Sta6: 3, 5, 9
	tStaLink6 := &nb.StationLinkInfo{
		Ecgi: &nb.ECGI{
			Plmnid: plmnid,
			Ecid:   "0006",
		},
		NeighborECGI: []*nb.ECGI{
			{
				Plmnid: plmnid,
				Ecid:   "0003",
			},
			{
				Plmnid: plmnid,
				Ecid:   "0005",
			},
			{
				Plmnid: plmnid,
				Ecid:   "0009",
			},
		},
	}

	// Sta7: 4, 8
	tStaLink7 := &nb.StationLinkInfo{
		Ecgi: &nb.ECGI{
			Plmnid: plmnid,
			Ecid:   "0007",
		},
		NeighborECGI: []*nb.ECGI{
			{
				Plmnid: plmnid,
				Ecid:   "0004",
			},
			{
				Plmnid: plmnid,
				Ecid:   "0008",
			},
		},
	}

	// Sta8: 5, 7, 9
	tStaLink8 := &nb.StationLinkInfo{
		Ecgi: &nb.ECGI{
			Plmnid: plmnid,
			Ecid:   "0008",
		},
		NeighborECGI: []*nb.ECGI{
			{
				Plmnid: plmnid,
				Ecid:   "0005",
			},
			{
				Plmnid: plmnid,
				Ecid:   "0007",
			},
			{
				Plmnid: plmnid,
				Ecid:   "0009",
			},
		},
	}

	// Sta9: 6, 8
	tStaLink9 := &nb.StationLinkInfo{
		Ecgi: &nb.ECGI{
			Plmnid: plmnid,
			Ecid:   "0009",
		},
		NeighborECGI: []*nb.ECGI{
			{
				Plmnid: plmnid,
				Ecid:   "0006",
			},
			{
				Plmnid: plmnid,
				Ecid:   "0008",
			},
		},
	}
	tStaLinks[tStaLink1.Ecgi.String()] = tStaLink1
	tStaLinks[tStaLink2.Ecgi.String()] = tStaLink2
	tStaLinks[tStaLink3.Ecgi.String()] = tStaLink3
	tStaLinks[tStaLink4.Ecgi.String()] = tStaLink4
	tStaLinks[tStaLink5.Ecgi.String()] = tStaLink5
	tStaLinks[tStaLink6.Ecgi.String()] = tStaLink6
	tStaLinks[tStaLink7.Ecgi.String()] = tStaLink7
	tStaLinks[tStaLink8.Ecgi.String()] = tStaLink8
	tStaLinks[tStaLink9.Ecgi.String()] = tStaLink9

	// UELinkInfo
	tUeInfo1 := &nb.UEInfo{
		Crnti: "00001",
		Ecgi: &nb.ECGI{
			Plmnid: plmnid,
			Ecid:   "0008",
		},
		Imsi: "315010999900001",
	}
	tUeInfo2 := &nb.UEInfo{
		Crnti: "00002",
		Ecgi: &nb.ECGI{
			Plmnid: plmnid,
			Ecid:   "0008",
		},
		Imsi: "315010999900002",
	}
	tUeInfo3 := &nb.UEInfo{
		Crnti: "00003",
		Ecgi: &nb.ECGI{
			Plmnid: plmnid,
			Ecid:   "0008",
		},
		Imsi: "315010999900003",
	}
	tUeInfo4 := &nb.UEInfo{
		Crnti: "00004",
		Ecgi: &nb.ECGI{
			Plmnid: plmnid,
			Ecid:   "0005",
		},
		Imsi: "315010999900004",
	}
	tUeInfo5 := &nb.UEInfo{
		Crnti: "00005",
		Ecgi: &nb.ECGI{
			Plmnid: plmnid,
			Ecid:   "0005",
		},
		Imsi: "315010999900005",
	}
	tUeInfo6 := &nb.UEInfo{
		Crnti: "00006",
		Ecgi: &nb.ECGI{
			Plmnid: plmnid,
			Ecid:   "0005",
		},
		Imsi: "315010999900006",
	}

	tUes[tUeInfo1.Imsi] = tUeInfo1
	tUes[tUeInfo2.Imsi] = tUeInfo2
	tUes[tUeInfo3.Imsi] = tUeInfo3
	tUes[tUeInfo4.Imsi] = tUeInfo4
	tUes[tUeInfo5.Imsi] = tUeInfo5
	tUes[tUeInfo6.Imsi] = tUeInfo6
	txPwrSta1 := nb.StationPowerOffset_PA_DB_0
	txPwrSta2 := nb.StationPowerOffset_PA_DB_0
	txPwrSta3 := nb.StationPowerOffset_PA_DB_0
	txPwrSta4 := nb.StationPowerOffset_PA_DB_0
	txPwrSta5 := nb.StationPowerOffset_PA_DB_0
	txPwrSta6 := nb.StationPowerOffset_PA_DB_0
	txPwrSta7 := nb.StationPowerOffset_PA_DB_0
	txPwrSta8 := nb.StationPowerOffset_PA_DB_0
	txPwrSta9 := nb.StationPowerOffset_PA_DB_0

	testResult, _ := MLBDecisionMaker(tStas, tStaLinks, tUes, &threshold)

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
