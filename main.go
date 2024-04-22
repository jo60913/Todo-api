package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jo60913/Todo-api/routers"
)

func main() {
	server := gin.Default()
	routers.InitRouters(server)

	server.Run()
}

type UpdateNotification struct {
	Name   string `json:"name"`
	Switch bool   `json:"switch"`
}
