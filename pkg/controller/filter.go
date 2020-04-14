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

package controller

import (
	mastershipstore "github.com/onosproject/onos-ric/pkg/store/mastership"
)

// Filter filters individual events for a node
// Each time an event is received from a watcher, the filter has the option to discard the request by
// returning false.
type Filter interface {
	// Accept indicates whether to accept the given object
	Accept(Request) bool
}

// MastershipFilter activates a controller on acquiring mastership
// The MastershipFilter requires a DeviceResolver to extract a device ID from each request. Given a device
// ID, the MastershipFilter rejects any requests for which the local node is not the master for the device.
type MastershipFilter struct {
	Store    mastershipstore.Store
	Resolver MastershipKeyResolver
}

// Accept accepts the given ID if the local node is the master
func (f *MastershipFilter) Accept(request Request) bool {
	key, err := f.Resolver.Resolve(request)
	if err != nil {
		return false
	}
	election, err := f.Store.GetElection(key)
	if err != nil {
		return false
	}
	master, err := election.IsMaster()
	if err != nil {
		return false
	}
	return master
}

var _ Filter = &MastershipFilter{}

// MastershipKeyResolver resolves a mastership key from a request
type MastershipKeyResolver interface {
	// Resolve resolves a mastership key
	Resolve(request Request) (mastershipstore.Key, error)
}
