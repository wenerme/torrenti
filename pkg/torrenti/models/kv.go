package models

import (
	"encoding/json"

	"gorm.io/gorm/clause"

	"gorm.io/datatypes"
)

type KV struct {
	Model
	Type       string `gorm:"uniqueIndex:kvs_type_key"`
	Key        string `gorm:"uniqueIndex:kvs_type_key"`
	Value      string
	Data       datatypes.JSON
	Attributes datatypes.JSON
}

func (KV) ConflictColumns() []clause.Column {
	return []clause.Column{{Name: "type"}, {Name: "key"}}
}

func (kv *KV) SetData(v interface{}) (err error) {
	kv.Data, err = json.Marshal(v)
	return
}

func (kv *KV) MustSetData(v interface{}) *KV {
	if err := kv.SetData(v); err != nil {
		panic(err)
	}
	return kv
}
