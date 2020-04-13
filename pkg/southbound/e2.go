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

package southbound

import (
	"context"
	"fmt"
	"github.com/onosproject/onos-lib-go/pkg/southbound"
	"github.com/onosproject/onos-ric/pkg/store/mastership"
	"sync"
	"time"

	topodevice "github.com/onosproject/onos-topo/api/device"

	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-ric/api/sb"
	"github.com/onosproject/onos-ric/api/sb/e2ap"
	"github.com/onosproject/onos-ric/api/sb/e2sm"
	"google.golang.org/grpc"
)

var log = logging.GetLogger("southbound")

//MapHOEventMeasuredRIC ...
var MapHOEventMeasuredRIC map[string]HOEventMeasuredRIC

// MutexMapHOEventMeasuredRIC ...
var MutexMapHOEventMeasuredRIC sync.RWMutex

//ChanHOEvent ...
var ChanHOEvent chan HOEventMeasuredRIC

// TelemetryUpdateHandler - a function for the session to write back to manager without the import cycle
type TelemetryUpdateHandler = func(e2ap.RicIndication)

// RicControlResponseHandler - ricControlResponse messages to manager without the import cycle
type RicControlResponseHandler = func(e2ap.RicControlResponse)

// ControlUpdateHandler - a function for the session to write back to manager without the import cycle
type ControlUpdateHandler = func(e2ap.RicIndication)

// Session is responsible for managing connections to and interactions with the RAN southbound.
type Session struct {
	Device   topodevice.Device
	Ecgi     sb.ECGI
	election mastership.Election

	ricControlRequestChan  chan e2ap.RicControlRequest
	ricControlResponseChan chan e2ap.RicControlResponse
	controlIndications     chan e2ap.RicIndication
	telemetryIndications   chan e2ap.RicIndication

	RicControlResponseHandlerFunc RicControlResponseHandler
	ControlUpdateHandlerFunc      ControlUpdateHandler
	TelemetryUpdateHandlerFunc    TelemetryUpdateHandler
	EnableMetrics                 bool
}

// HOEventMeasuredRIC is struct including UE ID and its eNB ID
type HOEventMeasuredRIC struct {
	Timestamp   time.Time
	Crnti       string
	DestPlmnID  string
	DestECID    string
	ElapsedTime int64
}

// NewSession creates a new southbound session controller.
func NewSession(ecgi sb.ECGI, device topodevice.Device, mastershipStore mastership.Store) (*Session, error) {
	log.Infof("Creating Session for %v at %s", ecgi, device.Address)

	election, err := mastershipStore.GetElection(mastership.NewKey(ecgi.PlmnId, ecgi.Ecid))
	if err != nil {
		return nil, err
	}

	session := &Session{
		Device:                 device,
		Ecgi:                   ecgi,
		election:               election,
		ricControlRequestChan:  make(chan e2ap.RicControlRequest),
		ricControlResponseChan: make(chan e2ap.RicControlResponse),
		controlIndications:     make(chan e2ap.RicIndication),
		telemetryIndications:   make(chan e2ap.RicIndication),
	}
	if err := session.listen(); err != nil {
		return nil, err
	}
	return session, nil
}

// listen starts the session listening for mastership
func (s *Session) listen() error {
	ch := make(chan mastership.State)
	err := s.election.Watch(ch)
	if err != nil {
		return err
	}
	go func() {
		var errCh chan error
		for state := range ch {
			if state.Master == s.election.NodeID() && errCh == nil {
				errCh = make(chan error)
				s.connect(errCh)
			} else if state.Master != s.election.NodeID() && errCh != nil {
				s.disconnect(errCh)
				errCh = nil
			}
		}
	}()
	return nil
}

// connect connects the session to the device
func (s *Session) connect(errCh chan error) {
	MapHOEventMeasuredRIC = make(map[string]HOEventMeasuredRIC)
	go s.manageConnections(errCh)
	go s.recvTelemetryUpdates(errCh)
	go s.recvUpdates(errCh)
}

// disconnect disconnects the session from the device
func (s *Session) disconnect(errCh chan error) {
	// Close the error channel to disconnect the session
	close(errCh)
}

