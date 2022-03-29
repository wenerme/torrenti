package plugin

import (
	"plugin"
)

var plugins []*plugin.Plugin

type LoadPluginOptions struct {
	Dir string
}

func LoadPlugins(o LoadPluginOptions) error {
	// prefer hcp plugin system
	return nil
}
