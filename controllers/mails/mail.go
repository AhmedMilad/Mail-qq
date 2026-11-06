package mails

import (
	"encoding/json"
	"log"
	"mail-qq/queues"
	"mail-qq/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListEmails(c *gin.Context) {

	var resp []utils.Response = []utils.Response{
		{
			Message: "Here is your message",
		},
	}

	c.JSON(http.StatusOK, resp)
}

func SendEmail(c *gin.Context) {

	var reqData utils.SendEmailRequest

	if err := c.ShouldBind(&reqData); err != nil {
		// TODO handle validation message response
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := json.Marshal(reqData)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := queues.Publish(data); err != nil {
		log.Printf("could not queue email: %s", err)
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "the mail queue is unavailable, please try again later"})
		return
	}

	var resp utils.Response = utils.Response{
		Message: "The message has been sent successfully.",
	}

	c.JSON(http.StatusAccepted, resp)
}
