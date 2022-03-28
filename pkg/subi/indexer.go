package subi

import (
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/wenerme/torrenti/pkg/torrenti/models"
	"github.com/wenerme/torrenti/pkg/torrenti/util"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type NewIndexerOptions struct {
	DB *gorm.DB
}

func NewIndexer(o NewIndexerOptions) (*Indexer, error) {
	if o.DB == nil {
		return nil, errors.New("db is nil")
	}
	idx := &Indexer{DB: o.DB}
	if err := idx.DB.Migrator().AutoMigrate(
		models.SubtitleRef{},
		models.SubtitleContent{},
	); err != nil {
		return nil, err
	}

	return idx, nil
}

type Indexer struct {
	DB *gorm.DB
}

func (idx *Indexer) Index(f *util.File) (err error) {
	data, err := f.ReadAll()
	if err != nil {
		return err
	}
	ch := util.ContentHashBytes(data)
	content := &models.SubtitleContent{
		Hash:     ch,
		Ext:      filepath.Ext(f.Name()),
		Size:     int(f.Length),
		RawBytes: data,
	}
	{
		ret := idx.DB.Clauses(clause.OnConflict{
			Columns:   content.ConflictColumns(),
			DoNothing: true,
		}).Create(&content)
		if err = errors.Wrap(ret.Error, "create SubtitleContent"); err != nil {
			return
		}
	}
	sref := &models.SubtitleRef{
		ContentHash: ch,
		Filename:    f.Name(),
		URL:         f.URL,
	}
	{
		ret := idx.DB.Clauses(clause.OnConflict{
			Columns:   sref.ConflictColumns(),
			DoNothing: true,
		}).Create(&sref)
		if err = errors.Wrap(ret.Error, "create SubtitleRef"); err != nil {
			return
		}
	}
	return
}
