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

	mlbappmanager "github.com/onosproject/onos-ran/pkg/apps/onos-ran-mlb/manager"
	log "k8s.io/klog"
)

// The main entry point.
func main() {
	onosranaddr := flag.String("onosranaddr", "localhost:5150", "address:port of the ONOS RAN subsystem")
	loadthresh := flag.Float64("threshold", 1, "Threshold for MLB [0, 1] (e.g., 0.5 means 50%)")
	period := flag.Int64("period", 10000, "Period to run MLB procedure [ms]")

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

	log.Info("Starting MLB Application")

	appMgr, err := mlbappmanager.NewManager()
	if err != nil {
		log.Fatal("Unable to load MLB Application: ", err)
	} else {
		appMgr.SB.ONOSRANAddr = onosranaddr
		appMgr.SB.LoadThresh = loadthresh
		appMgr.SB.Period = period
		appMgr.Run()
	}

	// If this application requires gRPC server, should be called here.
}
