package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"github.com/jo60913/Todo-api/routers"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

func main() {
	log.SetOutput(os.Stderr)
	log.Print("開始執行")
	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatal("Error loading .env file", envErr)
	}
	firebaseAdminSdk := os.Getenv("FIREBASE_ADMIN_SDK")

	sa := option.WithCredentialsJSON([]byte(firebaseAdminSdk))
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

	// 每3秒执行一次
	// s.Every(3).Seconds().Do(TriggerFCM)
	s.Every(1).Day().At("08:00").Do(TriggerFCM)
	s.StartBlocking()
}

func TriggerFCM() {
	fmt.Println("八點執行FCM")
}
