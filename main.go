package main

import (
	"context"
	"fmt"
	"time"

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"github.com/jo60913/Todo-api/routers"
	"google.golang.org/api/option"
)

func main() {
	sa := option.WithCredentialsFile("todo-app-firebase-adminsdk.json")
	app, newAppErr := firebase.NewApp(context.Background(), nil, sa)
	if newAppErr != nil {
		fmt.Println("firebase.NewApp錯誤")
	}

	client, err := app.Firestore(context.Background())
	if err != nil {
		fmt.Println("firestore登入錯誤", err.Error())
	}

	defer client.Close()
	server := gin.Default()
	routers.InitRouters(server, client)
	go SetFCMSetting()
	server.Run()
}

type UpdateNotification struct {
	Name   string `json:"name"`
	Switch bool   `json:"switch"`
}

func SetFCMSetting() {
	timezone, _ := time.LoadLocation("Asia/Taipei")
	s := gocron.NewScheduler(timezone)

	s.Every(1).Day().At("08:00").Do(TriggerFCM)
	s.StartBlocking()
}

func TriggerFCM() {
	//ToDo八點要執行的FCM
}
