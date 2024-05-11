package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	firebase "firebase.google.com/go"
	. "github.com/tbxark/g4vercel"
	"google.golang.org/api/option"
)

var (
	firebaseSdkAdmin = os.Getenv("FIREBASE_ADMIN_SDK")
)

func Handler(w http.ResponseWriter, r *http.Request) {
	server := New()
	log.Println("firestore admin 內容", firebaseSdkAdmin)
	sa := option.WithCredentialsJSON([]byte(firebaseSdkAdmin))
	app, newAppErr := firebase.NewApp(context.Background(), nil, sa)
	if newAppErr != nil {
		log.Fatal("firebase.NewApp錯誤", newAppErr)
		return
	}

	client, err := app.Firestore(context.Background())
	if err != nil {
		log.Fatal("firestore登入錯誤", err.Error())
		return
	}
	defer client.Close()

	server.POST("/update/notification", func(ctx *Context) {
		var notificationUpdate NotificationUpdate
		err := json.NewDecoder(ctx.Req.Body).Decode(&notificationUpdate)
		log.Println("/update/notification ", "UserToken : "+notificationUpdate.UserToken)
		if err != nil {
			log.Fatal("錯誤", err.Error())
			return
		}

		ctx.JSON(http.StatusOK, H{
			"ErrorMsg":  "",
			"ErrorFlag": "0",
		})
	})

	server.POST("/get/notification", func(ctx *Context) {
		var notificationGet NotificationGet
		err := json.NewDecoder(ctx.Req.Body).Decode(&notificationGet)

		log.Println("get/notification ", "UserToken : "+notificationGet.UserToken)
		if err != nil {
			log.Fatal("get/notification錯誤", "json轉換錯誤 "+err.Error())
			ctx.JSON(http.StatusOK, H{
				"ErrorMsg":  "欄位錯誤",
				"ErrorFlag": "3",
			})
			return
		}
		_, readErr := client.Collection(notificationGet.UserToken).Doc("notification").Get(context.Background())
		if readErr != nil { //沒有notificati
			log.Fatal("get/notification 尚未新增", "新增FCM欄位"+readErr.Error())
			_, addErr := client.Collection(notificationGet.UserToken).Doc("notification").Create(context.Background(), map[string]interface{}{
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

		notificationDoc := client.Collection(notificationGet.UserToken).Doc("notification")
		getvalue, getDocError := notificationDoc.Get(context.Background())

		if getDocError != nil {
			log.Fatal("get/notification notification欄位找不到", getDocError.Error())
			ctx.JSON(http.StatusOK, H{
				"ErrorMsg":  "找不到文件",
				"ErrorFlag": "2",
			})
			return
		}

		fcmValue, getFcmValueError := getvalue.DataAt("FCM")

		if getFcmValueError != nil {
			log.Fatal("get/notification ", "找不到FCM屬性")
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
	server.Handle(w, r)
}

type NotificationUpdate struct {
	UserToken         string `json:"UserToken"`
	NotificationValue bool   `json:"NotificationValue"`
}

type NotificationGet struct {
	UserToken string `json:"UserToken"`
}
