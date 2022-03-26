package indexer

import "github.com/wenerme/torrenti/pkg/indexer/util"

var IndexerContextKey = util.ContextKey[*Indexer]{"indexer"}
