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
	"net/http"
	"sync"
	"time"

	"github.com/onosproject/onos-lib-go/pkg/logging"

	"github.com/onosproject/onos-ric/api/sb"
	"github.com/onosproject/onos-ric/pkg/manager"
	"github.com/onosproject/onos-ric/pkg/southbound"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var log = logging.GetLogger("exporter")

var hoLatencyHistogram prometheus.Histogram

// RunRICExposer runs Prometheus exposer
func RunRICExposer(mgr *manager.Manager) {
	southbound.ChanHOEvent = make(chan southbound.HOEventMeasuredRIC)
	initHOHistogram(mgr)
	exposeCtrUpdateInfo(mgr)
	exposeHOLatency()
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
			wg.Add(2) // because there are four goroutines

			var listRNIBCell []prometheus.Counter
			var listRNIBUEAdmReq []prometheus.Counter

			go func() {
				listRNIBCell = exposeCellConfig(mgr)
				defer wg.Done()
			}()
			go func() {
				listRNIBUEAdmReq = exposeUEAdmRequests(mgr)
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

func exposeHOLatency() {
	go func() {
		for {
			for e := range southbound.ChanHOEvent {
				tmp := promauto.NewCounter(prometheus.CounterOpts{
					Name: "ho_events",
					ConstLabels: prometheus.Labels{
						"timestamp": e.Timestamp.Format(time.StampNano),
						"crnti":     e.Crnti,
						"plmnid":    e.DestPlmnID,
						"ecid":      e.DestECID,
					},
				})
				eTime := float64(e.ElapsedTime)
				tmp.Add(eTime)
				hoLatencyHistogram.Observe(eTime)
			}
		}
	}()
}

func initHOHistogram(mgr *manager.Manager) {
	hoLatencyHistogram = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "ho_histogram",
			Buckets: prometheus.ExponentialBuckets(1e3, 1.5, 20),
		},
	)
	prometheus.MustRegister(hoLatencyHistogram)
}
