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

package monitor

import (
	topodevice "github.com/onosproject/onos-topo/api/device"
)

type topoEventHandler = func(chan *topodevice.ListResponse)

// NewTopoMonitorBuilder builds a topology event monitor builder
func NewTopoMonitorBuilder() TopoMonitorBuilder {
	return &TopoMonitor{}
}

// TopoMonitorBuilder a topology monitor builder interface
type TopoMonitorBuilder interface {
	SetTopoChannel(chan *topodevice.ListResponse) TopoMonitorBuilder
	Build() TopoMonitor
}

// TopoEventHandler topology event handler
func (s *TopoMonitor) TopoEventHandler(eventHandler topoEventHandler) {
	go eventHandler(s.topoChannel)
}

// TopoMonitor topology monitor
type TopoMonitor struct {
	topoChannel chan *topodevice.ListResponse
}

// TopoChannel returns topo channel
func (s *TopoMonitor) TopoChannel() chan *topodevice.ListResponse {
	return s.topoChannel
}

// SetTopoChannel sets topo channel
func (s *TopoMonitor) SetTopoChannel(topoChannel chan *topodevice.ListResponse) TopoMonitorBuilder {
	s.topoChannel = topoChannel
	return s
}

// Build builds a topo monitor object
func (s *TopoMonitor) Build() TopoMonitor {
	return TopoMonitor{
		topoChannel: s.topoChannel,
	}
}

// Close closes the underlying channel
func (s *TopoMonitor) Close() {
	close(s.topoChannel)
}
