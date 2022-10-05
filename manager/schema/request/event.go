package request

import "github.com/hyperjiang/gin-skeleton/model"

type CreateEvent struct {
	Event        model.Event         `json:"event" binding:"required"`
	EventRewards []model.EventReward `json:"event_rewards" binding:"required"`
}

type UpdateEvent struct {
	Name               string `json:"event"`
	Category           string `json:"category"`
	EventConfigUri     string `json:"event_config_uri"`
	MaxNumParticipants int    `json:"max_num_participants"`
	AllowTemporaryHold int    `json:"allow_temporary_hold"`
}

type CheckEventParticipant struct {
	UserId uint `json:"user_id" binding:"required"`
}

type UpdateScore struct {
	Score float32 `json:"score"`
}
