package model

import "time"

type EventParticipant struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	Score     float32   `json:"score"`
	WalletId  *uint     `json:"wallet_id"`
	EventId   *int      `json:"event_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (EventParticipant) TableName() string {
	return "event_participants"
}

func (o *EventParticipant) Create() error {
	db := DB().Create(o)

	if db.Error != nil {
		return db.Error
	} else if db.RowsAffected == 0 {
		return ErrKeyConflict
	}

	return nil
}
