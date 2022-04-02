package web

import (
	"context"
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
	DB *gorm.DB
}

func NewWebServiceServer(conf NewWebServiceServerOptions) webv1.WebServiceServer {
	return &webServiceServer{DB: conf.DB}
}

type webServiceServer struct {
	webv1.UnimplementedWebServiceServer
	DB *gorm.DB
}

func (s *webServiceServer) ListTorrentRef(ctx context.Context, req *webv1.ListTorrentRefRequest) (resp *webv1.ListTorrentRefResponse, err error) {
	resp = &webv1.ListTorrentRefResponse{}
	var out []*models.MetaFile
	if err = s.DB.Limit(100).Order("created_at desc").Preload("Torrent").Find(&out).Error; err != nil {
		return
	}
	resp.Items = lo.Map(out, toTorrentRef)
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
