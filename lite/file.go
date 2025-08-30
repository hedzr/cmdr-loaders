package lite

import (
	"context"
	"os"
	"path"
	"strings"

	"github.com/hedzr/cmdr-loaders/common"
	logz "github.com/hedzr/logg/slog"

	"github.com/hedzr/store"
	"github.com/hedzr/store/codecs/json"
	"github.com/hedzr/store/codecs/toml"
	"github.com/hedzr/store/providers/file"

	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/is/dir"
)

// NewConfigFileLoader returns a new instance to load local config files.
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
func NewConfigFileLoader(opts ...Opt) *conffileloader {
	s := &conffileloader{confDFolderName: confSubFolderName, dot: true, writeBack: true}
	for _, opt := range opts {
		if opt != nil {
			opt(s)
		}
	}
	s.initOnce()
	return s
}

type Opt func(s *conffileloader)

func WithFolderMap(m map[string][]*common.Item) Opt {
	return func(s *conffileloader) {
		s.folderMap = m
	}
}

func WithFolderMapSubcategory(categoryName string, items ...*common.Item) Opt {
	return func(s *conffileloader) {
		if s.folderMap == nil {
			s.folderMap = make(map[string][]*common.Item)
		}
		if _, ok := s.folderMap[categoryName]; !ok {
			s.folderMap[categoryName] = make([]*common.Item, 0, len(items))
		}
		s.folderMap[categoryName] = append(s.folderMap[categoryName], items...)
	}
}

func WithConfDFolderName(name string) Opt {
	return func(s *conffileloader) {
		s.confDFolderName = name
	}
}

func WithAlternateDotPrefix(dotPrefix bool) Opt {
	return func(s *conffileloader) {
		s.dot = dotPrefix
	}
}

func WithAlternateWriteBack(b bool) Opt {
	return func(s *conffileloader) {
		s.writeBack = b
	}
}

// WithMoreSuffixCodecs adds more suffix codecs to the loader.
// The default codecs are: "toml", "json".
//
// The suffix is the file extension without the leading dot, like "yaml", "json", "toml", etc.
//
// For example:
//
//	import "github.com/hedzr/store/codecs/yaml"
//	lite.WithMoreSuffixCodecs(
//	    map[string]func()store.Codec{ "yaml": func() store.Codec { return yaml.New() } },
//	}
//
// If you want to add a codec for a file extension that is not supported by default,
// you can use this function to add it.
//
// You may also override the default codecs by using this function.
func WithMoreCodecs(descibers map[string]func() store.Codec) Opt {
	return func(s *conffileloader) {
		for suffix, getter := range descibers {
			s.suffixCodecMap[suffix] = getter
		}
	}
}

type conffileloader struct {
	folderMap       map[string][]*common.Item
	suffixCodecMap  map[string]func() store.Codec
	confDFolderName string

	dot       bool
	writeBack bool

	loaded cli.LoadedSources
}

func (w *conffileloader) LoadedSources() cli.LoadedSources {
	return w.loaded
}

func (w *conffileloader) Save(ctx context.Context) (err error) {
	for _, class := range []string{common.Primary, common.Secondary, common.Alternative} {
		for _, str := range w.folderMap[class] {
			if h := str.WriteBackHandler(); str.Hit() && str.WriteBack && h != nil {
				// logz.InfoContext(ctx, "Write-Back", "str", str.concreteFile)
				err = h.Save(ctx)
			}
		}
	}
	return
}

func (w *conffileloader) Load(ctx context.Context, app cli.App) (err error) {
	if w.loaded == nil {
		w.loaded = make(cli.LoadedSources)
	}

	// var conf = app.Store()

	// cwd := dir.GetCurrentDir()
	// logz.DebugContext(ctx, "conffileloader.Load()", "cwd", cwd)

	var found bool
	for _, class := range []string{common.Primary, common.Secondary, common.Alternative} {
		for _, it := range w.folderMap[class] {
			folderEx := os.ExpandEnv(it.Folder)
			logz.VerboseContext(ctx, "loading config files from Folder", "class", class, "Folder", it.Folder, "Folder-expanded", folderEx)
			if !dir.FileExists(folderEx) {
				continue
			}

			found, err = w.loadAppConfig(ctx, class, folderEx, it, app)

			if root := path.Join(folderEx, w.confDFolderName); it.Recursive && found && dir.FileExists(root) {
				found, err = w.loadSubDir(ctx, class, root, app)
			}
		}
	}

	// logz.Verbose("Store.Dump")
	// logz.Verbose(conf.Dump())
	return
}

func (w *conffileloader) LoadFile(ctx context.Context, filename string, app cli.App) (err error) {
	return w.loadConfigFile(ctx, filename, path.Ext(filename), &common.Item{Watch: true, WriteBack: false}, app)
}

