package model

import (
	"database/sql"
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Task struct {
	ID        uuid.UUID    `gorm:"type:uuid;primary_key;"`
	Name      string       `gorm:"type:string;size:256;not null"`
	Image     string       `gorm:"type:string;size:256;not null"`
	Namespace string       `gorm:"type:string;size:64;not null"`
	Runtime   string       `gorm:"type:string;size:32;not null"`
	Status    string       `gorm:"type:string;size:32;not null"`
	Script    string       `gorm:"type:string;size:1000;not null"`
	Result    string       `gorm:"type:string;size:1000"`
	CreatedAt time.Time    `gorm:"type:TIMESTAMP with time zone;not null"`
	UpdatedAt sql.NullTime `gorm:"type:TIMESTAMP with time zone;null"`
}

func (task *Task) BeforeCreate(tx *gorm.DB) (err error) {
	uuid, err := uuid.NewV4()
	if err != nil {
		return
	}
	task.ID = uuid
	task.CreatedAt = time.Now().UTC()
	return
}

func (task *Task) BeforeUpdate(tx *gorm.DB) (err error) {
	task.UpdatedAt = sql.NullTime{Time: time.Now().UTC(), Valid: true}
	return
}
