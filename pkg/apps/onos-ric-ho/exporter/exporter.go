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

package hoappexporter

import (
	"fmt"
	"net/http"
	"time"

	"github.com/onosproject/onos-lib-go/pkg/logging"

	hoappsouthbound "github.com/onosproject/onos-ric/pkg/apps/onos-ric-ho/southbound"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var log = logging.GetLogger("ho", "exporter")

// RunHOExposer runs Prometheus exporter
func RunHOExposer(sb *hoappsouthbound.HOSessions) {
	exposeHOInfo(sb)
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":7001", nil)
	if err != nil {
		log.Error(err)
	}
}

// exposeHOInfo is the function to expose all HO info - UELinkInfo, num(UELinKInfo), and num(HOEvents)
func exposeHOInfo(sb *hoappsouthbound.HOSessions) {
	go func() {
		for {
			listRNIB := exposeUELinkInfo(sb)
			numListRNIB := exposeNumUELinkInfo(sb)
			hoEventList := exposeHOEventListInfo(sb)
			numHOEventsCounter := exposeNumHOEvents(sb)

			time.Sleep(1000 * time.Millisecond)
			prometheus.Unregister(numListRNIB)
			for i := 0; i < len(listRNIB); i++ {
				prometheus.Unregister(listRNIB[i])
			}
			for i := 0; i < len(hoEventList); i++ {
				prometheus.Unregister(hoEventList[i])
			}
			prometheus.Unregister(numHOEventsCounter)
		}
	}()
}

// exposeUELinkInfo is the function to expose UELinkInfo
func exposeUELinkInfo(sb *hoappsouthbound.HOSessions) []prometheus.Counter {
	var listRNIB []prometheus.Counter
	for _, e := range sb.GetUELinkInfo() {
		tmp := promauto.NewCounter(prometheus.CounterOpts{
			Name: "hoapp_ue_link_info",
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

		listRNIB = append(listRNIB, tmp)
	}
	return listRNIB
}

// exposeNumUELinkInfo is the function to expose the number of UELinkInfo
func exposeNumUELinkInfo(sb *hoappsouthbound.HOSessions) prometheus.Counter {
	numListRNIB := promauto.NewCounter(prometheus.CounterOpts{
		Name: "hoapp_num_ue_link",
	})
	numListRNIB.Add(float64(len(sb.GetUELinkInfo())))

	return numListRNIB
}

// exposeNumHOEvents is the function to expose the number of HO events
func exposeNumHOEvents(sb *hoappsouthbound.HOSessions) prometheus.Counter {
	numHOEventsCounter := promauto.NewCounter(prometheus.CounterOpts{
		Name: "hoapp_num_ho_events",
	})
	numHOEventsCounter.Add(float64(sb.NumHOEvents))

	return numHOEventsCounter
}

func exposeHOEventListInfo(sb *hoappsouthbound.HOSessions) []prometheus.Counter {
	var listHOEvents []prometheus.Counter
	for _, e := range sb.HOEventStore {
		tmp := promauto.NewCounter(prometheus.CounterOpts{
			Name: "hoapp_ho_event_info",
			ConstLabels: prometheus.Labels{
				"timestamp": fmt.Sprintf("%d-%d-%d %d:%d:%d", e.TimeStamp.Year(), e.TimeStamp.Month(), e.TimeStamp.Day(), e.TimeStamp.Hour(), e.TimeStamp.Minute(), e.TimeStamp.Second()),
				"crnti":     e.CRNTI,
				"srcplmnid": e.SrcPlmnID,
				"srcecid":   e.SrcEcid,
				"dstplmnid": e.DstPlmnID,
				"dstecid":   e.DstEcid,
			},
		})
		tmp.Add(float64(e.ElapsedTime))
		listHOEvents = append(listHOEvents, tmp)
	}
	return listHOEvents
}
