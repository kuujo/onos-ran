// Code generated by onos-ric-generate. DO NOT EDIT.

package {{ .Impl.Package.Name }}

import (
	"github.com/gogo/protobuf/proto"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-ric/api/store/message"
	"github.com/onosproject/onos-ric/pkg/config"
	messagestore "github.com/onosproject/onos-ric/pkg/store/message"
	timestore "github.com/onosproject/onos-ric/pkg/store/time"
	{{ .DataType.Package.Alias }} {{ .DataType.Package.Path | quote }}
)

var log = logging.GetLogger("store", {{ .Impl.Package.Name | quote }})

func NewDistributedStore(config config.Config, timeStore timestore.Store) ({{ .Interface.Type.Name }}, error) {
	log.Info("Creating distributed store")
	messageStore, err := messagestore.NewDistributedStore({{ .Impl.Storage.Primitive | quote }}, config, {{ .Impl.Storage.Type | lower | quote }}, timeStore)
	if err != nil {
		return nil, err
	}
	return &{{ .Impl.Type.Name }}{
		messageStore: messageStore,
	}, nil
}

func NewLocalStore(timeStore timestore.Store) ({{ .Interface.Type.Name }}, error) {
	messageStore, err := messagestore.NewLocalStore({{ .Impl.Storage.Primitive | quote }}, timeStore)
	if err != nil {
		return nil, err
	}
	return &{{ .Impl.Type.Name }}{
		messageStore: messageStore,
	}, nil
}

var _ {{ .Interface.Type.Name }} = &{{ .Impl.Type.Name }}{}

type {{ .Impl.Type.Name }} struct {
	messageStore messagestore.Store
}

func (s *{{ .Impl.Type.Name }}) Get(id ID, opts ...GetOption) (*{{ .DataType.Package.Alias }}.{{ .DataType.Name }}, error) {
	messageOpts := make([]messagestore.GetOption, len(opts))
	for i, opt := range opts {
		messageOpts[i] = opt
	}
	entry, err := s.messageStore.Get(messagestore.Key(id), messageOpts...)
	if err != nil {
		return nil, err
	} else if entry == nil {
		return nil, nil
	}
	return decode(entry)
}

func (s *{{ .Impl.Type.Name }}) Put(message *{{ .DataType.Package.Alias }}.{{ .DataType.Name }}) error {
	entry, err := encode(message)
	if err != nil {
		return err
	}
	return s.messageStore.Put(getKey(message), entry)
}

func (s *{{ .Impl.Type.Name }}) Delete(id ID, opts ...DeleteOption) error {
	messageOpts := make([]messagestore.DeleteOption, len(opts))
	for i, opt := range opts {
		messageOpts[i] = opt
	}
	return s.messageStore.Delete(messagestore.Key(id), messageOpts...)
}

func (s *{{ .Impl.Type.Name }}) List(ch chan<- {{ .DataType.Package.Alias }}.{{ .DataType.Name }}) error {
	entryCh := make(chan message.MessageEntry)
	if err := s.messageStore.List(entryCh); err != nil {
		return err
	}
	go func() {
		defer close(ch)
		for entry := range entryCh {
			if message, err := decode(&entry); err == nil {
				ch <- *message
			}
		}
	}()
	return nil
}

func (s *{{ .Impl.Type.Name }}) Watch(ch chan<- {{ .DataType.Package.Alias }}.{{ .DataType.Name }}, opts ...WatchOption) error {
	messageOpts := make([]messagestore.WatchOption, len(opts))
	for i, opt := range opts {
		messageOpts[i] = opt
	}

	watchCh := make(chan message.MessageEntry)
	if err := s.messageStore.Watch(watchCh, messageOpts...); err != nil {
		return err
	}
	go func() {
		defer close(ch)
		for entry := range watchCh {
			if message, err := decode(&entry); err == nil {
				ch <- *message
			}
		}
	}()
	return nil
}

func (s *{{ .Impl.Type.Name }}) Clear() error {
	return s.messageStore.Clear()
}

func (s *{{ .Impl.Type.Name }}) Close() error {
	return s.messageStore.Close()
}

func decode(message *message.MessageEntry) (*{{ .DataType.Package.Alias }}.{{ .DataType.Name }}, error) {
	m := &{{ .DataType.Package.Alias }}.{{ .DataType.Name }}{}
	if err := proto.Unmarshal(message.Bytes, m); err != nil {
		return nil, err
	}
	return m, nil
}

func encode(m *{{ .DataType.Package.Alias }}.{{ .DataType.Name }}) (*message.MessageEntry, error) {
	bytes, err := proto.Marshal(m)
	if err != nil {
		return nil, err
	}
	return &message.MessageEntry{
		Bytes: bytes,
	}, nil
}

func toMessageRevision(revision Revision) messagestore.Revision {
	return messagestore.Revision{
		Term:      messagestore.Term(revision.Term),
		Timestamp: messagestore.Timestamp(revision.Timestamp),
	}
}

func getKey(message *{{ .DataType.Package.Alias }}.{{ .DataType.Name }}) messagestore.Key {
	return messagestore.Key(message.GetID())
}