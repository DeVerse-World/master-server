package model

import "time"

type EventRewardOwnership struct {
	ID          uint      `gorm:"primary_key" json:"id"`
	Amount      int       `json:"amount"`
	MintedNftId *uint     `json:"minted_nft_id"`
	EventId     *uint     `json:"event_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (EventRewardOwnership) TableName() string {
	return "event_rewards"
}

func (o *EventRewardOwnership) Create() error {
	db := DB().Create(o)

	if db.Error != nil {
		return db.Error
	} else if db.RowsAffected == 0 {
		return ErrKeyConflict
	}

	return nil
}
