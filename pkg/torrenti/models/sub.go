package models

import "gorm.io/gorm/clause"

type Subtitle struct {
	Model
	ContentHash string `gorm:"unique"`
	Lang        string
	Lang2       string
}

type SubtitleContent struct {
	Model
	Hash     string `gorm:"unique"`
	Ext      string
	Size     int
	RawBytes []byte
}

type SubtitleRef struct {
	Model
	ContentHash string `gorm:"uniqueIndex:subtitle_refs_content_hash_filename_url_idx"`
	Filename    string `gorm:"uniqueIndex:subtitle_refs_content_hash_filename_url_idx"`
	URL         string `gorm:"uniqueIndex:subtitle_refs_content_hash_filename_url_idx"`
}

func (SubtitleContent) ConflictColumns() []clause.Column {
	return []clause.Column{{Name: "hash"}}
}

func (SubtitleRef) ConflictColumns() []clause.Column {
	return []clause.Column{{Name: "content_hash"}, {Name: "filename"}, {Name: "url"}}
}
