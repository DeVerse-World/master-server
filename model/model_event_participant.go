package model

import (
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
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

func (o *EventParticipant) Update() error {
	err := DB().Model(&o).Save(o).Error

	return err
}

func (o *EventParticipant) Delete() error {
	db := DB().Delete(EventParticipant{}, "id = ?", o.ID)

	if db.Error != nil {
		return db.Error
	}

	return nil
}

func (o *EventParticipant) GetParticipant(userId uint, eventId int) error {
	err := DB().
		Where("user_id=?", userId).
		Where("event_id=?", eventId).
		First(o).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrDataNotFound
	}

	return err
}

func GetEventParticipants(eventId uint) ([]EventParticipant, error) {
	var participants []EventParticipant
	err := DB().Find(&participants, "event_id=?", eventId).Error
	return participants, err
}
