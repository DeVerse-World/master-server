package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type SystemSetting struct {
	ID                uint      `gorm:"primary_key" json:"id"`
	KeyName           string    `json:"key_name"`
	Value             string    `json:"value"`
	Category          string    `json:"category"`
	ObjectReferenceId string    `json:"object_reference_id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (SystemSetting) TableName() string {
	return "system_settings"
}

func (s *SystemSetting) GetByInfo(key string, category string, objectReferenceId string) error {
	err := DB().Where("key_name = ?", key).Where("category = ?", category).Where("object_reference_id = ?", objectReferenceId).First(s).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrDataNotFound
	}
	return err
}

func (s *SystemSetting) Create() error {
	db := DB().Create(s)

	if db.Error != nil {
		return db.Error
	} else if db.RowsAffected == 0 {
		return ErrKeyConflict
	}

	return nil
}

func (s *SystemSetting) Update() error {
	err := DB().Model(&s).Save(s).Error
	return err
}
