package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
)

var apiUrl = "https://todo-api"

func main() {
	// log.SetOutput(os.Stderr)
	// log.Print("開始執行")
	// envErr := godotenv.Load()
	// if envErr != nil {
	// 	log.Fatal("Error loading .env file", envErr)
	// }
	// firebaseAdminSdk := os.Getenv("FIREBASE_ADMIN_SDK")

	// sa := option.WithCredentialsFile("todo-app-firebase-adminsdk.json")
	// app, newAppErr := firebase.NewApp(context.Background(), nil, sa)
	// if newAppErr != nil {
	// 	fmt.Println("firebase.NewApp錯誤")
	// }

	// client, err := app.Firestore(context.Background())
	// if err != nil {
	// 	fmt.Println("firestore登入錯誤", err.Error())
	// }

	// defer client.Close()
	// server := gin.Default()
	// routers.InitRouters(server, client)
	// go SetFCMSetting()
	// server.Run()
	router := gin.Default()
	router.POST("/api", func(c *gin.Context) {
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
			"ErrorMsg":  "成功" + notificationUpdate.UserToken,
			"ErrorFlag": "1",
		})

		// _, readErr := client.Collection(notificationUpdate.UserToken).Doc("notification").Get(context.Background())
		// if readErr != nil { //沒有notification時新增
		// 	_, addErr := client.Collection(notificationUpdate.UserToken).Doc("notification").Create(context.Background(), map[string]interface{}{
		// 		"FCM": notificationUpdate.NotificationValue,
		// 	})
		// 	if addErr != nil {
		// 		c.JSON(http.StatusOK, gin.H{
		// 			"ErrorMsg":  "新增notification時錯誤",
		// 			"ErrorFlag": "2",
		// 		})
		// 		return
		// 	}
		// 	c.JSON(http.StatusOK, gin.H{
		// 		"ErrorMsg":  "",
		// 		"ErrorFlag": "0",
		// 	})
		// 	return
		// }

		// _, updateErr := client.Collection(notificationUpdate.UserToken).Doc("notification").Update(context.Background(), []firestore.Update{
		// 	{Path: "FCM", Value: notificationUpdate.NotificationValue},
		// })

		// if updateErr != nil {
		// 	c.JSON(http.StatusOK, gin.H{
		// 		"ErrorMsg":  "修改時發生錯誤",
		// 		"ErrorFlag": "2",
		// 	})
		// 	return
		// }

		// c.JSON(http.StatusOK, gin.H{
		// 	"ErrorMsg":  "",
		// 	"ErrorFlag": "0",
		// })

	})
	router.Run()
}

type UpdateNotification struct {
	Name   string `json:"name"`
	Switch bool   `json:"switch"`
}

func SetFCMSetting() {
	timezone, _ := time.LoadLocation("Asia/Taipei")
	s := gocron.NewScheduler(timezone)

	// 每3秒执行一次
	// s.Every(3).Seconds().Do(TriggerFCM)
	s.Every(1).Day().At("08:00").Do(TriggerFCM)
	s.StartBlocking()
}

func TriggerFCM() {
	fmt.Println("八點執行FCM")
}

type NotificationUpdate struct {
	UserToken         string `json:"UserToken"`
	NotificationValue bool   `json:"NotificationValue"`
}
