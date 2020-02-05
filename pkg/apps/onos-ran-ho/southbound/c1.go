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

// Package hoappsouthbound is the southbound of HO application running over ONOS RAN subsystem.
package hoappsouthbound

import (
	log "k8s.io/klog"
)

// HOSessions is responsible for mapping connections to and interactions with the Northbound of ONOS RAN subsystem.
type HOSessions struct {
	ONOSRANAddr *string
	//client              nb.C1InterfaceServiceClient
	//ueInfoResponses     chan nb.UEInfo
	//ueLinkInfoResponses chan nb.UELinkInfo
	//hoRequest           chan nb.HandOverRequest
}

// NewSession creates a new southbound sessions of HO application.
func NewSession() (*HOSessions, error) {
	log.Info("Creating HOSessions")
	return &HOSessions{}, nil
}

// Run starts the southbound control loop for handover.
func (m *HOSessions) Run() {
	// To-Do: Add control loop.
}
