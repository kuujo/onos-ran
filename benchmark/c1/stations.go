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

package c1

import (
	"context"
	"io"
	"time"

	"github.com/onosproject/helmit/pkg/benchmark"
	"github.com/onosproject/onos-ric/api/nb"
)

// BenchmarkGetStations tests get stations
func (s *BenchmarkSuite) BenchmarkGetStations(b *benchmark.Benchmark) error {
	request := nb.StationListRequest{Subscribe: false}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	stream, err := s.client.ListStations(ctx, &request)
	if err != nil {
		return err
	}
	for {
		_, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
	}
	return nil
}