func (s *Session) recvTelemetryUpdates(errCh chan error) {
	for {
		select {
		case update := <-s.telemetryIndications:
			s.TelemetryUpdateHandlerFunc(update)
		case <-errCh:
			return
		}
	}
}

func (s *Session) recvUpdates(errCh chan error) {
	for {
		select {
		case update := <-s.controlIndications:
			s.ControlUpdateHandlerFunc(update)
		case <-errCh:
			return
		}
	}
}

// SendRicControlRequest sends the specified RicControlRequest on the control channel
func (s *Session) SendRicControlRequest(req e2ap.RicControlRequest) error {
	mastership, err := s.election.GetState()
	if err != nil {
		return err
	}
	if mastership.Master != s.election.NodeID() {
		return fmt.Errorf("not the master for %s", s.Device.Address)
	}

	switch req.GetHdr().GetMessageType() {
	case sb.MessageType_HO_REQUEST:
		if s.EnableMetrics {
			tC := time.Now()
			go s.updateHOEventMeasuredRIC(req, tC)
		}
	default:
	}

	s.ricControlRequestChan <- req
	return nil
}

func (s *Session) manageConnections(masterCh chan error) {
	errCh := make(chan error)
	for {
		conn, err := s.createConnection(errCh)
		if err == nil {
			for {
				select {
				case <-masterCh:
					if conn != nil {
						_ = conn.Close()
					}
					return
				case <-errCh:
					break
				}
			}
		}
		select {
		case <-masterCh:
			return
		case <-time.After(5 * time.Second):
		}
	}
}

func (s *Session) createConnection(errCh chan<- error) (*grpc.ClientConn, error) {
	// Attempt to create connection to the simulator
	log.Infof("Connecting to simulator...%s with context", s.Device.Address)
	opts := []grpc.DialOption{
		grpc.WithStreamInterceptor(southbound.RetryingStreamClientInterceptor(100 * time.Millisecond)),
	}
	conn, err := southbound.Connect(context.Background(), string(s.Device.Address), s.Device.TLS.Cert, s.Device.TLS.Key, opts...)
	if err != nil {
		return nil, err
	}
	client := e2ap.NewE2APClient(conn)
	go s.ricControl(client, errCh)
	go s.ricSubscribe(client, errCh)
	return conn, nil
}

func (s *Session) ricControl(client e2ap.E2APClient, errors chan<- error) {
	stream, err := client.RicControl(context.Background())
	if err != nil {
		errors <- err
		return
	}
	waitc := make(chan error)
	go func() {
		for {
			resp, err := stream.Recv()
			if err != nil {
				close(waitc)
				return
			}
			//log.Infof("Got ricControlResponse messageType %d from %s", resp.Hdr.MessageType, s.EndPoint)
			s.ricControlResponse(resp)
		}
	}()

	cellConfigReq := e2ap.RicControlRequest{
		Hdr: &e2sm.RicControlHeader{
			MessageType: sb.MessageType_CELL_CONFIG_REQUEST,
		},
		Msg: &e2sm.RicControlMessage{
			S: &e2sm.RicControlMessage_CellConfigRequest{
				CellConfigRequest: &sb.CellConfigRequest{
					Ecgi: &sb.ECGI{
						PlmnId: s.Ecgi.GetPlmnId(),
						Ecid:   s.Ecgi.GetEcid()},
				},
			},
		},
	}

	err = stream.Send(&cellConfigReq)
	if err != nil {
		waitc <- err
		return
	}

	for {
		req := <-s.ricControlRequestChan
		err = stream.Send(&req)
		if err != nil {
			errors <- err
			return
		}
	}
}

func (s *Session) ricControlResponse(update *e2ap.RicControlResponse) {
	msgType := update.GetHdr().GetMessageType()
	switch msgType {
	case sb.MessageType_CELL_CONFIG_REPORT:
		msg := update.GetMsg().GetCellConfigReport()
		log.Infof("%s CellConfigReport plmnid:%s, ecid:%s", s.Device.Address, msg.GetEcgi().GetPlmnId(), msg.GetEcgi().GetEcid())
	default:
		log.Fatalf("%s ControlResponse has unexpected type %T", s.Device.Address, msgType)
	}
	s.ricControlResponseChan <- *update
}

