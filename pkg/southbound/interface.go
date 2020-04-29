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

package southbound

import (
	"time"

	topodevice "github.com/onosproject/onos-topo/api/device"

	"github.com/onosproject/onos-ric/api/sb"
	e2ap "github.com/onosproject/onos-ric/api/sb/e2ap"
)

//ChanHOEvent ...
var ChanHOEvent chan HOEventMeasuredRIC

// HOEventMeasuredRIC is struct including UE ID and its eNB ID
type HOEventMeasuredRIC struct {
	Timestamp   time.Time
	Crnti       string
	DestPlmnID  string
	DestECID    string
	ElapsedTime int64
}

// TelemetryUpdateHandler - a function for the session to write back to manager without the import cycle
type TelemetryUpdateHandler = func(e2ap.RicIndication)

// RicControlResponseHandler - ricIndication messages to manager without the import cycle
type RicControlResponseHandler = func(e2ap.RicIndication)

// ControlUpdateHandler - a function for the session to write back to manager without the import cycle
type ControlUpdateHandler = func(e2ap.RicIndication)

// E2 ...
type E2 interface {
	Run(ecgi sb.ECGI, endPoint sb.Endpoint,
		tls topodevice.TlsConfig, creds topodevice.Credentials,
		storeRicControlResponse RicControlResponseHandler,
		storeControlUpdate ControlUpdateHandler,
		storeTelemetry TelemetryUpdateHandler,
		enableMetrics bool,
		e2Chan chan E2)
	RemoteAddress() sb.Endpoint
	Setup() error
	RRMConfig(pa sb.XICICPA) error
	UeHandover(crnti []string, srcEcgi sb.ECGI, dstEcgi sb.ECGI) error
	L2MeasConfig(l2ReportInterval *sb.L2MeasConfig) error
}
