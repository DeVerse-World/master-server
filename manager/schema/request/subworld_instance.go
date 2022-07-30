package request

import "github.com/hyperjiang/gin-skeleton/model"

type CreateSubworldInstance struct {
	SubworldInstance model.SubworldInstance `json:"subworld_instance" binding:"required"`
}

type UpdateSubworldInstance struct {
	HostName          string `json:"host_name"`
	Region            string `json:"region"`
	MaxPlayers        int    `json:"max_players"`
	NumCurrentPlayers int    `json:"num_current_players"`
	InstancePort      string `json:"instance_port"`
	BeaconPort        string `json:"beacon_port"`
}
