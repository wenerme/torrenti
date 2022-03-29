package services

import (
	"context"
	"os"

	common "github.com/wenerme/torrenti/pkg/apis/indexer/common"
	torrentiv1 "github.com/wenerme/torrenti/pkg/apis/indexer/torrenti/v1"
	"github.com/wenerme/torrenti/pkg/torrenti"
	"github.com/wenerme/torrenti/pkg/torrenti/util"
)

type TorrentIndexerServer struct {
	Indexer *torrenti.Indexer
	torrentiv1.UnimplementedTorrentIndexServiceServer
}

func (i *TorrentIndexerServer) Stat(ctx context.Context, req *torrentiv1.StatRequest) (resp *torrentiv1.StatResponse, err error) {
	stat, err := i.Indexer.Stat(ctx)
	if err != nil {
		return
	}
	resp = &torrentiv1.StatResponse{
		Stat: &torrentiv1.Stat{
			MetaCount:            stat.MetaCount,
			MetaSize:             stat.MetaSize,
			TorrentCount:         stat.TorrentCount,
			TorrentFileCount:     stat.TorrentFileCount,
			TorrentFileTotalSize: stat.TorrentFileTotalSize,
		},
	}
	return
}

func (i *TorrentIndexerServer) IndexTorrent(ctx context.Context, request *torrentiv1.IndexTorrentRequest) (resp *torrentiv1.IndexTorrentResponse, err error) {
	file := toFile(request.GetFile())
	_, err = i.Indexer.IndexTorrent(ctx, &torrenti.Torrent{
		File: file,
		H:    request.GetHash(),
	})
	return nil, err
}

func toFile(f *common.File) *util.File {
	out := &util.File{
		Path:     f.GetPath(),
		Length:   f.GetLength(),
		FileMode: os.FileMode(f.GetFileMode()),
		Modified: f.Modified.AsTime(),
		Internal: f,
		Data:     f.Data,
		URL:      f.GetUrl(),
		Response: nil,
	}
	return out
}
