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

package main

import (
	"flag"

	mlbappexporter "github.com/onosproject/onos-ric/pkg/apps/onos-ric-mlb/exporter"

	"github.com/onosproject/onos-lib-go/pkg/logging"

	mlbappmanager "github.com/onosproject/onos-ric/pkg/apps/onos-ric-mlb/manager"
)

var log = logging.GetLogger("main")

// The main entry point.
func main() {
	onosricEndpoint := flag.String("onosricaddr", "onos-ric:5150", "endpoint:port of the ONOS RIC subsystem")
	loadthresh := flag.Float64("threshold", 1, "Threshold for MLB [0, 1] (e.g., 0.5 means 50%)")
	period := flag.Int64("period", 10000, "Period to run MLB procedure [ms]")
	enableMetrics := flag.Bool("enableMetrics", true, "Enable gathering of metrics for Prometheus")
	loglevel := flag.String("loglevel", "warn", "Initial log level - debug, info, warn, error")

	// TODO: Replace this when config.yaml is available to set initial conditions
	initialLogLevel := logging.WarnLevel
	switch *loglevel {
	case "debug":
		initialLogLevel = logging.DebugLevel
	case "info":
		initialLogLevel = logging.InfoLevel
	case "warn":
		initialLogLevel = logging.WarnLevel
	case "error":
		initialLogLevel = logging.ErrorLevel
	}
	flag.Parse()
	log.Infof("logs level: %s", initialLogLevel)

	log.SetLevel(initialLogLevel)
	logging.GetLogger("main").SetLevel(initialLogLevel)
	logging.GetLogger("mlb").SetLevel(initialLogLevel)
	logging.GetLogger("mlb", "manager").SetLevel(initialLogLevel)
	logging.GetLogger("mlb", "southbound").SetLevel(initialLogLevel)
	logging.GetLogger("mlb", "exporter").SetLevel(initialLogLevel)
	log.Info("Starting MLB Application")

	appMgr, err := mlbappmanager.NewManager()
	if err != nil {
		log.Fatal("Unable to load MLB Application: ", err)
	}

	if *enableMetrics {
		log.Warn("Starting MLB Exporter")
		go mlbappexporter.RunMLBExposer(appMgr.SB)
	}

	appMgr.SB.ONOSRICAddr = onosricEndpoint
	appMgr.SB.LoadThresh = loadthresh
	appMgr.SB.Period = period
	appMgr.SB.EnableMetric = *enableMetrics
	appMgr.Run()
}
