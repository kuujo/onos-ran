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

	"github.com/onosproject/onos-lib-go/pkg/logging"

	hoappexporter "github.com/onosproject/onos-ric/pkg/apps/onos-ric-ho/exporter"
	hoappmanager "github.com/onosproject/onos-ric/pkg/apps/onos-ric-ho/manager"
)

var log = logging.GetLogger("main")

// The main entry point.
func main() {
	onosricEndpoint := flag.String("onosricaddr", "onos-ric:5150", "endpoint:port of the ONOS RIC subsystem")
	enableMetrics := flag.Bool("enableMetrics", true, "Enable gathering of metrics for Prometheus")
	flag.Parse()

	log.Info("Starting HO Application")

	appMgr, err := hoappmanager.NewManager()
	if err != nil {
		log.Fatal("Unable to load HO Application: ", err)
	}

	if *enableMetrics {
		log.Info("Starting HO Exporter")
		go hoappexporter.RunHOExposer(appMgr.SB)
	}

	appMgr.SB.ONOSRICAddr = onosricEndpoint
	appMgr.SB.EnableExporter = *enableMetrics
	appMgr.Run()
}
