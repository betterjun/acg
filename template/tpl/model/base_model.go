package model

import (
	"reflect"
	"time"

	"gorm.io/gorm"
)

type ModelBase struct {
	ID        uint           `gorm:"primary_key" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func isNil(i interface{}) bool {
	//defer func() {
	//	recover()
	//}()
	//vi := reflect.ValueOf(i)
	//return vi.IsNil()

	vi := reflect.ValueOf(i)
	if vi.Kind() == reflect.Ptr {
		return vi.IsNil()
	}
	return false
}
