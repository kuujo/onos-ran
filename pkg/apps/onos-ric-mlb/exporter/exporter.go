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

package mlbappexporter

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/onosproject/onos-lib-go/pkg/logging"
	mlbappsouthbound "github.com/onosproject/onos-ric/pkg/apps/onos-ric-mlb/southbound"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var log = logging.GetLogger("mlb", "exporter")

// RunMLBExposer runs Prometheus exporter
func RunMLBExposer(sb *mlbappsouthbound.MLBSessions) {
	exposeMLBInfo(sb)
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":7002", nil)
	if err != nil {
		log.Error(err)
	}
}

// exposeMLBInfo exposes MLB Information - UELink info, Cell info, and MLB events
func exposeMLBInfo(sb *mlbappsouthbound.MLBSessions) {
	go exposeMLBEvents(sb)
	go func() {
		for {
			var wg sync.WaitGroup
			wg.Add(2)

			var listUEInfo []prometheus.Counter
			var listCellInfo []prometheus.Counter

			go func() {
				listCellInfo = exposeCellInfo(sb)
				defer wg.Done()
			}()
			go func() {
				listUEInfo = exposeUEInfo(sb)
				defer wg.Done()
			}()
			wg.Wait()

			time.Sleep(1000 * time.Millisecond)
			for i := 0; i < len(listUEInfo); i++ {
				prometheus.Unregister(listUEInfo[i])
			}
			for i := 0; i < len(listCellInfo); i++ {
				prometheus.Unregister(listCellInfo[i])
			}
		}
	}()
}

// exposeCellInfo exposes cell infomration
func exposeCellInfo(sb *mlbappsouthbound.MLBSessions) []prometheus.Counter {
	var listRNIBCell []prometheus.Counter
	for _, e := range sb.RNIBCellInfo {
		tmp := promauto.NewCounter(prometheus.CounterOpts{
			Name: "mlbapp_cell_config_info",
			ConstLabels: prometheus.Labels{
				"plmnid":       e.GetEcgi().GetPlmnid(),
				"ecid":         e.GetEcgi().GetEcid(),
				"maxnumconues": fmt.Sprintf("%d", e.GetMaxNumConnectedUes()),
			},
		})
		listRNIBCell = append(listRNIBCell, tmp)
	}
	return listRNIBCell
}

func exposeUEInfo(sb *mlbappsouthbound.MLBSessions) []prometheus.Counter {
	var listUEInfo []prometheus.Counter
	for _, e := range sb.UEInfoList {
		tmp := promauto.NewCounter(prometheus.CounterOpts{
			Name: "mlbapp_ue_info",
			ConstLabels: prometheus.Labels{
				"imsi":   e.GetImsi(),
				"crnti":  e.GetCrnti(),
				"plmnid": e.GetEcgi().GetPlmnid(),
				"ecid":   e.GetEcgi().GetEcid(),
			},
		})
		listUEInfo = append(listUEInfo, tmp)
	}
	return listUEInfo
}

func exposeMLBEvents(sb *mlbappsouthbound.MLBSessions) {
	for {
		if sb.MLBEventChan == nil {
			continue
		}
		for e := range sb.MLBEventChan {
			tmp := promauto.NewCounter(prometheus.CounterOpts{
				Name: "mlbapp_mlb_events",
				ConstLabels: prometheus.Labels{
					"plmnid":    e.MLBReq.GetEcgi().GetPlmnid(),
					"ecgi":      e.MLBReq.GetEcgi().GetEcid(),
					"pa":        e.MLBReq.GetOffset().String(),
					"timestamp": e.EndTime.Format(time.StampNano),
				},
			})
			tmp.Add(float64(e.EndTime.Sub(e.StartTime).Microseconds()))
		}
	}
}
