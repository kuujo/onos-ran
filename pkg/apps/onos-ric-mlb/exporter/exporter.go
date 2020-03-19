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
	go func() {
		for {
			var wg sync.WaitGroup
			wg.Add(3)

			var listUELinkInfo []prometheus.Counter
			var listCellInfo []prometheus.Counter
			var listMLBEventInfo []prometheus.Counter

			go func() {
				listUELinkInfo = exposeUELinkInfo(sb)
				defer wg.Done()
			}()
			go func() {
				listCellInfo = exposeCellInfo(sb)
				defer wg.Done()
			}()
			go func() {
				listMLBEventInfo = exposeMLBEventInfo(sb)
				defer wg.Done()
			}()
			wg.Wait()

			time.Sleep(1000 * time.Millisecond)
			for i := 0; i < len(listUELinkInfo); i++ {
				prometheus.Unregister(listUELinkInfo[i])
			}
			for i := 0; i < len(listCellInfo); i++ {
				prometheus.Unregister(listCellInfo[i])
			}
			for i := 0; i < len(listMLBEventInfo); i++ {
				prometheus.Unregister(listMLBEventInfo[i])
			}
		}
	}()
}

// exposeUELinkInfo exposes UE link information
func exposeUELinkInfo(sb *mlbappsouthbound.MLBSessions) []prometheus.Counter {
	var listRNIBUELink []prometheus.Counter
	for _, e := range sb.GetUELinkInfo() {
		tmp := promauto.NewCounter(prometheus.CounterOpts{
			Name: "mlbapp_ue_link_info",
			ConstLabels: prometheus.Labels{
				"crnti":       e.GetCrnti(),
				"serv_plmnid": e.GetEcgi().GetPlmnid(),
				"serv_ecid":   e.GetEcgi().GetEcid(),
				"n1_plmnid":   e.GetChannelQualities()[0].GetTargetEcgi().GetPlmnid(),
				"n1_ecid":     e.GetChannelQualities()[0].GetTargetEcgi().GetEcid(),
				"n1_cqi":      fmt.Sprintf("%d", e.GetChannelQualities()[0].GetCqiHist()),
				"n2_plmnid":   e.GetChannelQualities()[1].GetTargetEcgi().GetPlmnid(),
				"n2_ecid":     e.GetChannelQualities()[1].GetTargetEcgi().GetEcid(),
				"n2_cqi":      fmt.Sprintf("%d", e.GetChannelQualities()[1].GetCqiHist()),
				"n3_plmnid":   e.GetChannelQualities()[2].GetTargetEcgi().GetPlmnid(),
				"n3_ecid":     e.GetChannelQualities()[2].GetTargetEcgi().GetEcid(),
				"n3_cqi":      fmt.Sprintf("%d", e.GetChannelQualities()[2].GetCqiHist()),
			},
		})
		listRNIBUELink = append(listRNIBUELink, tmp)
	}
	return listRNIBUELink
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

// exposeMLBEventInfo exposes MLB events
func exposeMLBEventInfo(sb *mlbappsouthbound.MLBSessions) []prometheus.Counter {
	var listMLBEvents []prometheus.Counter
	for _, e := range sb.MLBEventStore {
		tmp := promauto.NewCounter(prometheus.CounterOpts{
			Name: "mlbapp_mlb_event_info",
			ConstLabels: prometheus.Labels{
				"timestamp": fmt.Sprintf("%d-%d-%d %d:%d:%d.%d", e.TimeStamp.Year(), e.TimeStamp.Month(), e.TimeStamp.Day(), e.TimeStamp.Hour(), e.TimeStamp.Minute(), e.TimeStamp.Second(), e.TimeStamp.Nanosecond()),
				"plmnid":    e.PlmnID,
				"ecid":      e.Ecid,
				"maxnumues": fmt.Sprintf("%d", e.MaxNumUes),
				"numues":    fmt.Sprintf("%d", e.NumUes),
				"pa":        fmt.Sprintf("%d", e.Pa),
			},
		})
		tmp.Add(float64(e.ElapsedTime))
		listMLBEvents = append(listMLBEvents, tmp)
	}
	return listMLBEvents
}
