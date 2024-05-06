package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Handler(c *gin.Context) {
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
		"ErrorMsg":  "成功執行",
		"ErrorFlag": "0",
	})
}

type NotificationUpdate struct {
	UserToken         string `json:"UserToken"`
	NotificationValue bool   `json:"NotificationValue"`
}

type NotificationGet struct {
	UserToken string `json:"UserToken"`
}
