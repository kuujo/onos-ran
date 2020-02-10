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
	// for Handover
	m.getListUELinks()

	conn.Close()
}

// getListUELinks gets the list of link between each UE and serving/neighbor stations, and call sendHandoverTrigger if HO is necessary.
func (m *HOSessions) getListUELinks() {
	stream, err := m.client.ListUELinks(context.Background(), &nb.UELinkListRequest{})

	if err != nil {
		return
	}

	numRoutines := 0
	joinChan := make(chan int32)
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
		log.Infof("UE Link(plmnid:%s,ecid:%s,crnti:%s): STA1(ecid:%s,cqi:%d), STA2(ecid:%s,cqi:%d), STA3(ecid:%s,cqi:%d)",
			ueInfo.GetEcgi().GetPlmnid(), ueInfo.GetEcgi().GetEcid(), ueInfo.GetCrnti(),
			ueInfo.GetChannelQualities()[0].GetTargetEcgi().GetEcid(), ueInfo.GetChannelQualities()[0].GetCqiHist(),
			ueInfo.GetChannelQualities()[1].GetTargetEcgi().GetEcid(), ueInfo.GetChannelQualities()[1].GetCqiHist(),
			ueInfo.GetChannelQualities()[2].GetTargetEcgi().GetEcid(), ueInfo.GetChannelQualities()[2].GetCqiHist())
		go m.sendHandoverTrigger(hoapphandover.HODecisionMaker(ueInfo), joinChan)
		numRoutines++
	}

	for i := 0; i < numRoutines; i++ {
		<-joinChan
	}
}

// sendHanmdoverTrigger sends handover trigger to appropriate stations.
func (m *HOSessions) sendHandoverTrigger(hoReq nb.HandOverRequest, joinChan chan int32) {

	// HODecisionMaker function returns nb.HandOverRequest{} when serving stations is the best one
	// No need to trigger handover because serving station is the best one
	if hoReq.GetDstStation() == nil && hoReq.GetSrcStation() == nil {
		log.Info("No need to trigger HO")
		joinChan <- 0
		return
	}
	log.Infof("HO %s from %s to %s", hoReq.GetCrnti(), hoReq.GetSrcStation().GetEcid(), hoReq.GetDstStation().GetEcid())
	_, err := m.client.TriggerHandOver(context.Background(), &hoReq)

	if err != nil {
		log.Error(err)
	}
	joinChan <- 1
}
