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

	hoappsouthbound "github.com/onosproject/onos-ric/pkg/apps/onos-ric-ho/southbound"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "k8s.io/klog"
)

// RunHOExposer runs Prometheus exposer
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
			numHOEventsCounter := exposeNumHOEvents(sb)

			time.Sleep(1000 * time.Millisecond)
			prometheus.Unregister(numListRNIB)
			for i := 0; i < len(listRNIB); i++ {
				prometheus.Unregister(listRNIB[i])
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
			Name: "ue_link_info",
			ConstLabels: prometheus.Labels{
				"time":        time.Now().String(),
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
		err := prometheus.Register(tmp)
		if err != nil {
			log.Error(err)
		}
		listRNIB = append(listRNIB, tmp)
	}
	return listRNIB
}

// exposeNumUELinkInfo is the function to expose the number of UELinkInfo
func exposeNumUELinkInfo(sb *hoappsouthbound.HOSessions) prometheus.Counter {
	numListRNIB := promauto.NewCounter(prometheus.CounterOpts{
		Name: "num_ue_link",
	})
	numListRNIB.Add(float64(len(sb.GetUELinkInfo())))
	err := prometheus.Register(numListRNIB)
	if err != nil {
		log.Error(err)
	}
	return numListRNIB
}

// exposeNumHOEvents is the function to expose the number of HO events
func exposeNumHOEvents(sb *hoappsouthbound.HOSessions) prometheus.Counter {
	numHOEventsCounter := promauto.NewCounter(prometheus.CounterOpts{
		Name: "num_ho_events",
	})
	numHOEventsCounter.Add(float64(sb.NumHOEvents))
	err := prometheus.Register(numHOEventsCounter)
	if err != nil {
		log.Error(err)
	}
	return numHOEventsCounter
}
