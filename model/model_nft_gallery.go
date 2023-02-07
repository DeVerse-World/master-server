package model

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

var nftGalleryValidCategories = []string{"Personal", "Contract"}

type NftGallery struct {
	ID                 uint      `gorm:"primary_key" json:"id"`
	DisplayName        string    `json:"display_name"`
	OwnerAddress       string    `json:"owner_address"`
	Category           string    `json:"category"`
	CollectionAddress  string    `json:"collectionAddress"`
	Chain              string    `json:"chain"`
	IsAutoFetch        int       `json:"is_auto_fetch"`
	SubworldTemplateId *uint     `json:"subworld_template_id"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func (NftGallery) TableName() string {
	return "nft_galleries"
}

func (e *NftGallery) GetById(id int) error {
	err := DB().Where("id=?", id).First(e).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrDataNotFound
	}

	return err
}

func (e *NftGallery) Create() error {
	if !e.isValidCategory() {
		return errors.New(
			"invalid gallery category, supported " + strings.Join(nftGalleryValidCategories, ","),
		)
	}
	db := DB().Create(e)

	if db.Error != nil {
		return db.Error
	} else if db.RowsAffected == 0 {
		return ErrKeyConflict
	}

	return nil
}

func (a *NftGallery) Delete() error {
	db := DB().Delete(NftGallery{}, "id = ?", a.ID)

	if db.Error != nil {
		return db.Error
	}

	return nil
}

func (e *NftGallery) isValidCategory() bool {
	for _, category := range nftGalleryValidCategories {
		if e.Category == category {
			return true
		}
	}
	return false
}

func GetAll() ([]NftGallery, error) {
	var nftGalleries []NftGallery
	err := DB().Find(&nftGalleries).Error
	return nftGalleries, err
}

func GetAllBySubworldTemplateId(id int) ([]NftGallery, error) {
	var nftGalleries []NftGallery
	err := DB().Where("subworld_template_id=?", id).Find(&nftGalleries).Error
	return nftGalleries, err
}
