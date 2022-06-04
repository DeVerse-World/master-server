package model

import "time"

type EventReward struct {
	ID          uint      `gorm:"primary_key" json:"id"`
	Supply      int       `json:"supply"`
	MintedNftId *uint     `json:"minted_nft_id"`
	EventId     *uint     `json:"event_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (EventReward) TableName() string {
	return "event_rewards"
}

func (o *EventReward) Create() error {
	db := DB().Create(o)

	if db.Error != nil {
		return db.Error
	} else if db.RowsAffected == 0 {
		return ErrKeyConflict
	}

	return nil
}
