package recipes

import (
	"github.com/influx6/box/recipes/exec"
	"github.com/influx6/faux/context"
	"github.com/influx6/faux/metrics"
	"github.com/influx6/faux/ops"
)

// GenericRunner will run necessary commands to update apt-get for a debian/ubuntu system.
type GenericRunner struct {
	Command   string
	DoWithCmd exec.CommanderOption
}

// Exec executes giving recipe for building giving docker image for the provided binary.
func (gn GenericRunner) Exec(ctx context.CancelContext, m metrics.Metrics) error {
	cmd := exec.New(exec.Command(gn.Command), exec.Async())

	if gn.DoWithCmd != nil {
		gn.DoWithCmd(cmd)
	}

	return cmd.Exec(ctx, m)
}

//===============================================================================================================

// MultiRunner will run necessary commands to update apt-get for a debian/ubuntu system.
type MultiRunner struct {
	Pre  []ops.Op
	Then ops.Op
	Post []ops.Op
}

// Exec executes giving spells in a before-now-after sequence.
func (gn MultiRunner) Exec(ctx context.CancelContext, m metrics.Metrics) error {
	for _, item := range gn.Pre {
		if err := item.Exec(ctx, m); err != nil {
			return err
		}
	}

	if gn.Then != nil {
		if err := gn.Then.Exec(ctx, m); err != nil {
			return err
		}
	}

	for _, item := range gn.Post {
		if err := item.Exec(ctx, m); err != nil {
			return err
		}
	}

	return nil
}
