package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Handler(r *gin.Engine) {
	api := r.Group("/api")
	api.GET("/get", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"ErrorMsg":  "欄位錯誤",
			"ErrorFlag": "3",
		})
	})
}
