package app

import (
	"artemkv.net/journey2/reststats"
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

	// TODO: now simply returns input
	toSuccess(c, action)

	// stats
	reststats.CountRequestByEndpoint("action")
	reststats.UpdateResponseStats(c.Writer.Status())
}
