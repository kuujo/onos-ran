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
	EnableMetric bool
	client       nb.C1InterfaceServiceClient
	RNIBCellInfo []*nb.StationInfo
	UEInfoList   []*nb.UEInfo
	MLBEventChan chan *MLBEvent
}

// MLBEvent is responsible for representing each MLB event
type MLBEvent struct {
	MLBReq    nb.RadioPowerRequest
	StartTime time.Time
	EndTime   time.Time
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

	tStart := time.Now()

	// get all R-NIB information
	var wg sync.WaitGroup

	var stationLinkInfoList []nb.StationLinkInfo
	m.RNIBCellInfo = nil
	m.UEInfoList = nil

	// Fork three go-routines.
	wg.Add(3)
	go func() {
		m.RNIBCellInfo = m.getListStations()
		defer wg.Done()
	}()
	go func() {
		stationLinkInfoList = m.getListStationLinks()
		defer wg.Done()
	}()
	go func() {
		m.UEInfoList = m.getListUEs()
		defer wg.Done()
	}()

	// Wait until all go-routines join.
	wg.Wait()

	// if R-NIB (UELink) is not old one, start MLB procedure
	// otherwise, skip this timeslot, because MLBDecisionMaker was already run before
	mlbReqs, _ := mlbapploadbalance.MLBDecisionMaker(m.RNIBCellInfo, stationLinkInfoList, m.UEInfoList, m.LoadThresh)
	tEnd := time.Now()

	for _, req := range *mlbReqs {
		m.sendRadioPowerOffset(req)
	}

	if m.EnableMetric {
		m.MLBEventChan = make(chan *MLBEvent)
		for _, req := range *mlbReqs {
			event := &MLBEvent{
				StartTime: tStart,
				EndTime:   tEnd,
				MLBReq:    req,
			}
			m.MLBEventChan <- event
		}
		close(m.MLBEventChan)
	}
}

// getListStations gets list of stations from ONOS RAN subsystem.
func (m *MLBSessions) getListStations() []*nb.StationInfo {
	var stationInfoList []*nb.StationInfo
	stream, err := m.client.ListStations(context.Background(), &nb.StationListRequest{})

	if err != nil {
		log.Error(err)
		return nil
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
		stationInfoList = append(stationInfoList, stationInfo)
	}
	return stationInfoList
}

// getListStationLinks gets list of the relationship among stations from ONOS RAN subsystem.
func (m *MLBSessions) getListStationLinks() []nb.StationLinkInfo {
	var stationLinkInfoList []nb.StationLinkInfo
	stream, err := m.client.ListStationLinks(context.Background(), &nb.StationLinkListRequest{})

	if err != nil {
		log.Error(err)
		return nil
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
		stationLinkInfoList = append(stationLinkInfoList, *stationLinkInfo)
	}
	return stationLinkInfoList
}

func (m *MLBSessions) getListUEs() []*nb.UEInfo {
	var ueInfoList []*nb.UEInfo
	stream, err := m.client.ListUEs(context.Background(), &nb.UEListRequest{})

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
		ueInfoList = append(ueInfoList, ueInfo)
	}
	return ueInfoList
}

// sendRadioPowerOffset sends power offset to appropriate stations.
func (m *MLBSessions) sendRadioPowerOffset(mlbReq nb.RadioPowerRequest) {

	log.Infof("MLB: plmnid:%s,ecid:%s,offset:%s", mlbReq.GetEcgi(), mlbReq.GetEcgi(), mlbReq.GetOffset().String())

	_, err := m.client.SetRadioPower(context.Background(), &mlbReq)

	if err != nil {
		log.Error(err)
	}
}
