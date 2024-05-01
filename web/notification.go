package web

import (
	"net/http"

	"context"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

func UpdateNotification(c *gin.Context, client *firestore.Client) {
	var notificationUpdate NotificationUpdate
	err := c.ShouldBindJSON(&notificationUpdate)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"ErrorMsg":  "欄位錯誤",
			"ErrorFlag": "3",
		})
		return
	}

	_, readErr := client.Collection(notificationUpdate.UserToken).Doc("notification").Get(context.Background())
	if readErr != nil { //沒有notification時新增
		_, addErr := client.Collection(notificationUpdate.UserToken).Doc("notification").Create(context.Background(), map[string]interface{}{
			"FCM": notificationUpdate.NotificationValue,
		})
		if addErr != nil {
			c.JSON(http.StatusOK, gin.H{
				"ErrorMsg":  "新增notification時錯誤",
				"ErrorFlag": "2",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"ErrorMsg":  "",
			"ErrorFlag": "0",
		})
		return
	}

	_, updateErr := client.Collection(notificationUpdate.UserToken).Doc("notification").Update(context.Background(), []firestore.Update{
		{Path: "FCM", Value: notificationUpdate.NotificationValue},
	})

	if updateErr != nil {
		c.JSON(http.StatusOK, gin.H{
			"ErrorMsg":  "修改時發生錯誤",
			"ErrorFlag": "2",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ErrorMsg":  "",
		"ErrorFlag": "0",
	})
}

func GetNotification(c *gin.Context, client *firestore.Client) {
	var notificationGet NotificationGet
	err := c.ShouldBindJSON(&notificationGet)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"ErrorMsg":  "欄位錯誤",
			"ErrorFlag": "3",
		})
		return
	}

	_, readErr := client.Collection(notificationGet.UserToken).Doc("notification").Get(context.Background())
	if readErr != nil { //沒有notification時新增
		_, addErr := client.Collection(notificationGet.UserToken).Doc("notification").Create(context.Background(), map[string]interface{}{
			"FCM": true,
		})
		if addErr != nil {
			c.JSON(http.StatusOK, gin.H{
				"ErrorMsg":  "新增notification時錯誤",
				"ErrorFlag": "2",
			})
			return
		}
	}

	notificationDoc := client.Collection(notificationGet.UserToken).Doc("notification")
	getvalue, getDocError := notificationDoc.Get(context.Background())

	if getDocError != nil {
		c.JSON(http.StatusOK, gin.H{
			"ErrorMsg":  "找不到文件",
			"ErrorFlag": "2",
		})
		return
	}

	fcmValue, getFcmValueError := getvalue.DataAt("FCM")

	if getFcmValueError != nil {
		c.JSON(http.StatusOK, gin.H{
			"ErrorMsg":  "找不到FCM屬性",
			"ErrorFlag": "2",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ErrorMsg":  "",
		"ErrorFlag": "0",
		"Data":      fcmValue,
	})
}

type NotificationUpdate struct {
	UserToken         string `json:"UserToken"`
	NotificationValue bool   `json:"NotificationValue"`
}

type NotificationGet struct {
	UserToken string `json:"UserToken"`
}
