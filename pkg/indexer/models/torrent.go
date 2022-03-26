package models

import (
	"time"

	"gorm.io/datatypes"

	"gorm.io/gorm"
)

type MetaFile struct {
	Model
	Filename     string
	ContentHash  string `gorm:"unique"`
	TorrentHash  string `gorm:"index"`
	CreatedBy    string
	CreationDate int64
	Comment      string
	Encoding     string
	Size         int64
	SourceURL    *string `gorm:"index"`
	Raw          datatypes.JSON
	RawBytes     []byte
}

type TorrentFile struct {
	Model
	TorrentHash string `gorm:"uniqueIndex:torrent_files_torrent_hash_path"`
	Size        int64
	Path        string `gorm:"uniqueIndex:torrent_files_torrent_hash_path"`
	Filename    string
	Ext         string
}

type Torrent struct {
	Model
	Hash       string `gorm:"unique"`
	Name       string
	TotalSize  int64
	FileCount  int
	PieceCount int
	IsDir      bool
	InfoBytes  []byte
}

type Model struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	DB *gorm.DB `gorm:"-" mapstructure:"-" json:"-" yaml:"-"`
}

func (model *Model) GetModel() *Model {
	return model
}

func (model *Model) GetDB() *gorm.DB {
	return model.DB
}

func (model *Model) AfterFind(tx *gorm.DB) (err error) {
	model.DB = tx
	return
}

func (model *Model) BeforeSave(tx *gorm.DB) (err error) {
	model.DB = tx
	return
}
