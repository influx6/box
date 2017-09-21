package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/influx6/faux/context"
	"github.com/influx6/faux/ops"
	"github.com/moby/moby/client"
)

// NetworkInspect returns a new NetworkInspectOp instance to be executed on the client.
func (d *DockerCaster) NetworkInspect(netOp types.NetworkInspectOptions) (*NetworkInspectOp, error) {
	var spell NetworkInspectOp

	spell.netOp = netOp

	return &spell, nil
}

// NetworkInspectOptions defines a function type to modify internal fields of the NetworkInspectOp.
type NetworkInspectOptions func(*NetworkInspectOp)

// NetworkInspectResponseCallback defines a function type for NetworkInspectOp response.
type NetworkInspectResponseCallback func(types.NetworkResource) error

// NetworkInspectOp defines a structure which implements the Op interface
// for executing of docker based commands for NetworkInspect.
type NetworkInspectOp struct {
	client *client.Client

	netOp types.NetworkInspectOptions
}

// Op returns a object implementing the ops.Op interface.
func (cm *NetworkInspectOp) Op(callback NetworkInspectResponseCallback) ops.Op {
	return &onceNetworkInspectOp{spell: cm, callback: cb}
}

type onceNetworkInspectOp struct {
	callback NetworkInspectResponseCallback
	spell    *NetworkInspectOp
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceNetworkInspectOp) Exec(ctx context.CancelContext) error {
	return cm.spell.Exec(ctx, cm.callback)
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *NetworkInspectOp) Exec(ctx context.CancelContext, callback NetworkInspectResponseCallback) error {
	// Execute client NetworkInspect method.
	ret0, err := cm.client.NetworkInspect(cm.netOp)
	if err != nil {
		return err
	}

	if callback != nil {
		return callback(ret0)
	}

	return nil
}
