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
	"context"
	"io"
	"time"
	"unsafe"

	"github.com/onosproject/onos-ran/api/nb"
	hoapphandover "github.com/onosproject/onos-ran/pkg/apps/onos-ran-ho/handover"
	hoappservice "github.com/onosproject/onos-ran/pkg/apps/onos-ran-ho/service"
	"google.golang.org/grpc"
	log "k8s.io/klog"
)

// HOSessions is responsible for mapping connections to and interactions with the Northbound of ONOS RAN subsystem.
type HOSessions struct {
	ONOSRANAddr *string
	client      nb.C1InterfaceServiceClient
}

// NewSession creates a new southbound sessions of HO application.
func NewSession() (*HOSessions, error) {
	log.Info("Creating HOSessions")
	return &HOSessions{}, nil
}

// Run starts the southbound control loop for handover.
func (m *HOSessions) Run() {

	log.Info("Started HO App Manager")
	go m.manageConnections()
	for {
		time.Sleep(100 * time.Second)
	}
}

// manageConnections handles connections between HO App and ONOS RAN subsystem.
func (m *HOSessions) manageConnections() {
	log.Infof("Connecting to ONOS RAN controller...%s", *m.ONOSRANAddr)

	for {
		// Attempt to create connection to the simulator
		opts := []grpc.DialOption{
			grpc.WithInsecure(),
			grpc.WithBlock(),
		}

		conn, err := hoappservice.Connect(*m.ONOSRANAddr, opts...)
		if err == nil {
			// If successful, manage this connection and don't return until it is
			// no longer valid and all related resources have been properly cleaned-up.
			m.manageConnection(conn)
		}
		time.Sleep(1000 * time.Millisecond) // need to be in 10ms - 100ms
	}
}

// manageConnection is responsible for managing a single connection between HO App and ONOS RAN subsystem.
func (m *HOSessions) manageConnection(conn *grpc.ClientConn) {
	m.client = nb.NewC1InterfaceServiceClient(conn)

	if m.client == nil {
		return
	}

	//log.Info("Connected to ONOS RAN subsystem")

	// for test -> will be removed
	go m.getListStations()

	// for test -> will be removed
	go m.getListStationLinks()

	// for Handover
	m.getListUELinks()

	conn.Close()

	//log.Info("Disconnected from ONOS RAN subsystem")

}

// getListStations gets list of stations from ONOS RAN subsystem.
// It is for test: it is not really necessary for handover application, but necessary for test whether gRPC client in HO app works well.
// After test, it will be totally moved to MLB app
func (m *HOSessions) getListStations() {
	stream, err := m.client.ListStations(context.Background(), &nb.StationListRequest{})

	if err != nil {
		log.Error(err)
		return
	}
	for {
		stationInfo, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Error(err)
			break
		}

		// For debugging
		log.Infof("Received station info: ECGI (PLMNID: %s and ECID: %s) and MaxNumConnectedUEs: %d",
			stationInfo.GetEcgi().GetPlmnid(), stationInfo.GetEcgi().GetEcid(), stationInfo.GetMaxNumConnectedUes())
	}
}

// getListStationLinks gets list of the relationship among stations from ONOS RAN subsystem.
// It is for test: it is not really necessary for handover applicationm, but necessary for test whether gRPC client in HO app works well.
// After test, it will be totally moved to MLB app
func (m *HOSessions) getListStationLinks() {
	stream, err := m.client.ListStationLinks(context.Background(), &nb.StationLinkListRequest{})

	if err != nil {
		log.Error(err)
		return
	}

	for {
		stationLinkInfo, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Error(err)
			break
		}

		// For debugging
		log.Infof("The station (PLMNID: %s and ECID: %s) has %d neighbor stations",
			stationLinkInfo.GetEcgi().GetPlmnid(), stationLinkInfo.GetEcgi().GetEcid(), len(stationLinkInfo.GetNeighborECGI()))
	}
}

// getListUELinks gets the list of link between each UE and serving/neighbor stations, and call sendHandoverTrigger if HO is necessary.
func (m *HOSessions) getListUELinks() {
	stream, err := m.client.ListUELinks(context.Background(), &nb.UELinkListRequest{})

	if err != nil {
		return
	}

	for {
		ueInfo, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Error(err)
			break
		}

		// analyze UEInfo and call sendHandoverTrigger if handover is necessary.
		go m.sendHandoverTrigger(hoapphandover.HODecisionMaker(ueInfo))
	}
}

func (m *HOSessions) sendHandoverTrigger(hoReq nb.HandOverRequest) {

	if unsafe.Sizeof(hoReq) == 0 {
		return
	}

	_, err := m.client.TriggerHandOver(context.Background(), &hoReq)

	if err != nil {
		log.Error(err)
	}
}
