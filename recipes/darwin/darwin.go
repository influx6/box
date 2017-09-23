package darwin

import (
	"errors"

	"github.com/influx6/box"
	"github.com/influx6/faux/context"
	"github.com/influx6/faux/metrics"
	"github.com/influx6/faux/ops"
)

var (
	_ = box.RegisterJSON("darwin", func() ops.Op {
		return &DarwinProvisioner{}
	})
)

// DarwinProvisioner handles the implementation details for setuping on
// box on a system.
type DarwinProvisioner struct {
	OSName string
}

// Exec implements the box.Spell system.
func (dw *DarwinProvisioner) Exec(ctx context.CancelContext, metric metrics.Metrics) error {
	return errors.New("darwin/osx is not supported yet")
}
