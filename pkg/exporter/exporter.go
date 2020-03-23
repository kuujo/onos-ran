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

package exporter

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/onosproject/onos-lib-go/pkg/logging"

	"github.com/onosproject/onos-ric/api/sb"
	"github.com/onosproject/onos-ric/pkg/manager"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var log = logging.GetLogger("exporter")

// RunRICExposer runs Prometheus exposer
func RunRICExposer(mgr *manager.Manager) {
	exposeCtrUpdateInfo(mgr)
	exposeTelemetryInfo(mgr)
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":7000", nil)
	if err != nil {
		log.Error(err)
	}
}

// exposeCtrUpdateInfo exposes Control Update info
func exposeCtrUpdateInfo(mgr *manager.Manager) {
	go func() {
		for {
			var wg sync.WaitGroup
			wg.Add(3) // because there are four goroutines

			var listRNIBCell []prometheus.Counter
			var listRNIBUEAdmReq []prometheus.Counter
			var listRNIBUERel []prometheus.Counter
			var listHOEvents []prometheus.Counter

			go func() {
				listRNIBCell = exposeCellConfig(mgr)
				defer wg.Done()
			}()
			go func() {
				listRNIBUEAdmReq = exposeUEAdmRequests(mgr)
				defer wg.Done()
			}()
			go func() {
				listRNIBUERel = exposeUEReleaseInd(mgr)
				defer wg.Done()
			}()
			wg.Wait()

			time.Sleep(1000 * time.Millisecond)
			for i := 0; i < len(listRNIBCell); i++ {
				prometheus.Unregister(listRNIBCell[i])
			}
			for i := 0; i < len(listRNIBUEAdmReq); i++ {
				prometheus.Unregister(listRNIBUEAdmReq[i])
			}
			for i := 0; i < len(listRNIBUERel); i++ {
				prometheus.Unregister(listRNIBUERel[i])
			}
			for i := 0; i < len(listHOEvents); i++ {
				prometheus.Unregister(listHOEvents[i])
			}
		}
	}()
}

// exposeCellConfig exposes CellConfig info
func exposeCellConfig(mgr *manager.Manager) []prometheus.Counter {
	var listCellConfigMsgs []prometheus.Counter

	allControlUpdates, _ := mgr.GetControlUpdates()
	for i := 0; i < len(allControlUpdates); i++ {
		switch allControlUpdates[i].S.(type) {
		case *sb.ControlUpdate_CellConfigReport:
			tmp := promauto.NewCounter(prometheus.CounterOpts{
				Name: "cell_config_info",
				ConstLabels: prometheus.Labels{
					"plmnid":                  allControlUpdates[i].GetCellConfigReport().GetEcgi().GetPlmnId(),
					"ecid":                    allControlUpdates[i].GetCellConfigReport().GetEcgi().GetEcid(),
					"pci":                     string(allControlUpdates[i].GetCellConfigReport().GetPci()),
					"earfcndl":                allControlUpdates[i].GetCellConfigReport().GetEarfcnDl(),
					"earfcnul":                allControlUpdates[i].GetCellConfigReport().GetEarfcnUl(),
					"rbsperttidl":             string(allControlUpdates[i].GetCellConfigReport().GetRbsPerTtiDl()),
					"rbsperttiul":             string(allControlUpdates[i].GetCellConfigReport().GetRbsPerTtiUl()),
					"numtxantenna":            string(allControlUpdates[i].GetCellConfigReport().GetNumTxAntenna()),
					"duplexmode":              allControlUpdates[i].GetCellConfigReport().GetDuplexMode(),
					"maxnumconnectedues":      string(allControlUpdates[i].GetCellConfigReport().GetMaxNumConnectedUes()),
					"maxnumconnectedbearers":  string(allControlUpdates[i].GetCellConfigReport().GetMaxNumConnectedBearers()),
					"maxnumuesschedperrttidl": string(allControlUpdates[i].GetCellConfigReport().GetMaxNumUesSchedPerTtiDl()),
					"maxnumuesschedperttiul":  string(allControlUpdates[i].GetCellConfigReport().GetMaxNumUesSchedPerTtiUl()),
					"dlfsschedenable":         allControlUpdates[i].GetCellConfigReport().GetDlfsSchedEnable(),
				},
			})
			listCellConfigMsgs = append(listCellConfigMsgs, tmp)
		default:
		}
	}
	return listCellConfigMsgs
}

