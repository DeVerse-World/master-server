package model

import "time"

type eventStage string

const (
	UNSTARTED   eventStage = "Unstarted"
	IN_PROGRESS eventStage = "InProgress"
	PAUSED      eventStage = "Paused"
	FINISHED    eventStage = "Finished"
)

type Event struct {
	ID                 uint      `gorm:"primary_key" json:"id"`
	Name               string    `json:"name"`
	EventConfigUri     string    `json:"event_config_uri"`
	MaxNumParticipants string    `json:"max_num_participants"`
	AllowTemporaryHold int       `json:"allow_temporary_hold"`
	Stage              string    `json:"stage"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func (Event) TableName() string {
	return "minted_nfts"
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
