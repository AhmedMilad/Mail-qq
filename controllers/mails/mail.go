package mails

import (
	"encoding/json"
	"log"
	"mail-qq/database"
	"mail-qq/models"
	"mail-qq/queues"
	"mail-qq/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListEmails(c *gin.Context) {

	var emailLog []models.EmailLog

	if err := database.Conn.Order("id desc").Find(&emailLog).Error; err != nil {
		message := "could not fetch the email log"
		log.Printf(message+" : %s", err)

		c.JSON(http.StatusInternalServerError, utils.Response{
			Message: message,
		})
		return
	}

	var logs []utils.EmailLogResponse

	for _, log := range emailLog {

		logs = append(logs, utils.EmailLogResponse{
			ID:        log.ID,
			Sender:    log.From,
			Recipient: log.To,
			Body:      log.Body,
			Subject:   log.Subject,
		})
	}

	var resp utils.Response = utils.Response{
		Message: "Fetched email log successfully.",
		Data:    logs,
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
