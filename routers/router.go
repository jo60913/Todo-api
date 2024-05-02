package routers

import (
	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

func InitRouters(r *gin.Engine, client *firestore.Client) {
	InitApi(r, client)
}
