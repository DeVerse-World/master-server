package model

import (
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type SubworldTemplateTag struct {
	ID                 uint      `gorm:"primary_key" json:"id"`
	TagName            string    `json:"tag_name"`
	UpdatedAt          time.Time `json:"updated_at"`
	CreatedAt          time.Time `json:"created_at"`
	SubworldTemplateID uint      `json:"subworld_template_id"`
}

func (SubworldTemplateTag) TableName() string {
	return "subworld_template_tags"
}

func (m *SubworldTemplateTag) GetById(id uint) error {
	err := DB().Where("id=?", id).First(m).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrDataNotFound
	}

	return err
}

func (m *SubworldTemplateTag) Create() error {
	db := DB().Create(m)

	var mysqlErr *mysql.MySQLError
	if errors.As(db.Error, &mysqlErr) && mysqlErr.Number == DbDuplicateEntryCode {
		return ErrKeyConflict
	} else if db.Error != nil {
		return db.Error
	}

	return nil
}

func (m *SubworldTemplateTag) Update() error {
	err := DB().Model(&m).Save(m).Error

	return err
}

func (m *SubworldTemplateTag) Delete() error {
	db := DB().Delete(SubworldTemplateTag{}, "id = ?", m.ID)

	if db.Error != nil {
		return db.Error
	}

	return nil
}

func GetAllSubworldTemplateTags(subworldTemplateID uint) ([]SubworldTemplateTag, error) {
	var rs []SubworldTemplateTag
	err := DB().
		Where("subworld_template_id = ?", subworldTemplateID).
		Find(&rs).Error
	return rs, err
}
