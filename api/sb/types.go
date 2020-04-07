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

package sb

import (
	"fmt"
)

// Endpoint - a type for K8s endpoints
type Endpoint string

// GetID returns a unique identifier for the message
func (m *TelemetryMessage) GetID() string {
	var ecgi ECGI
	var crnti string
	switch m.MessageType {
	case MessageType_RADIO_MEAS_REPORT_PER_UE:
		ecgi = *m.GetRadioMeasReportPerUE().GetEcgi()
		crnti = m.GetRadioMeasReportPerUE().GetCrnti()
	case MessageType_RADIO_MEAS_REPORT_PER_CELL:
		ecgi = *m.GetRadioMeasReportPerCell().GetEcgi()
	case MessageType_SCHED_MEAS_REPORT_PER_UE:
	case MessageType_SCHED_MEAS_REPORT_PER_CELL:
		ecgi = *m.GetSchedMeasReportPerCell().GetEcgi()
	}
	return NewID(m.GetMessageType(), ecgi.PlmnId, ecgi.Ecid, crnti)
}

// GetID returns a unique identifier for the message
func (m *ControlUpdate) GetID() string {
	var ecgi ECGI
	var crnti string
	switch m.MessageType {
	case MessageType_RRM_CONFIG_STATUS:
		ecgi = *m.GetRRMConfigStatus().GetEcgi()
	case MessageType_UE_ADMISSION_REQUEST:
		ecgi = *m.GetUEAdmissionRequest().GetEcgi()
		crnti = m.GetUEAdmissionRequest().GetCrnti()
	case MessageType_UE_ADMISSION_STATUS:
		ecgi = *m.GetUEAdmissionStatus().GetEcgi()
		crnti = m.GetUEAdmissionStatus().GetCrnti()
	case MessageType_UE_CONTEXT_UPDATE:
		ecgi = *m.GetUEContextUpdate().GetEcgi()
		crnti = m.GetUEContextUpdate().GetCrnti()
	case MessageType_BEARER_ADMISSION_REQUEST:
		ecgi = *m.GetBearerAdmissionRequest().GetEcgi()
		crnti = m.GetBearerAdmissionStatus().GetCrnti()
	case MessageType_BEARER_ADMISSION_STATUS:
		ecgi = *m.GetBearerAdmissionStatus().GetEcgi()
		crnti = m.GetBearerAdmissionStatus().GetCrnti()
	}
	return NewID(m.GetMessageType(), ecgi.PlmnId, ecgi.Ecid, crnti)
}

// NewID creates a new message identifier
func NewID(messageType MessageType, plmnid, ecid, crnti string) string {
	return fmt.Sprintf("%s:%s:%s:%s", messageType, plmnid, ecid, crnti)
}
