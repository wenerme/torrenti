package models

import (
	"time"

	"gorm.io/gorm"
)

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
