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

// Package hoappmanager is the main coordinator.
package hoappmanager

import (
	"github.com/onosproject/onos-lib-go/pkg/logging"
	hoappsouthbound "github.com/onosproject/onos-ric/pkg/apps/onos-ric-ho/southbound"
)

var log = logging.GetLogger("ho", "manager")

var appMgr AppManager

// AppManager is a single point of entry for Handover application running over ONOS-RIC subsystem.
type AppManager struct {
	SB *hoappsouthbound.HOSessions
}

// NewManager initializes the manager of Handover application.
func NewManager() (*AppManager, error) {
	log.Info("Creating HO App Manager")

	appSb, err := hoappsouthbound.NewSession()
	if err != nil {
		return nil, err
	}

	appMgr = AppManager{appSb}
	return &appMgr, nil
}

// Run starts background tasks associated with the manager.
func (m *AppManager) Run() {
	log.Info("Starting HO App Manager")
	m.SB.Run()
}

// Close kills the channels and manager related objects.
func (m *AppManager) Close() {
	log.Info("Closing HO App Manager")
}

// GetManager returns the initialized and running manager instance.
// Should be called only after NewManager and Run are done.
func GetManager() *AppManager {
	return &appMgr
}
