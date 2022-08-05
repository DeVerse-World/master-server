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

func (ctrl *EventController) GetAll(c *gin.Context) {
	const (
		success = "Get All Events successfully"
		failed  = "Get All Events unsuccessfully"
	)

	events, err := model.GetAllEvents()
	if err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}
	JSONReturn(c, http.StatusOK, success, gin.H{
		"events": events,
	})
}

func (ctrl *EventController) Get(c *gin.Context) {
	const (
		success = "Get Event successfully"
		failed  = "Get Event unsuccessfully"
	)

	var event model.Event
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}
	if err := event.GetById(id); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}
	participants, err := model.GetEventParticipants(event.ID)
	if err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}
	JSONReturn(c, http.StatusOK, success, gin.H{
		"event":        event,
		"participants": participants,
	})

}

func (ctrl *EventController) Create(c *gin.Context) {
	const (
		success = "Create Event successfully"
		failed  = "Create Event unsuccessfully"
	)

	authU, _ := jwt.HandleUserCookie(c.Writer, c.Request)
	if authU == nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, errors.New("can't find token"))
		return
	}
	var user model.User
	user.GetUserById(strconv.FormatUint(uint64(authU.ID), 10))

	var req requestSchema.CreateEvent
	if err := c.BindJSON(&req); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	var event = req.Event
	event.Stage = model.EVENT_UNSTARTED
	event.UserId = &user.ID
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

func (ctrl *EventController) Update(c *gin.Context) {
	const (
		success = "Update Event successfully"
		failed  = "Update Event unsuccessfully"
	)

	var req requestSchema.UpdateEvent
	if err := c.BindJSON(&req); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	authU, _ := jwt.HandleUserCookie(c.Writer, c.Request)
	if authU == nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, errors.New("can't find token"))
		return
	}
	var user model.User
	user.GetUserById(strconv.FormatUint(uint64(authU.ID), 10))

	var event model.Event
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}
	if err := event.GetById(id); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	if *event.UserId != user.ID {
		abortWithStatusError(c, http.StatusBadRequest, failed, errors.New("unauthorized to update other's event"))
		return
	}

	event.Name = req.Name
	event.EventConfigUri = req.EventConfigUri
	event.MaxNumParticipants = req.MaxNumParticipants
	event.AllowTemporaryHold = req.AllowTemporaryHold
	if err := event.Update(); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	JSONReturn(c, http.StatusOK, success, gin.H{
		"event": event,
	})
}

func (ctrl *EventController) Delete(c *gin.Context) {
	const (
		success = "Delete Event successfully"
		failed  = "Delete Event unsuccessfully"
	)

	authU, _ := jwt.HandleUserCookie(c.Writer, c.Request)
	if authU == nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, errors.New("can't find token"))
		return
	}
	var user model.User
	user.GetUserById(strconv.FormatUint(uint64(authU.ID), 10))

	var event model.Event
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}
	if err := event.GetById(id); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	if *event.UserId != user.ID {
		abortWithStatusError(c, http.StatusBadRequest, failed, errors.New("unauthorized to delete other's event"))
		return
	}
	if err := event.Delete(); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	JSONReturn(c, http.StatusOK, success, gin.H{})
}

func (ctrl *EventController) Start(c *gin.Context) {
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
	if event.Stage != model.EVENT_UNSTARTED && event.Stage != model.EVENT_PAUSED {
		abortWithStatusError(c, http.StatusBadRequest, failed, errors.New("can only start un-started/ paused event"))
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

func (ctrl *EventController) Pause(c *gin.Context) {
	const (
		success = "Pause Event successfully"
		failed  = "Pause Event unsuccessfully"
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
	if event.Stage != model.EVENT_IN_PROGRESS {
		abortWithStatusError(c, http.StatusBadRequest, failed, errors.New("can only pause in-progress event"))
		return
	}
	event.Stage = model.EVENT_PAUSED
	if err := event.Update(); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	JSONReturn(c, http.StatusOK, success, gin.H{
		"event": event,
	})
}

func (ctrl *EventController) Stop(c *gin.Context) {
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

func (ctrl *EventController) Join(c *gin.Context) {
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
	if event.Stage != model.EVENT_IN_PROGRESS {
		abortWithStatusError(c, http.StatusBadRequest, failed, errors.New("event not at join-able stage"))
		return
	}

	authU, _ := jwt.HandleUserCookie(c.Writer, c.Request)
	if authU == nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, errors.New("can't find token"))
		return
	}
	var user model.User
	user.GetUserById(strconv.FormatUint(uint64(authU.ID), 10))

	// TODO: CHeck if exceed max_num_participants
	var eventParticipant model.EventParticipant
	eventParticipant.EventId = &eventId
	eventParticipant.UserId = &user.ID
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

func (ctrl *EventController) UpdateScore(c *gin.Context) {
	const (
		success = "Update Participant Score successfully"
		failed  = "Update Participant Score unsuccessfully"
	)
	authU, _ := jwt.HandleUserCookie(c.Writer, c.Request)
	if authU == nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, errors.New("can't find token"))
		return
	}
	var user model.User
	user.GetUserById(strconv.FormatUint(uint64(authU.ID), 10))

	var eventIdStr = c.Param("id")
	eventId, err := strconv.Atoi(eventIdStr)
	if err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	var req requestSchema.UpdateScore
	if err := c.BindJSON(&req); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	var eventParticipant model.EventParticipant
	if err := eventParticipant.GetParticipant(user.ID, eventId); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}
	eventParticipant.Score = req.Score
	if err := eventParticipant.Update(); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}
	JSONReturn(c, http.StatusOK, success, gin.H{
		"eventParticipant": eventParticipant,
	})
}

func (ctrl *EventController) Rejoin(c *gin.Context) {
	const (
		success = "Rejoin Event successfully"
		failed  = "Rejoin Event unsuccessfully"
	)

	authU, _ := jwt.HandleUserCookie(c.Writer, c.Request)
	if authU == nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, errors.New("can't find token"))
		return
	}
	var user model.User
	user.GetUserById(strconv.FormatUint(uint64(authU.ID), 10))

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
	if event.Stage != model.EVENT_IN_PROGRESS {
		abortWithStatusError(c, http.StatusBadRequest, failed, errors.New("event not at join-able stage"))
		return
	}

	var eventParticipant model.EventParticipant
	if err := eventParticipant.GetParticipant(user.ID, eventId); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}
	eventParticipant.Score = 0
	if err := eventParticipant.Update(); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	JSONReturn(c, http.StatusOK, success, gin.H{
		"eventParticipant": eventParticipant,
	})
}

func (ctrl *EventController) Exit(c *gin.Context) {
	const (
		success = "Exit Event successfully"
		failed  = "Exit Event unsuccessfully"
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

	authU, _ := jwt.HandleUserCookie(c.Writer, c.Request)
	if authU == nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, errors.New("can't find token"))
		return
	}
	var user model.User
	user.GetUserById(strconv.FormatUint(uint64(authU.ID), 10))

	var eventParticipant model.EventParticipant
	if err := eventParticipant.GetParticipant(user.ID, eventId); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	if err := eventParticipant.Delete(); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	JSONReturn(c, http.StatusOK, success, gin.H{})
}
