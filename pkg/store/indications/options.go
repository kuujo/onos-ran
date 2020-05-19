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

// LookupOption is an indication store lookup option
type LookupOption interface {
	applyLookup(options *lookupOptions)
}

// lookupOptions is a struct of indication lookup options
type lookupOptions struct{}

// RecordOption is an indication store record option
type RecordOption interface {
	applyRecord(options *recordOptions)
}

// recordOptions is a struct of indication record options
type recordOptions struct{}

// DiscardOption is an indication store discard option
type DiscardOption interface {
	applyDiscard(options *discardOptions)
}

// discardOptions is a struct of indication discard options
type discardOptions struct{}

// WatchOption is a message store watch option
type WatchOption interface {
	applyWatch(options *watchOptions)
}

// watchOptions is a struct of message store watch options
type watchOptions struct {
	replay bool
}

// WithReplay returns a watch option that replays existing messages
func WithReplay() WatchOption {
	return &watchReplayOption{
		replay: true,
	}
}

// watchReplayOption is an option for configuring whether to replay message on watch calls
type watchReplayOption struct {
	replay bool
}

func (o *watchReplayOption) applyWatch(options *watchOptions) {
	options.replay = o.replay
}
