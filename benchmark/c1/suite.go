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

package c1

import (
	"github.com/onosproject/onos-ric/api/nb"
	"github.com/onosproject/onos-test/pkg/benchmark"
	"github.com/onosproject/onos-test/pkg/onit/env"
	"github.com/onosproject/onos-test/pkg/onit/setup"
)

// BenchmarkSuite is an onos-ric gNMI benchmark suite
type BenchmarkSuite struct {
	benchmark.Suite
	simulator env.SimulatorEnv
	client    nb.C1InterfaceServiceClient
}

// SetupSuite :: benchmark
func (s *BenchmarkSuite) SetupSuite(c *benchmark.Context) {
	setup.Atomix()
	setup.Database().Raft()
	setup.RAN().SetReplicas(1)
	setup.SetupOrDie()
}

// SetupBenchmark :: benchmark
func (s *BenchmarkSuite) SetupBenchmark(c *benchmark.Context) {
	ranSim := env.Simulators().
		New().
		SetImage("onosproject/ric-simulator:latest").
		SetPort(5150).
		SetName("ran-simulator").
		SetAddDevice(false).
		AddOrDie()

	s.simulator = ranSim
	client, err := env.RAN().NewRANC1ServiceClient()
	if err != nil {
		panic(err)
	}
	s.client = client
}

// TearDownBenchmark :: benchmark
func (s *BenchmarkSuite) TearDownBenchmark(c *benchmark.Context) {
	s.simulator.RemoveOrDie()
}
