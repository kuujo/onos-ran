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
	ONOSRICAddr    *string
	LoadThresh     *float64
	Period         *int64
	EnableMetric   bool
	client         nb.C1InterfaceServiceClient
	RNIBCellInfo   []*nb.StationInfo
	UEInfoList     []*nb.UEInfo
	RNIBCellMap    map[string]*nb.StationInfo
	UEInfoMap      map[string]*nb.UEInfo
	StationLinkMap map[string]*nb.StationLinkInfo
}

// ChanMLBEvent is a go channel to pass MLBEvent to exporter
var ChanMLBEvent chan MLBEvent

// MLBEvent is responsible for representing each MLB event
type MLBEvent struct {
	MLBReqs   []nb.RadioPowerRequest
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

	// Init maps
	m.RNIBCellMap = make(map[string]*nb.StationInfo)
	m.UEInfoMap = make(map[string]*nb.UEInfo)
	m.StationLinkMap = make(map[string]*nb.StationLinkInfo)

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
		time.Sleep(100 * time.Millisecond) // Re-dial timer
	}
}

// manageConnection is responsible for managing a sigle connection between MLB App and ONOS RAN subsystem.
func (m *MLBSessions) manageConnection(conn *grpc.ClientConn) {
	m.client = nb.NewC1InterfaceServiceClient(conn)

	if m.client == nil {
		return
	}

	// run MLB procedure
	if m.EnableMetric {
		errMLB := m.runMLBProcedureWithExporter()
		if errMLB != nil {
			log.Error(errMLB)
		}
	} else {
		errMLB := m.runMLBProcedure()
		if errMLB != nil {
			log.Error(errMLB)
		}
	}

	conn.Close()
}

func (m *MLBSessions) runMLBProcedureWithExporter() error {
	staCh := make(chan *nb.StationInfo)
	staLinkCh := make(chan *nb.StationLinkInfo)
	ueCh := make(chan *nb.UEInfo)

	if err := m.watchStations(staCh); err != nil {
		return err
	}
	if err := m.watchStationLinks(staLinkCh); err != nil {
		return err
	}
	if err := m.watchUEs(ueCh); err != nil {
		return err
	}

	go m.collectNetworkView(staCh, staLinkCh, ueCh)

	// run MLB algorithm periodically
	for {
		tStart := time.Now()
		mlbReqs, _ := mlbapploadbalance.MLBDecisionMaker(m.RNIBCellMap, m.StationLinkMap, m.UEInfoMap, m.LoadThresh)
		tEnd := time.Now()

		numEvents := 0
		for _, req := range *mlbReqs {
			if req.Offset != nb.StationPowerOffset_PA_DB_0 {
				numEvents++
			}
			m.sendRadioPowerOffset(req)
		}

		if numEvents > 0 {
			tmp := MLBEvent{
				MLBReqs:   *mlbReqs,
				StartTime: tStart,
				EndTime:   tEnd,
			}
			ChanMLBEvent <- tmp
		}
		time.Sleep(time.Duration(*m.Period) * time.Millisecond)
	}
}

func (m *MLBSessions) runMLBProcedure() error {
	staCh := make(chan *nb.StationInfo)
	staLinkCh := make(chan *nb.StationLinkInfo)
	ueCh := make(chan *nb.UEInfo)

	if err := m.watchStations(staCh); err != nil {
		return err
	}
	if err := m.watchStationLinks(staLinkCh); err != nil {
		return err
	}
	if err := m.watchUEs(ueCh); err != nil {
		return err
	}

	go m.collectNetworkView(staCh, staLinkCh, ueCh)

	// run MLB algorithm periodically
	for {
		mlbReqs, _ := mlbapploadbalance.MLBDecisionMaker(m.RNIBCellMap, m.StationLinkMap, m.UEInfoMap, m.LoadThresh)

		for _, req := range *mlbReqs {
			m.sendRadioPowerOffset(req)
		}

		time.Sleep(time.Duration(*m.Period) * time.Millisecond)
	}
}

func (m *MLBSessions) collectNetworkView(staCh chan *nb.StationInfo, staLinkCh chan *nb.StationLinkInfo, ueCh chan *nb.UEInfo) {
	go m.updateStaInfo(staCh)
	go m.updateStaLinkInfo(staLinkCh)
	go m.updateUEInfo(ueCh)
}

func (m *MLBSessions) updateStaInfo(staCh chan *nb.StationInfo) {
	for sta := range staCh {
		mlbapploadbalance.RNIBCellMapMutex.Lock()
		m.RNIBCellMap[m.getStaID(sta)] = sta
		mlbapploadbalance.RNIBCellMapMutex.Unlock()
	}
}

func (m *MLBSessions) updateStaLinkInfo(staLinkCh chan *nb.StationLinkInfo) {
	for staLink := range staLinkCh {
		mlbapploadbalance.StaLinkMapMutex.Lock()
		m.StationLinkMap[m.getStaLinkID(staLink)] = staLink
		mlbapploadbalance.StaLinkMapMutex.Unlock()
	}
}

func (m *MLBSessions) updateUEInfo(ueCh chan *nb.UEInfo) {
	for ue := range ueCh {
		mlbapploadbalance.UEMapMutex.Lock()
		m.UEInfoMap[m.getUeID(ue)] = ue
		mlbapploadbalance.UEMapMutex.Unlock()
	}
}

func (m *MLBSessions) getStaID(sta *nb.StationInfo) string {
	return sta.GetEcgi().String()
}

func (m *MLBSessions) getStaLinkID(staLink *nb.StationLinkInfo) string {
	return staLink.GetEcgi().String()
}

func (m *MLBSessions) getUeID(ue *nb.UEInfo) string {
	return ue.GetImsi()
}

func (m *MLBSessions) watchStations(ch chan<- *nb.StationInfo) error {
	stream, err := m.client.ListStations(context.Background(), &nb.StationListRequest{
		Subscribe: true,
	})

	if err != nil {
		log.Errorf("Failed to get stream: %s", err)
		return err
	}

	go func() {
		for {
			staInfo, err := stream.Recv()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Error(err)
			}
			ch <- staInfo
		}
	}()

	return nil
}

func (m *MLBSessions) watchStationLinks(ch chan<- *nb.StationLinkInfo) error {
	stream, err := m.client.ListStationLinks(context.Background(), &nb.StationLinkListRequest{
		Subscribe: true,
	})

	if err != nil {
		log.Errorf("Failed to get stream: %s", err)
		return err
	}

	go func() {
		for {
			staLinkInfo, err := stream.Recv()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Error(err)
			}
			ch <- staLinkInfo
		}
	}()

	return nil
}

func (m *MLBSessions) watchUEs(ch chan<- *nb.UEInfo) error {
	stream, err := m.client.ListUEs(context.Background(), &nb.UEListRequest{
		Subscribe: true,
	})

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
			}
			ch <- ueInfo
		}
	}()

	return nil
}

// sendRadioPowerOffset sends power offset to appropriate stations.
func (m *MLBSessions) sendRadioPowerOffset(mlbReq nb.RadioPowerRequest) {

	log.Infof("MLB: plmnid:%s,ecid:%s,offset:%s", mlbReq.GetEcgi(), mlbReq.GetEcgi(), mlbReq.GetOffset().String())

	_, err := m.client.SetRadioPower(context.Background(), &mlbReq)

	if err != nil {
		log.Error(err)
	}
}
