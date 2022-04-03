package search

import (
	"context"
	"path/filepath"
	"strings"
	"time"

	"github.com/blugelabs/bluge"
	"github.com/blugelabs/bluge/search"
	"github.com/pkg/errors"
)

type NewServiceOptions struct {
	DataDir string
}

func NewService(opts NewServiceOptions) (s *Service, err error) {
	s = &Service{
		Torrent: &CollectionIndex{},
	}
	{
		config := bluge.DefaultConfig(filepath.Join(opts.DataDir, "torrent"))
		s.Torrent.Writer, err = bluge.OpenWriter(config)
		if err != nil {
			return
		}
		s.Torrent.Reader, err = s.Torrent.Writer.Reader()
		if err != nil {
			return
		}
	}
	return
}

type Service struct {
	Torrent *CollectionIndex
}

type SearchRequest struct {
	QueryString string
	Query       bluge.Query
	Limit       int
	Offset      int
	Orders      []string
}

type SearchResponse struct {
	Docs     []*DocumentMatch
	Count    int
	Duration time.Duration
	MaxScore float64
}
type DocumentMatch struct {
	ID        string
	Score     float64
	Locations search.FieldTermLocationMap
}

type Torrent struct {
	ID              string
	Size            int64
	MetaFileName    string
	TorrentFileName string
	CreatedAt       time.Time
}

const (
	TorrentFieldMetaFileName    = "file_name"
	TorrentFieldTorrentFileName = "torrent_file_name"
	docFieldSize                = "size"
	docFieldCreatedAt           = "created_at"
)

func (m *Torrent) Document() *bluge.Document {
	doc := bluge.NewDocument(m.ID)

	if m.MetaFileName != "" {
		doc.AddField(bluge.NewTextField(TorrentFieldMetaFileName, strings.TrimSuffix(m.MetaFileName, ".torrent")).WithAnalyzer(filenameAnalyzer).HighlightMatches())
	}
	if m.TorrentFileName != "" {
		doc.AddField(bluge.NewTextField(TorrentFieldTorrentFileName, m.TorrentFileName).WithAnalyzer(filenameAnalyzer).HighlightMatches())
	}
	if m.Size != 0 {
		doc.AddField(bluge.NewNumericField(docFieldSize, float64(m.Size)))
	}
	if !m.CreatedAt.IsZero() {
		doc.AddField(bluge.NewDateTimeField(docFieldCreatedAt, m.CreatedAt))
	}
	return doc
}

func (s *Service) IndexTorrent(ctx context.Context, v []*Torrent) (err error) {
	batch := bluge.NewBatch()
	for _, t := range v {
		doc := t.Document()
		id := doc.ID()
		if len(id.Term()) == 0 {
			return errors.New("invalid indexing: empty id")
		}
		batch.Update(id, doc)
	}
	err = s.Torrent.Writer.Batch(batch)
	return
}

func (s *Service) SearchTorrent(ctx context.Context, req *SearchRequest) (resp *SearchResponse, err error) {
	if req.Limit <= 0 {
		req.Limit = 100
	}

	query := req.Query
	if query == nil {
		query = bluge.NewBooleanQuery().
			AddShould(bluge.NewMatchQuery(req.QueryString).SetField(TorrentFieldTorrentFileName)).
			AddShould(bluge.NewMatchQuery(req.QueryString).SetField(TorrentFieldMetaFileName))
	}

	r := bluge.NewTopNSearch(req.Limit, query).SetFrom(req.Offset).WithStandardAggregations()
	r = r.IncludeLocations()

	if len(req.Orders) == 0 {
		req.Orders = []string{"-created_at"}
	}
	sorts := append([]string{"_score"}, req.Orders...)
	r = r.SortBy(sorts)

	iterator, err := s.Torrent.Reader.Search(ctx, r)
	if err != nil {
		return
	}

	agg := iterator.Aggregations()
	resp = &SearchResponse{
		Count:    int(agg.Count()),
		Duration: agg.Duration(),
		MaxScore: agg.Metric("max_score"),
	}

	var doc *search.DocumentMatch
	for {
		doc, err = iterator.Next()
		if err != nil {
			return
		}
		if doc == nil {
			break
		}

		o := &DocumentMatch{
			Locations: doc.Locations,
			Score:     doc.Score,
		}
		err = doc.VisitStoredFields(func(field string, value []byte) bool {
			if field == "_id" {
				o.ID = string(value)
			} else {
				// fixme 暂时不考虑存储的字段
				panic("found unexpected field: " + field)
			}
			return true
		})
		if err != nil {
			return
		}
		resp.Docs = append(resp.Docs, o)
	}

	return
}

type CollectionIndex struct {
	Name   string
	Reader *bluge.Reader
	Writer *bluge.Writer
}

func MergeHTMLMark(s string) string {
	return strings.ReplaceAll(s, "</mark><mark>", "")
}
