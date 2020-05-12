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

// GetOption is a message store get option
type GetOption interface {
	applyGet(options *getOptions)
}

// getOptions is a struct of message get options
type getOptions struct{}

// PutOption is a message store put option
type PutOption interface {
	applyPut(options *putOptions)
}

// putOptions is a struct of message put options
type putOptions struct {
	revision Revision
}

// DeleteOption is a message store delete option
type DeleteOption interface {
	applyDelete(options *deleteOptions)
}

// deleteOptions is a struct of message delete options
type deleteOptions struct {
	revision Revision
}

// UpdateOption is an option supporting both Put and Delete
type UpdateOption interface {
	PutOption
	DeleteOption
}

// IfRevision returns a delete option that deletes the message only if its revision matches the given revision
func IfRevision(revision Revision) UpdateOption {
	return &updateRevisionOption{
		revision: revision,
	}
}

type updateRevisionOption struct {
	revision Revision
}

func (o *updateRevisionOption) applyPut(options *putOptions) {
	options.revision = o.revision
}

func (o *updateRevisionOption) applyDelete(options *deleteOptions) {
	options.revision = o.revision
}

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
