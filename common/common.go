package common

import (
	"context"

	"github.com/hedzr/cmdr/v2/cli"
)

type Item struct {
	// In a Folder, we try to stat() '$APP.yaml' or with another suffix.
	// But if Dot is true, '.$APP.yaml' will be stat() and loaded.
	Folder    string
	Dot       bool // prefix '.' to the filename?
	Recursive bool // following 'conf.d' subdirectory?
	Watch     bool // enable watching routine?
	WriteBack bool // write-back to "alternative config" file?

	// NoFlattenKeys bool // don't flatten keys (a flattened key looks like: "app.logging.days = 7")

	hit              bool // this item is valid and the config file loaded?
	writeBackHandler WriteBackHandler
	concreteFile     string
}

func (s *Item) Hit() bool                          { return s.hit }
func (s *Item) WriteBackHandler() WriteBackHandler { return s.writeBackHandler }
func (s *Item) ConcreteFile() string               { return s.concreteFile }

func (s *Item) SetHit(b bool)                          { s.hit = b }
func (s *Item) SetWriteBackHandler(h WriteBackHandler) { s.writeBackHandler = h }
func (s *Item) SetConcreteFile(str string)             { s.concreteFile = str }

type WriteBackHandler interface {
	Save(ctx context.Context) error
}

type Loadable interface {
	Load(ctx context.Context, app cli.App) (err error)
}

type SingleFileLoadable interface {
	LoadFile(ctx context.Context, filename string, app cli.App) (err error)
}

type QueryLoadedSources interface {
	LoadedSources() cli.LoadedSources
}

const (
	Primary     = "primary"
	Secondary   = "secondary"
	Alternative = "alternative"

	// confSubFolderName = "conf.d"
)

func DefaultFolderMap(dot, writeBack bool) map[string][]*Item {
	return map[string][]*Item{
		// Primary configs, which define the baseline of app config, are generally
		// bundled with application release.
		// App installer will dispatch primary config files to the standard directory
		// position. It's `/etc/$APP/` on linux, or `/usr/loca/etc/$app` on macOS by
		// Homebrew.
		// For debugging easier in developing, we also check `./ci/etc/$app`.
		Primary: {
			{Folder: "/etc/$APP", Recursive: true, Watch: true},
			{Folder: "/usr/local/etc/$APP", Recursive: true, Watch: true},
			{Folder: "/opt/homebrew/etc/$APP", Recursive: true, Watch: true},
			{Folder: "/usr/lib/$APP", Recursive: true, Watch: true},
			{Folder: "./ci/etc/$APP", Recursive: true, Watch: true},
			{Folder: "../ci/etc/$APP", Recursive: true, Watch: true},
			{Folder: "../../ci/etc/$APP", Recursive: true, Watch: true},
		},
		// Secondary configs, which may make some patches on the baseline if necessary.
		// On linux and macOS, it can be `~/.$app` or `~/.config/$app` (`XDG_CONFIG_DIR`).
		Secondary: {
			{Folder: "$HOME/.cmdrrc", Recursive: false, Watch: true, WriteBack: false},
			{Folder: "./.cmdrrc", Recursive: false, Watch: true, WriteBack: false},
			{Folder: "$HOME/.$APP", Recursive: true, Watch: true},
			{Folder: "$CONFIG_DIR/$APP", Recursive: true, Watch: true},
			{Folder: "./ci/config/$APP", Recursive: true, Watch: true},
			{Folder: "../ci/config/$APP", Recursive: true, Watch: true},
			{Folder: "../../ci/config/$APP", Recursive: true, Watch: true},
		},
		// Alternative config, which is live config, can be read and written.
		// Application, such as cmdr-based, reads primary config on startup, and
		// patches it with secondary config, and updates these configs with
		// alternative config finally.
		// At application terminating, the changes can be written back to alternative
		// config.
		Alternative: {
			{Folder: ".", Dot: dot, Recursive: false, Watch: true, WriteBack: writeBack},
		},
	}
}
