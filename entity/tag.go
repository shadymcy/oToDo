package entity

import "github.com/google/uuid"

type Tag struct {
	ID         string
	Title      string `json:"title"`
	CountTotal uint   `json:"count_total"`
	CountTodo  uint   `json:"count_todo"`
	CountDone  uint   `json:"count_done"`

	UserId uuid.UUID `json:"user_id"`
	User   User      `json:"-"`
}
