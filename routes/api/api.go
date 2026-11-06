package api

import (
	"github.com/gin-gonic/gin"
)

func LoadAPIRoutes(r *gin.Engine) {

	api := r.Group("/api")
	{
		LoadMailerRoutes(api)
	}
}
