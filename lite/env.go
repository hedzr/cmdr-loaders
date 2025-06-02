package lite

import (
	"context"
	"strings"

	logz "github.com/hedzr/logg/slog"
	"github.com/hedzr/store"
	"github.com/hedzr/store/providers/env"

	"github.com/hedzr/cmdr/v2/cli"
)

// NewEnvVarLoader return a new instance to load current
// executive environment.
//
// The key-values in environment will be loaded into Store as
// key-value pairs.
//
// For example,
//
//	app = cmdr.New().
//	    Info("demo-app", "0.3.1").
//	    Author("your-name")
//	if err := app.Run(
//	    cmdr.WithStore(store.New()),
//	    cmdr.WithExternalLoaders(
//	      local.NewConfigFileLoader(),
//	      local.NewEnvVarLoader(),
//	    ),
//	    cmdr.WithForceDefaultAction(true), // true for debug in developing time
//	); err != nil {
//	    logz.Error("Application Error:", "err", err)
//	}
//
// EnvVar-loader finds out envvars which has appname prefix (eg "TINY_APP_"),
// and strips the prefix and parses the key-value into Store.
//
// For example, an app named "tiny-app", these following vars present:
//
//	TINY_APP_COOL_S1=1
//	TINY_APP_COOL_S2=wow
//
// After cmdr loaded (before parsing), the envvars will be mapped into store as:
//
//	cool.               # mapped from TINY_APP_xxx
//	  s1  -> 1
//	  s2  -> "wow"
//	app.cmd. ...        # another top level node, mapping from app command system
//	app. ...            # default top level node, loading from external config files
//
// What difference with [cli.Config.AutoEnv]?
//
// AutoEnv is a builtin feature which maps APP_(subcmds)_(flag) env var value as
// cmdr [cli.Flag]'s default value.
//
// For example, APP_JUMP_TO_FULL=1 will map to (root)/{Cmd:jump}/{Cmd:to}/{Flag:full} as true.
func NewEnvVarLoader() *envvarloader {
	return &envvarloader{}
}

type envvarloader struct{}

func (w *envvarloader) Load(ctx context.Context, app cli.App) (err error) {
	conf := app.Store()
	prefix := conf.Prefix()
	name := strings.Replace(app.Name(), "-", "_", -1)
	namePrefixed := name
	if !strings.HasSuffix(namePrefixed, "_") {
		namePrefixed = namePrefixed + "_"
	}
	_, err = conf.Load(ctx,
		store.WithProvider(env.New(
			env.WithStorePrefix(prefix),                // load into store, at position `prefix`, eg. "app.envvars"
			env.WithPrefix(namePrefixed, namePrefixed), // only these envvars if has `name`_ prefix, eg. "TINY_APP_"
			env.WithLowerCase(true),                    // convert names to lowercase
			env.WithUnderlineToDot(true),               // split to (dotted-) tree with delimiter '_'
		)),
		// store.WithStorePrefix("app.cmd"),
	)
	if err == nil {
		logz.Verbose("envvars loaded")
	}
	return
}
