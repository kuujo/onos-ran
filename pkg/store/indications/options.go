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

package indications

// SubscribeOption is a message store subscribe option
type SubscribeOption interface {
	applySubscribe(options *subscribeOptions)
}

// subscribeOptions is a struct of message store subscribe options
type subscribeOptions struct {
	replay bool
}

// WithReplay returns a subscribe option that replays existing messages
func WithReplay() SubscribeOption {
	return &subscribeReplayOption{
		replay: true,
	}
}

// subscribeReplayOption is an option for configuring whether to replay message on subscribe calls
type subscribeReplayOption struct {
	replay bool
}

func (o *subscribeReplayOption) applySubscribe(options *subscribeOptions) {
	options.replay = o.replay
}