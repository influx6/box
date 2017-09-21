package linux

import (
	"fmt"

	"github.com/influx6/box"
	"github.com/influx6/box/recipes/exec/osinfo"
	"github.com/influx6/faux/context"
	"github.com/influx6/faux/ops"
)

var (
	_ = box.RegisterJSON("linux", func() ops.Op {
		return &LinuxProvisioner{}
	})
)

// LinuxProvisioner handles the implementation details for setuping on
// box on a system.
type LinuxProvisioner struct {
	OSName string
}

// Exec implements the box.Spell system.
func (dw *LinuxProvisioner) Exec(ctx context.CancelContext) error {
	osInfo, err := osinfo.OSInfo(ctx)
	if err != nil {
		return err
	}

	fmt.Printf("%#v\n", osInfo)
	return nil
}
