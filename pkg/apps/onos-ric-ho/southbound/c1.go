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

	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-lib-go/pkg/southbound"

	"github.com/onosproject/onos-ric/api/nb"
	hoapphandover "github.com/onosproject/onos-ric/pkg/apps/onos-ric-ho/handover"
	"google.golang.org/grpc"
)

var log = logging.GetLogger("ho", "southbound")

// HOSessions is responsible for mapping connections to and interactions with the Northbound of ONOS RAN subsystem.
type HOSessions struct {
	ONOSRICAddr *string
	client      nb.C1InterfaceServiceClient
}

// NewSession creates a new southbound session of HO application.
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
	log.Infof("Connecting to ONOS RAN controller...%s", *m.ONOSRICAddr)

	for {
		// Attempt to create connection to the RIC
		conn, err := southbound.Connect(context.Background(), *m.ONOSRICAddr, "", "")
		if err != nil {
			log.Errorf("Failed to connect: %s", err)
			continue
		}
		log.Infof("Connected to %s", *m.ONOSRICAddr)
		// If successful, manage this connection and don't return until it is
		// no longer valid and all related resources have been properly cleaned-up.
		m.manageConnection(conn)
		time.Sleep(100 * time.Millisecond) // need to be in 10ms - 100ms
	}
}

// manageConnection is responsible for managing a single connection between HO App and ONOS RAN subsystem.
func (m *HOSessions) manageConnection(conn *grpc.ClientConn) {
	m.client = nb.NewC1InterfaceServiceClient(conn)
	if m.client == nil {
		log.Error("Unable to get gRPC NewC1InterfaceServiceClient")
		return
	}

	defer conn.Close()
	// run Handover procedure
	err := m.runHandoverProcedure()
	log.Fatalf("Failed to watch UELinks %s. Closing", err.Error())
}

// runHandoverProcedure runs entire handover procedure - getting UELinkInfo, making decision, and sending trigger messages.
func (m *HOSessions) runHandoverProcedure() error {
	ch := make(chan *nb.UELinkInfo)
	if err := m.watchUELinks(ch); err != nil {
		return err
	}
	for ueLink := range ch { // Block here and wait for UELink stream
		// HO procedure 2. get requirement messages
		hoReq := hoapphandover.HODecisionMaker(ueLink)
		// HO procedure 3. send trigger message
		if hoReq != nil {
			if err := m.sendHandoverTrigger(hoReq); err != nil {
				log.Errorf("Error triggering HO event %s", err.Error())
			}
		}
	}
	return nil
}

// getListUELinks gets the list of link between each UE and serving/neighbor stations, and call sendHandoverTrigger if HO is necessary.
func (m *HOSessions) watchUELinks(ch chan<- *nb.UELinkInfo) error {
	stream, err := m.client.ListUELinks(context.Background(), &nb.UELinkListRequest{Subscribe: true, NoReplay: true})
	if err != nil {
		log.Errorf("Failed to get stream: %s", err)
		return err
	}

	go func() {
		for {
			ueInfo, err := stream.Recv()

			if err == io.EOF {
				break
			} else if err != nil {
				log.Error(err)
				break
			}

			ch <- ueInfo
		}
	}()
	return nil
}

// sendHanmdoverTrigger sends handover trigger to appropriate stations.
func (m *HOSessions) sendHandoverTrigger(hoReq *nb.HandOverRequest) error {
	log.Infof("HO %s(%s,%s) from %s,%s to %s,%s", hoReq.GetCrnti(), hoReq.GetSrcStation().GetPlmnid(), hoReq.GetSrcStation().GetEcid(),
		hoReq.GetSrcStation().GetPlmnid(), hoReq.GetSrcStation().GetEcid(),
		hoReq.GetDstStation().GetPlmnid(), hoReq.GetDstStation().GetEcid())
	_, err := m.client.TriggerHandOver(context.Background(), hoReq)
	return err
}
