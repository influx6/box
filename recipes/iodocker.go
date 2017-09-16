package recipes

import "github.com/influx6/box"

// DockerIORecipe will build giving binary file associated into a docker image.
type DockerIORecipe struct {
	BinaryPath string
}

// Exec executes giving recipe for building giving docker image for the provided binary.
func (d *DockerIORecipe) Exec(ctx box.CancelContext) error {

	return nil
}
