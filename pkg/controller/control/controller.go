package control

import (
	"fmt"
	"github.com/onosproject/onos-ric/pkg/controller"
	"github.com/onosproject/onos-ric/pkg/southbound"
	"github.com/onosproject/onos-ric/pkg/store/control"
	"github.com/onosproject/onos-ric/pkg/store/mastership"
	"strings"
)

const queueSize = 100

// NewController creates a new control request controller
func NewController(mastershipStore mastership.Store, controlStore control.Store, sessionManager *southbound.SessionManager) *controller.Controller {
	c := controller.NewController("RicControlRequest")
	c.Filter(&controller.MastershipFilter{
		Store:    mastershipStore,
		Resolver: &Resolver{},
	})
	c.Partition(&Partitioner{})
	c.Watch(&Watcher{
		controlStore: controlStore,
	})
	c.Reconcile(&Reconciler{
		controlStore:   controlStore,
		sessionManager: sessionManager,
	})
	return c
}

// Partitioner is a WorkPartitioner for control requests
type Partitioner struct{}

// Partition returns the device as a partition key
func (p *Partitioner) Partition(request controller.Request) (controller.PartitionKey, error) {
	parts := strings.Split(string(request.ID), ":")
	return controller.PartitionKey(fmt.Sprintf("%s:%s", parts[1], parts[2])), nil
}

var _ controller.WorkPartitioner = &Partitioner{}

// Resolver is a MastershipResolver that resolves mastership keys from requests
type Resolver struct{}

// Resolve resolves a device ID from a device change ID
func (r *Resolver) Resolve(request controller.Request) (mastership.Key, error) {
	parts := strings.Split(string(request.ID), ":")
	return mastership.Key(fmt.Sprintf("%s:%s", parts[1], parts[2])), nil
}

type Reconciler struct {
	controlStore   control.Store
	sessionManager *southbound.SessionManager
}

func (r *Reconciler) Reconcile(id controller.Request) (controller.Result, error) {
	request, err := r.controlStore.Get(control.ID(id.ID))
	if err != nil {
		return controller.Result{}, err
	} else if request == nil {
		return controller.Result{}, nil
	}

	session, err := r.sessionManager.GetSession(*request.GetHdr().Ecgi)
	if err != nil {
		return controller.Result{}, err
	}
	err = session.SendRicControlRequest(*request)
	if err != nil {
		return controller.Result{}, err
	}
	return controller.Result{}, nil
}

type Watcher struct {
	controlStore control.Store
}

// Start starts the network change watcher
func (w *Watcher) Start(ch chan<- controller.Request) error {
	configCh := make(chan control.Event, queueSize)
	err := w.controlStore.Watch(configCh, control.WithReplay())
	if err != nil {
		return err
	}

	go func() {
		for request := range configCh {
			ch <- controller.Request{
				ID: controller.ID(request.ID),
			}
		}
		close(ch)
	}()
	return nil
}

// Stop stops the network change watcher
func (w *Watcher) Stop() {
}

var _ controller.Watcher = &Watcher{}