func (w *conffileloader) loadAppConfig(ctx context.Context, class, folderExpanded string, it *common.Item, app cli.App) (found bool, err error) {
	rootCmd := app.RootCommand()

	// if the folderExpanded is a regular file, load it directly
	if isfile, _ := dir.IsRegularFile(folderExpanded); isfile {
		err = w.loadConfigFile(ctx, folderExpanded, path.Ext(folderExpanded), it, app)
		if err == nil {
			found = true
			w.add(false, class, folderExpanded)
			logz.VerboseContext(ctx, "config file loaded", "file", folderExpanded)
		}
		return
	}

	// or loop the files in this folder to find one
	err = dir.ForFileMax(folderExpanded, 0, 1,
		func(depth int, dirName string, fi os.DirEntry) (stop bool, err error) {
			baseName, ext, appName := dir.Basename(fi.Name()), dir.Ext(fi.Name()), rootCmd.AppName
			if it.Dot {
				appName = "." + appName
			}
			if baseName != appName {
				return
			}

			// logz.VerboseContext(ctx, "loading config file", "dir", dirName, "file", fi.Name())
			file := path.Join(dirName, fi.Name())
			err = w.loadConfigFile(ctx, file, ext, it, app)
			if err == nil {
				logz.VerboseContext(ctx, "config file loaded", "file", file)
				found, stop = true, true
				w.add(false, class, file)
			}
			return
		})
	return
}

func (w *conffileloader) loadConfigFile(ctx context.Context, filename, ext string, it *common.Item, app cli.App) (err error) {
	logz.VerboseContext(ctx, "try loading config file", "file", filename)
	strings.TrimPrefix(ext, ".")
	if codec, ok := w.suffixCodecMap[ext]; ok {
		// if ext == "" {
		// 	x, _ := os.ReadFile(filename)
		// 	logz.DebugContext(ctx, "FILE CONTENT", "file", filename, "content", string(x))
		// }
		var wr common.WriteBackHandler
		conf := app.Store()
		wr, err = conf.Load(ctx,
			// store.WithStorePrefix("app.yaml"),
			// store.WithPosition("app"),
			store.WithCodec(codec()),
			store.WithProvider(file.New(filename,
				file.WithWatchEnabled(it.Watch),
				file.WithWriteBackEnabled(it.WriteBack),
				// file.WithoutFlattenKeys(it.NoFlattenKeys),
			)),
		)
		if err == nil {
			if it.WriteBack && wr != nil {
				it.SetWriteBackHandler(wr)
				it.SetHit(true)
			}
			it.SetConcreteFile(filename)
		}
	}
	return
}

func (w *conffileloader) loadSubDir(ctx context.Context, class, root string, app cli.App) (found bool, err error) {
	err = dir.ForFile(root,
		func(depth int, dirName string, fi os.DirEntry) (stop bool, err error) {
			ext := dir.Ext(fi.Name())
			strings.TrimPrefix(ext, ".")

			if codec, ok := w.suffixCodecMap[ext]; ok {
				filename := path.Join(dirName, fi.Name())
				_, err = app.Store().Load(ctx,
					// store.WithStorePrefix("app.yaml"),
					// store.WithPosition("app"),
					store.WithCodec(codec()),
					store.WithProvider(file.New(filename)),
				)
				if err == nil {
					logz.VerboseContext(ctx, "conf.d file loaded", "file", filename)
					found, stop = true, false
					w.add(true, class, filename)
				}
			}
			return
		})
	return
}

// SetAlternativeConfigFile adds a user-specified config file into
// alternative list.
//
// Generally user can specify a prefer config file from command-line
// by option `--config file'.
//
// The point is, an Alternative config file is writeable: its content
// will be refreshed at end of invocation on a cmdr-app. The feature
// is called 'write-back' in cmdr.
//
// By this token, there is only one Alternative config file in the list.
func (w *conffileloader) SetAlternativeConfigFile(file string) {
	// w.folderMap[Alternative] = append(w.folderMap[Alternative], &Item{Folder: file, Watch: true})
	w.folderMap[common.Alternative] = []*common.Item{{Folder: file, Watch: true}}
}

func (w *conffileloader) add(subdir bool, class, file string) {
	if _, ok := w.loaded[class]; !ok {
		w.loaded[class] = new(cli.LoadedSource)
	}
	if subdir {
		w.loaded[class].Children = append(w.loaded[class].Children, file)
	} else {
		w.loaded[class].Main = append(w.loaded[class].Main, file)
	}
}

func (w *conffileloader) initOnce() {
	if w.folderMap == nil {
		w.folderMap = common.DefaultFolderMap(w.dot, w.writeBack)
	}
	if w.suffixCodecMap == nil {
		w.suffixCodecMap = map[string]func() store.Codec{
			"toml": func() store.Codec { return toml.New() },
			// "yaml":       func() store.Codec { return yaml.New() },
			// "yml":        func() store.Codec { return yaml.New() },
			"json": func() store.Codec { return json.New() },
			// "hjson":      func() store.Codec { return hjson.New() },
			// "hcl":        func() store.Codec { return hcl.New() },
			// "nestedtext": func() store.Codec { return nestext.New() },
			// "txt":        func() store.Codec { return nestext.New() },
			// "conf":       func() store.Codec { return nestext.New() },
			// "":           func() store.Codec { return nestext.New() },
		}
	}
}

const (
	confSubFolderName = "conf.d"
)
