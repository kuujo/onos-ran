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
Package onos-ran is the main entry point to the ONOS RAN subsystem.

Arguments

-caPath <the location of a CA certificate>

-keyPath <the location of a client private key>

-certPath <the location of a client certificate>


See ../../docs/run.md for how to run the application.
*/
package main

import (
	"flag"
	"github.com/onosproject/onos-ran/pkg/manager"
	"github.com/onosproject/onos-ran/pkg/northbound/c1"
	"github.com/onosproject/onos-ran/pkg/service"
	log "k8s.io/klog"
)

// The main entry point
func main() {
	caPath := flag.String("caPath", "", "path to CA certificate")
	keyPath := flag.String("keyPath", "", "path to client private key")
	certPath := flag.String("certPath", "", "path to client certificate")
	simulator := flag.String("simulator", "", "address:port of the RAN simulator")

	//lines 93-109 are implemented according to
	// https://github.com/kubernetes/klog/blob/master/examples/coexist_glog/coexist_glog.go
	// because of libraries importing glog. With glog import we can't call log.InitFlags(nil) as per klog readme
	// thus the alsologtostderr is not set properly and we issue multiple logs.
	// Calling log.InitFlags(nil) throws panic with error `flag redefined: log_dir`
	err := flag.Set("alsologtostderr", "true")
	if err != nil {
		log.Error("Cant' avoid double Error logging ", err)
	}
	flag.Parse()

	klogFlags := flag.NewFlagSet("klog", flag.ExitOnError)
	log.InitFlags(klogFlags)

	// Sync the glog and klog flags.
	flag.CommandLine.VisitAll(func(f1 *flag.Flag) {
		f2 := klogFlags.Lookup(f1.Name)
		if f2 != nil {
			value := f1.Value.String()
			_ = f2.Value.Set(value)
		}
	})
	log.Info("Starting onos-ran")

	mgr, err := manager.NewManager()
	if err != nil {
		log.Fatal("Unable to load onos-ran ", err)
	} else {
		mgr.SB.Simulator = simulator
		mgr.Run()
		err = startServer(*caPath, *keyPath, *certPath)
		if err != nil {
			log.Fatal("Unable to start onos-ran ", err)
		}
	}
}

// Creates gRPC server and registers various services; then serves.
func startServer(caPath string, keyPath string, certPath string) error {
	s := service.NewServer(service.NewServerConfig(caPath, keyPath, certPath))
	s.AddService(c1.Service{})

	return s.Serve(func(started string) {
		log.Info("Started NBI on ", started)
	})
}
