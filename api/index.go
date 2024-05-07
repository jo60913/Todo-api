package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	router := gin.Default()
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "hellow",
		})
	})

	router.POST("/hello", func(c *gin.Context) {
		var notificationUpdate NotificationUpdate
		err := c.ShouldBindJSON(&notificationUpdate)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"ErrorMsg":  "欄位錯誤",
				"ErrorFlag": "3",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"ErrorMsg":  notificationUpdate.NotificationValue,
			"ErrorFlag": "0",
		})
	})
	router.ServeHTTP(w, r)
	router.Run()
}
