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
	"io"
	"reflect"
	"sync"
	"time"

	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-lib-go/pkg/southbound"

	"github.com/onosproject/onos-ric/api/nb"
	mlbapploadbalance "github.com/onosproject/onos-ric/pkg/apps/onos-ric-mlb/mlb"
	"google.golang.org/grpc"
)

var log = logging.GetLogger("mlb", "southbound")

// MLBSessions is responsible for mapping connnections to and interactions with the northbound of ONOS-RAN subsystem.
type MLBSessions struct {
	ONOSRICAddr  *string
	LoadThresh   *float64
	Period       *int64
	client       nb.C1InterfaceServiceClient
	prevRNIB     []*nb.UELinkInfo
	RNIBCellInfo []*nb.StationInfo
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
	log.Infof("Connecting to ONOS RAN controllers...%s", *m.ONOSRICAddr)

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
		time.Sleep(time.Duration(*m.Period) * time.Millisecond)
	}
}

// manageConnection is responsible for managing a sigle connection between MLB App and ONOS RAN subsystem.
func (m *MLBSessions) manageConnection(conn *grpc.ClientConn) {
	m.client = nb.NewC1InterfaceServiceClient(conn)

	if m.client == nil {
		return
	}

	// run MLB procedure
	m.runMLBProcedure()

	conn.Close()
}

func (m *MLBSessions) runMLBProcedure() {

	// get all R-NIB information
	var wg sync.WaitGroup

	var stationLinkInfoList []nb.StationLinkInfo
	var ueLinkInfoList []*nb.UELinkInfo
	m.RNIBCellInfo = nil

	// Fork three go-routines.
	wg.Add(3)
	go m.getListStations(&m.RNIBCellInfo, &wg)
	go m.getListStationLinks(&stationLinkInfoList, &wg)
	go func() {
		ueLinkInfoList = m.getListUELinks(ueLinkInfoList)
		defer wg.Done()
	}()

	// Wait until all go-routines join.
	wg.Wait()

	// if R-NIB (UELink) is not old one, start MLB procedure
	// otherwise, skip this timeslot, because MLBDecisionMaker was already run before
	if m.prevRNIB == nil || !m.isEqualUeLinkLists(m.prevRNIB, ueLinkInfoList) {
		mlbReqs, _ := mlbapploadbalance.MLBDecisionMaker(m.RNIBCellInfo, stationLinkInfoList, ueLinkInfoList, m.LoadThresh)
		for _, req := range *mlbReqs {
			m.sendRadioPowerOffset(req)
		}
		m.prevRNIB = ueLinkInfoList
	}
}

// isEqualUeLinkList checks whether the recently received UELinkInfoList and the previously received UELinkInfoList are equivalent.
func (m *MLBSessions) isEqualUeLinkLists(pList []*nb.UELinkInfo, cList []*nb.UELinkInfo) bool {
	if len(pList) == len(cList) && m.containUeLinkLists(pList, cList) {
		return true
	}
	return false
}

// containUeLinkLists checks whether the recently received UELinkInfoList is the subset of the previously received UELinkInfoList.
func (m *MLBSessions) containUeLinkLists(pList []*nb.UELinkInfo, cList []*nb.UELinkInfo) bool {
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

// getListStations gets list of stations from ONOS RAN subsystem.
func (m *MLBSessions) getListStations(stationInfoList *[]*nb.StationInfo, wg *sync.WaitGroup) {
	stream, err := m.client.ListStations(context.Background(), &nb.StationListRequest{})

	if err != nil {
		log.Error(err)
		defer wg.Done()
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
		*stationInfoList = append(*stationInfoList, stationInfo)
	}
	defer wg.Done()
}

// getListStationLinks gets list of the relationship among stations from ONOS RAN subsystem.
func (m *MLBSessions) getListStationLinks(stationLinkInfoList *[]nb.StationLinkInfo, wg *sync.WaitGroup) {
	stream, err := m.client.ListStationLinks(context.Background(), &nb.StationLinkListRequest{})

	if err != nil {
		log.Error(err)
		defer wg.Done()
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
		*stationLinkInfoList = append(*stationLinkInfoList, *stationLinkInfo)
	}
	defer wg.Done()
}

// getListUELinks gets the list of link between each UE and serving/neighbor stations.
func (m *MLBSessions) getListUELinks(ueLinkInfoList []*nb.UELinkInfo) []*nb.UELinkInfo {
	stream, err := m.client.ListUELinks(context.Background(), &nb.UELinkListRequest{})

	if err != nil {
		log.Error(err)
		return nil
	}
	for {
		ueInfo, err := stream.Recv()

		if err == io.EOF {
			break
		} else if err != nil {
			log.Error(err)
			break
		}
		ueLinkInfoList = append(ueLinkInfoList, ueInfo)
	}
	return ueLinkInfoList
}

// sendRadioPowerOffset sends power offset to appropriate stations.
func (m *MLBSessions) sendRadioPowerOffset(mlbReq nb.RadioPowerRequest) {

	log.Infof("MLB: plmnid:%s,ecid:%s,offset:%s", mlbReq.GetEcgi(), mlbReq.GetEcgi(), mlbReq.GetOffset().String())

	_, err := m.client.SetRadioPower(context.Background(), &mlbReq)

	if err != nil {
		log.Error(err)
	}
}

// GetUELinkInfo gets a list of UELinkInfo
func (m *MLBSessions) GetUELinkInfo() []*nb.UELinkInfo {
	return m.prevRNIB
}
