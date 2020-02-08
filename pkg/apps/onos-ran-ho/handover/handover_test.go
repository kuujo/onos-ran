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
	"testing"

	"github.com/onosproject/onos-ran/api/nb"
	"github.com/stretchr/testify/assert"
)

// Define sample UELinkInfo messages

// Case 1. Serving station is the best one
func TestHODecisionMakerCase1(t *testing.T) {
	var cqV1 []*nb.ChannelQuality
	cqV1S := nb.ChannelQuality{
		TargetEcgi: &nb.ECGI{
			Plmnid: "315010",
			Ecid:   "0001",
		},
		CqiHist: 15,
	}
	cqV1 = append(cqV1, &cqV1S)
	cqV1N1 := nb.ChannelQuality{
		TargetEcgi: &nb.ECGI{
			Plmnid: "315010",
			Ecid:   "0002",
		},
		CqiHist: 12,
	}
	cqV1 = append(cqV1, &cqV1N1)
	cqV1N2 := nb.ChannelQuality{
		TargetEcgi: &nb.ECGI{
			Plmnid: "315010",
			Ecid:   "0003",
		},
		CqiHist: 11,
	}
	cqV1 = append(cqV1, &cqV1N2)
	cqV1N3 := nb.ChannelQuality{
		TargetEcgi: &nb.ECGI{
			Plmnid: "315010",
			Ecid:   "0004",
		},
		CqiHist: 1,
	}
	cqV1 = append(cqV1, &cqV1N3)
	value1 := nb.UELinkInfo{
		Crnti: "1",
		Ecgi: &nb.ECGI{
			Plmnid: "315010",
			Ecid:   "0001",
		},
		ChannelQualities: cqV1,
	}
	hoReqV1 := HODecisionMaker(&value1)
	assert.Nil(t, hoReqV1.GetDstStation()) // nil means "no need to trigger handover"
	assert.Nil(t, hoReqV1.GetSrcStation()) // nil means "no need to trigger handover"
}

// Case 2-1. 1st neighbor station is the best one
func TestHODecisionMakerCase2Dot1(t *testing.T) {

	v2TargetECID := "0002"
	var cqV2 []*nb.ChannelQuality
	cqV2S := nb.ChannelQuality{
		TargetEcgi: &nb.ECGI{
			Plmnid: "315010",
			Ecid:   "0001",
		},
		CqiHist: 10,
	}
	cqV2 = append(cqV2, &cqV2S)
	cqV2N1 := nb.ChannelQuality{
		TargetEcgi: &nb.ECGI{
			Plmnid: "315010",
			Ecid:   "0002",
		},
		CqiHist: 12,
	}
	cqV2 = append(cqV2, &cqV2N1)
	cqV2N2 := nb.ChannelQuality{
		TargetEcgi: &nb.ECGI{
			Plmnid: "315010",
			Ecid:   "0003",
		},
		CqiHist: 11,
	}
	cqV2 = append(cqV2, &cqV2N2)
	cqV2N3 := nb.ChannelQuality{
		TargetEcgi: &nb.ECGI{
			Plmnid: "315010",
			Ecid:   "0004",
		},
		CqiHist: 1,
	}
	cqV2 = append(cqV2, &cqV2N3)
	value2 := nb.UELinkInfo{
		Crnti: "1",
		Ecgi: &nb.ECGI{
			Plmnid: "315010",
			Ecid:   "0001",
		},
		ChannelQualities: cqV2,
	}
	hoReqV2 := HODecisionMaker(&value2)
	assert.Equal(t, hoReqV2.GetDstStation().GetEcid(), v2TargetECID)
}

// Case 2-2. 2nd neighbor station is the best one
func TestHODecisionMakerCase2Dot2(t *testing.T) {

	v2TargetECID := "0003"
	var cqV2 []*nb.ChannelQuality
	cqV2S := nb.ChannelQuality{
		TargetEcgi: &nb.ECGI{
			Plmnid: "315010",
			Ecid:   "0001",
		},
		CqiHist: 10,
	}
	cqV2 = append(cqV2, &cqV2S)
	cqV2N1 := nb.ChannelQuality{
		TargetEcgi: &nb.ECGI{
			Plmnid: "315010",
			Ecid:   "0002",
		},
		CqiHist: 9,
	}
	cqV2 = append(cqV2, &cqV2N1)
	cqV2N2 := nb.ChannelQuality{
		TargetEcgi: &nb.ECGI{
			Plmnid: "315010",
			Ecid:   "0003",
		},
		CqiHist: 11,
	}
	cqV2 = append(cqV2, &cqV2N2)
	cqV2N3 := nb.ChannelQuality{
		TargetEcgi: &nb.ECGI{
			Plmnid: "315010",
			Ecid:   "0004",
		},
		CqiHist: 1,
	}
	cqV2 = append(cqV2, &cqV2N3)
	value2 := nb.UELinkInfo{
		Crnti: "1",
		Ecgi: &nb.ECGI{
			Plmnid: "315010",
			Ecid:   "0001",
		},
		ChannelQualities: cqV2,
	}
	hoReqV2 := HODecisionMaker(&value2)
	assert.Equal(t, hoReqV2.GetDstStation().GetEcid(), v2TargetECID)
}

// Case 2-3. 3rd neighbor station is the best one
func TestHODecisionMakerCase2Dot3(t *testing.T) {

	v2TargetECID := "0004"
	var cqV2 []*nb.ChannelQuality
	cqV2S := nb.ChannelQuality{
		TargetEcgi: &nb.ECGI{
			Plmnid: "315010",
			Ecid:   "0001",
		},
		CqiHist: 10,
	}
	cqV2 = append(cqV2, &cqV2S)
	cqV2N1 := nb.ChannelQuality{
		TargetEcgi: &nb.ECGI{
			Plmnid: "315010",
			Ecid:   "0002",
		},
		CqiHist: 9,
	}
	cqV2 = append(cqV2, &cqV2N1)
	cqV2N2 := nb.ChannelQuality{
		TargetEcgi: &nb.ECGI{
			Plmnid: "315010",
			Ecid:   "0003",
		},
		CqiHist: 11,
	}
	cqV2 = append(cqV2, &cqV2N2)
	cqV2N3 := nb.ChannelQuality{
		TargetEcgi: &nb.ECGI{
			Plmnid: "315010",
			Ecid:   "0004",
		},
		CqiHist: 15,
	}
	cqV2 = append(cqV2, &cqV2N3)
	value2 := nb.UELinkInfo{
		Crnti: "1",
		Ecgi: &nb.ECGI{
			Plmnid: "315010",
			Ecid:   "0001",
		},
		ChannelQualities: cqV2,
	}
	hoReqV2 := HODecisionMaker(&value2)
	assert.Equal(t, hoReqV2.GetDstStation().GetEcid(), v2TargetECID)
}
