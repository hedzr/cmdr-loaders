package loaders

import (
	"github.com/hedzr/cmdr-loaders/local"
	"github.com/hedzr/cmdr/v2"
	"github.com/hedzr/cmdr/v2/cli"
)

// Create provides a concise interface to create an cli app easily.
func Create(appName, version, author, desc string, opts ...cli.Opt) cmdr.Creator {
	return cmdr.Create(appName, version, author, desc, append([]cli.Opt{
		// import "github.com/hedzr/cmdr-loaders/local" to get in advanced external loading features
		cmdr.WithExternalLoaders(
			local.NewConfigFileLoader(
				local.WithAlternateWriteBack(false), // disable write-back the modified state into alternative config file
				local.WithAlternateDotPrefix(false), // use '<app-name>.toml' instead of '.<app-name>.toml'
			),
			local.NewEnvVarLoader(),
		),
	}, opts...)...)
}
