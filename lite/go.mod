module github.com/hedzr/cmdr-loaders/lite

go 1.23.0

toolchain go1.23.3

// replace github.com/hedzr/cmdr/v2 => ../cmdr

// replace gopkg.in/hedzr/errors.v3 => ../../24/libs.errors

require (
	github.com/hedzr/cmdr/v2 v2.1.43
	github.com/hedzr/evendeep v1.3.43
	github.com/hedzr/is v0.8.43
	github.com/hedzr/logg v0.8.43
	github.com/hedzr/store v1.3.43
	github.com/hedzr/store/codecs/json v1.3.43
	github.com/hedzr/store/codecs/toml v1.3.43
	github.com/hedzr/store/providers/env v1.3.43
	github.com/hedzr/store/providers/file v1.3.43
	gopkg.in/hedzr/errors.v3 v3.3.5
)

require (
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	golang.org/x/exp v0.0.0-20250620022241-b7579e27df2b // indirect
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/term v0.33.0 // indirect
)
