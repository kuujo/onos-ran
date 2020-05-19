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

package e2

import (
	"context"
	"errors"
	"github.com/onosproject/onos-lib-go/pkg/southbound"
	"github.com/onosproject/onos-ric/api/sb"
	"github.com/onosproject/onos-ric/api/sb/e2ap"
	"github.com/onosproject/onos-ric/api/sb/e2sm"
	"github.com/onosproject/onos-ric/pkg/store/device"
	"github.com/onosproject/onos-ric/pkg/store/mastership"
	"google.golang.org/grpc"
	"io"
	"sync"
	"time"
)

// newSession creates a new E2 session
func newSession(device device.Device, election mastership.Election) (*Session, error) {
	session := &Session{
		Device:   &device,
		election: election,
		closeCh:  make(chan struct{}),
	}
	err := session.open()
	if err != nil {
		return nil, err
	}
	return session, nil
}

// Session is an E2 session
type Session struct {
	Device   *device.Device
	election mastership.Election
	state    *mastership.State
	conn     *grpc.ClientConn
	stream   e2ap.E2AP_RicChanClient
	recvCh   chan<- e2ap.RicIndication
	closeCh  chan struct{}
	mu       sync.RWMutex
}

// open opens the E2 session
func (s *Session) open() error {
	ch := make(chan mastership.State)
	err := s.election.Watch(ch)
	if err != nil {
		return err
	}
	go func() {
		connected := false
		state, _ := s.election.GetState()
		if state != nil {
			s.mu.Lock()
			s.state = state
			s.mu.Unlock()
			if state.Master == s.election.NodeID() {
				err := s.connect()
				if err != nil {
					log.Error(err)
				} else {
					connected = true
				}
			}

		}

		for {
			select {
			case state := <-ch:
				s.mu.Lock()
				s.state = &state
				s.mu.Unlock()
				if state.Master == s.election.NodeID() && !connected {
					err := s.connect()
					if err != nil {
						log.Error(err)
					} else {
						connected = true
					}
				} else if state.Master != s.election.NodeID() && connected {
					err := s.disconnect()
					if err != nil {
						log.Error(err)
					} else {
						connected = false
					}
				}
			case <-s.closeCh:
				return
			}
		}
	}()
	return nil
}

// sendRequest sends a request on the session
func (s *Session) sendRequest(request *e2ap.RicControlRequest) error {
	s.mu.RLock()
	state := s.state
	stream := s.stream
	s.mu.RUnlock()

	if state == nil || state.Master != s.election.NodeID() {
		return errors.New("not the master")
	}

	if stream == nil {
		return errors.New("not connected")
	}
	return stream.Send(request)
}

// subscribe subscribes to southbound messages arriving via the session
func (s *Session) subscribe(ch chan<- e2ap.RicIndication) error {
	s.mu.Lock()
	s.recvCh = ch
	s.mu.Unlock()
	return nil
}

// connect connects the session to the device
func (s *Session) connect() error {
	log.Infof("Connecting to device %s", s.Device.Address)
	opts := []grpc.DialOption{
		grpc.WithStreamInterceptor(southbound.RetryingStreamClientInterceptor(100 * time.Millisecond)),
	}
	conn, err := southbound.Connect(context.Background(), s.Device.Address, s.Device.TLS.Cert, s.Device.TLS.Key, opts...)
	if err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	oldConn := s.conn
	if oldConn != nil {
		oldConn.Close()
	}
	s.conn = conn

	client := e2ap.NewE2APClient(conn)
	stream, err := client.RicChan(context.Background())
	if err != nil {
		conn.Close()
		return err
	}
	s.stream = stream

	err = s.setup(stream)
	if err != nil {
		return err
	}

	err = s.l2MeasConfig(stream)
	if err != nil {
		return err
	}

	go s.recvRicIndications(stream)
	return nil
}

// setup sends an initial cell config request
func (s *Session) setup(stream e2ap.E2AP_RicChanClient) error {
	setupRequest := &e2ap.RicControlRequest{
		Hdr: &e2sm.RicControlHeader{
			MessageType: sb.MessageType_CELL_CONFIG_REQUEST,
		},
		Msg: &e2sm.RicControlMessage{
			S: &e2sm.RicControlMessage_CellConfigRequest{
				CellConfigRequest: &sb.CellConfigRequest{},
			},
		},
	}
	return stream.Send(setupRequest)
}

// l2MeasConfig sends an initial l2MeasConfig request
func (s *Session) l2MeasConfig(stream e2ap.E2AP_RicChanClient) error {
	// TODO: Make report interval configurable
	l2MeasConfig := &sb.L2MeasConfig{
		RadioMeasReportPerUe:   sb.L2MeasReportInterval_MS_10,
		RadioMeasReportPerCell: sb.L2MeasReportInterval_MS_10,
		SchedMeasReportPerUe:   sb.L2MeasReportInterval_MS_10,
		SchedMeasReportPerCell: sb.L2MeasReportInterval_MS_10,
		PdcpMeasReportPerUe:    sb.L2MeasReportInterval_MS_10,
	}

	l2MeasConfigRequest := &e2ap.RicControlRequest{
		Hdr: &e2sm.RicControlHeader{
			MessageType: sb.MessageType_L2_MEAS_CONFIG,
		},
		Msg: &e2sm.RicControlMessage{
			S: &e2sm.RicControlMessage_L2MeasConfig{
				L2MeasConfig: l2MeasConfig,
			},
		},
	}
	return stream.Send(l2MeasConfigRequest)
}

// recvRicIndications receives RicIndication messages on the given stream
func (s *Session) recvRicIndications(stream e2ap.E2AP_RicChanClient) {
	for {
		response, err := stream.Recv()
		if err == io.EOF {
			return
		} else if err != nil {
			log.Errorf("An error was received from %s", s.Device.Address, err)
		} else {
			s.mu.RLock()
			ch := s.recvCh
			s.mu.RUnlock()
			if ch != nil {
				ch <- *response
			}
		}
	}
}

// disconnect disconnects the session from the device
func (s *Session) disconnect() error {
	log.Infof("Disconnecting from device %s", s.Device.Address)
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.conn == nil {
		return nil
	}

	if s.stream != nil {
		_ = s.stream.CloseSend()
	}

	err := s.conn.Close()
	s.conn = nil
	s.stream = nil
	return err
}

// Close closes the E2 session
func (s *Session) Close() {
	err := s.disconnect()
	if err != nil {
		log.Error(err)
	}
	close(s.closeCh)
}
