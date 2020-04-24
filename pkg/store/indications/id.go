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

package indications

import (
	"github.com/onosproject/onos-ric/api/sb"
	"github.com/onosproject/onos-ric/api/sb/e2ap"
	"github.com/onosproject/onos-ric/pkg/store/message"
)

// NewID creates a new telemetry store ID
func NewID(messageType sb.MessageType, plmnidn, ecid, crnti string) ID {
	return ID{
		Partition: PartitionKey(message.NewPartitionKey(plmnidn, ecid)),
		Key:       Key(message.NewID(messageType, plmnidn, ecid, crnti)),
	}
}

// GetID gets the telemetry store ID for the given message
func GetID(message *e2ap.RicIndication) ID {
	var ecgi sb.ECGI
	var crnti string
	msgType := message.GetHdr().GetMessageType()
	switch msgType {
	case sb.MessageType_RADIO_MEAS_REPORT_PER_UE:
		ecgi = *message.GetMsg().GetRadioMeasReportPerUE().GetEcgi()
		crnti = message.GetMsg().GetRadioMeasReportPerUE().GetCrnti()
	case sb.MessageType_RADIO_MEAS_REPORT_PER_CELL:
		ecgi = *message.GetMsg().GetRadioMeasReportPerCell().GetEcgi()
	case sb.MessageType_CELL_CONFIG_REPORT:
		ecgi = *message.GetMsg().GetCellConfigReport().GetEcgi()
	case sb.MessageType_UE_ADMISSION_REQUEST:
		ecgi = *message.GetMsg().GetUEAdmissionRequest().GetEcgi()
		crnti = message.GetMsg().GetUEAdmissionRequest().GetCrnti()
	}
	return NewID(msgType, ecgi.PlmnId, ecgi.Ecid, crnti)
}