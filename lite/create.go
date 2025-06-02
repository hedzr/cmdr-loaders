package lite

import (
	"github.com/hedzr/cmdr/v2"
	"github.com/hedzr/cmdr/v2/cli"
)

// Create provides a concise interface to create an cli app easily.
func Create(appName, version, author, desc string, opts ...cli.Opt) cmdr.Creator {
	return cmdr.Create(appName, version, author, desc, append([]cli.Opt{
		// import "github.com/hedzr/cmdr-loaders/lite" to get in advanced external loading features
		cmdr.WithExternalLoaders(
			NewConfigFileLoader(
				WithAlternateWriteBack(false), // disable write-back the modified state into alternative config file
				WithAlternateDotPrefix(false), // use '<app-name>.toml' instead of '.<app-name>.toml'
			),
			NewEnvVarLoader(),
		),
	}, opts...)...)
}
