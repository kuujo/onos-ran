// Code generated by onos-ric-generate. DO NOT EDIT.

package control

import (
	e2ap "github.com/onosproject/onos-ric/api/sb/e2ap"
	messagestore "github.com/onosproject/onos-ric/pkg/store/message"
	"io"
)

type ID string

// Revision is a message revision
type Revision struct {
	Term      Term
	Timestamp Timestamp
}

type Term uint64

type Timestamp uint64

// Store is interface for store
type Store interface {
	io.Closer

	// Gets a message based on a given ID
	Get(ID, ...GetOption) (*e2ap.RicControlResponse, error)

	// Puts a message to the store
	Put(*e2ap.RicControlResponse) error

	// Removes a message from the store
	Delete(ID, ...DeleteOption) error

	// List all of the last up to date messages
	List(ch chan<- e2ap.RicControlResponse) error

	// Watch watches messages
	Watch(ch chan<- e2ap.RicControlResponse, opts ...WatchOption) error

	// Clear deletes all messages from the store
	Clear() error
}

// GetOption is a message store get option
type GetOption messagestore.GetOption

// WithRevision returns a GetOption that ensures the retrieved entry is newer than the given revision
func WithRevision(revision Revision) GetOption {
	return messagestore.WithRevision(toMessageRevision(revision))
}

// DeleteOption is a message store delete option
type DeleteOption messagestore.DeleteOption

// IfRevision returns a delete option that deletes the message only if its revision matches the given revision
func IfRevision(revision Revision) DeleteOption {
	return messagestore.IfRevision(toMessageRevision(revision))
}

// WatchOption is a store watch option
type WatchOption messagestore.WatchOption

// WithReplay returns a watch option that replays existing messages
func WithReplay() WatchOption {
	return messagestore.WithReplay()
}
