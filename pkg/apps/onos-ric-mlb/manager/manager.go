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

// Package mlbappmanager is the main coordinator.
package mlbappmanager

import (
	"github.com/onosproject/onos-lib-go/pkg/logging"
	mlbappsouthbound "github.com/onosproject/onos-ric/pkg/apps/onos-ric-mlb/southbound"
)

var log = logging.GetLogger("mlb", "manager")

var appMgr AppManager

// AppManager is a single point of entry for MLB application running over ONOS-RAN subsystem.
type AppManager struct {
	SB *mlbappsouthbound.MLBSessions
}

// NewManager initializes the manager of MLB application.
func NewManager() (*AppManager, error) {
	log.Info("Creating MLB App Manager")

	appSb, err := mlbappsouthbound.NewSession()
	if err != nil {
		return nil, err
	}

	appMgr = AppManager{appSb}
	return &appMgr, nil
}

// Run starts background tasks associated with the manager.
func (m *AppManager) Run() {
	log.Info("Starting MLB App manager")
	m.SB.Run()
}

// Close kills the channels and manager related objects.
func (m *AppManager) Close() {
	log.Info("Closing MLB App Manager")
}

// GetManager returns the initialized and running manager instance.
// Should be called only after NewManager and Run are done.
func GetManager() *AppManager {
	return &appMgr
}
