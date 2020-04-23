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
	"net/http"
	"time"

	"github.com/onosproject/onos-lib-go/pkg/logging"
	mlbappsouthbound "github.com/onosproject/onos-ric/pkg/apps/onos-ric-mlb/southbound"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var log = logging.GetLogger("mlb", "exporter")

var mlbLatencyHistogram prometheus.Histogram

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
	mlbappsouthbound.ChanMLBEvent = make(chan mlbappsouthbound.MLBEvent)
	initMLBHistogram()
	exposeMLBEvents(sb)
}

func exposeMLBEvents(sb *mlbappsouthbound.MLBSessions) {
	go func() {
		for {
			time.Sleep(1 * time.Second)
			for e := range mlbappsouthbound.ChanMLBEvent {
				mlbLatencyHistogram.Observe(float64(e.EndTime.Sub(e.StartTime).Microseconds()))
			}
		}
	}()
}

func initMLBHistogram() {
	mlbLatencyHistogram = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "mlbapp_mlb_histogram",
			Buckets: prometheus.ExponentialBuckets(1e3, 1.5, 20),
		},
	)
	prometheus.MustRegister(mlbLatencyHistogram)
}
