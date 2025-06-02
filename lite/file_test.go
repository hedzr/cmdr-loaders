package lite

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"

	errorsv3 "gopkg.in/hedzr/errors.v3"

	"github.com/hedzr/evendeep"
	"github.com/hedzr/evendeep/diff"
	"github.com/hedzr/store"

	"github.com/hedzr/cmdr/v2/cli"
)

func TestDiff(t *testing.T) {
	a, b := 1, int64(1)
	delta, equal := evendeep.DeepDiff(a, b)
	t.Logf("delta: %v, equal: %v", delta, equal) // ""

	delta, equal = evendeep.DeepDiff([]int{3, 0, 9}, []int{9, 3, 0}, diff.WithSliceOrderedComparison(true))
	t.Logf("delta: %v, equal: %v", delta, equal) // ""

	delta, equal = evendeep.DeepDiff([]int{3, 0, 9}, []int{9, 3, 0}, diff.WithSliceOrderedComparison(false))
	t.Logf("delta: %v, equal: %v", delta, equal) // ""

	delta, equal = evendeep.DeepDiff([]int{3, 0}, []int{9, 3, 0}, diff.WithSliceOrderedComparison(true))
	t.Logf("delta: %v", delta) // "added: [0] = 9\n"

	delta, equal = evendeep.DeepDiff([]int{3, 0}, []int{9, 3, 0})
	t.Logf("delta: %v", delta)
}

