package control

import (
	"github.com/atomix/go-client/pkg/client/map"
	"github.com/onosproject/onos-ric/api/sb/e2ap"
	messagestore "github.com/onosproject/onos-ric/pkg/store/message"
	"io"
)

type ID string

// Revision is a message revision
type Revision uint64

// Event is a store event
type Event struct {
	// Type is the event type
	Type EventType
	// ID is the message identifier
	ID ID
	// Message is the event message
	Message e2ap.RicControlRequest
}

// EventType is a store event type
type EventType messagestore.EventType

const (
	EventNone   EventType = EventType(messagestore.EventNone)
	EventInsert EventType = EventType(messagestore.EventInsert)
	EventUpdate EventType = EventType(messagestore.EventUpdate)
	EventDelete EventType = EventType(messagestore.EventDelete)
)

// Store is interface for store
type Store interface {
	io.Closer

	// Gets a message based on a given ID
	Get(ID, ...GetOption) (*e2ap.RicControlRequest, error)

	// Puts a message to the store
	Put(ID, *e2ap.RicControlRequest, ...PutOption) error

	// Removes a message from the store
	Delete(ID, ...DeleteOption) error

	// List all of the last up to date messages
	List(ch chan<- e2ap.RicControlRequest) error

	// Watch watches the store for changes
	Watch(ch chan<- Event, opts ...WatchOption) error

	// Clear deletes all messages from the store
	Clear() error
}

// GetOption is a message store get option
type GetOption interface {
	getGetOption() _map.GetOption
}

// PutOption is a message store put option
type PutOption interface {
	getPutOption() _map.PutOption
}

// DeleteOption is a message store delete option
type DeleteOption interface {
	getRemoveOption() _map.RemoveOption
}

// WatchOption is a store watch option
type WatchOption interface {
	getWatchOption() _map.WatchOption
}

// WithReplay returns a watch option that replays existing messages
func WithReplay() WatchOption {
	return &replayOption{
		replay: true,
	}
}

type replayOption struct {
	replay bool
}

func (o *replayOption) getWatchOption() _map.WatchOption {
	if o.replay {
		return _map.WithReplay()
	}
	return nil
}
