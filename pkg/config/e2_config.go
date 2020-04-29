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

package config

import (
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-ric/api/sb"
	"github.com/onosproject/onos-ric/pkg/southbound"
)

var log = logging.GetLogger("config")

type e2Config struct {
	e2Chan chan southbound.E2 // channel on which new E2 sessions are advertised
}

// Run fetches the initial RicConfig from onos-config at startup
// and subsequently keeps the cache updated with streaming updates
// from onos-config (assuming thats how onos-config works).
func (c e2Config) Run(e2Chan chan southbound.E2) E2Config {
	log.Infof("Starting E2 config")
	c.e2Chan = e2Chan
	// TBD - go routine to load config from onos-config

	go c.processE2Updates()

	// go c.processOnosConfigUpdates()

	return c
}

func (c e2Config) processE2Updates() {
	for e2 := range c.e2Chan {
		err := e2.Setup()
		if err != nil {
			log.Errorf("Setup procedure failed %s", err.Error())
		}
		err = e2.L2MeasConfig(c.getL2MeasConfig(e2.RemoteAddress()))
		if err != nil {
			log.Errorf("L2MeasConfig procedure failed %s", err.Error())
		}
	}
}

func (c e2Config) Stop() {
	close(c.e2Chan)
}

// getL2MeasReportConfig returns the l2MeasConfig
func (c e2Config) getL2MeasConfig(addr sb.Endpoint) *sb.L2MeasConfig {
	// TBD - Remove hard-coded values and Fetch config from onos-config
	return &sb.L2MeasConfig{
		RadioMeasReportPerUe:   sb.L2MeasReportInterval_MS_10,
		RadioMeasReportPerCell: sb.L2MeasReportInterval_MS_10,
		SchedMeasReportPerUe:   sb.L2MeasReportInterval_MS_10,
		SchedMeasReportPerCell: sb.L2MeasReportInterval_MS_10,
		PdcpMeasReportPerUe:    sb.L2MeasReportInterval_MS_10,
	}
}

// NewE2Config ...
func NewE2Config() E2Config {
	return &e2Config{}
}
