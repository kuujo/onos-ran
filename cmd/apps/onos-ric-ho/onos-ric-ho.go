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

	hoappexporter "github.com/onosproject/onos-ric/pkg/apps/onos-ric-ho/exporter"
	hoappmanager "github.com/onosproject/onos-ric/pkg/apps/onos-ric-ho/manager"
	log "k8s.io/klog"
)

// The main entry point.
func main() {
	onosricaddr := flag.String("onosricaddr", "localhost:5150", "address:port of the ONOS RIC subsystem")

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

	log.Info("Starting HO Application")

	appMgr, err := hoappmanager.NewManager()

	log.Info("Starting HO Exposer")
	go hoappexporter.RunHOExposer(appMgr.SB)

	if err != nil {
		log.Fatal("Unable to load HO Application: ", err)
	} else {
		appMgr.SB.ONOSRICAddr = onosricaddr
		appMgr.Run()
	}

	// If this application requires gRPC server, should be called here.
}
