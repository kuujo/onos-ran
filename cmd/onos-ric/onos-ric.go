// Copyright 2019-present Open Networking Foundation.
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

/*
Package onos-ric is the main entry point to the ONOS RAN Intelligent Controller.

Arguments

-caPath <the location of a CA certificate>

-keyPath <the location of a client private key>

-certPath <the location of a client certificate>


See ../../docs/run.md for how to run the application.
*/
package main

import (
	"flag"
	"github.com/onosproject/onos-ric/pkg/exporter"

	"github.com/onosproject/onos-lib-go/pkg/certs"
	"github.com/onosproject/onos-lib-go/pkg/logging"

	service "github.com/onosproject/onos-lib-go/pkg/northbound"

	"github.com/onosproject/onos-ric/pkg/manager"
	"github.com/onosproject/onos-ric/pkg/northbound/c1"
)

var log = logging.GetLogger("main")

// The main entry point
func main() {
	caPath := flag.String("caPath", "", "path to CA certificate")
	keyPath := flag.String("keyPath", "", "path to client private key")
	certPath := flag.String("certPath", "", "path to client certificate")
	topoEndpoint := flag.String("topoEndpoint", "onos-topo:5150", "topology service endpoint")
	enableMetrics := flag.Bool("enableMetrics", true, "Enable gathering of metrics for Prometheus")
	flag.Parse()

	log.Info("Starting onos-ric")

	opts, err := certs.HandleCertArgs(keyPath, certPath)
	if err != nil {
		log.Fatal(err)
	}

	mgr, err := manager.NewManager(*topoEndpoint, opts)

	if *enableMetrics {
		log.Info("Starting ONOS-RIC Exposer")
		go exporter.RunRICExposer(mgr)
	}

	if err != nil {
		log.Fatal("Unable to load onos-ric ", err)
	} else {
		mgr.Run()
		err = startServer(*caPath, *keyPath, *certPath)
		if err != nil {
			log.Fatal("Unable to start onos-ric ", err)
		}
	}
	mgr.Close()
}

// Creates gRPC server and registers various services; then serves.
func startServer(caPath string, keyPath string, certPath string) error {
	s := service.NewServer(service.NewServerConfig(caPath, keyPath, certPath, 5150, true))
	s.AddService(c1.Service{})
	s.AddService(logging.Service{})

	return s.Serve(func(started string) {
		log.Info("Started NBI on ", started)
	})
}
