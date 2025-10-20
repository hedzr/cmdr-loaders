module github.com/hedzr/cmdr-loaders/lite

go 1.24.0

toolchain go1.24.5

// replace github.com/hedzr/cmdr/v2 => ../cmdr

// replace gopkg.in/hedzr/errors.v3 => ../../24/libs.errors

require (
	github.com/hedzr/cmdr-loaders v1.3.60
	github.com/hedzr/cmdr/v2 v2.1.61
	github.com/hedzr/evendeep v1.3.61
	github.com/hedzr/is v0.8.61
	github.com/hedzr/logg v0.8.61
	github.com/hedzr/store v1.3.61
	github.com/hedzr/store/codecs/json v1.3.61
	github.com/hedzr/store/codecs/toml v1.3.61
	github.com/hedzr/store/providers/env v1.3.61
	github.com/hedzr/store/providers/file v1.3.61
	gopkg.in/hedzr/errors.v3 v3.3.5
)

require (
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	golang.org/x/exp v0.0.0-20251017212417-90e834f514db // indirect
	golang.org/x/net v0.46.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/term v0.36.0 // indirect
)
