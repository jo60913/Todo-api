package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	. "github.com/tbxark/g4vercel"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

var (
	firebaseSdkAdmin = os.Getenv("FIREBASE_ADMIN_SDK")
	firebasefcmkey   = os.Getenv("TODO_API_FIREBASE_FCM_KEY")
	fcmHeader        = os.Getenv("FCM_HEADER")
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
			_, addErr := client.Collection(firstlogin.UserID).Doc("notification").Create(context.Background(), map[string]interface{}{
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

	server.POST("/notification/fcm", func(ctx *Context) {
		authHeader := ctx.Req.Header.Get("FCMHeader")
		if authHeader != fcmHeader {
			ctx.JSON(http.StatusUnauthorized, H{
				"ErrorMsg":  "請設置header",
				"ErrorFlag": "2",
			})
			return
		}
		cctx := context.Background()
		collections := client.Collections(cctx)
		for {
			collectionRef, err := collections.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Fatalf("Failed to get collection: %v", err)
			}

			log.Printf("Collection: %s", collectionRef.ID)
			collection := client.Collection(collectionRef.ID)

			log.Printf("取得notification資料")
			notificationDoc := collection.Doc("notification")
			getvalue, getDocError := notificationDoc.Get(context.Background())
			if getDocError != nil {
				log.Printf("取得 notfication時 錯誤")
				continue
			}

			var nyData FcmInfo
			log.Printf("轉換notification的值")
			fcmError := getvalue.DataTo(&nyData)
			if fcmError != nil {
				log.Printf("轉換notification錯誤", fcmError)
			}
			if nyData.FcmValue {
				log.Printf("開始傳送fcm token 為", nyData.FCMToken)
				taskInfo := getTaskCount(collection)
				if hasIncompleteTesk(taskInfo) {
					//有代辦事項時傳送
					hasIncompleteTodos(nyData.FCMToken, taskInfo.inCompleteCount, taskInfo.totalCount)
				} else {
					// 沒有訊息時傳送
					getNoToDoListMessage(nyData.FCMToken)
				}
			}

		}

	})

	server.Handle(w, r)
}

func getNoToDoListMessage(FCMToken string) {
	data := map[string]interface{}{
		"to": FCMToken,
		"notification": map[string]string{
			"body":  "點擊通知新增事項",
			"title": "美好的一天開始 8點新增事項",
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("getNoToDoListMessage轉換JSON錯誤:", err)
		return
	}
	sendNotficationToUser(jsonData, FCMToken)
}

func hasIncompleteTodos(FCMToken string, inCompleteCount int, totalCount int) {
	fmt.Println("未完成數 :" + strconv.Itoa(inCompleteCount) + " 總數" + strconv.Itoa(totalCount))
	var successRate float32 = 0
	if float32(totalCount)-float32(inCompleteCount) > 0 {
		completedCount := totalCount - inCompleteCount
		successRate = float32(completedCount) / float32(totalCount)
	}
	data := map[string]interface{}{
		"to": FCMToken,
		"notification": map[string]string{
			"body":  fmt.Sprintf("點擊查看未完成任務 完成率為%.1f%%", successRate*100),
			"title": "加油 目前還有" + strconv.Itoa(inCompleteCount) + "個未完成任務",
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("hasIncompleteTodos轉換JSON錯誤:", err)
		return
	}
	sendNotficationToUser(jsonData, FCMToken)
}

func hasIncompleteTesk(taskinfo TaskInfo) bool {
	return taskinfo.inCompleteCount > 0
}

func getTaskCount(collection *firestore.CollectionRef) TaskInfo {
	docsRefs, docsErr := collection.Documents(context.Background()).GetAll()
	inCompleteCount := 0
	totalCount := 0
	if docsErr != nil {
		return TaskInfo{
			inCompleteCount: inCompleteCount,
			totalCount:      totalCount,
		}
	}

	for _, element := range docsRefs {
		if element.Ref.ID == "notification" {
			continue
		}

		fmt.Println("element : " + element.Ref.ID)
		todoEntry := element.Ref.Collection("todo-entries")
		todoEntryDocs, todoentryErr := todoEntry.Documents(context.Background()).GetAll()
		if todoentryErr != nil {
			continue
		}

		for _, todoEntryItemElement := range todoEntryDocs {
			todoEntryData := todoEntryItemElement.Data()
			isDone := todoEntryData["isDone"].(bool)
			totalCount += 1
			if !isDone {
				inCompleteCount += 1
			}
		}
	}
	fmt.Println("未完成 ：", inCompleteCount)
	fmt.Println("總數", totalCount)
	return TaskInfo{
		inCompleteCount: inCompleteCount,
		totalCount:      totalCount,
	}
}

func sendNotficationToUser(message []byte, FCMToken string) {
	log.Printf("執行sendNotifcation方法 ", FCMToken)

	url := "https://fcm.googleapis.com/fcm/send"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(message))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", firebasefcmkey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("FCM推播錯誤", err)
		return
	}
	defer resp.Body.Close()
	log.Printf("FCM推播結果", resp.Status)
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

type FcmInfo struct {
	FcmValue bool   `firestore:"FCM"`
	FCMToken string `firestore:"FCMToken"` // in millions
}
type TaskInfo struct {
	inCompleteCount int
	totalCount      int
}
