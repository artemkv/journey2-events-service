package app

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type actionIncomingData struct {
	AccId  string `json:"acc" binding:"required"`
	AppId  string `json:"aid" binding:"required"`
	UserId string `json:"uid" binding:"required"`
	Action string `json:"act" binding:"required"`
	Param  string `json:"par"`
}

type actionOutgoingData struct {
	AccId  string `json:"acc" binding:"required"`
	AppId  string `json:"aid" binding:"required"`
	UserId string `json:"uid" binding:"required"`
	Action string `json:"act" binding:"required"`
	Param  string `json:"par"`
	Date   string `json:"dts"`
}

func handlePostAction(c *gin.Context) {
	var actionIn actionIncomingData
	if err := c.ShouldBindJSON(&actionIn); err != nil {
		toBadRequest(c, err)
		return
	}

	// TODO: this is temporary solution to avoid high billing
	if !IsWhitelisted(actionIn.AppId) {
		toBadRequest(c, fmt.Errorf("invalid aid: %s", actionIn.AppId))
		return
	}

	actionOut := constructActionOut(&actionIn)

	msgId, err := EnqueueAction(actionOut)
	if err != nil {
		toInternalServerError(c, err.Error())
		return
	}

	toSuccess(c, msgId)
}

func constructActionOut(actionIn *actionIncomingData) *actionOutgoingData {
	return &actionOutgoingData{
		AccId:  actionIn.AccId,
		AppId:  actionIn.AppId,
		UserId: actionIn.UserId,
		Action: actionIn.Action,
		Param:  actionIn.Param,
		Date:   time.Now().UTC().Format(time.RFC3339),
	}
}
