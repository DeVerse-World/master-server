package model

import (
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type SubworldInstance struct {
	ID                 uint      `gorm:"primary_key" json:"id"`
	HostName           string    `json:"host_name"`
	Region             string    `json:"region"`
	MaxPlayers         int       `json:"max_players"`
	NumCurrentPlayers  int       `json:"num_current_players"`
	InstancePort       string    `json:"instance_port"`
	BeaconPort         string    `json:"beacon_port"`
	SubworldTemplateId *uint     `json:"subworld_template_id"`
	HostId             *uint     `json:"host_id"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func (SubworldInstance) TableName() string {
	return "subworld_instances"
}

func GetSubworldInstancesFromHost(hostId int) ([]SubworldInstance, error) {
	var sis []SubworldInstance
	err := DB().
		Where("host_id=?", hostId).
		Find(&sis).Error
	return sis, err
}

func GetAllSubworldInstances() ([]SubworldInstance, error) {
	var sis []SubworldInstance
	err := DB().
		Find(&sis).Error
	return sis, err
}

func (st *SubworldInstance) GetById(id int) error {
	err := DB().Where("id=?", id).First(st).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrDataNotFound
	}

	return err
}

func (st *SubworldInstance) Create() error {
	db := DB().Create(st)

	var mysqlErr *mysql.MySQLError
	if errors.As(db.Error, &mysqlErr) && mysqlErr.Number == DbDuplicateEntryCode {
		return ErrKeyConflict
	} else if db.Error != nil {
		return db.Error
	}

	return nil
}

func (st *SubworldInstance) Update() error {
	err := DB().Model(&st).Save(st).Error

	return err
}

func (st *SubworldInstance) Delete() error {
	db := DB().Delete(SubworldInstance{}, "id = ?", st.ID)

	if db.Error != nil {
		return db.Error
	}

	return nil
}
