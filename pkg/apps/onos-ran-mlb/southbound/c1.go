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

package mlbappsouthbound

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/onosproject/onos-ran/api/nb"
	mlbapploadbalance "github.com/onosproject/onos-ran/pkg/apps/onos-ran-mlb/mlb"
	mlbappservice "github.com/onosproject/onos-ran/pkg/apps/onos-ran-mlb/service"
	"google.golang.org/grpc"
	log "k8s.io/klog"
)

// MLBSessions is responsible for mapping connnections to and interactions with the northbound of ONOS-RAN subsystem.
type MLBSessions struct {
	ONOSRANAddr *string
	LoadThresh  *float64
	client      nb.C1InterfaceServiceClient
}

// NewSession creates a new southbound session of MLB application.
func NewSession() (*MLBSessions, error) {
	log.Info("Creating MLBSessions")
	return &MLBSessions{}, nil
}

// Run starts the southbound control loop for mobility load balancing.
func (m *MLBSessions) Run() {
	log.Info("Started MLB App Manager")

	m.manageConnections()
}

// manageConnections handles connnections between MLB App and ONOS RAN subsystem.
func (m *MLBSessions) manageConnections() {
	log.Infof("Connecting to ONOS RAN controllers...%s", *m.ONOSRANAddr)

	for {
		// Attempt to create connection to the simulator
		opts := []grpc.DialOption{
			grpc.WithInsecure(),
			grpc.WithBlock(),
		}
		conn, err := mlbappservice.Connect(*m.ONOSRANAddr, opts...)
		if err == nil {
			// If successful, manage this connection and don't return until it is
			// no longer valid and all related resources have been properly cleaned-up.
			m.manageConnection(conn)
		}
		time.Sleep(1000 * time.Millisecond)
	}
}

// manageConnection is responsible for managing a sigle connection between MLB App and ONOS RAN subsystem.
func (m *MLBSessions) manageConnection(conn *grpc.ClientConn) {
	m.client = nb.NewC1InterfaceServiceClient(conn)

	if m.client == nil {
		return
	}

	joinChan := make(chan int32)

	var stationInfoList []nb.StationInfo
	var stationLinkInfoList []nb.StationLinkInfo
	var ueLinkInfoList []nb.UELinkInfo

	go m.getListStations(&stationInfoList, joinChan)
	go m.getListStationLinks(&stationLinkInfoList, joinChan)
	go m.getListUELinks(&ueLinkInfoList, joinChan)

	// wait until above three go-routine is terminated
	for i := 0; i < 3; i++ {
		<-joinChan
	}

	for _, s := range stationInfoList {
		log.Infof("STA: plmnid:%s,ecid:%s,maxUEs:%d", s.GetEcgi().GetPlmnid(), s.GetEcgi().GetEcid(), s.GetMaxNumConnectedUes())
	}

	for _, sl := range stationLinkInfoList {
		tmpStaInfo := fmt.Sprintf("STA: plmnid:%s,ecid:%s", sl.GetEcgi().GetPlmnid(), sl.GetEcgi().GetEcid())
		tmpNStaInfo := ""
		for _, e := range sl.GetNeighborECGI() {
			tmpNStaInfo = fmt.Sprintf("%s\tNSTA: plmnid:%s,ecid:%s", tmpNStaInfo, e.GetPlmnid(), e.GetEcid())
		}
		log.Info(tmpStaInfo)
	}

	for _, u := range ueLinkInfoList {
		log.Infof("UE: plmnid:%s,ecid:%s,crnti:%s", u.GetEcgi().GetPlmnid(), u.GetEcgi().GetEcid(), u.GetCrnti())
	}

	mlbReqs := mlbapploadbalance.MLBDecisionMaker(stationInfoList, stationLinkInfoList, ueLinkInfoList, m.LoadThresh)
	for _, req := range *mlbReqs {
		m.sendRadioPowerOffset(req)
	}
}

// getListStations gets list of stations from ONOS RAN subsystem.
func (m *MLBSessions) getListStations(stationInfoList *[]nb.StationInfo, joinChan chan int32) {
	stream, err := m.client.ListStations(context.Background(), &nb.StationListRequest{})

	if err != nil {
		log.Error(err)
		joinChan <- 0
		return
	}
	for {
		stationInfo, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Error(err)
			break
		}
		// For debugging
		*stationInfoList = append(*stationInfoList, *stationInfo)
	}
	joinChan <- 1
}

// getListStationLinks gets list of the relationship among stations from ONOS RAN subsystem.
func (m *MLBSessions) getListStationLinks(stationLinkInfoList *[]nb.StationLinkInfo, joinChan chan int32) {
	stream, err := m.client.ListStationLinks(context.Background(), &nb.StationLinkListRequest{})

	if err != nil {
		log.Error(err)
		joinChan <- 0
		return
	}
	for {
		stationLinkInfo, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Error(err)
			break
		}
		// For debugging
		printNeighbors := ""
		for _, n := range stationLinkInfo.GetNeighborECGI() {
			printNeighbors = printNeighbors + " PLMNID: " + n.GetPlmnid() + ",ECID: " + n.GetEcid() + "\t"
		}
		*stationLinkInfoList = append(*stationLinkInfoList, *stationLinkInfo)
	}
	joinChan <- 1
}

// getListUELinks gets the list of link between each UE and serving/neighbor stations.
func (m *MLBSessions) getListUELinks(ueLinkInfoList *[]nb.UELinkInfo, joinChan chan int32) {
	stream, err := m.client.ListUELinks(context.Background(), &nb.UELinkListRequest{})

	if err != nil {
		log.Error(err)
		joinChan <- 0
		return
	}
	for {
		ueInfo, err := stream.Recv()

		if err == io.EOF {
			break
		} else if err != nil {
			log.Error(err)
			break
		}
		*ueLinkInfoList = append(*ueLinkInfoList, *ueInfo)
	}
	joinChan <- 1
}

// sendRadioPowerOffset sends power offset to appropriate stations.
func (m *MLBSessions) sendRadioPowerOffset(mlbReq nb.RadioPowerRequest) {

	log.Infof("MLB: plmnid:%s,ecid:%s,offset:%s", mlbReq.GetEcgi(), mlbReq.GetEcgi(), mlbReq.GetOffset().String())

	_, err := m.client.SetRadioPower(context.Background(), &mlbReq)

	if err != nil {
		log.Error(err)
	}
}
