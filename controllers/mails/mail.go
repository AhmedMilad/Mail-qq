package mails

import (
	"encoding/json"
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

	// publish to the mailing queue
	queues.Publish(data)

	var resp utils.Response = utils.Response{
		Message: "The message has been sent successfully.",
	}

	c.JSON(http.StatusOK, resp)
}
