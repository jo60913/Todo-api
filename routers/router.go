package routers

import "github.com/gin-gonic/gin"

func InitRouters(r *gin.Engine) {
	InitApi(r)
}
