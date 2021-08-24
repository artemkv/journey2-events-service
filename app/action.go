package app

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type actionData struct {
	AppId  string `json:"aid" binding:"required"`
	UserId string `json:"uid" binding:"required"`
	Action string `json:"act" binding:"required"`
	Param  string `json:"par"`
}

func handlePostAction(c *gin.Context) {
	var action actionData
	if err := c.ShouldBindJSON(&action); err != nil {
		toBadRequest(c, err)
		return
	}

	// TODO: this is temporary solution to avoid high billing
	if !IsWhitelisted(action.AppId) {
		toBadRequest(c, fmt.Errorf("invalid aid: %s", action.AppId))
		return
	}

	msgId, err := EnqueueAction(action)
	if err != nil {
		toInternalServerError(c, err.Error())
		return
	}

	toSuccess(c, msgId)
}
