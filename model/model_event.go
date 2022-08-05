package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type eventStage string

const (
	EVENT_UNSTARTED   eventStage = "Unstarted"
	EVENT_IN_PROGRESS eventStage = "InProgress"
	EVENT_PAUSED      eventStage = "Paused"
	EVENT_FINISHED    eventStage = "Finished"
)

type Event struct {
	ID                 uint       `gorm:"primary_key" json:"id"`
	Name               string     `json:"name"`
	EventConfigUri     string     `json:"event_config_uri"`
	MaxNumParticipants int        `json:"max_num_participants"`
	AllowTemporaryHold int        `json:"allow_temporary_hold"`
	Stage              eventStage `json:"stage"`
	UserId             *uint      `json:"user_id"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

func (Event) TableName() string {
	return "events"
}

func GetUserCreatedEvents(userID uint) ([]Event, error) {
	var events []Event
	err := DB().Find(&events, "user_id = ?", userID).Error
	return events, err
}

func (e *Event) GetById(id int) error {
	err := DB().Where("id=?", id).First(e).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrDataNotFound
	}

	return err
}

func (e *Event) Create() error {
	db := DB().Create(e)

	if db.Error != nil {
		return db.Error
	} else if db.RowsAffected == 0 {
		return ErrKeyConflict
	}

	return nil
}

func (e *Event) Update() error {
	err := DB().Model(&e).Save(e).Error

	return err
}

func (a *Event) Delete() error {
	db := DB().Delete(Event{}, "id = ?", a.ID)

	if db.Error != nil {
		return db.Error
	}

	return nil
}

func GetAllEvents() ([]Event, error) {
	var events []Event
	err := DB().Find(&events).Error
	return events, err
}