func (s *Session) ricSubscribe(client e2ap.E2APClient, errors chan<- error) {
	stream, err := client.RicSubscription(context.Background())
	if err != nil {
		errors <- err
		return
	}

	waitc := make(chan error)
	go func() {
		for {
			ind, err := stream.Recv()
			if err != nil {
				close(waitc)
				return
			}
			log.Infof("%s indication messageType %d", s.Device.Address, ind.GetHdr().GetMessageType())
			msgType := ind.GetHdr().GetMessageType()
			switch msgType {
			case sb.MessageType_RADIO_MEAS_REPORT_PER_UE:
				if s.EnableMetrics {
					tT := time.Now()
					go s.addReceivedHOEventMeasuredRIC(ind, tT) // to monitor HO delay]
				}
				s.telemetryIndications <- *ind
			case sb.MessageType_UE_ADMISSION_REQUEST, sb.MessageType_UE_RELEASE_IND:
				s.controlIndications <- *ind
			default:
				log.Fatalf("%s indication has unexpected type %T", s.Device.Address, msgType)
			}
		}
	}()
}

func (s *Session) addReceivedHOEventMeasuredRIC(update *e2ap.RicIndication, tT time.Time) {
	servStationID := update.GetMsg().GetRadioMeasReportPerUE().GetEcgi()
	numNeighborCells := len(update.GetMsg().GetRadioMeasReportPerUE().GetRadioReportServCells())
	bestStationID := update.GetMsg().GetRadioMeasReportPerUE().GetRadioReportServCells()[0].GetEcgi()
	bestCQI := update.GetMsg().GetRadioMeasReportPerUE().GetRadioReportServCells()[0].GetCqiHist()[0]

	for i := 1; i < numNeighborCells; i++ {
		tmpCQI := update.GetMsg().GetRadioMeasReportPerUE().GetRadioReportServCells()[i].GetCqiHist()[0]
		if bestCQI < tmpCQI {
			bestStationID = update.GetMsg().GetRadioMeasReportPerUE().GetRadioReportServCells()[i].GetEcgi()
			bestCQI = tmpCQI
		}
	}

	if servStationID.GetEcid() != bestStationID.GetEcid() || servStationID.GetPlmnId() != bestStationID.GetPlmnId() {
		tmpHOEventMeasuredRIC := HOEventMeasuredRIC{
			Timestamp:  tT,
			Crnti:      update.GetMsg().GetRadioMeasReportPerUE().GetCrnti(),
			DestPlmnID: bestStationID.GetPlmnId(),
			DestECID:   bestStationID.GetEcid(),
		}
		key := fmt.Sprintf("%s:%s:%s:%s:%s", update.GetMsg().GetRadioMeasReportPerUE().GetCrnti(), servStationID.GetPlmnId(), servStationID.GetEcid(), bestStationID.GetPlmnId(), bestStationID.GetEcid())
		MutexMapHOEventMeasuredRIC.Lock()
		_, f := MapHOEventMeasuredRIC[key]
		if !f {
			MapHOEventMeasuredRIC[key] = tmpHOEventMeasuredRIC
		}
		MutexMapHOEventMeasuredRIC.Unlock()
	}
}

func (s *Session) updateHOEventMeasuredRIC(req e2ap.RicControlRequest, tC time.Time) {
	key := fmt.Sprintf("%s:%s:%s:%s:%s", req.GetMsg().GetHORequest().GetCrnti(),
		req.GetMsg().GetHORequest().GetEcgiS().GetPlmnId(),
		req.GetMsg().GetHORequest().GetEcgiS().GetEcid(),
		req.GetMsg().GetHORequest().GetEcgiT().GetPlmnId(),
		req.GetMsg().GetHORequest().GetEcgiT().GetEcid())
	MutexMapHOEventMeasuredRIC.Lock()
	v, f := MapHOEventMeasuredRIC[key]
	if !f {
		log.Warnf("Telemetry message is missing when calculating HO latency - %s", key)
	} else {
		v.ElapsedTime = tC.Sub(v.Timestamp).Microseconds()
		ChanHOEvent <- v
		delete(MapHOEventMeasuredRIC, key)
	}
	MutexMapHOEventMeasuredRIC.Unlock()
}
