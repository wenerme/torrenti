package web

import (
	"context"
	"strings"

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
		FileName:  in.Name,
		Hash:      in.Hash,
		Magnet:    in.Hash,
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
