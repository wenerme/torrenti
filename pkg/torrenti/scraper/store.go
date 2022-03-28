package scraper

import (
	"github.com/wenerme/torrenti/pkg/torrenti/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Store struct {
	DB *gorm.DB
}

func (s *Store) IsScraped(url string) (visited bool, err error) {
	kv, err := s.getVisit(url)
	if kv != nil {
		visited = kv.Value == visitScraped
	}
	return
}

func (s *Store) setVisit(u string, v string) (err error) {
	kv := &models.KV{
		Type:  typeVisit,
		Key:   u,
		Value: v,
	}
	conflict := clause.OnConflict{
		Columns: kv.ConflictColumns(),
		Where: clause.Where{Exprs: []clause.Expression{
			clause.Neq{Column: clause.Column{Table: "excluded", Name: "value"}, Value: kv.Value},
		}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "updated_at"}),
	}
	return s.DB.Clauses(conflict).Create(kv).Error
}

func (s *Store) getVisit(u string) (out *models.KV, err error) {
	err = s.DB.Where(models.KV{Type: "Visit", Key: u}).Find(&out).Error
	if gorm.ErrRecordNotFound == err {
		err = nil
	}
	return
}

func (s *Store) MarkVisiting(u string) error {
	kv := &models.KV{
		Type:  typeVisit,
		Key:   u,
		Value: visitVisiting,
	}
	conflict := clause.OnConflict{
		Columns:   kv.ConflictColumns(),
		DoNothing: true,
	}
	return s.DB.Clauses(conflict).Create(kv).Error
}

const typeVisit = "Visit"
const (
	visitScraped  = "Scraped"
	visitVisiting = "Visiting"
)

func (s *Store) MarkScraped(u string) error {
	return s.setVisit(u, visitScraped)
}
