package scraper

import (
	"github.com/wenerme/torrenti/pkg/torrenti/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Store struct {
	DB *gorm.DB
}

func (s *Store) IsVisited(url string) (found bool, err error) {
	n := int64(0)
	err = s.DB.Model(models.KV{}).Where(&models.KV{
		Type: "Visit",
		Key:  url,
	}).Count(&n).Error
	found = n > 0
	return
}

func (s *Store) Visit(u string) error {
	kv := &models.KV{
		Type: "Visit",
		Key:  u,
	}
	return s.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "type"}, {Name: "key"}},
		DoNothing: true,
	}).Create(kv).Error
}
