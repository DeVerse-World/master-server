package model

import (
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type Avatar struct {
	ID             uint      `gorm:"primary_key" json:"id"`
	Name           string    `json:"name"`
	PreprocessUrl  string    `json:"preprocess_url"`
	PostprocessUrl string    `json:"postprocess_url"`
	UserId         *uint     `json:"user_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (Avatar) TableName() string {
	return "avatars"
}

func GetUserAvatars(userID uint) ([]Avatar, error) {
	var avatars []Avatar
	err := DB().Find(&avatars, "user_id = ?", userID).Error
	return avatars, err
}

func (a *Avatar) GetById(id int) error {
	err := DB().Where("id=?", id).First(a).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrDataNotFound
	}

	return err
}

func (a *Avatar) Create() error {
	db := DB().Create(a)

	var mysqlErr *mysql.MySQLError
	if errors.As(db.Error, &mysqlErr) && mysqlErr.Number == DbDuplicateEntryCode {
		return ErrKeyConflict
	} else if db.Error != nil {
		return db.Error
	}

	return nil
}

func (a *Avatar) Update() error {
	err := DB().Model(&a).Save(a).Error

	return err
}

func (a *Avatar) Delete() error {
	db := DB().Delete(Avatar{}, "id = ?", a.ID)

	if db.Error != nil {
		return db.Error
	}

	return nil
}
