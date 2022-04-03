package torrenti

import "github.com/wenerme/torrenti/pkg/torrenti/util"

var IndexerContextKey = util.ContextKey[*Service]{"torrenti.Indexer"}
