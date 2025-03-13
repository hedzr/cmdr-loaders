# Loaders for cmdr/v2

![Go](https://github.com/hedzr/cmdr-loaders/workflows/Go/badge.svg)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/hedzr/cmdr-loaders)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/hedzr/cmdr-loaders.svg?label=release)](https://github.com/hedzr/cmdr-loaders/releases)

Local configuration file loaders for various file formats, such as YAML, TOML, HCL, and much more.

This is an addon library especially for [cmdr/v2](https://github.com/hedzr/cmdr).

The typical app is [cmdr-test/examples/large](https://github.com/hedzr/cmdr-tests/blob/master/examples/large).

![image-20241111141228632](https://cdn.jsdelivr.net/gh/hzimg/blog-pics@master/upgit/2024/11/20241111_1731305562.png)

A tiny app using `cmdr/v2` and `cmdr-loaders` is:

```go
package main

import (
	"context"
	"os"

	loaders "github.com/hedzr/cmdr-loaders"
	"github.com/hedzr/cmdr/v2"
	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/examples"
	"github.com/hedzr/cmdr/v2/examples/blueprint/cmd"
	"github.com/hedzr/cmdr/v2/examples/devmode"
	"github.com/hedzr/cmdr/v2/pkg/logz"
)

const (
	appName = "blueprint"
	desc    = `a good blueprint for you.`
	version = cmdr.Version
	author  = `The Examples Authors`
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	app := prepareApp(cmd.Commands...) // define your own commands implementations with cmd/*.go
	if err := app.Run(ctx); err != nil {
		logz.ErrorContext(ctx, "Application Error:", "err", err) // stacktrace if in debug mode/build
		os.Exit(app.SuggestRetCode())
	} else if rc := app.SuggestRetCode(); rc != 0 {
		os.Exit(rc)
	}
}

func prepareApp(commands ...cli.CmdAdder) cli.App {
	return loaders.Create(
		appName, version, author, desc,
		append([]cli.Opt{
			cmdr.WithAutoEnvBindings(true),  // default it's false
			cmdr.WithSortInHelpScreen(true), // default it's false
			// cmdr.WithDontGroupInHelpScreen(false), // default it's false
			// cmdr.WithForceDefaultAction(false),
		})...,
	).
		// importing devmode package and run its init():
		With(func(app cli.App) { logz.Debug("in dev mode?", "mode", devmode.InDevelopmentMode()) }).
		WithBuilders(
			examples.AddHeadLikeFlagWithoutCmd, // add a `--line` option, feel free to remove it.
			examples.AddToggleGroupFlags,       //
			examples.AddTypedFlags,             //
			examples.AddKilobytesFlag,          //
			examples.AddValidArgsFlag,          //
		).
		WithAdders(commands...).
		Build()
}
```

and `cmd/...` at:

- [cmd/](https://github.com/hedzr/cmdr-tests/blob/master/examples/blueprint/cmd/)

See also:

- [cmdr/v2](https://github.com/hedzr/cmdr)
- [blueprint app - cmdr-tests](https://github.com/hedzr/cmdr-tests/blob/master/examples/blueprint/)

## History

See full list in [CHANGELOG](https://github.com/hedzr/cmdr-loaders/blob/master/CHANGELOG)

## License

Apache 2.0
