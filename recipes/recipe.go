package recipes

import (
	"github.com/influx6/box"
	"github.com/influx6/faux/context"
)

// MultiRunner will run necessary commands to update apt-get for a debian/ubuntu system.
type MultiRunner struct {
	Pre  []box.Spell
	Then box.Spell
	Post []box.Spell
}

// Exec executes giving spells in a before-now-after sequence.
func (gn MultiRunner) Exec(ctx context.CancelContext) error {
	for _, item := range gn.Pre {
		if err := item.Exec(ctx); err != nil {
			return err
		}
	}

	if gn.Then != nil {
		if err := gn.Then.Exec(ctx); err != nil {
			return err
		}
	}

	for _, item := range gn.Post {
		if err := item.Exec(ctx); err != nil {
			return err
		}
	}

	return nil
}
