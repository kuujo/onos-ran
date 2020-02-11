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

package nb

import (
	"context"
	"errors"
	"github.com/onosproject/onos-ran/api/nb"
	"github.com/onosproject/onos-test/pkg/onit/env"
	"io"
	"time"
)

// waitForSimulator polls until the simulator is responding properly.
// the can take a while, allow a minute before giving up.
func waitForSimulator() error {
	const sleepPeriodSeconds = 2
	const tries = 30

	for i := 1; i <= tries; i++ {
		client, clientErr := env.RAN().NewRANC1ServiceClient()

		if clientErr != nil {
			return clientErr
		}

		request := &nb.StationListRequest{
			Subscribe: false,
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		stations, stationsErr := client.ListStations(ctx, request)

		if stationsErr != nil {
			cancel()
			return stationsErr
		}

		_, pollError := stations.Recv()
		cancel()
		if pollError != nil && pollError != io.EOF {
			return pollError
		}
		if pollError == nil {
			return nil
		}
		time.Sleep(sleepPeriodSeconds * time.Second)
	}

	return errors.New("simulator never responded properly")
}
