package model

import (
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type SubworldTemplate struct {
	ID                       uint      `gorm:"primary_key" json:"id"`
	FileName                 string    `json:"file_name"`
	DisplayName              string    `json:"display_name"`
	LevelIpfsUri             string    `json:"level_ipfs_uri"`
	LevelCentralizedUri      string    `json:"level_centralized_uri"`
	ThumbnailCentralizedUri  string    `json:"thumbnail_centralized_uri"`
	DerivativeUri            string    `json:"derivative_uri"`
	ParentSubworldTemplateId *uint     `json:"parent_subworld_template_id"`
	CreatorId                *uint     `json:"creator_id"`
	CreatedAt                time.Time `json:"created_at"`
	UpdatedAt                time.Time `json:"updated_at"`
}

func (SubworldTemplate) TableName() string {
	return "subworld_templates"
}

func GetRootFromCreator(creatorId int) ([]SubworldTemplate, error) {
	var sts []SubworldTemplate
	err := DB().
		Where("parent_subworld_template_id is null").
		Where("creator_id=?", creatorId).
		Find(&sts).Error
	return sts, err
}

func GetAllRoot() ([]SubworldTemplate, error) {
	var sts []SubworldTemplate
	err := DB().
		Where("parent_subworld_template_id is null").
		Find(&sts).Error
	return sts, err
}

func GetDerivFromCreator(rootId int, creatorId int) ([]SubworldTemplate, error) {
	var sts []SubworldTemplate
	err := DB().
		Where("parent_subworld_template_id = ?", rootId).
		Where("creator_id=?", creatorId).
		Find(&sts).Error
	return sts, err
}

func GetAllDeriv(rootId uint) ([]SubworldTemplate, error) {
	var sts []SubworldTemplate
	err := DB().
		Where("parent_subworld_template_id = ?", rootId).
		Find(&sts).Error
	return sts, err
}

func (st *SubworldTemplate) GetById(id int) error {
	err := DB().Where("id=?", id).First(st).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrDataNotFound
	}

	return err
}

func (st *SubworldTemplate) Create() error {
	db := DB().Create(st)

	var mysqlErr *mysql.MySQLError
	if errors.As(db.Error, &mysqlErr) && mysqlErr.Number == DbDuplicateEntryCode {
		return ErrKeyConflict
	} else if db.Error != nil {
		return db.Error
	}

	return nil
}

func (st *SubworldTemplate) Update() error {
	err := DB().Model(&st).Save(st).Error

	return err
}

func (st *SubworldTemplate) Delete() error {
	db := DB().Delete(SubworldTemplate{}, "id = ?", st.ID)

	if db.Error != nil {
		return db.Error
	}

	return nil
}
