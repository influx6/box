package docker

import (
	"errors"

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
