//go:build !cgo

package indexer

import "github.com/rs/zerolog/log"

func _LoadPlugins(o LoadPluginOptions) error {
	log.Debug().Msg("cgo disabled, plugin loading disabled")
	return nil
}