func TestWorkerS_Pre(t *testing.T) {
	app, ww := cleanApp(t,
		cli.WithStore(store.New()),
		cli.WithExternalLoaders(NewConfigFileLoader(), NewEnvVarLoader()),
		cli.WithArgs("--debug"),
		cli.WithHelpScreenWriter(&discardP{}),
	)

	// app := buildDemoApp()
	// ww := postBuild(app)
	// ww.InitGlobally()
	// assert.EqualTrue(t, ww.Ready())
	//
	// ww.ForceDefaultAction = true
	// ww.wrHelpScreen = &discardP{}
	// ww.wrDebugScreen = os.Stdout
	// ww.wrHelpScreen = os.Stdout

	// ww.setArgs([]string{"--debug"})
	// ww.Config.Store = store.New()
	// ww.Loaders = []cli.Loader{NewConfigFileLoader(), NewEnvVarLoader()}

	_ = app
	ctx := context.Background()
	err := ww.Run(ctx,
		cli.WithTasksBeforeParse(func(ctx context.Context, cmd cli.Cmd, runner cli.Runner, extras ...any) (err error) {
			cmd.Root().SelfAssert()
			t.Logf("root.SelfAssert() passed. runner = %v", runner)
			return
		}),
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestWorkerS_Pre_v1(t *testing.T) {
	app, ww := cleanApp(t,
		cli.WithStore(store.New()),
		cli.WithExternalLoaders(NewConfigFileLoader(), NewEnvVarLoader()),
		cli.WithArgs("--debug", "-v"),
		cli.WithHelpScreenWriter(&discardP{}),
	)

	// ww.setArgs([]string{"--debug", "-v"})
	// ww.Config.Store = store.New()
	// ww.Loaders = []cli.Loader{loaders.NewConfigFileLoader(), loaders.NewEnvVarLoader()}
	_ = app
	ctx := context.Background()
	err := ww.Run(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestWorkerS_Pre_v3(t *testing.T) {
	app, ww := cleanApp(t,
		cli.WithStore(store.New()),
		cli.WithExternalLoaders(NewConfigFileLoader(), NewEnvVarLoader()),
		cli.WithArgs("--debug", "-vv", "--verbose"),
		cli.WithHelpScreenWriter(&discardP{}),
	)

	// ww.setArgs([]string{"--debug", "-vv", "--verbose"})
	// ww.Config.Store = store.New()
	// ww.Loaders = []cli.Loader{loaders.NewConfigFileLoader(), loaders.NewEnvVarLoader()}
	_ = app
	ctx := context.Background()
	err := ww.Run(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestWorkerS_Parse(t *testing.T) { //nolint:revive
	aTaskBeforeRun := func(ctx context.Context, cmd cli.Cmd, runner cli.Runner, extras ...any) (err error) { return } //nolint:revive
	ctx := context.Background()

	for i, c := range []struct {
		args     string
		verifier cli.Task
		opts     []cli.Opt
	}{
		{},
		{args: "m unk snd cool", verifier: func(ctx context.Context, cmd cli.Cmd, runner cli.Runner, extras ...any) (err error) { //nolint:revive
			runner.DumpErrors(os.Stdout)
			errParsed := runner.Error()
			if !regexp.MustCompile(`UNKNOWN (CmdS|Flag) FOUND:?`).MatchString(errParsed.Error()) {
				t.Log("expect 'UNKNOWN CmdS FOUND' error, but not matched.") // "unk"
			}
			return /* errParsed */
		}, opts: []cli.Opt{cli.WithUnmatchedAsError(true)}},

		{args: "m snd -n -wn cool fog --pp box", verifier: func(ctx context.Context, cmd cli.Cmd, runner cli.Runner, extras ...any) (err error) { //nolint:revive
			errParsed := runner.Error()
			if !regexp.MustCompile(`UNKNOWN (CmdS|Flag) FOUND:?`).MatchString(errParsed.Error()) {
				t.Log("expect 'UNKNOWN Flag FOUND' error, but not matched.") // "--pp"
			}
			hitTest(runner, extras, "dry-run", 2)
			hitTest(runner, extras, "wet-run", 1)
			argsAre(runner, extras, "cool", "fog")
			return nil /* errParsed */
		}},

		// general commands and flags
		{args: "jump to --full -f --dry-run", verifier: func(ctx context.Context, cmd cli.Cmd, runner cli.Runner, extras ...any) (err error) { //nolint:revive
			hitTest(runner, extras, "full", 2)
			hitTest(runner, extras, "dry-run", 1)
			errParsed := runner.Error()
			return errParsed
		}},
		// compact flags
		{args: "-qvqDq gen --debug sh --zsh -b -Dwmann --dry-run",
			verifier: func(ctx context.Context, cmd cli.Cmd, runner cli.Runner, extras ...any) (err error) { //nolint:revive
				hitTest(runner, extras, "quiet", 3)
				hitTest(runner, extras, "debug", 3)
				hitTest(runner, extras, "verbose", 1)
				hitTest(runner, extras, "wet-run", 1)
				hitTest(runner, extras, "dry-run", 2)
				errParsed := runner.Error()
				return errParsed
			}},

		// general, unknown cmd/flg errors
		{args: "m snd --help"},
		{args: "m unk snd cool", verifier: func(ctx context.Context, cmd cli.Cmd, runner cli.Runner, extras ...any) (err error) { //nolint:revive
			errParsed := runner.Error()
			if !regexp.MustCompile(`UNKNOWN (CmdS|Flag) FOUND:?`).MatchString(errParsed.Error()) {
				t.Log("expect 'UNKNOWN CmdS FOUND' error, but not matched.") // "unk"
			}
			return /* errParsed */
		}, opts: []cli.Opt{cli.WithUnmatchedAsError(true)}},
		{args: "m snd -n -wn cool fog --pp box", verifier: func(ctx context.Context, cmd cli.Cmd, runner cli.Runner, extras ...any) (err error) { //nolint:revive
			errParsed := runner.Error()
			if !regexp.MustCompile(`UNKNOWN (CmdS|Flag) FOUND:?`).MatchString(errParsed.Error()) {
				t.Log("expect 'UNKNOWN Flag FOUND' error, but not matched.") // "--pp"
			}
			hitTest(runner, extras, "dry-run", 2)
			hitTest(runner, extras, "wet-run", 1)
			argsAre(runner, extras, "cool", "fog")
			return nil /* errParsed */
		}},

		// headLike
		{args: "server start -f -129", verifier: func(ctx context.Context, cmd cli.Cmd, runner cli.Runner, extras ...any) (err error) { //nolint:revive
			hitTest(runner, extras, "foreground", 1)
			hitTest(runner, extras, "head", 1)
			hitTest(runner, extras, "tail", 0)
			valTest(runner, extras, "head", 129) //nolint:revive
			errParsed := runner.Error()
			return errParsed
		}},

		// toggle group
		{args: "generate shell --bash --zsh -p", verifier: func(ctx context.Context, cmd cli.Cmd, runner cli.Runner, extras ...any) (err error) { //nolint:revive
			if f := getMatchedFlag(extras, "shell"); f != nil {
				assertEqual(t, f.MatchedTG().MatchedTitle, "powershell")
			}
			errParsed := runner.Error()
			return errParsed
		}},

		// valid args
		{args: "server start -e apple -e zig", verifier: func(ctx context.Context, cmd cli.Cmd, runner cli.Runner, extras ...any) (err error) { //nolint:revive
			valTest(runner, extras, "enum", "zig")
			errParsed := runner.Error()
			return errParsed
		}},

		// parsing slice (-cr foo,bar,noz), compact flag with value (-mmt3)
		{args: "ms t modify -2 -cr foo,bar,noz -nfool -mmi3", verifier: func(ctx context.Context, cmd cli.Cmd, runner cli.Runner, extras ...any) (err error) { //nolint:revive
			hitTest(runner, extras, "money", 1)
			valTest(runner, extras, "both", true)
			valTest(runner, extras, "clear", true)
			valTest(runner, extras, "name", "fool")
			valTest(runner, extras, "id", "3")
			valTest(runner, extras, "remove", []string{"foo", "bar", "noz"})
			errParsed := runner.Error()
			return errParsed
		}},

		// parsing slice (-cr foo,bar,noz), compact flag with value (-mmt3)
		// merge/append to slice
		{args: "ms t modify -2 -cr foo,bar,noz -n fool -mmr 1", verifier: func(ctx context.Context, cmd cli.Cmd, runner cli.Runner, extras ...any) (err error) { //nolint:revive
			hitTest(runner, extras, "money", 1)
			valTest(runner, extras, "both", true)
			valTest(runner, extras, "clear", true)
			valTest(runner, extras, "name", "fool")
			valTest(runner, extras, "remove", []string{"foo", "bar", "noz", "1"})
			errParsed := runner.Error()
			return errParsed
		}},

		// ~~tree
		{args: "ms t t --tree", verifier: func(ctx context.Context, cmd cli.Cmd, runner cli.Runner, extras ...any) (err error) { //nolint:revive
			errParsed := runner.Error()
			if errorsv3.Is(errParsed, cli.ErrUnmatchedFlag) {
				t.Log("ErrUnmatchedFlag FOUND, that's expecting.")
			}
			return errParsed
		}},

		// ~~tree 2
		{args: "ms t t ~~tree", verifier: func(ctx context.Context, cmd cli.Cmd, runner cli.Runner, extras ...any) (err error) { //nolint:revive
			errParsed := runner.Error()
			if errorsv3.Is(errParsed, cli.ErrUnmatchedFlag) {
				t.Fatal("ErrUnmatchedFlag FOUND, that's NOT expecting.")
			}
			if f := getMatchedFlag(extras, "tree"); f != nil {
				if ms := getMatchedState(extras, f); ms != nil {
					if ms.DblTilde == false {
						t.Fatal("expecting DblTilde is true but fault.")
					}
				}
			}
			return errParsed
		}},

		{args: "ms t t -K -2 -cun foo,bar,noz", verifier: func(ctx context.Context, cmd cli.Cmd, runner cli.Runner, extras ...any) (err error) { //nolint:revive
			hitTest(runner, extras, "insecure", 1)
			valTest(runner, extras, "insecure", true)
			valTest(runner, extras, "both", true)
			valTest(runner, extras, "clear", true)
			valTest(runner, extras, "unset", []string{"foo", "bar", "noz"})
			errParsed := runner.Error()
			return errParsed
		}},

		{},
		{},
		{},
		{},
		{},
		{},
	} {
		if c.args == "" && c.verifier == nil {
			continue
		}

		t.Log()
		t.Log()
		t.Log()
		t.Logf("--------------- test #%d: Parsing %q\n", i, c.args)

		app, ww := cleanApp(t,
			cli.WithStore(store.New()),
			cli.WithExternalLoaders(NewConfigFileLoader(), NewEnvVarLoader()),
			cli.WithArgs(append([]string{"demo-app"}, strings.Split(c.args, " ")...)...),
			cli.WithHelpScreenWriter(&discardP{}),

			cli.WithTasksBeforeParse(func(ctx context.Context, cmd cli.Cmd, runner cli.Runner, extras ...any) (err error) {
				return
			}, c.verifier),
			cli.WithTasksBeforeRun(func(ctx context.Context, cmd cli.Cmd, runner cli.Runner, extras ...any) (err error) {
				// err = runner.Error()
				return
			}, aTaskBeforeRun),
		)
		// ww.Config.Store = store.New()
		// ww.Config.Loaders = []cli.Loader{loaders.NewConfigFileLoader(), loaders.NewEnvVarLoader()}

		// ww.setArgs(append([]string{app.Name()}, strings.Split(c.args, " ")...))
		// ww.tasksAfterParse = []taskAfterParse{c.verifier}
		// ww.Config.TasksBeforeRun = []cli.Task{aTaskBeforeRun}
		err := ww.Run(ctx, c.opts...) // withTasksBeforeRun(taskBeforeRun),withTasksAfterParse(c.verifier))
		// err := app.Run()
		if err != nil {
			_ = app
			t.Fatal(err)
		}
	}
}

func argsAre(runner cli.Runner, extras []any, list ...string) {
	args, _ := getPositionalArgs(extras), runner
	if !reflect.DeepEqual(args, list) {
		cc, _ := getLastCmd(extras), runner
		panic(fmt.Sprintf("expect positional args are %v but got %v (for cmd %v)", list, args, cc))
	}
}

func hitTest(runner cli.Runner, extras []any, longTitle string, times int) {
	cc, _ := getLastCmd(extras), runner
	ctx := context.Background()
	if f := cc.FindFlagBackwards(ctx, longTitle); f == nil {
		panic(fmt.Sprintf("can't found flag: %q", longTitle))
	} else if f.GetTriggeredTimes() != times {
		panic(fmt.Sprintf("expect hit times is %d but got %d (for flag %v)", times, f.GetTriggeredTimes(), f))
	}
}

func valTest(runner cli.Runner, extras []any, longTitle string, val any) {
	cc, _ := getLastCmd(extras), runner
	ctx := context.Background()
	if f := cc.FindFlagBackwards(ctx, longTitle); f == nil {
		panic(fmt.Sprintf("can't found flag: %q", longTitle))
	} else if !reflect.DeepEqual(f.DefaultValue(), val) {
		panic(fmt.Sprintf("expect flag's value is '%v' but got '%v' (for flag %v)", val, f.DefaultValue(), f))
	}
}

func getMatchedCommand(extras []any, longTitle string) (cc *cli.CmdS) {
	if len(extras) > 0 {
		if ctx, ok := extras[0].(interface {
			MatchedCommand(longTitle string) (cc *cli.CmdS)
		}); ok {
			cc = ctx.MatchedCommand(longTitle)
		}
	}
	return
}

func getCommandMatchedState(extras []any, c *cli.CmdS) (ms *cli.MatchState) {
	if len(extras) > 0 {
		if ctx, ok := extras[0].(interface {
			CommandMatchedState(c *cli.CmdS) (m *cli.MatchState)
		}); ok {
			ms = ctx.CommandMatchedState(c)
		}
	}
	return
}

func getMatchedFlag(extras []any, longTitle string) (ff *cli.Flag) {
	if len(extras) > 0 {
		if ctx, ok := extras[0].(interface {
			MatchedFlag(longTitle string) (ff *cli.Flag)
		}); ok {
			ff = ctx.MatchedFlag(longTitle)
		}
	}
	return
}

func getMatchedState(extras []any, f *cli.Flag) (ms *cli.MatchState) {
	if len(extras) > 0 {
		if ctx, ok := extras[0].(interface {
			FlagMatchedState(f *cli.Flag) (m *cli.MatchState)
		}); ok {
			ms = ctx.FlagMatchedState(f)
		}
	}
	return
}

func getLastCmd(extras []any) (cc *cli.CmdS) {
	if len(extras) > 0 {
		if ctx, ok := extras[0].(interface{ LastCmd() *cli.CmdS }); ok {
			cc = ctx.LastCmd()
		}
	}
	return
}

func getPositionalArgs(extras []any) (args []string) {
	if len(extras) > 0 {
		if ctx, ok := extras[0].(interface{ PositionalArgs() []string }); ok {
			args = ctx.PositionalArgs()
		}
	}
	return
}

func assertEqual(t *testing.T, expect, actual any, msgs ...any) {
	if expect != actual {
		t.Fatalf("expecting %v but got %v", actual, expect)
	}
	_ = msgs
}
