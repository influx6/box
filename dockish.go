package dockish

import (
	"errors"
	"sync"

	"github.com/moby/moby/client"
)

// errors
var (
	ErrNoFilesystemProvided       = errors.New("No Filesystem was provided")
	ErrNoDockerClientProvided     = errors.New("No docker client provided")
	ErrNoImageBuildOptionProvided = errors.New("No types.ImageBuildOptions was provided")
)

// CancelContext defines a type which provides Done signal for cancelling operations.
type CancelContext interface {
	Done() <-chan struct{}
}

// CnclContext defines a struct to implement the CancelContext.
type CnclContext struct {
	close chan struct{}
	once  sync.Once
}

// NewCnclContext returns a new instance of the CnclContext.
func NewCnclContext() *CnclContext {
	return &CnclContext{close: make(chan struct{})}
}

// Close closes the internal channel of the contxt
func (cn *CnclContext) Close() {
	cn.once.Do(func() {
		close(cn.close)
	})
}

// Done returns a channel to signal ending of op.
// It implements the CancelContext.
func (cn *CnclContext) Done() <-chan struct{} {
	return cn.close
}

// Spell defines an interface which expose an exec method.
type Spell interface {
	Exec(CancelContext) error
}

// DockerCaster provides the central structure that provides methods for executing different
// operations on a docker client. It instantiates the docker client and
// passes all necessary details to different spells.
type DockerCaster struct {
	client *client.Client
}

// New returns a new instance of a DockerCaster.
func New(client *client.Client) *DockerCaster {
	return &DockerCaster{
		client: client,
	}
}
