package local

import (
	"context"

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
func NewEnvVarLoader() *envvarloader {
	return &envvarloader{}
}

type envvarloader struct{}

func (w *envvarloader) Load(ctx context.Context, app cli.App) (err error) {
	conf := app.Store()
	name := app.Name()
	_, err = conf.Load(ctx,
		store.WithProvider(env.New(
			env.WithStorePrefix(conf.Prefix()),
			env.WithPrefix(name+"_", name+"_"),
			env.WithLowerCase(true),
			env.WithUnderlineToDot(true),
		)),
		// store.WithStorePrefix("app.cmd"),
	)
	if err == nil {
		logz.Verbose("envvars loaded")
	}
	return
}
