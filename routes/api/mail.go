package api

import (
	"github.com/gin-gonic/gin"
	mailController "mail-qq/controllers/mails"
)

func LoadMailerRoutes(r *gin.RouterGroup) {

	mails := r.Group("/mails")
	{
		mails.GET("/", mailController.ListEmails)
		mails.POST("/", mailController.SendEmail)
	}
}
