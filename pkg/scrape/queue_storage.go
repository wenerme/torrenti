package scrape

import (
	"time"

	"github.com/gocolly/colly/v2/queue"
	"github.com/rs/zerolog/log"
	"github.com/tidwall/gjson"
	"github.com/wenerme/torrenti/pkg/torrenti/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ queue.Storage = (*QueueStorage)(nil)

type QueueRequest struct {
	models.Model
	URL      string `gorm:"unique"`
	Depth    int
	PulledAt *time.Time
	Raw      []byte
}

type QueueStorage struct {
	DB *gorm.DB
}

func (q *QueueStorage) Init() error {
	return q.DB.AutoMigrate(QueueRequest{})
}

func (q *QueueStorage) AddRequest(bytes []byte) error {
	r := &QueueRequest{
		URL:   gjson.GetBytes(bytes, "URL").String(),
		Depth: int(gjson.GetBytes(bytes, "Depth").Int()),
		Raw:   bytes,
	}
	log.Debug().Str("url", r.URL).Msg("add request")
	return q.DB.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "url"}},
			DoNothing: true,
		}).
		Create(r).Error
}

func (q *QueueStorage) GetRequest() (out []byte, err error) {
	var req QueueRequest
	db := q.DB

	// now := time.Now()
	// err = db.
	//	Model(&req).
	//	Where("id = (?)", db.Model(QueueRequest{}).Where("pulled_at IS NULL").Limit(1).Select("id")).
	//	Clauses(clause.Returning{}).
	//	Updates(QueueRequest{PulledAt: &now}).Error

	err = db.
		Model(QueueRequest{}).
		Clauses(clause.Returning{}).
		Where("id = (?)", db.Model(QueueRequest{}).Limit(1).Select("id")).
		Delete(&req).Error
	return req.Raw, err
}

func (q *QueueStorage) QueueSize() (out int, err error) {
	var count int64
	err = q.DB.
		Model(QueueRequest{}).
		// Where("pulled_at is null").
		Count(&count).Error
	out = int(count)
	return
}
