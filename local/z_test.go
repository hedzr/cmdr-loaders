package local

import (
	"context"
	"os"
	"testing"

	"github.com/hedzr/cmdr/v2/builder"
	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/cli/examples"
	"github.com/hedzr/cmdr/v2/cli/worker"
)

func cleanApp(t *testing.T, opts ...cli.Opt) (app cli.App, ww cli.Runner) { //nolint:revive
	opts = append(opts,
		// cli.WithHelpScreenWriter(os.Stdout),
		cli.WithDebugScreenWriter(os.Stdout),
		cli.WithForceDefaultAction(true),
		cli.WithTasksBeforeParse(func(ctx context.Context, cmd *cli.Command, runner cli.Runner, extras ...any) (err error) {
			_, _, _ = cmd, runner, extras
			return
		}),
	)
	app = buildDemoApp(opts...)
	ww = postBuild(app)
	if r, ok := ww.(interface{ InitGlobally(ctx context.Context) }); ok {
		r.InitGlobally(context.TODO())
	}
	if !ww.Ready() {
		t.Fatalf("not ready")
	}

	// assertTrue(t, ww.Ready())
	// ww.wrHelpScreen = &discardP{}
	// if helpScreen {
	// 	ww.wrHelpScreen = os.Stdout
	// }
	// ww.wrDebugScreen = os.Stdout
	// ww.ForceDefaultAction = true
	// ww.tasksAfterParse = []taskAfterParse{func(w *workerS, ctx *parseCtx, errParsed error) (err error) { return }} //nolint:revive

	// ww.setArgs([]string{"--debug"})
	// err := ww.Run(withTasksBeforeParse(func(root *cli.RootCommand, runner cli.Runner) (err error) {
	// 	root.SelfAssert()
	// 	t.Logf("root.SelfAssert() passed.")
	// 	return
	// }))
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// ww.TasksBeforeParse = nil
	return
}

func buildDemoApp(opts ...cli.Opt) (app cli.App) { //nolint:revive
	cfg := cli.NewConfig(opts...)

	// cfg := cli.New(cli.WithStore(store.New()))

	w := worker.New(cfg)

	app = builder.New(w).
		Info("demo-app", "0.3.1").
		Author("hedzr")

	b := app.Cmd("jump").
		Description("jump is a demo command").
		Examples(`jump example`).
		Deprecated(`since: v0.9.1`).
		Hidden(false)

	b1 := b.Cmd("to").
		Description("to command").
		Examples(``).
		Deprecated(``).
		Hidden(false).
		OnAction(func(ctx context.Context, cmd *cli.Command, args []string) (err error) { //nolint:revive
			return // handling command action here
		})

	b1.Flg("full", "f").
		Default(false).
		Description("full command").
		Build()
	b1.Build()

	b.Flg("empty", "e").
		With(func(b cli.FlagBuilder) {
			b.Default(false).
				Description("empty command")
		})

	b.Build()

	app.Flg("dry-run", "n").
		Default(false).
		Description("run all but without committing").
		Build()

	app.Flg("wet-run", "w").
		Default(false).
		Description("run all but with committing").
		Build() // no matter even if you're adding the duplicated one.

	b = app.Cmd("consul", "c").
		Description("command set for consul operations")
	b.Flg("data-center", "dc", "datacenter").
		// Description("set data-center").
		Default("dc-1").
		Build()
	b.Build()

	b = app.Cmd("server")
	examples.AttachServerCommand(b)

	b = app.Cmd("kv")
	examples.AttachKvCommand(b)

	b = app.Cmd("ms")
	examples.AttachMsCommand(b)

	b = app.Cmd("more")
	examples.AttachMoreCommandsForTest(b, true)

	b = app.Cmd("display", "da").
		Description("command set for display adapter operations")

	b1 = b.Cmd("voodoo", "vd").
		Description("command set for voodoo operations")
	b1.Flg("data-center", "dc", "datacenter").
		Default("dc-1").
		Build()
	b1.Build()

	b2 := b.Cmd("nvidia", "nv").
		Description("command set for nvidia operations")
	b2.Flg("data-center", "dc", "datacenter").
		Default("dc-1").
		Build()
	b2.Build()

	b3 := b.Cmd("amd", "amd").
		Description("command set for AMD operations")
	b3.Flg("data-center", "dc", "datacenter").
		Default("dc-1").
		Build()
	b3.Build()

	b.Build()

	return
}

func postBuild(app cli.App, args ...string) (ww cli.Runner) { //nolint:revive,unparam
	if sr, ok := app.(interface{ Worker() cli.Runner }); ok {
		ww = sr.Worker()
		if ww1, ok := sr.Worker().(interface {
			SetRoot(root *cli.RootCommand, args []string)
		}); ok {
			if r, ok := app.(interface{ Root() *cli.RootCommand }); ok {
				r.Root().EnsureTree(context.TODO(), app, r.Root())
				ww1.SetRoot(r.Root(), app.Args())
				_ = args
			}
		}
	}
	return
}

type discardP struct{}

func (*discardP) Write([]byte) (n int, err error) { return }

func (*discardP) WriteString(string) (n int, err error) { return }
