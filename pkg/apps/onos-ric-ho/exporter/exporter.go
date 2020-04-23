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
	"github.com/onosproject/onos-lib-go/pkg/logging"
	hoappsouthbound "github.com/onosproject/onos-ric/pkg/apps/onos-ric-ho/southbound"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

var log = logging.GetLogger("ho", "exporter")

var hoLatencyHistogram prometheus.Histogram

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
	hoappsouthbound.ChanHOEvent = make(chan hoappsouthbound.HOEvent)
	initHOHistogram()
	exposeHOLatency(sb)
}

// exposeUELinkInfo is the function to expose UELinkInfo
func exposeHOLatency(sb *hoappsouthbound.HOSessions) {
	go func() {
		for {
			time.Sleep(1 * time.Second)
			for e := range hoappsouthbound.ChanHOEvent {
				hoLatencyHistogram.Observe(float64(e.ElapsedTime))
			}
		}
	}()
}

func initHOHistogram() {
	hoLatencyHistogram = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "hoapp_ho_histogram",
			Buckets: prometheus.ExponentialBuckets(1e3, 1.5, 20),
		},
	)
	prometheus.MustRegister(hoLatencyHistogram)
}
