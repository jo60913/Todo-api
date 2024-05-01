package routers

import (
	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/jo60913/Todo-api/web"
)

func InitApi(r *gin.Engine, client *firestore.Client) {
	api := r.Group("/api")
	api.POST("/update/notification", func(c *gin.Context) {
		web.UpdateNotification(c, client)
	})

	api.POST("/get/notification", func(c *gin.Context) {
		web.GetNotification(c, client)
	})
}
