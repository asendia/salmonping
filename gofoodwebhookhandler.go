package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// historyHandler godoc
//
// @Summary		Show salmon ping history
// @Description	get ping history based on query string params
// @Tags		ping
// @Accept		json
// @Produce		json
// @Security	GofoodSignature
// @Param		request	body	GofoodWebhookPayload	true	"Webhook Payload sent by Gofood server"
// @Success		200	{object}	DefaultResponse
// @Failure		400	{object}	DefaultErrorResponse
// @Router		webhook/gofood	[post]
func gofoodWebhookHandler(c *gin.Context) {
	var payload GofoodWebhookPayload
	err := c.BindJSON(&payload)
	if err != nil {
		log := DefaultErrorResponse{
			Error:   err.Error(),
			Level:   "error",
			Message: "Error parsing body",
		}
		logJson(log.JSON())
		c.JSON(http.StatusBadRequest, log)
		return
	}
	log := DefaultErrorResponse{
		Header:  c.Request.Header,
		Level:   "info",
		Message: "Gofood Webhook received",
		Payload: payload,
	}
	logJson(log.JSON())

	c.JSON(http.StatusOK, DefaultResponse{
		Message: "ok",
	})
}
