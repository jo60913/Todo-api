package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	. "github.com/tbxark/g4vercel"
	"google.golang.org/api/option"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// _, filename, _, _ := runtime.Caller(0)
	// currentDir := filepath.Dir(filename)
	// currentfiles, currenterr := os.ReadDir(currentDir)
	// fmt.Print("當前路徑", currentfiles)
	// if currenterr != nil {
	// 	fmt.Println("讀取目錄失敗:", currenterr)
	// 	return
	// }
	// // 獲取當前路徑的上一層目錄
	// parentDir := filepath.Dir(currentDir)

	// // 列出上一層目錄下的所有檔案和目錄
	// files, err := os.ReadDir(parentDir)
	// if err != nil {
	// 	fmt.Println("讀取目錄失敗:", err)
	// 	return
	// }

	// for _, file := range files {
	// 	log.Fatal("檔名", file.Name())
	// 	fmt.Println(file.Name())
	// }

	// data, readErr := ioutil.ReadFile(fileName)
	// if readErr != nil {
	// 	log.Fatal("讀取json檔案錯誤", readErr)
	// }
	// log.Fatal(string(data))

	server := New()
	sa := option.WithCredentialsFile("todo-app-firebase-adminsdk.json")
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
		_, readErr := client.Collection(notificationGet.UserToken).Doc("notification").Get(context.Background())
		if readErr != nil { //沒有notification時新增
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
			ctx.JSON(http.StatusOK, H{
				"ErrorMsg":  "找不到文件",
				"ErrorFlag": "2",
			})
			return
		}

		fcmValue, getFcmValueError := getvalue.DataAt("FCM")

		if getFcmValueError != nil {
			ctx.JSON(http.StatusOK, H{
				"ErrorMsg":  "找不到FCM屬性",
				"ErrorFlag": "2",
			})
			return
		}

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
