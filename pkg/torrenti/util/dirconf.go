package util

import (
	"path/filepath"

	"github.com/adrg/xdg"
)

type DirConf struct {
	RootDir   string `env:"ROOT_DIR"`
	DataDir   string `env:"DATA_DIR"`
	CacheDir  string `env:"CACHE_DIR"`
	ConfigDir string `env:"CONFIG_DIR"`
}

func (conf *DirConf) InitDirConf(name string) {
	if conf.RootDir != "" {
		conf.DataDir = defaultTo(conf.DataDir, filepath.Join(conf.RootDir, "data"))
		conf.CacheDir = defaultTo(conf.CacheDir, filepath.Join(conf.RootDir, "cache"))
		conf.ConfigDir = defaultTo(conf.ConfigDir, filepath.Join(conf.RootDir, "config"))
	} else {
		conf.DataDir = defaultTo(conf.DataDir, filepath.Join(xdg.DataHome, name))
		conf.CacheDir = defaultTo(conf.CacheDir, filepath.Join(xdg.CacheHome, name))
		conf.ConfigDir = defaultTo(conf.ConfigDir, filepath.Join(xdg.ConfigHome, name))
	}
}

func defaultTo(v string, def string) string {
	if v == "" {
		return def
	}
	return v
}
