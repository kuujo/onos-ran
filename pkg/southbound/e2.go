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
	"time"

	topodevice "github.com/onosproject/onos-topo/api/device"

	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-lib-go/pkg/southbound"

	"github.com/onosproject/onos-ric/api/sb"
	e2ap "github.com/onosproject/onos-ric/api/sb/e2ap"
	"google.golang.org/grpc"
)

var log = logging.GetLogger("southbound")

// TelemetryUpdateHandler - a function for the session to write back to manager without the import cycle
type TelemetryUpdateHandler = func(sb.TelemetryMessage)

// ControlUpdateHandler - a function for the session to write back to manager without the import cycle
type ControlUpdateHandler = func(sb.ControlUpdate)

// Session is responsible for managing connections to and interactions with the RAN southbound.
type Session struct {
	EndPoint sb.Endpoint
	Ecgi     sb.ECGI
	client   e2ap.E2APClient

	ricControlRequestChan chan e2ap.RicControlRequest
	controlResponses      chan sb.ControlResponse
	controlUpdates        chan sb.ControlUpdate

	telemetryUpdates chan sb.TelemetryMessage

	ControlUpdateHandlerFunc   ControlUpdateHandler
	TelemetryUpdateHandlerFunc TelemetryUpdateHandler
}

// HOEventMeasuredRIC is struct including UE ID and its eNB ID
type HOEventMeasuredRIC struct {
	Timestamp   time.Time
	Crnti       string
	DestPlmnID  string
	DestECID    string
	ElapsedTime int64
}

// ListHOEventMeasuredRIC is the list of HOEvent which is measured at ONOS RIC
var ListHOEventMeasuredRIC []HOEventMeasuredRIC

// NewSession creates a new southbound session controller.
func NewSession(ecgi sb.ECGI, endPoint sb.Endpoint) (*Session, error) {
	log.Infof("Creating Session for %v at %s", ecgi, endPoint)

	return &Session{
		EndPoint:              endPoint,
		Ecgi:                  ecgi,
		ricControlRequestChan: make(chan e2ap.RicControlRequest),
		controlResponses:      make(chan sb.ControlResponse),
		controlUpdates:        make(chan sb.ControlUpdate),
		telemetryUpdates:      make(chan sb.TelemetryMessage),
	}, nil
}

// Run starts the southbound control loop.
func (s *Session) Run(tls topodevice.TlsConfig, creds topodevice.Credentials) {
	go s.manageConnections(tls, creds)
	go s.recvTelemetryUpdates()
	go s.recvUpdates()
}

func (s *Session) recvTelemetryUpdates() {
	for update := range s.telemetryUpdates {
		s.TelemetryUpdateHandlerFunc(update)
	}
}

func (s *Session) recvUpdates() {
	for update := range s.controlUpdates {
		s.ControlUpdateHandlerFunc(update)
	}
}

// SendRicControlRequest sends the specified RicControlRequest on the control channel
func (s *Session) SendRicControlRequest(req e2ap.RicControlRequest) error {
	s.ricControlRequestChan <- req
	return nil
}

func (s *Session) manageConnections(tls topodevice.TlsConfig, creds topodevice.Credentials) {
	for {
		// Attempt to create connection to the simulator
		log.Infof("Connecting to simulator...%s with context", s.EndPoint)
		connection, err := southbound.Connect(context.Background(), string(s.EndPoint), tls.GetCert(), tls.GetKey())
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

	go s.handleControl(errors)

	go s.handleTelemetry(errors)

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
	//	waitc := make(chan error)
	//	go func() {
	//		for {
	//			resp, err := stream.Recv()
	//			if err != nil {
	//				close(waitc)
	//				return
	//			}
	//			log.Infof("Got messageType %d from %s", resp.Hdr.MessageType, s.EndPoint)
	//			s.ricControlResponse(resp)
	//		}
	//	}()
	for {
		req := <-s.ricControlRequestChan
		err := stream.Send(&req)
		if err != nil {
			errors <- err
			return
		}
	}
}

func (s *Session) handleControl(errors chan error) {

	cellConfigRequest := &sb.ControlResponse{
		MessageType: sb.MessageType_CELL_CONFIG_REQUEST,
		S: &sb.ControlResponse_CellConfigRequest{
			CellConfigRequest: &sb.CellConfigRequest{
				Ecgi: &sb.ECGI{PlmnId: "test", Ecid: "test"},
			},
		},
	}

	stream, err := s.client.SendControl(context.Background())
	if err != nil {
		errors <- err
		return
	}

	waitc := make(chan error)
	go func() {
		for {
			update, err := stream.Recv()
			if err != nil {
				close(waitc)
				return
			}
			log.Infof("Got messageType %d from %s", update.MessageType, s.EndPoint)
			s.processControlUpdate(update)
		}
	}()

	if err := stream.Send(cellConfigRequest); err != nil {
		log.Fatalf("Failed to send a note: %v", err)
	}

	go s.ricControl(errors)

	for {
		response := <-s.controlResponses
		err := stream.Send(&response)
		if err != nil {
			waitc <- err
			return
		}
	}
}

/*
func (s *Session) ricControlResponse(update *e2ap.RicControlResponse) {
	// Nothing to do here yet
}
*/

func (s *Session) processControlUpdate(update *sb.ControlUpdate) {
	switch x := update.S.(type) {
	case *sb.ControlUpdate_CellConfigReport:
		log.Infof("%s CellConfigReport plmnid:%s, ecid:%s", s.EndPoint, x.CellConfigReport.Ecgi.PlmnId, x.CellConfigReport.Ecgi.Ecid)
	case *sb.ControlUpdate_UEAdmissionRequest:
		log.Infof("%s UEAdmissionRequest plmnid:%s, ecid:%s, crnti:%s", s.EndPoint, x.UEAdmissionRequest.Ecgi.PlmnId, x.UEAdmissionRequest.Ecgi.Ecid, x.UEAdmissionRequest.Crnti)
	case *sb.ControlUpdate_UEReleaseInd:
		log.Infof("%s UEReleaseInd plmnid:%s, ecid:%s, crnti:%s", s.EndPoint, x.UEReleaseInd.Ecgi.PlmnId, x.UEReleaseInd.Ecgi.Ecid, x.UEReleaseInd.Crnti)
	default:
		log.Fatalf("%s Control update has unexpected type %T", s.EndPoint, x)
	}
	s.controlUpdates <- *update
}

func (s *Session) handleTelemetry(errors chan error) {
	l2MeasConfig := &sb.L2MeasConfig{
		Ecgi: &sb.ECGI{PlmnId: "test", Ecid: "test"},
	}

	stream, err := s.client.SendTelemetry(context.Background(), l2MeasConfig)
	if err != nil {
		errors <- err
		return
	}

	waitc := make(chan error)
	go func() {
		for {
			update, err := stream.Recv()
			if err != nil {
				close(waitc)
				return
			}
			log.Infof("%s telemetry messageType %d", s.EndPoint, update.MessageType)
			s.processTelemetryUpdate(update)
		}
	}()
}

func (s *Session) processTelemetryUpdate(update *sb.TelemetryMessage) {
	switch x := update.S.(type) {
	case *sb.TelemetryMessage_RadioMeasReportPerUE:
	default:
		log.Fatalf("%s Telemetry update has unexpected type %T", s.EndPoint, x)
	}
	s.telemetryUpdates <- *update
}
