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

package nb

import (
	"github.com/onosproject/helmit/pkg/helm"
	"github.com/onosproject/helmit/pkg/test"
)

// TestSuite is the primary onos-ric test suite
type TestSuite struct {
	test.Suite
}

// SetupTestSuite sets up the onos-ric northbound test suite
func (s *TestSuite) SetupTestSuite() error {
	err := helm.Chart("kubernetes-controller", "https://charts.atomix.io").
		Release("onos-ric-atomix").
		Set("scope", "Namespace").
		Install(true)
	if err != nil {
		return err
	}

	err = helm.Chart("raft-storage-controller", "https://charts.atomix.io").
		Release("onos-ric-raft").
		Set("scope", "Namespace").
		Install(true)
	if err != nil {
		return err
	}

	err = helm.Chart("cache-storage-controller", "https://charts.atomix.io").
		Release("onos-ric-cache").
		Set("scope", "Namespace").
		Install(true)
	if err != nil {
		return err
	}

	sdran := helm.Chart("sd-ran").
		Release("sd-ran").
		Set("global.store.controller", "onos-ric-atomix-kubernetes-controller:5679").
		Set("import.onos-gui.enabled", false).
		Set("import.nem-monitoring.enabled", false).
		Set("onos-ric.service.external.nodePort", 0).
		Set("onos-ric-ho.service.external.nodePort", 0).
		Set("onos-ric-mlb.service.external.nodePort", 0)
	return sdran.Install(true)
}
