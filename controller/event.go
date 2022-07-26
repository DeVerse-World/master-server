package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/hyperjiang/gin-skeleton/manager/jwt"
	requestSchema "github.com/hyperjiang/gin-skeleton/manager/schema/request"
	"github.com/hyperjiang/gin-skeleton/model"
)

type EventController struct {
}

func NewEventController() *EventController {
	return &EventController{}
}

func (ctrl *EventController) CreateEvent(c *gin.Context) {
	const (
		success = "Create Event successfully"
		failed  = "Create Event unsuccessfully"
	)

	var req requestSchema.CreateEvent
	if err := c.BindJSON(&req); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	authW, _ := jwt.HandleUserCookie(c.Writer, c.Request)
	if authW == nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, errors.New("can't find token"))
		return
	}
	var wallet model.Wallet
	wallet.GetWalletByAddress(authW.Address)

	var event = req.Event
	event.Stage = model.EVENT_UNSTARTED
	event.WalletId = &wallet.ID
	if err := event.Create(); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	var eventRewards = req.EventRewards
	for i := 0; i < len(eventRewards); i++ {
		var eventReward = &eventRewards[i]
		eventReward.EventId = &event.ID
		if err := eventReward.Create(); err != nil {
			abortWithStatusError(c, http.StatusBadRequest, failed, err)
			return
		}
	}

	JSONReturn(c, http.StatusOK, success, gin.H{
		"event":         event,
		"event_rewards": eventRewards,
	})
}

func (ctrl *EventController) StartEvent(c *gin.Context) {
	const (
		success = "Start Event successfully"
		failed  = "Start Event unsuccessfully"
	)
	var eventIdStr = c.Param("id")
	eventId, err := strconv.Atoi(eventIdStr)
	if err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	var event model.Event
	if err := event.GetById(eventId); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}
	event.Stage = model.EVENT_IN_PROGRESS
	if err := event.Update(); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	JSONReturn(c, http.StatusOK, success, gin.H{
		"event": event,
	})
}

func (ctrl *EventController) StopEvent(c *gin.Context) {
	const (
		success = "Stop Event successfully"
		failed  = "Stop Event unsuccessfully"
	)
	var eventIdStr = c.Param("id")
	eventId, err := strconv.Atoi(eventIdStr)
	if err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	var event model.Event
	if err := event.GetById(eventId); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}
	event.Stage = model.EVENT_FINISHED
	if err := event.Update(); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	JSONReturn(c, http.StatusOK, success, gin.H{
		"event": event,
	})
}

func (ctrl *EventController) JoinEvent(c *gin.Context) {
	const (
		success = "Join Event successfully"
		failed  = "Join Event unsuccessfully"
	)
	var eventIdStr = c.Param("id")
	eventId, err := strconv.Atoi(eventIdStr)
	if err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	var event model.Event
	if err := event.GetById(eventId); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}
	if event.Stage == model.EVENT_UNSTARTED || event.Stage == model.EVENT_FINISHED {
		abortWithStatusError(c, http.StatusBadRequest, failed, errors.New("event not at join-able stage"))
		return
	}

	authW, _ := jwt.HandleUserCookie(c.Writer, c.Request)
	if authW == nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, errors.New("can't find token"))
		return
	}

	// TODO: CHeck if exceed max_num_participants
	var wallet model.Wallet
	wallet.GetWalletByAddress(authW.Address)
	var eventParticipant model.EventParticipant
	eventParticipant.EventId = &eventId
	eventParticipant.WalletId = &wallet.ID
	eventParticipant.Score = 0
	err = eventParticipant.Create()
	if err == model.ErrKeyConflict {
		abortWithStatusError(c, http.StatusMethodNotAllowed, failed, err)
		return
	} else if err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}
	JSONReturn(c, http.StatusOK, success, gin.H{
		"event_participant": eventParticipant,
	})
}
