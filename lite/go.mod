module github.com/hedzr/cmdr-loaders/lite

go 1.24.0

toolchain go1.24.5

// replace github.com/hedzr/cmdr/v2 => ../cmdr

// replace gopkg.in/hedzr/errors.v3 => ../../24/libs.errors

require (
	github.com/hedzr/cmdr-loaders v1.3.55
	github.com/hedzr/cmdr/v2 v2.1.60
	github.com/hedzr/evendeep v1.3.60
	github.com/hedzr/is v0.8.60
	github.com/hedzr/logg v0.8.60
	github.com/hedzr/store v1.3.60
	github.com/hedzr/store/codecs/json v1.3.60
	github.com/hedzr/store/codecs/toml v1.3.60
	github.com/hedzr/store/providers/env v1.3.60
	github.com/hedzr/store/providers/file v1.3.60
	gopkg.in/hedzr/errors.v3 v3.3.5
)

require (
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	golang.org/x/exp v0.0.0-20250911091902-df9299821621 // indirect
	golang.org/x/net v0.44.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
	golang.org/x/term v0.35.0 // indirect
)
