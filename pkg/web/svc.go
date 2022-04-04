package web

import (
	"context"
	"strings"

	"github.com/blugelabs/bluge/search/highlight"
	"github.com/rs/zerolog/log"
	"github.com/wenerme/torrenti/pkg/torrenti"

	"github.com/wenerme/torrenti/pkg/search"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/samber/lo"
	webv1 "github.com/wenerme/torrenti/pkg/apis/media/web/v1"
	"github.com/wenerme/torrenti/pkg/scrape/handlers"
	"github.com/wenerme/torrenti/pkg/torrenti/models"
	"github.com/wenerme/torrenti/pkg/torrenti/util/nilx"
	"github.com/wenerme/torrenti/pkg/torrenti/util/protox"
	"gorm.io/gorm"
)

type NewWebServiceServerOptions struct {
	DB     *gorm.DB
	Search *search.Service
}

func NewWebServiceServer(conf NewWebServiceServerOptions) webv1.WebServiceServer {
	return &webServiceServer{DB: conf.DB, Search: conf.Search}
}

type webServiceServer struct {
	webv1.UnimplementedWebServiceServer
	DB     *gorm.DB
	Search *search.Service
}

func (s *webServiceServer) SearchTorrentRef(ctx context.Context, req *webv1.SearchTorrentRefRequest) (resp *webv1.SearchTorrentRefResponse, err error) {
	if req.Limit <= 0 || req.Limit > 200 {
		req.Limit = 100
	}
	if req.Offset < 0 {
		req.Offset = 0
	}
	if req.Search == "" {
		return resp, status.Error(codes.InvalidArgument, "query is empty")
	}
	sr, err := s.Search.SearchTorrent(context.Background(), &search.SearchRequest{
		QueryString: req.Search,
		Limit:       int(req.Limit),
		Offset:      int(req.Offset),
	})
	if err != nil {
		return
	}
	resp = &webv1.SearchTorrentRefResponse{
		Items:    nil,
		Total:    int32(sr.Count),
		Duration: int32(sr.Duration.Milliseconds()),
	}
	ids := lo.Map(sr.Docs, func(t *search.DocumentMatch, i int) string {
		return t.ID
	})
	docs := lo.Map(sr.Docs, func(t *search.DocumentMatch, i int) *torrenti.TorrentSearchMatch {
		return &torrenti.TorrentSearchMatch{
			Match: t,
		}
	})
	byID := lo.KeyBy(docs, func(t *torrenti.TorrentSearchMatch) string {
		return t.Match.ID
	})

	var out []*models.MetaFile
	err = s.DB.Model(models.MetaFile{}).
		Where("torrent_hash in (?)", ids).
		Select([]string{"filename", "content_hash", "torrent_hash", "creation_date", "comment", "created_by"}).
		Preload("Torrent", func(db *gorm.DB) *gorm.DB {
			return db.Select([]string{"name", "hash", "total_file_size", "file_count", "is_dir"})
		}).
		Find(&out).Error
	if err != nil {
		return
	}
	hi := highlight.NewHTMLHighlighter()
	for _, v := range out {
		doc := byID[v.TorrentHash]
		doc.Model = v

		if doc.Model.Torrent != nil {
			doc.HighlightTorrentName = hi.BestFragment(doc.Match.Locations[search.TorrentFieldTorrentFileName], []byte(doc.Model.Torrent.Name))
		}
		doc.HighlightFileName = hi.BestFragment(doc.Match.Locations[search.TorrentFieldMetaFileName], []byte(doc.Model.Filename))

		doc.HighlightTorrentName = search.MergeHTMLMark(doc.HighlightTorrentName)
		doc.HighlightFileName = search.MergeHTMLMark(doc.HighlightFileName)
	}

	resp.Items = make([]*webv1.SearchTorrentRef, 0, len(out))
	for _, v := range docs {
		vv := toItem(v)
		if vv != nil {
			resp.Items = append(resp.Items, vv)
		}
	}
	return
}

func toItem(t *torrenti.TorrentSearchMatch) *webv1.SearchTorrentRef {
	if t.Model == nil {
		log.Warn().Str("hash", t.Match.ID).Msg("no model")
		return nil
	}

	return &webv1.SearchTorrentRef{
		Item:                 toTorrentRef(t.Model, 0),
		HighlightFileName:    t.HighlightFileName,
		HighlightTorrentName: t.HighlightTorrentName,
	}
}

func (s *webServiceServer) ListTorrentRef(ctx context.Context, req *webv1.ListTorrentRefRequest) (resp *webv1.ListTorrentRefResponse, err error) {
	resp = &webv1.ListTorrentRefResponse{}
	se := strings.TrimSpace(req.Search)
	pageSize := 100
	offset := int(req.GetPage()) * pageSize
	query := s.DB.Preload("Torrent")

	var sr *search.SearchResponse
	if se != "" && s.Search != nil {
		sr, err = s.Search.SearchTorrent(ctx, &search.SearchRequest{
			QueryString: se,
			Limit:       pageSize,
			Offset:      offset,
		})
		if err != nil {
			return
		}
	} else {
		query = query.Limit(pageSize).Order("created_at desc")
	}

	var out []*models.MetaFile
	if sr != nil {
		query = query.Where("torrent_hash in (?)", lo.Map(sr.Docs, func(t *search.DocumentMatch, i int) string {
			return t.ID
		}))
	}
	if err = query.Find(&out).Error; err != nil {
		return
	}
	resp.Items = lo.Map(out, toTorrentRef)
	if sr != nil {
		m := make(map[string]int, len(sr.Docs))
		for i, v := range sr.Docs {
			m[v.ID] = i
		}
		slices.SortFunc(resp.Items, func(i, j *webv1.TorrentRef) bool {
			return m[i.FileHash] < m[j.FileHash]
		})
	}
	return
}

func (s *webServiceServer) GetTorrent(ctx context.Context, req *webv1.GetTorrentRequest) (resp *webv1.GetTorrentResponse, err error) {
	var out *models.Torrent
	err = s.DB.Where(models.Torrent{Hash: req.GetHash()}).Find(&out).Error
	if err == gorm.ErrRecordNotFound {
		err = status.Errorf(codes.NotFound, "torrent not found")
		return
	}
	resp = &webv1.GetTorrentResponse{
		Item: toTorrent(out),
	}
	return
}

func toTorrent(in *models.Torrent) (out *webv1.Torrent) {
	if in == nil {
		return nil
	}
	out = &webv1.Torrent{
		FileName: in.Name,
		Hash:     in.Hash,
		// Magnet:    in.Hash,
		FileSize:  in.TotalFileSize,
		FileCount: int32(in.FileCount),
		Ext:       "",
		IsDir:     in.IsDir,
	}
	if !in.IsDir {
		out.Ext = handlers.Ext(in.Name)
	}
	return out
}

func toTorrentRef(in *models.MetaFile, idx int) *webv1.TorrentRef {
	if in == nil {
		return nil
	}

	return &webv1.TorrentRef{
		FileName:    in.Filename,
		FileHash:    in.ContentHash,
		Referer:     nilx.NilToEmpty(in.Referer),
		Comment:     in.Comment,
		CreatedBy:   in.CreatedBy,
		CreatedAt:   protox.UnixToTimestamp(in.CreationDate),
		TorrentHash: in.TorrentHash,
		Torrent:     toTorrent(in.Torrent),
	}
}
