module github.com/hedzr/cmdr-loaders/lite

go 1.23.0

toolchain go1.23.3

// replace github.com/hedzr/cmdr/v2 => ../cmdr

// replace gopkg.in/hedzr/errors.v3 => ../../24/libs.errors

require (
	github.com/hedzr/cmdr/v2 v2.1.23
	github.com/hedzr/evendeep v1.3.25
	github.com/hedzr/is v0.8.25
	github.com/hedzr/logg v0.8.25
	github.com/hedzr/store v1.3.25
	github.com/hedzr/store/codecs/json v1.3.25
	github.com/hedzr/store/codecs/toml v1.3.25
	github.com/hedzr/store/providers/env v1.3.25
	github.com/hedzr/store/providers/file v1.3.25
	gopkg.in/hedzr/errors.v3 v3.3.5
)

require (
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	golang.org/x/exp v0.0.0-20250606033433-dcc06ee1d476 // indirect
	golang.org/x/net v0.41.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/term v0.32.0 // indirect
)
