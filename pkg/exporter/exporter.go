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
	"math"
	"net/http"
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
	initHOHistogram()
	//exposeCtrUpdateInfo(mgr)
	exposeHOLatency()
	exposeStdevLoad()
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":7000", nil)
	if err != nil {
		log.Error(err)
	}
}

func exposeHOLatency() {
	go func() {
		for {
			time.Sleep(1 * time.Second)
			for e := range southbound.ChanHOEvent {
				hoLatencyHistogram.Observe(float64(e.ElapsedTime))
			}
		}
	}()
}

func initHOHistogram() {
	hoLatencyHistogram = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "ho_histogram",
			Buckets: prometheus.ExponentialBuckets(1e3, 1.5, 20),
		},
	)
	prometheus.MustRegister(hoLatencyHistogram)
}

func exposeStdevLoad() {
	totalNumUEsGuage := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "total_num_ues",
	})
	totalNumBSGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "total_num_bses",
	})
	stdDevGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "stdev_ues",
	})
	go func() {
		for {
			time.Sleep(1 * time.Second)

			tNow := time.Now()
			allIndication, _ := manager.GetManager().GetIndications()

			tmpMap := make(map[sb.ECGI]int)
			numTotalUEs := 0
			numBSes := 0

			for i := 0; i < len(allIndication); i++ {
				switch allIndication[i].GetHdr().GetMessageType() {
				case sb.MessageType_CELL_CONFIG_REPORT:
					tmpMap[*allIndication[i].GetMsg().GetCellConfigReport().GetEcgi()] = 0
					numBSes++
				case sb.MessageType_UE_ADMISSION_REQUEST:
					numTotalUEs++
					if _, ok := tmpMap[*allIndication[i].GetMsg().GetUEAdmissionRequest().GetEcgi()]; ok {
						tmpMap[*allIndication[i].GetMsg().GetUEAdmissionRequest().GetEcgi()]++
					}
				default:
				}
			}

			avgNumUEs := float64(numTotalUEs) / float64(numBSes)
			var stdev float64

			for _, v := range tmpMap {
				stdev += math.Pow(float64(v)-avgNumUEs, 2)
			}

			stdev = math.Sqrt(stdev / float64(numBSes))

			// expose Total Num UEs
			go func() {
				historyTotalNumUEsCounter := promauto.NewCounter(prometheus.CounterOpts{
					Name: "history_total_num_ues",
					ConstLabels: prometheus.Labels{
						"timestamp":      tNow.Format(time.StampNano),
						"timestamp_nano": fmt.Sprintf("%d", tNow.UnixNano()),
					},
				})
				totalNumUEsGuage.Set(float64(numTotalUEs))
				historyTotalNumUEsCounter.Add(float64(numTotalUEs))
			}()

			// expose Total Num BSes
			go func() {
				historyTotalNumBSCounter := promauto.NewCounter(prometheus.CounterOpts{
					Name: "history_total_num_bses",
					ConstLabels: prometheus.Labels{
						"timestamp":      tNow.Format(time.StampNano),
						"timestamp_nano": fmt.Sprintf("%d", tNow.UnixNano()),
					},
				})
				totalNumBSGauge.Set(float64(numBSes))
				historyTotalNumBSCounter.Add(float64(numBSes))
			}()

			// expose Stdev(Connected UEs for each BS)
			go func() {
				historyStdevCounter := promauto.NewCounter(prometheus.CounterOpts{
					Name: "history_stdev_ues",
					ConstLabels: prometheus.Labels{
						"timestamp":      tNow.Format(time.StampNano),
						"timestamp_nano": fmt.Sprintf("%d", tNow.UnixNano()),
					},
				})
				stdDevGauge.Set(stdev)
				historyStdevCounter.Add(stdev)
			}()
		}
	}()
}
