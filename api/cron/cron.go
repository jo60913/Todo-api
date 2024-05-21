package api

import (
	"log"
	"net/http"
	// . "github.com/tbxark/g4vercel"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	log.Println("執行定時任務 5/20")
	// server := New()
	// server.POST("/cron/get", func(ctx *Context) {
	// 	log.Println("執行定時任務 /cron/get api")
	// 	ctx.JSON(http.StatusOK, H{
	// 		"ErrorMsg":  "定時任務執行完成",
	// 		"ErrorFlag": "0",
	// 	})
	// })
	// server.Handle(w, r)
}

// func handler(w http.ResponseWriter, r *http.Request) {
// 	log.Println("執行定時任務 方法 5/12 17:18設定 5/13 7:30執行")
// 	server := New()
// 	server.POST("/cron/get", func(ctx *Context) {
// 		log.Println("執行定時任務 /cron/get api")
// 		ctx.JSON(http.StatusOK, H{
// 			"ErrorMsg":  "定時任務執行完成",
// 			"ErrorFlag": "0",
// 		})
// 	})
// 	server.Handle(w, r)
// }
