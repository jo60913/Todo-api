package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	. "github.com/tbxark/g4vercel"
	"google.golang.org/api/option"
)

var (
	firebaseSdkAdmin = os.Getenv("FIREBASE_ADMIN_SDK")
)

func Handler(w http.ResponseWriter, r *http.Request) {
	server := New()
	sa := option.WithCredentialsJSON([]byte(firebaseSdkAdmin))
	app, newAppErr := firebase.NewApp(context.Background(), nil, sa)
	if newAppErr != nil {
		log.Println("firebase.NewApp錯誤", newAppErr)
		return
	}

	client, err := app.Firestore(context.Background())
	if err != nil {
		log.Println("firestore登入錯誤", err.Error())
		return
	}
	defer client.Close()

	server.POST("/update/notification", func(ctx *Context) {
		var notificationUpdate NotificationUpdate
		err := json.NewDecoder(ctx.Req.Body).Decode(&notificationUpdate)
		log.Println("/update/notification ", "UserToken : "+notificationUpdate.UserID)
		if err != nil {
			log.Println("/update/notification 傳入參數錯誤", err.Error())
			ctx.JSON(http.StatusOK, H{
				"ErrorMsg":  "參數錯誤",
				"ErrorFlag": "1",
			})
			return
		}

		_, readErr := client.Collection(notificationUpdate.UserID).Doc("notification").Get(context.Background())
		if readErr != nil { //沒有notification時新增
			log.Println("/update/notification 找notification時錯誤", readErr.Error())
			_, addErr := client.Collection(notificationUpdate.UserID).Doc("notification").Create(context.Background(), map[string]interface{}{
				"FCM": notificationUpdate.NotificationValue,
			})
			if addErr != nil {
				log.Println("/update/notification 新增notification時錯誤", addErr.Error())
				ctx.JSON(http.StatusOK, H{
					"ErrorMsg":  "新增notification時錯誤",
					"ErrorFlag": "2",
				})
				return
			}
			log.Println("/update/notification 新增notification成功")
			ctx.JSON(http.StatusOK, H{
				"ErrorMsg":  "",
				"ErrorFlag": "0",
			})
			return
		}

		_, updateErr := client.Collection(notificationUpdate.UserID).Doc("notification").Update(context.Background(), []firestore.Update{
			{Path: "FCM", Value: notificationUpdate.NotificationValue},
		})

		if updateErr != nil {
			log.Println("/update/notification 更新FCM時錯誤" + updateErr.Error())
			ctx.JSON(http.StatusOK, H{
				"ErrorMsg":  "修改時發生錯誤",
				"ErrorFlag": "2",
			})
			return
		}

		log.Println("/update/notification 更新notification成功")
		ctx.JSON(http.StatusOK, H{
			"ErrorMsg":  "",
			"ErrorFlag": "0",
			"Data":      notificationUpdate.NotificationValue,
		})
	})

	server.POST("/get/notification", func(ctx *Context) {
		var notificationGet NotificationGet
		err := json.NewDecoder(ctx.Req.Body).Decode(&notificationGet)

		log.Println("get/notification ", "UserToken : "+notificationGet.UserID)
		if err != nil {
			log.Println("get/notification錯誤", "json轉換錯誤 "+err.Error())
			ctx.JSON(http.StatusOK, H{
				"ErrorMsg":  "欄位錯誤",
				"ErrorFlag": "3",
			})
			return
		}
		_, readErr := client.Collection(notificationGet.UserID).Doc("notification").Get(context.Background())
		if readErr != nil { //沒有notificati
			log.Println("get/notification 尚未新增", "新增FCM欄位"+readErr.Error())
			_, addErr := client.Collection(notificationGet.UserID).Doc("notification").Create(context.Background(), map[string]interface{}{
				"FCM": true,
			})
			if addErr != nil {
				ctx.JSON(http.StatusOK, H{
					"ErrorMsg":  "新增notification時錯誤",
					"ErrorFlag": "2",
				})
				return
			}
		}

		notificationDoc := client.Collection(notificationGet.UserID).Doc("notification")
		getvalue, getDocError := notificationDoc.Get(context.Background())

		if getDocError != nil {
			log.Println("get/notification notification欄位找不到", getDocError.Error())
			ctx.JSON(http.StatusOK, H{
				"ErrorMsg":  "找不到文件",
				"ErrorFlag": "2",
			})
			return
		}

		fcmValue, getFcmValueError := getvalue.DataAt("FCM")

		if getFcmValueError != nil {
			log.Println("get/notification ", "找不到FCM屬性")
			ctx.JSON(http.StatusOK, H{
				"ErrorMsg":  "找不到FCM屬性",
				"ErrorFlag": "2",
			})
			return
		}

		log.Println("get/notification ", "成功")
		ctx.JSON(http.StatusOK, H{
			"ErrorMsg":  "",
			"ErrorFlag": "0",
			"Data":      fcmValue,
		})
	})

	server.POST("/update/firstlogin", func(ctx *Context) {
		var firstlogin FirstLogin
		err := json.NewDecoder(ctx.Req.Body).Decode(&firstlogin)

		log.Println("update/firstlogin ", "UserToken : "+firstlogin.UserToken)
		if err != nil {
			log.Println("update/firstlogin", "json轉換錯誤 "+err.Error())
			ctx.JSON(http.StatusOK, H{
				"ErrorMsg":  "欄位錯誤",
				"ErrorFlag": "3",
			})
			return
		}
		_, readErr := client.Collection(firstlogin.UserID).Doc("notification").Get(context.Background())
		if readErr != nil { //沒有notification
			log.Println("update/firstlogin 尚未新增", "新增FCM欄位"+readErr.Error())
			_, addErr := client.Collection(firstlogin.UserToken).Doc("notification").Create(context.Background(), map[string]interface{}{
				"FCM":      true,
				"FCMToken": firstlogin.UserToken,
			})
			if addErr != nil {
				ctx.JSON(http.StatusOK, H{
					"ErrorMsg":  "新增notification FCM時錯誤",
					"ErrorFlag": "2",
				})
				return
			}
		}

		notificationDoc := client.Collection(firstlogin.UserID).Doc("notification")
		getvalue, getDocError := notificationDoc.Get(context.Background())

		if getDocError != nil {
			log.Println("update/firstlogin notification欄位找不到", getDocError.Error())
			ctx.JSON(http.StatusOK, H{
				"ErrorMsg":  "找不到文件",
				"ErrorFlag": "2",
			})
			return
		}

		fcmValue, getFcmValueError := getvalue.DataAt("FCM")

		if getFcmValueError != nil {
			log.Println("update/firstlogin ", "找不到FCM屬性")
			ctx.JSON(http.StatusOK, H{
				"ErrorMsg":  "找不到FCM屬性",
				"ErrorFlag": "2",
			})
			return
		}

		log.Println("update/firstlogin ", "成功")
		ctx.JSON(http.StatusOK, H{
			"ErrorMsg":  "",
			"ErrorFlag": "0",
			"Data":      fcmValue,
		})
	})

	server.Handle(w, r)
}

type NotificationUpdate struct {
	UserID            string `json:"UserID"`
	NotificationValue bool   `json:"NotificationValue"`
}

type NotificationGet struct {
	UserID string `json:"UserID"`
}

type FirstLogin struct {
	UserID    string `json:"UserID"`
	UserToken string `json:"UserToken"`
}
