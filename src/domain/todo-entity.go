package domain

import (
	"time"

	"gorm.io/gorm"
)

type ToDoEntity struct {
	gorm.Model
	UserId      uint   `gorm:"index:idx_user_id"`
	Title       string `gorm:"not null;type:varchar(255)"`
	Description string
	IsComplete  bool `gorm:"default:false"`
	DueDate     time.Time
}
