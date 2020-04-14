package control

import (
	"context"
	"github.com/atomix/go-client/pkg/client/map"
	"github.com/atomix/go-client/pkg/client/primitive"
	"github.com/atomix/go-client/pkg/client/util/net"
	"github.com/gogo/protobuf/proto"
	"github.com/onosproject/onos-lib-go/pkg/atomix"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-ric/api/sb/e2ap"
	"github.com/onosproject/onos-ric/pkg/config"
	"time"
)

var log = logging.GetLogger("store", "control")

const primitiveName = "control"
const requestTimeout = 15 * time.Second
const retryInterval = time.Second

// NewDistributedStore creates a new distributed message store
func NewDistributedStore(config config.Config) (Store, error) {
	log.Info("Creating distributed message store")
	database, err := atomix.GetDatabase(config.Atomix, config.Atomix.GetDatabase(atomix.DatabaseTypeCache))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	messages, err := database.GetMap(ctx, primitiveName)
	if err != nil {
		return nil, err
	}
	return &store{
		dist: messages,
	}, nil
}

// NewLocalStore returns a new local message store
func NewLocalStore() (Store, error) {
	_, address := atomix.StartLocalNode()
	return newLocalStore(address)
}

// newLocalStore creates a new local message store
func newLocalStore(address net.Address) (Store, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	session, err := primitive.NewSession(ctx, primitive.Partition{ID: 1, Address: address})
	if err != nil {
		return nil, err
	}
	primitiveName := primitive.Name{
		Namespace: "local",
		Name:      primitiveName,
	}
	messages, err := _map.New(context.Background(), primitiveName, []*primitive.Session{session})
	if err != nil {
		return nil, err
	}

	return &store{
		dist: messages,
	}, nil
}

var _ Store = &store{}

type store struct {
	dist _map.Map
}

func (s *store) Get(id ID, opts ...GetOption) (*e2ap.RicControlRequest, error) {
	getOpts := make([]_map.GetOption, len(opts))
	for i, opt := range opts {
		getOpts[i] = opt.getGetOption()
	}
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	entry, err := s.dist.Get(ctx, string(id), getOpts...)
	if err != nil {
		return nil, err
	} else if entry == nil {
		return nil, nil
	}
	return decode(entry.Value)
}

func (s *store) Put(id ID, message *e2ap.RicControlRequest, opts ...PutOption) error {
	putOpts := make([]_map.PutOption, len(opts))
	for i, opt := range opts {
		putOpts[i] = opt.getPutOption()
	}
	bytes, err := encode(message)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	_, err = s.dist.Put(ctx, string(id), bytes, putOpts...)
	return err
}

func (s *store) Delete(id ID, opts ...DeleteOption) error {
	removeOpts := make([]_map.RemoveOption, len(opts))
	for i, opt := range opts {
		removeOpts[i] = opt.getRemoveOption()
	}
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	_, err := s.dist.Remove(ctx, string(id), removeOpts...)
	return err
}

func (s *store) List(ch chan<- e2ap.RicControlRequest) error {
	entryCh := make(chan *_map.Entry)
	if err := s.dist.Entries(context.Background(), entryCh); err != nil {
		return err
	}
	go func() {
		defer close(ch)
		for entry := range entryCh {
			if message, err := decode(entry.Value); err == nil {
				ch <- *message
			}
		}
	}()
	return nil
}

func (s *store) Watch(ch chan<- Event, opts ...WatchOption) error {
	messageOpts := make([]_map.WatchOption, len(opts))
	for i, opt := range opts {
		messageOpts[i] = opt.getWatchOption()
	}

	watchCh := make(chan *_map.Event)
	if err := s.dist.Watch(context.Background(), watchCh, messageOpts...); err != nil {
		return err
	}
	go func() {
		defer close(ch)
		for event := range watchCh {
			if message, err := decode(event.Entry.Value); err == nil {
				event := Event{
					Type:    EventType(event.Type),
					ID:      GetID(message),
					Message: *message,
				}
				ch <- event
			}
		}
	}()
	return nil
}

func (s *store) Clear() error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	return s.dist.Clear(ctx)
}

func (s *store) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	return s.dist.Close(ctx)
}

func decode(bytes []byte) (*e2ap.RicControlRequest, error) {
	message := &e2ap.RicControlRequest{}
	if err := proto.Unmarshal(bytes, message); err != nil {
		return nil, err
	}
	return message, nil
}

func encode(message *e2ap.RicControlRequest) ([]byte, error) {
	return proto.Marshal(message)
}
