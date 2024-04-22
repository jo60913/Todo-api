package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jo60913/Todo-api/web"
)

func InitApi(r *gin.Engine) {
	api := r.Group("/api")
	api.POST("/login", web.Login)
}
