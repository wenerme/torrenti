package indexer

import (
	"context"

	indexerv1 "github.com/wenerme/torrenti/pkg/apis/media/indexer/v1"
	"github.com/wenerme/torrenti/pkg/torrenti/util"
	"github.com/wenerme/torrenti/pkg/torrenti/util/protou"
)

type service struct {
	indexerv1.UnimplementedIndexServiceServer
	s *Service
}

func (svc *service) Index(ctx context.Context, req *indexerv1.IndexRequest) (resp *indexerv1.IndexResponse, err error) {
	resp = &indexerv1.IndexResponse{}
	var r *IndexResponse
	r, err = svc.s.Index(ctx, &IndexRequest{
		File: protou.ToFile(req.GetFile()),
		URL:  req.Url,
	})
	if r != nil {
	}
	return
}

type IndexRequest struct {
	File *util.File
	URL  string
}

type (
	IndexResponse struct{}
	Service       struct{}
)

func (s *Service) Index(ctx context.Context, req *IndexRequest) (resp *IndexResponse, err error) {
	return
}
