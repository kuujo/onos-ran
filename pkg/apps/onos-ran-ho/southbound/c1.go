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
	"reflect"
	"time"

	"github.com/onosproject/onos-ran/api/nb"
	"github.com/onosproject/onos-ran/pkg/apps/onos-ran-ho/handover"
	"github.com/onosproject/onos-ran/pkg/apps/onos-ran-ho/service"
	"google.golang.org/grpc"
	log "k8s.io/klog"
)

// HOSessions is responsible for mapping connections to and interactions with the Northbound of ONOS RAN subsystem.
type HOSessions struct {
	ONOSRANAddr *string
	client      nb.C1InterfaceServiceClient
	prevRNIB    []*nb.UELinkInfo
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
		time.Sleep(100 * time.Millisecond) // need to be in 10ms - 100ms
	}
}

// manageConnection is responsible for managing a single connection between HO App and ONOS RAN subsystem.
func (m *HOSessions) manageConnection(conn *grpc.ClientConn) {
	m.client = nb.NewC1InterfaceServiceClient(conn)
	if m.client == nil {
		return
	}
	// run Handover procedure
	m.runHandoverProcedure()

	conn.Close()
}

// runHandoverProcedure runs entire handover procedure - getting UELinkInfo, making decision, and sending trigger messages.
func (m *HOSessions) runHandoverProcedure() {
	ch := make(chan []*nb.UELinkInfo)
	if err := m.watchUELinks(ch); err != nil {
		return
	}
	for ueLinkList := range ch {
		// if R-NIB (UELink) is not old one, start HO procedure
		// otherwise, skip this timeslot, because HODecisionMaker was already run before
		if m.prevRNIB == nil || !m.isEqualUeLinkLists(m.prevRNIB, ueLinkList) {
			// compare previous and current UELinkList and pick new or different UELinks
			newUeLinks := m.getNewUeLinks(m.prevRNIB, ueLinkList)
			// HO procedure 2. get requirement messages
			hoReqs := hoapphandover.HODecisionMaker(newUeLinks)
			// HO procedure 3. send trigger message
			m.sendHandoverTrigger(hoReqs)
			// Update RNIB in HO App.
			m.prevRNIB = ueLinkList
		}
	}
}

// getNewUeLinks gets new or modified UELinks in the recently received UELinkInfoList compared with the previously received UELinkInfoList.
func (m *HOSessions) getNewUeLinks(pList []*nb.UELinkInfo, cList []*nb.UELinkInfo) []*nb.UELinkInfo {
	var newUeLinks []*nb.UELinkInfo
	for i := 0; i < len(cList); i++ {
		for j := 0; j < len(pList); j++ {
			if reflect.DeepEqual(cList[i], pList[j]) {
				break
			}
			if j == len(pList)-1 {
				newUeLinks = append(newUeLinks, cList[i])
				// analyze UEInfo and call sendHandoverTrigger if handover is necessary.
				log.Infof("UE Link(plmnid:%s,ecid:%s,crnti:%s): STA1(ecid:%s,cqi:%d), STA2(ecid:%s,cqi:%d), STA3(ecid:%s,cqi:%d)",
					cList[i].GetEcgi().GetPlmnid(), cList[i].GetEcgi().GetEcid(), cList[i].GetCrnti(),
					cList[i].GetChannelQualities()[0].GetTargetEcgi().GetEcid(), cList[i].GetChannelQualities()[0].GetCqiHist(),
					cList[i].GetChannelQualities()[1].GetTargetEcgi().GetEcid(), cList[i].GetChannelQualities()[1].GetCqiHist(),
					cList[i].GetChannelQualities()[2].GetTargetEcgi().GetEcid(), cList[i].GetChannelQualities()[2].GetCqiHist())
			}
		}
	}
	return newUeLinks
}

// isEqualUeLinkList checks whether the recently received UELinkInfoList and the previously received UELinkInfoList are equivalent.
func (m *HOSessions) isEqualUeLinkLists(pList []*nb.UELinkInfo, cList []*nb.UELinkInfo) bool {
	if len(pList) == len(cList) && m.containUeLinkLists(pList, cList) {
		return true
	}
	return false
}

// containUeLinkLists checks whether the recently received UELinkInfoList is the subset of the previously received UELinkInfoList.
func (m *HOSessions) containUeLinkLists(pList []*nb.UELinkInfo, cList []*nb.UELinkInfo) bool {
	for i := 0; i < len(cList); i++ {
		for j := 0; j < len(pList); j++ {
			if reflect.DeepEqual(cList[i], pList[j]) {
				break
			}
			if j == len(pList)-1 {
				return false
			}
		}
	}
	return true
}

// getListUELinks gets the list of link between each UE and serving/neighbor stations, and call sendHandoverTrigger if HO is necessary.
func (m *HOSessions) watchUELinks(ch chan<- []*nb.UELinkInfo) error {
	ueLinks := make(map[ueLinkID]*nb.UELinkInfo)

	stream, err := m.client.ListUELinks(context.Background(), &nb.UELinkListRequest{Subscribe: true})
	if err != nil {
		return nil
	}

	go func() {
		var rxUeLinkInfoList []*nb.UELinkInfo
		for {
			ueInfo, err := stream.Recv()

			if err == io.EOF {
				break
			} else if err != nil {
				log.Error(err)
				break
			}

			id := ueLinkID{
				PlmnID: ueInfo.Ecgi.Plmnid,
				Ecid:   ueInfo.Ecgi.Ecid,
				Crnti:  ueInfo.Crnti,
			}
			ueLinks[id] = ueInfo

			rxUeLinkInfoList = make([]*nb.UELinkInfo, 0, len(ueLinks))
			for _, link := range ueLinks {
				rxUeLinkInfoList = append(rxUeLinkInfoList, link)
			}
			ch <- rxUeLinkInfoList
		}
	}()
	return nil
}

type ueLinkID struct {
	PlmnID string
	Ecid   string
	Crnti  string
}

// sendHanmdoverTrigger sends handover trigger to appropriate stations.
func (m *HOSessions) sendHandoverTrigger(hoReqs []*nb.HandOverRequest) {
	for i := 0; i < len(hoReqs); i++ {
		log.Infof("HO %s(%s,%s) from %s,%s to %s,%s", hoReqs[i].GetCrnti(), hoReqs[i].GetSrcStation().GetPlmnid(), hoReqs[i].GetSrcStation().GetEcid(),
			hoReqs[i].GetSrcStation().GetPlmnid(), hoReqs[i].GetSrcStation().GetEcid(),
			hoReqs[i].GetDstStation().GetPlmnid(), hoReqs[i].GetDstStation().GetEcid())
		_, err := m.client.TriggerHandOver(context.Background(), hoReqs[i])
		if err != nil {
			log.Error(err)
		}
	}
}
