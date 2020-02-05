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

// Package manager is is the main coordinator for the ONOS RAN subsystem.
package manager

import (
	"github.com/onosproject/onos-ran/api/sb"
	"github.com/onosproject/onos-ran/pkg/southbound"
	"github.com/onosproject/onos-ran/pkg/store/ran"
	log "k8s.io/klog"
)

var mgr Manager

// NewManager initializes the RAN subsystem.
func NewManager() (*Manager, error) {
	log.Info("Creating Manager")

	db, err := ran.NewRanStore()
	if err != nil {
		return nil, err
	}

	sb, err := southbound.NewSessions()
	if err != nil {
		return nil, err
	}

	mgr = Manager{db, sb}
	return &mgr, nil

}

// Manager single point of entry for the topology system.
type Manager struct {
	store ran.Store
	SB    *southbound.Sessions
}

// GetControlUpdates gets a control update based on a given ID
func (m *Manager) GetControlUpdates() ([]sb.ControlUpdate, error) {
	return m.store.List(), nil
}

// Run starts a synchronizer based on the devices and the northbound services.
func (m *Manager) Run() {
	log.Info("Starting Manager")
	m.SB.Run()
}

//Close kills the channels and manager related objects
func (m *Manager) Close() {
	log.Info("Closing Manager")
}

// GetManager returns the initialized and running instance of manager.
// Should be called only after NewManager and Run are done.
func GetManager() *Manager {
	return &mgr
}
