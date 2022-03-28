package torrenti

import (
	"bytes"
	"context"
	"encoding/json"
	"path/filepath"
	"strings"

	"go.uber.org/multierr"

	"github.com/wenerme/torrenti/pkg/torrenti/util"

	"github.com/dustin/go-humanize"
	"github.com/rs/zerolog/log"

	"github.com/pkg/errors"
	"github.com/xgfone/bt/bencode"
	"gorm.io/gorm/clause"

	"github.com/wenerme/torrenti/pkg/torrenti/models"
	"gorm.io/gorm"
)

type Indexer struct {
	DB *gorm.DB
}

type NewIndexerOptions struct {
	DB *gorm.DB
}

func NewIndexer(o NewIndexerOptions) (*Indexer, error) {
	if o.DB == nil {
		return nil, errors.New("db is nil")
	}
	idx := &Indexer{DB: o.DB}
	if err := idx.DB.Migrator().AutoMigrate(
		models.MetaFile{},
		models.Torrent{},
		models.TorrentFile{},
	); err != nil {
		return nil, err
	}

	return idx, nil
}

type IndexTorrentStat struct {
	MetaCount            int64
	MetaSize             int64
	TorrentCount         int64
	TorrentFileCount     int64
	TorrentFileTotalSize int64
}

type IndexTorrentOptions struct {
	Stat  *IndexTorrentStat
	Force bool //  if already exists, force to re-index
}
type IndexTorrentRequest struct {
	File *util.File
	Hash string
}

func (idx *Indexer) Stat(ctx context.Context) (stat *IndexTorrentStat, err error) {
	db := idx.DB
	stat = &IndexTorrentStat{}
	err = multierr.Combine(
		db.Model(models.MetaFile{}).Count(&stat.MetaCount).Error,
		db.Model(models.MetaFile{}).Select("sum(size)").Scan(&stat.MetaSize).Error,
		db.Model(models.Torrent{}).Count(&stat.TorrentCount).Error,
		db.Model(models.TorrentFile{}).Count(&stat.TorrentFileCount).Error,
		db.Model(models.Torrent{}).Select("sum(total_file_size)").Scan(&stat.TorrentFileTotalSize).Error,
	)
	return
}

func (idx *Indexer) IndexTorrent(ctx context.Context, t *Torrent, opts ...func(o *IndexTorrentOptions)) (stat *IndexTorrentStat, err error) {
	o := &IndexTorrentOptions{
		Stat: &IndexTorrentStat{},
	}
	for _, f := range opts {
		f(o)
	}
	stat = o.Stat

	if t.Meta == nil {
		err = t.Load()
		if err != nil {
			return
		}
	}

	mi := t.Meta

	mf := models.MetaFile{
		Filename:     t.FileInfo.Name(),
		ContentHash:  util.ContentHashBytes(t.Data),
		TorrentHash:  t.Hash.String(),
		CreatedBy:    mi.CreatedBy,
		CreationDate: mi.CreationDate,
		Comment:      mi.Comment,
		Encoding:     mi.Encoding,
		Size:         t.FileInfo.Size(),
		SourceURL:    nilString(t.URL),
		Raw:          nil,
		RawBytes:     nil,
	}
	m := map[string]interface{}{}
	err = bencode.NewDecoder(bytes.NewReader(t.Data)).Decode(&m)
	err = errors.Wrap(err, "decode data")
	if err != nil {
		return
	}
	delete(m, "info")
	mf.Raw, err = json.Marshal(m)
	err = errors.Wrap(err, "json.Marshal data")
	if err != nil {
		return
	}
	{
		ret := idx.DB.Clauses(clause.OnConflict{
			Columns:   mf.ConflictColumns(),
			DoNothing: true,
		}).Create(&mf)
		if err = errors.Wrap(ret.Error, "save meta"); err != nil {
			return
		}
		stat.MetaCount += ret.RowsAffected
		log.Debug().
			Str("file", t.FileInfo.Name()).Str("size", humanize.Bytes(uint64(t.FileInfo.Size()))).
			Int64("affected", ret.RowsAffected).
			Msg("index meta")
		if ret.RowsAffected == 0 && !o.Force {
			return
		}
	}

	info, err := mi.Info()
	if err != nil {
		return
	}
	files := info.AllFiles()
	tt := models.Torrent{
		Hash:          t.Hash.String(),
		Name:          info.Name,
		TotalFileSize: info.TotalLength(),
		FileCount:     len(files),
		PieceCount:    info.CountPieces(),
		IsDir:         info.IsDir(),
		InfoBytes:     mi.InfoBytes,
	}
	{
		ret := idx.DB.Clauses(clause.OnConflict{
			Columns:   tt.ConflictColumns(),
			DoNothing: true,
		}).Create(&tt)
		if err = errors.Wrap(ret.Error, "save torrent"); err != nil {
			return
		}
		stat.TorrentCount += ret.RowsAffected
		stat.TorrentFileTotalSize += tt.TotalFileSize
		log.Debug().
			Str("name", info.Name).Int("files", tt.FileCount).Str("size", humanize.Bytes(uint64(tt.TotalFileSize))).
			Int64("affected", ret.RowsAffected).
			Msg("index torrent")
		if ret.RowsAffected == 0 && !o.Force {
			return
		}
	}

	for _, f := range files {
		tf := models.TorrentFile{
			TorrentHash: tt.Hash,
			Path:        strings.Join(f.Paths, "/"),
			Size:        f.Length,
		}

		if info.IsDir() {
			tf.Path = info.Name
		}

		tf.Filename = filepath.Base(tf.Path)
		tf.Ext = filepath.Ext(tf.Path)

		ret := idx.DB.Clauses(clause.OnConflict{
			Columns:   tf.ConflictColumns(),
			DoNothing: true,
		}).Create(&tf)
		if err = errors.Wrap(ret.Error, "save torrent file"); err != nil {
			return
		}
		stat.TorrentFileCount += ret.RowsAffected
		log.Trace().
			Str("file", tf.Filename).Str("size", humanize.Bytes(uint64(tf.Size))).
			Int64("affected", ret.RowsAffected).
			Msg("index torrent file")
	}

	return
}

func nilString(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}
