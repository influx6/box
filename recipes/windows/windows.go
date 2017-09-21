package windows

import (
	"errors"

	"github.com/influx6/box"
	"github.com/influx6/faux/context"
	"github.com/influx6/faux/ops"
)

var (
	_ = box.RegisterJSON("windows", func() ops.Op {
		return &WindowProvisioner{}
	})
)

// WindowProvisioner handles the implementation details for setuping on
// box on a system.
type WindowProvisioner struct {
	OSName string
}

// Exec implements the box.Spell system.
func (dw *WindowProvisioner) Exec(ctx context.CancelContext) error {
	return errors.New("window(adm64/i386) is not supported yet")
}
