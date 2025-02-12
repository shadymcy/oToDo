package entity

import (
	"time"

	"github.com/yzx9/otodo/otodo"
	"gorm.io/gorm"
)

type Entity struct {
	ID        int64          `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"createdAt" gorm:"not null"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (e *Entity) BeforeCreate(tx *gorm.DB) (err error) {
	e.ID = otodo.NewID()
	return
}
