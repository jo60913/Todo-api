package api

import (
	"encoding/json"
	"log"
	"net/http"

	. "github.com/tbxark/g4vercel"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	server := New()

	server.POST("/update/notification", func(ctx *Context) {
		var notificationUpdate NotificationUpdate
		err := json.NewDecoder(ctx.Req.Body).Decode(&notificationUpdate)
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
		if err != nil {
			ctx.JSON(http.StatusOK, H{
				"ErrorMsg":  "欄位錯誤",
				"ErrorFlag": "3",
			})
			return
		}
		ctx.JSON(http.StatusOK, H{
			"ErrorMsg":  "",
			"ErrorFlag": "1",
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
