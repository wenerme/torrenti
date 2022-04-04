package torrenti

import (
	"github.com/wenerme/torrenti/pkg/search"
	"github.com/wenerme/torrenti/pkg/torrenti/models"
)

type TorrentSearchMatch struct {
	Match                *search.DocumentMatch
	Model                *models.MetaFile
	HighlightFileName    string
	HighlightTorrentName string
}
