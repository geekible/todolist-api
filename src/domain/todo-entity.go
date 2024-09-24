package domain

import (
	"time"

	"gorm.io/gorm"
)

type ToDoEntity struct {
	gorm.Model
	Title       string `gorm:"not null;type:varchar(255)"`
	Description string
	IsComplete  bool `gorm:"default:false"`
	DueDate     time.Time
}
