package scrape

import (
	"github.com/wenerme/torrenti/pkg/torrenti/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type VisitStore struct {
	DB *gorm.DB
}

type VisitRecord struct {
	models.Model
	URL      string `gorm:"type:unique"`
	Visiting bool
	Scraped  bool
	File     bool
}

func (VisitRecord) ConflictColumns() []clause.Column {
	return []clause.Column{{Name: "url"}}
}

func (s *VisitStore) Init() error {
	return s.DB.AutoMigrate(VisitRecord{})
}

func (s *VisitStore) IsScraped(url string) (visited bool, err error) {
	var record VisitRecord
	err = s.DB.Where(VisitRecord{URL: url}).First(&record).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return record.Scraped, err
}

func (s *VisitStore) MarkVisiting(u string) error {
	vr := &VisitRecord{URL: u, Visiting: true}
	conflict := clause.OnConflict{
		Columns: vr.ConflictColumns(),
		Where: clause.Where{Exprs: []clause.Expression{
			clause.Neq{Column: clause.Column{Name: "visiting"}, Value: true},
		}},
		DoUpdates: clause.AssignmentColumns([]string{"visiting", "updated_at"}),
	}
	return s.DB.Clauses(conflict).Create(vr).Error
}

func (s *VisitStore) MarkScraped(u string) error {
	vr := &VisitRecord{URL: u, Scraped: true}
	conflict := clause.OnConflict{
		Columns: vr.ConflictColumns(),
		Where: clause.Where{Exprs: []clause.Expression{
			clause.Neq{Column: clause.Column{Name: "scraped"}, Value: true},
		}},
		DoUpdates: clause.AssignmentColumns([]string{"scraped", "updated_at"}),
	}
	return s.DB.Clauses(conflict).Create(vr).Error
}

func (s *VisitStore) MarkFile(u string) error {
	vr := &VisitRecord{URL: u, Scraped: true}
	conflict := clause.OnConflict{
		Columns: vr.ConflictColumns(),
		Where: clause.Where{Exprs: []clause.Expression{
			clause.Neq{Column: clause.Column{Name: "file"}, Value: true},
		}},
		DoUpdates: clause.AssignmentColumns([]string{"file", "updated_at"}),
	}
	return s.DB.Clauses(conflict).Create(vr).Error
}
