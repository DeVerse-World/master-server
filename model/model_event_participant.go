package model

import (
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
)

type EventParticipant struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	Score     float32   `json:"score"`
	UserId    *uint     `json:"user_id"`
	EventId   *int      `json:"event_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (EventParticipant) TableName() string {
	return "event_participants"
}

func (o *EventParticipant) Create() error {
	db := DB().Create(o)

	var mysqlErr *mysql.MySQLError
	if errors.As(db.Error, &mysqlErr) && mysqlErr.Number == DbDuplicateEntryCode {
		return ErrKeyConflict
	} else if db.Error != nil {
		return db.Error
	}

	return nil
}
