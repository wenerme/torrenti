package indexer

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/rs/zerolog/log"

	"github.com/pkg/errors"
	"github.com/xgfone/bt/bencode"
	"gorm.io/gorm/clause"

	"github.com/wenerme/torrenti/pkg/indexer/models"
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

func (idx *Indexer) IndexTorrent(t *Torrent) (err error) {
	mi := t.Meta

	mf := models.MetaFile{
		Path:         t.File,
		Filename:     t.FileInfo.Name(),
		ContentHash:  contentHash(t.Data),
		TorrentHash:  t.Hash.String(),
		CreatedBy:    mi.CreatedBy,
		CreationDate: mi.CreationDate,
		Comment:      mi.Comment,
		Encoding:     mi.Encoding,
		Size:         t.FileInfo.Size(),
		SourceURL:    nilString(t.URL),
		Raw:          nil,
		RawBytes:     t.Data,
	}
	m := map[string]interface{}{}
	err = bencode.NewDecoder(bytes.NewReader(t.Data)).Decode(&m)
	err = errors.Wrap(err, "decode data")
	if err != nil {
		return err
	}
	delete(m, "info")
	mf.Raw, err = json.Marshal(m)
	err = errors.Wrap(err, "json.Marshal data")
	if err != nil {
		return err
	}
	{
		ret := idx.DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "content_hash"}},
			DoUpdates: clause.AssignmentColumns([]string{"path", "size"}),
		}).Create(&mf)
		if err = ret.Error; err != nil {
			return err
		}
		log.Debug().
			Str("path", t.File).Str("size", humanize.Bytes(uint64(t.FileInfo.Size()))).
			Int64("affected", ret.RowsAffected).
			Msg("index meta")
	}

	info, err := mi.Info()
	if err != nil {
		return err
	}
	files := info.AllFiles()
	tt := models.Torrent{
		Hash:       t.Hash.String(),
		Name:       info.Name,
		TotalSize:  info.TotalLength(),
		FileCount:  len(files),
		PieceCount: info.CountPieces(),
		IsDir:      info.IsDir(),
		InfoBytes:  mi.InfoBytes,
	}
	{
		ret := idx.DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "hash"}},
			DoNothing: true,
		}).Create(&tt)
		if err = ret.Error; err != nil {
			return err
		}
		log.Debug().
			Str("name", info.Name).Int("files", tt.FileCount).Str("size", humanize.Bytes(uint64(tt.TotalSize))).
			Int64("affected", ret.RowsAffected).
			Msg("index torrent")
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
			Columns:   []clause.Column{{Name: "torrent_hash"}, {Name: "path"}},
			DoNothing: true,
		}).Create(&tf)
		if err = ret.Error; err != nil {
			return err
		}
		log.Trace().
			Str("file", tf.Filename).Str("size", humanize.Bytes(uint64(tf.Size))).
			Int64("affected", ret.RowsAffected).
			Msg("index torrent file")
	}

	return
}

func contentHash(v []byte) string {
	sum := sha256.Sum256(v)
	return hex.EncodeToString(sum[:])
}

func nilString(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}
