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
	"sync"
	"time"

	topodevice "github.com/onosproject/onos-topo/api/device"

	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-lib-go/pkg/southbound"

	"github.com/onosproject/onos-ric/api/sb"
	e2ap "github.com/onosproject/onos-ric/api/sb/e2ap"
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
	EndPoint sb.Endpoint
	Ecgi     sb.ECGI
	client   e2ap.E2APClient

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
func NewSession(ecgi sb.ECGI, endPoint sb.Endpoint) (*Session, error) {
	log.Infof("Creating Session for %v at %s", ecgi, endPoint)

	return &Session{
		EndPoint:               endPoint,
		Ecgi:                   ecgi,
		ricControlRequestChan:  make(chan e2ap.RicControlRequest),
		ricControlResponseChan: make(chan e2ap.RicControlResponse),
		controlIndications:     make(chan e2ap.RicIndication),
		telemetryIndications:   make(chan e2ap.RicIndication),
	}, nil
}

// Run starts the southbound control loop.
func (s *Session) Run(tls topodevice.TlsConfig, creds topodevice.Credentials) {
	MapHOEventMeasuredRIC = make(map[string]HOEventMeasuredRIC)
	go s.manageConnections(tls, creds)
	go s.recvTelemetryUpdates()
	go s.recvUpdates()
}

func (s *Session) recvTelemetryUpdates() {
	for update := range s.telemetryIndications {
		s.TelemetryUpdateHandlerFunc(update)
	}
}

func (s *Session) recvUpdates() {
	for update := range s.controlIndications {
		s.ControlUpdateHandlerFunc(update)
	}
}

// SendRicControlRequest sends the specified RicControlRequest on the control channel
func (s *Session) SendRicControlRequest(req e2ap.RicControlRequest) error {

	req.GetHdr().Ecgi = &s.Ecgi

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

func (s *Session) manageConnections(tls topodevice.TlsConfig, creds topodevice.Credentials) {
	for {
		// Attempt to create connection to the simulator
		log.Infof("Connecting to simulator...%s with context", s.EndPoint)
		opts := []grpc.DialOption{
			grpc.WithStreamInterceptor(southbound.RetryingStreamClientInterceptor(100 * time.Millisecond)),
		}
		connection, err := southbound.Connect(context.Background(), string(s.EndPoint), tls.GetCert(), tls.GetKey(), opts...)
		if err == nil {
			// If successful, manage this connection and don't return until it is
			// no longer valid and all related resources have been properly cleaned-up.
			s.manageConnection(connection)
		}
		time.Sleep(5 * time.Second)
	}
}

func (s *Session) manageConnection(connection *grpc.ClientConn) {
	// Offer the telemetry and control surfaces to the E2 devices
	s.client = e2ap.NewE2APClient(connection)
	if s.client == nil {
		return
	}
	log.Infof("Connected to simulator on %s", s.EndPoint)
	// Setup coordination channel
	errors := make(chan error)

	go s.ricControl(errors)

	go s.ricSubscribe(errors)

	// Wait for the first error on the coordination channel.
	for i := 0; i < 2; i++ {
		<-errors
	}

	// Clean-up the connection, forcing other consumers to terminate and clean-up.
	_ = connection.Close()

	// FIXME: should be done separately
	//close(s.responses)
	//close(s.updates)

	close(errors)
	log.Info("Disconnected from simulator %s", s.EndPoint)
}

func (s *Session) ricControl(errors chan error) {
	stream, err := s.client.RicControl(context.Background())
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
		log.Infof("%s CellConfigReport plmnid:%s, ecid:%s", s.EndPoint, msg.GetEcgi().GetPlmnId(), msg.GetEcgi().GetEcid())
	default:
		log.Fatalf("%s ControlResponse has unexpected type %T", s.EndPoint, msgType)
	}
	s.ricControlResponseChan <- *update
}

func (s *Session) ricSubscribe(errors chan error) {
	stream, err := s.client.RicSubscription(context.Background())
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
			log.Infof("%s indication messageType %d", s.EndPoint, ind.GetHdr().GetMessageType())
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
				log.Fatalf("%s indication has unexpected type %T", s.EndPoint, msgType)
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
