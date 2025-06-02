package lite

import (
	"github.com/hedzr/cmdr/v2/cli"
)

// PrepareApp will be used for cmdr-tests
func PrepareApp(appName, desc string, opts ...cli.Opt) func(adders ...cli.CmdAdder) (app cli.App) {
	return func(adders ...cli.CmdAdder) (app cli.App) {
		cs := Create(appName, version, author, desc, opts...)
		cs.With(func(app cli.App) {
			// another way to disable `cmdr.WithForceDefaultAction(true)` is using
			// env-var FORCE_RUN=1 (builtin already).
			app.Flg("no-default").
				Description("disable force default action").
				// Group(cli.UnsortedGroup).
				OnMatched(func(f *cli.Flag, position int, hitState *cli.MatchState) (err error) {
					if b, ok := hitState.Value.(bool); ok {
						// disable/enable the final state about 'force default action'
						f.Set().Set("app.force-default-action", b)
					}
					return
				}).
				Build()

			app.Flg("dry-run", "n").
				Default(false).
				Description("run all but without committing").
				Build()

			app.Flg("wet-run", "w").
				Default(false).
				Description("run all but with committing").
				Build() // no matter even if you're adding the duplicated one.

			for _, adder := range adders {
				adder.Add(app)
			}
		})

		cs.WithAdders(adders...)

		app = cs.Build()

		return
	}
}