// exposeUEAdmRequests exposes UEAdmissionRequest info
func exposeUEAdmRequests(mgr *manager.Manager) []prometheus.Counter {
	var listUEAdmRequestMsgs []prometheus.Counter
	allControlUpdates, _ := mgr.GetControlUpdates()
	for i := 0; i < len(allControlUpdates); i++ {
		switch allControlUpdates[i].S.(type) {
		case *sb.ControlUpdate_UEAdmissionRequest:
			tmp := promauto.NewCounter(prometheus.CounterOpts{
				Name: "ue_adm_req_info",
				ConstLabels: prometheus.Labels{
					"crnti":  allControlUpdates[i].GetUEAdmissionRequest().GetCrnti(),
					"plmnid": allControlUpdates[i].GetUEAdmissionRequest().Ecgi.GetPlmnId(),
					"ecid":   allControlUpdates[i].GetUEAdmissionRequest().GetEcgi().GetEcid(),
				},
			})
			listUEAdmRequestMsgs = append(listUEAdmRequestMsgs, tmp)
		default:
		}
	}
	return listUEAdmRequestMsgs
}

// exposeUEReleaseInd exposes UEReleaseInd info
func exposeUEReleaseInd(mgr *manager.Manager) []prometheus.Counter {
	var listUEReleaseIndMsgs []prometheus.Counter
	allControlUpdates, _ := mgr.GetControlUpdates()
	for i := 0; i < len(allControlUpdates); i++ {
		switch allControlUpdates[i].S.(type) {
		case *sb.ControlUpdate_UEReleaseInd:
			tmp := promauto.NewCounter(prometheus.CounterOpts{
				Name: "ue_rel_info",
				ConstLabels: prometheus.Labels{
					"crnti":  allControlUpdates[i].GetUEReleaseInd().GetCrnti(),
					"plmnid": allControlUpdates[i].GetUEReleaseInd().Ecgi.GetPlmnId(),
					"ecid":   allControlUpdates[i].GetUEReleaseInd().GetEcgi().GetEcid(),
				},
			})
			listUEReleaseIndMsgs = append(listUEReleaseIndMsgs, tmp)
		default:
		}
	}
	return listUEReleaseIndMsgs
}

// exposeTelemtryInfo is the function to expose Telemetry
// Currently, Telemetry is singular, but should be extended
func exposeTelemetryInfo(mgr *manager.Manager) {
	go func() {
		for {
			listRNIBUELink := exposeUELinkInfo(mgr)

			time.Sleep(1000 * time.Millisecond)
			for i := 0; i < len(listRNIBUELink); i++ {
				prometheus.Unregister(listRNIBUELink[i])
			}
		}
	}()
}

// exposeUELinkInfo exposes UELink info
func exposeUELinkInfo(mgr *manager.Manager) []prometheus.Counter {
	var listUELinkInfo []prometheus.Counter
	allTelemetry, _ := mgr.GetTelemetry()
	for i := 0; i < len(allTelemetry); i++ {
		switch allTelemetry[i].S.(type) {
		case *sb.TelemetryMessage_RadioMeasReportPerUE:
			tmp := promauto.NewCounter(prometheus.CounterOpts{
				Name: "ue_link_info",
				ConstLabels: prometheus.Labels{
					"crnti":       allTelemetry[i].GetRadioMeasReportPerUE().GetCrnti(),
					"serv_plmnid": allTelemetry[i].GetRadioMeasReportPerUE().GetEcgi().GetPlmnId(),
					"serv_ecid":   allTelemetry[i].GetRadioMeasReportPerUE().GetEcgi().GetEcid(),
					"n1_plmnid":   allTelemetry[i].GetRadioMeasReportPerUE().GetRadioReportServCells()[0].GetEcgi().GetPlmnId(),
					"n1_ecid":     allTelemetry[i].GetRadioMeasReportPerUE().GetRadioReportServCells()[0].GetEcgi().GetEcid(),
					"n1_cqi":      fmt.Sprintf("%d", allTelemetry[i].GetRadioMeasReportPerUE().GetRadioReportServCells()[0].GetCqiHist()[0]),
					"n2_plmnid":   allTelemetry[i].GetRadioMeasReportPerUE().GetRadioReportServCells()[1].GetEcgi().GetPlmnId(),
					"n2_ecid":     allTelemetry[i].GetRadioMeasReportPerUE().GetRadioReportServCells()[1].GetEcgi().GetEcid(),
					"n2_cqi":      fmt.Sprintf("%d", allTelemetry[i].GetRadioMeasReportPerUE().GetRadioReportServCells()[1].GetCqiHist()[0]),
					"n3_plmnid":   allTelemetry[i].GetRadioMeasReportPerUE().GetRadioReportServCells()[2].GetEcgi().GetPlmnId(),
					"n3_ecid":     allTelemetry[i].GetRadioMeasReportPerUE().GetRadioReportServCells()[2].GetEcgi().GetEcid(),
					"n3_cqi":      fmt.Sprintf("%d", allTelemetry[i].GetRadioMeasReportPerUE().GetRadioReportServCells()[2].GetCqiHist()[0]),
				},
			})
			listUELinkInfo = append(listUELinkInfo, tmp)
		default:
		}
	}
	return listUELinkInfo
}
