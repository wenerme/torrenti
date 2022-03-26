package torrenti

import (
	"plugin"
)

var plugins []*plugin.Plugin

type LoadPluginOptions struct {
	Dir string
}

func LoadPlugins(o LoadPluginOptions) error {
	return _LoadPlugins(o)
}
