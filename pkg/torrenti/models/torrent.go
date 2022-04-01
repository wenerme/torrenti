package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm/clause"
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
	Announce     string
	Size         int64   `gorm:"index"`
	Referer      *string `gorm:"index"`
	Raw          datatypes.JSON
	RawBytes     []byte

	Torrent *Torrent `gorm:"foreignKey:TorrentHash;references:Hash"`
}

type TorrentFile struct {
	Model
	TorrentHash string `gorm:"uniqueIndex:torrent_files_torrent_hash_path"`
	Size        int64  `gorm:"index"`
	Path        string `gorm:"uniqueIndex:torrent_files_torrent_hash_path"`
	Filename    string
	Ext         string
}

type Torrent struct {
	Model
	Hash          string `gorm:"unique"`
	Name          string
	TotalFileSize int64 `gorm:"index"`
	FileCount     int
	PieceCount    int
	IsDir         bool
	InfoBytes     []byte
}

type Tracker struct {
	Model
	URL      string `gorm:"unique"`
	Protocol string
}

func (Tracker) ConflictColumns() []clause.Column {
	return []clause.Column{{Name: "url"}}
}

func (Torrent) ConflictColumns() []clause.Column {
	return []clause.Column{{Name: "hash"}}
}

func (TorrentFile) ConflictColumns() []clause.Column {
	return []clause.Column{{Name: "torrent_hash"}, {Name: "path"}}
}

func (MetaFile) ConflictColumns() []clause.Column {
	return []clause.Column{{Name: "content_hash"}}
}
