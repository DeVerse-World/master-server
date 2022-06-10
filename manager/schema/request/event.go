package request

import "github.com/hyperjiang/gin-skeleton/model"

type CreateEvent struct {
	Event        model.Event         `json:"event" binding:"required"`
	EventRewards []model.EventReward `json:"event_rewards" binding:"required"`
}
