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
	"io"
	"time"

	"github.com/onosproject/onos-ran/api/sb"
	"github.com/onosproject/onos-ran/pkg/service"
	"google.golang.org/grpc"
	log "k8s.io/klog"
)

// Sessions is responsible for managing connections to and interactions with the RAN southbound.
type Sessions struct {
	Simulator *string
	client    sb.InterfaceServiceClient

	responses chan sb.ControlResponse
	updates   chan sb.ControlUpdate
}

// NewSessions creates a new southbound sessions controller.
func NewSessions() (*Sessions, error) {
	log.Info("Creating Sessions")
	return &Sessions{}, nil
}

// Run starts the southbound control loop.
func (m *Sessions) Run(updates chan sb.ControlUpdate, responses chan sb.ControlResponse) {
	// Kick off a go routine that manages the connection to the simulator
	m.updates = updates
	m.responses = responses
	go m.manageConnections()
}

// SendResponse sends the specified response on the control channel.
func (m *Sessions) SendResponse(response sb.ControlResponse) error {
	m.responses <- response
	return nil
}

func (m *Sessions) manageConnections() {
	for {
		// Attempt to create connection to the simulator
		opts := []grpc.DialOption{
			grpc.WithInsecure(),
			grpc.WithBlock(),
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		log.Infof("Connecting to simulator...%s with context", *m.Simulator)
		connection, err := service.Connect(ctx, *m.Simulator, opts...)
		cancel()
		if err == nil {
			// If successful, manage this connection and don't return until it is
			// no longer valid and all related resources have been properly cleaned-up.
			m.manageConnection(connection)
		}
		time.Sleep(5 * time.Second)
	}
}

func (m *Sessions) manageConnection(connection *grpc.ClientConn) {
	// Offer the telemetry and control surfaces to the E2 devices
	m.client = sb.NewInterfaceServiceClient(connection)
	if m.client == nil {
		return
	}

	log.Info("Connected to simulator")
	// Setup coordination channel
	errors := make(chan error)

	// FIXME: should be done separately
	//m.updates = make(chan sb.ControlUpdate)
	//m.responses = make(chan sb.ControlResponse)

	//go m.handleTelemetry(errors)
	go m.handleControl(errors)

	// Wait for the first error on the coordination channel.
	<-errors

	// Clean-up the connection, forcing other consumers to terminate and clean-up.
	_ = connection.Close()

	// FIXME: should be done separately
	//close(m.responses)
	//close(m.updates)

	close(errors)
	log.Info("Disconnected from simulator")
}

/*
func (m *Sessions) handleTelemetry(errors chan error) {
	stream, err := m.client.SendTelemetry(context.Background(), &sb.TelemetryRequest{})
	if err != nil {
		errors <- err
		return
	}

	waitc := make(chan error)
	go func() {
		for {
			telemetry, err := stream.Recv()
			if err != nil {
				close(waitc)
				return
			}
			m.processTelemetry(telemetry)
		}
	}()

	_ = stream.CloseSend()
	<-waitc
	errors <- io.EOF
}
*/

func (m *Sessions) handleControl(errors chan error) {

	cellConfigRequest := &sb.ControlResponse{
		MessageType: sb.MessageType_CELL_CONFIG_REQUEST,
		S: &sb.ControlResponse_CellConfigRequest{
			CellConfigRequest: &sb.CellConfigRequest{
				Ecgi: &sb.ECGI{PlmnId: "test", Ecid: "test"},
			},
		},
	}

	stream, err := m.client.SendControl(context.Background())
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
			log.Infof("Got messageType %d", update.MessageType)
			m.processControlUpdate(update)
		}
	}()

	if err := stream.Send(cellConfigRequest); err != nil {
		log.Fatalf("Failed to send a note: %v", err)
	}

	go func() {
		for {
			response := <-m.responses
			err := stream.Send(&response)
			if err != nil {
				waitc <- err
				return
			}
		}
	}()

	_ = stream.CloseSend()
	<-waitc
	errors <- io.EOF
}

/*
func (m *Sessions) processTelemetry(msg *sb.TelemetryMessage) {
	// TODO: process the telemetry message
	fmt.Println(msg)
}
*/

func (m *Sessions) processControlUpdate(update *sb.ControlUpdate) {
	switch x := update.S.(type) {
	case *sb.ControlUpdate_CellConfigReport:
		log.Infof("plmnid:%s, ecid:%s", x.CellConfigReport.Ecgi.PlmnId, x.CellConfigReport.Ecgi.Ecid)
	case *sb.ControlUpdate_UEAdmissionRequest:
		log.Infof("plmnid:%s, ecid:%s, crnti:%s", x.UEAdmissionRequest.Ecgi.PlmnId, x.UEAdmissionRequest.Ecgi.Ecid, x.UEAdmissionRequest.Crnti)
	default:
		log.Fatalf("ControlReport has unexpected type %T", x)
	}
	m.updates <- *update
}
