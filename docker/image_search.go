package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/influx6/faux/context"
	"github.com/influx6/faux/ops"
	"github.com/moby/moby/api/types/registry"
	"github.com/moby/moby/client"
)

// ImageSearch returns a new ImageSearchOp instance to be executed on the client.
func (d *DockerCaster) ImageSearch(searchOps types.ImageSearchOptions) (*ImageSearchOp, error) {
	var spell ImageSearchOp

	spell.searchOps = searchOps

	return &spell, nil
}

// ImageSearchOptions defines a function type to modify internal fields of the ImageSearchOp.
type ImageSearchOptions func(*ImageSearchOp)

// ImageSearchResponseCallback defines a function type for ImageSearchOp response.
type ImageSearchResponseCallback func([]registry.SearchResult) error

// ImageSearchOp defines a structure which implements the Op interface
// for executing of docker based commands for ImageSearch.
type ImageSearchOp struct {
	client *client.Client

	searchOps types.ImageSearchOptions
}

// Op returns a object implementing the ops.Op interface.
func (cm *ImageSearchOp) Op(callback ImageSearchResponseCallback) ops.Op {
	return &onceImageSearchOp{spell: cm, callback: cb}
}

type onceImageSearchOp struct {
	callback ImageSearchResponseCallback
	spell    *ImageSearchOp
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceImageSearchOp) Exec(ctx context.CancelContext) error {
	return cm.spell.Exec(ctx, cm.callback)
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *ImageSearchOp) Exec(ctx context.CancelContext, callback ImageSearchResponseCallback) error {
	// Execute client ImageSearch method.
	ret0, err := cm.client.ImageSearch(cm.searchOps)
	if err != nil {
		return err
	}

	if callback != nil {
		return callback(ret0)
	}

	return nil
}
