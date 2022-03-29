//go:build cgo

package plugin

import (
	"os"
	"path/filepath"
	"plugin"

	"github.com/rs/zerolog/log"
)

func _LoadPlugins(o LoadPluginOptions) error {
	return filepath.Walk(o.Dir, func(path string, info os.FileInfo, err error) error {
		if info == nil || info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".so" {
			log.Debug().Str("path", path).Msg("load plugin")
			p, err := plugin.Open(path)
			if err != nil {
				return err
			}
			plugins = append(plugins, p)
		}
		return nil
	})
}
