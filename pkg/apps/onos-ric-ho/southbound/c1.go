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

// ChanHOEvent is a go channel to pass HOEvent to exporter
var ChanHOEvent chan HOEvent

// HOEvent represents a single HO event
type HOEvent struct {
	TimeStamp   time.Time
	CRNTI       string
	SrcPlmnID   string
	SrcEcid     string
	DstPlmnID   string
	DstEcid     string
	ElapsedTime int64
}

// HOSessions is responsible for mapping connections to and interactions with the Northbound of ONOS RAN subsystem.
type HOSessions struct {
	ONOSRICAddr    *string
	EnableExporter bool
	HystCqi        int
	A3OffsetCqi    int
	TTTMs          int
	client         nb.C1InterfaceServiceClient
	hoReqChan      chan *nb.HandOverRequest
}

// NewSession creates a new southbound session of HO application.
func NewSession() (*HOSessions, error) {
	log.Info("Creating HOSessions")
	return &HOSessions{
		hoReqChan: make(chan *nb.HandOverRequest),
	}, nil
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
		opts := []grpc.DialOption{
			grpc.WithStreamInterceptor(southbound.RetryingStreamClientInterceptor(100 * time.Millisecond)),
		}
		conn, err := southbound.Connect(context.Background(), *m.ONOSRICAddr, "", "", opts...)
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

	go m.openHandoverTriggersStream(m.hoReqChan)
	var err error
	// run Handover procedure
	if m.EnableExporter {
		err = m.runHandoverProcedureWithExporter()
	} else {
		err = m.runHandoverProcedure()
	}
	log.Fatalf("Failed to watch UELinks %s. Closing", err.Error())
}

// runHandoverProcedure runs entire handover procedure - getting UELinkInfo, making decision, and sending trigger messages.
func (m *HOSessions) runHandoverProcedure() error {
	ch := make(chan *nb.UELinkInfo)
	if err := m.watchUELinks(ch); err != nil {
		return err
	}
	for ueLink := range ch { // Block here and wait for UELink stream
		go hoapphandover.HODecisionMakerWithHOParams(ueLink, m.hoReqChan, m.HystCqi, m.A3OffsetCqi, m.TTTMs)
	}
	return nil
}

// runHandoverProcedure runs entire handover procedure - getting UELinkInfo, making decision, and sending trigger messages.
func (m *HOSessions) runHandoverProcedureWithExporter() error {
	ch := make(chan *nb.UELinkInfo)
	if err := m.watchUELinks(ch); err != nil {
		return err
	}

	for ueLink := range ch { // Block here and wait for UELink stream
		go hoapphandover.HODecisionMakerWithHOParams(ueLink, m.hoReqChan, m.HystCqi, m.A3OffsetCqi, m.TTTMs)
		// TODO: add exporter code here:
	}
	return nil
}

// getListUELinks gets the list of link between each UE and serving/neighbor stations, and call sendHandoverTrigger if HO is necessary.
func (m *HOSessions) watchUELinks(ch chan<- *nb.UELinkInfo) error {
	stream, err := m.client.ListUELinks(context.Background(),
		&nb.UELinkListRequest{
			Subscribe: true,
			NoReplay:  true,
			Noimsi:    true,
		})
	if err != nil {
		log.Errorf("Failed to get stream: %s", err)
		return err
	}

	go func() {
		hoapphandover.InitA3EventMap()
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

func (m *HOSessions) openHandoverTriggersStream(reqChan chan *nb.HandOverRequest) {
	stream, err := m.client.TriggerHandOverStream(context.Background())
	if err != nil {
		log.Errorf("Error on opening HO Trigger Channel %v", err)
		return
	}
	for hoReq := range reqChan { // Block here until HORequest is received on channel
		log.Infof("HO %s(%s,%s) from %s,%s to %s,%s", hoReq.GetCrnti(), hoReq.GetSrcStation().GetPlmnid(), hoReq.GetSrcStation().GetEcid(),
			hoReq.GetSrcStation().GetPlmnid(), hoReq.GetSrcStation().GetEcid(),
			hoReq.GetDstStation().GetPlmnid(), hoReq.GetDstStation().GetEcid())
		if err := stream.Send(hoReq); err != nil {
			log.Errorf("Error on sending HORequest %v", err)
			return
		}
	}
}
