module github.com/hedzr/cmdr-loaders

go 1.23.0

toolchain go1.23.3

// replace github.com/hedzr/cmdr/v2 => ../cmdr

// replace gopkg.in/hedzr/errors.v3 => ../../24/libs.errors

require (
	github.com/hedzr/cmdr/v2 v2.1.48
	github.com/hedzr/evendeep v1.3.48
	github.com/hedzr/is v0.8.47
	github.com/hedzr/logg v0.8.48
	github.com/hedzr/store v1.3.48
	github.com/hedzr/store/codecs/hcl v1.3.48
	github.com/hedzr/store/codecs/hjson v1.3.48
	github.com/hedzr/store/codecs/json v1.3.48
	github.com/hedzr/store/codecs/nestext v1.3.48
	github.com/hedzr/store/codecs/toml v1.3.48
	github.com/hedzr/store/codecs/yaml v1.3.48
	github.com/hedzr/store/providers/env v1.3.48
	github.com/hedzr/store/providers/file v1.3.48
	gopkg.in/hedzr/errors.v3 v3.3.5
)

require (
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/apparentlymart/go-textseg/v15 v15.0.0 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/hashicorp/hcl/v2 v2.24.0 // indirect
	github.com/hjson/hjson-go/v4 v4.5.0 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/npillmayer/nestext v0.1.3 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/zclconf/go-cty v1.16.3 // indirect
	golang.org/x/exp v0.0.0-20250718183923-645b1fa84792 // indirect
	golang.org/x/mod v0.26.0 // indirect
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sync v0.16.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/term v0.33.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	golang.org/x/tools v0.35.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
