package services

import (
	"context"

	torrentiv12 "github.com/wenerme/torrenti/pkg/apis/media/torrenti/v1"
	"github.com/wenerme/torrenti/pkg/torrenti"
	"github.com/wenerme/torrenti/pkg/torrenti/util/protou"
)

type TorrentIndexerServer struct {
	Indexer *torrenti.Service
	torrentiv12.UnimplementedTorrentIndexServiceServer
}

func (i *TorrentIndexerServer) Stat(ctx context.Context, req *torrentiv12.StatRequest) (resp *torrentiv12.StatResponse, err error) {
	stat, err := i.Indexer.Stat(ctx)
	if err != nil {
		return
	}
	resp = &torrentiv12.StatResponse{
		Stat: &torrentiv12.Stat{
			MetaCount:            stat.MetaCount,
			MetaSize:             stat.MetaSize,
			TorrentCount:         stat.TorrentCount,
			TorrentFileCount:     stat.TorrentFileCount,
			TorrentFileTotalSize: stat.TorrentFileTotalSize,
		},
	}
	return
}

func (i *TorrentIndexerServer) IndexTorrent(ctx context.Context, request *torrentiv12.IndexTorrentRequest) (resp *torrentiv12.IndexTorrentResponse, err error) {
	file := protou.ToFile(request.GetFile())
	_, err = i.Indexer.IndexTorrent(ctx, &torrenti.Torrent{
		File: file,
		H:    request.GetHash(),
	})
	return nil, err
}
