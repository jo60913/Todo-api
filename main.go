package main

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
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

	server.Run()
}

type UpdateNotification struct {
	Name   string `json:"name"`
	Switch bool   `json:"switch"`
}
