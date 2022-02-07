package entity

import (
	"time"

	"github.com/google/uuid"
)

type TodoList struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name      string    `json:"name"`
	Deletable bool      `json:"deletable"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`

	UserID uuid.UUID `json:"user_id" gorm:"type:uuid;default:uuid_generate_v4()"`
	User   User      `json:"-"`
}
