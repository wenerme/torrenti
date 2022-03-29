package handlers

import (
	"context"

	"github.com/wenerme/torrenti/pkg/torrenti/util"
	"golang.org/x/exp/slices"
)

type Handler struct {
	Name       string
	Extensions []string
	OnFile     func(ctx context.Context, in *util.File, cb func(ctx context.Context, in *util.File) error) error
}

var subtitleExts = []string{".srt", ".ass", ".vtt", ".ssa", ".ttml", ".tml", ".sami", ".sbv", ".sub", ".rt", ".scc", ".dfxp", ".sup", ".idx", ".mks"}

func init() {
	slices.Sort(subtitleExts)
}

func IsSubtitleExt(ext string) bool {
	return util.BinarySearchContain(subtitleExts, ext)
}
