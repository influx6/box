package linux

import (
	"fmt"

	"github.com/influx6/box"
	"github.com/influx6/box/recipes/exec/osinfo"
	"github.com/influx6/faux/context"
	"github.com/influx6/faux/metrics"
	"github.com/influx6/faux/ops"

	_ "github.com/influx6/box/recipes/linux/ubuntu"
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
func (dw *LinuxProvisioner) Exec(ctx context.CancelContext, m metrics.Metrics) error {
	info, err := osinfo.OSInfo(ctx)
	if err != nil {
		return err
	}

	provisioner, err := box.CreateWithJSON(fmt.Sprintf("linux/%s", info.ID), map[string]interface{}{
		"os_info": info,
	})
	if err != nil {
		return fmt.Errorf("Linux Distro %q not supported: %+q", info.ID, err)
	}

	if err := provisioner.Exec(ctx, m); err != nil {
		return err
	}

	return nil
}
