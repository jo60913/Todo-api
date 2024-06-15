# Todo-api
給ToDo app 所使用的RESTful api

## Getting Started
Go版本為1.22.2

RESTful api 設置在api/index.go Handler方法
執行api前需要設置以下環境變量

1. FIREBASE_ADMIN_SDK => 需要到firebase/專案設定/服務帳戶/Firebase admin sdk 產生新的私密金鑰
2. TODO_API_FIREBASE_FCM_KEY => 需要到firebase/專案設定/雲端通信/Cloud Messaging API (舊版) 伺服器金鑰匙 並在開頭加上key=。如 key=伺服器金鑰匙
3. FCM_HEADER => 觸發fcm api時的header

### 如何部署

1. 安裝vercel
```
npm install -g vercel
```

2. 登入vercel 
```
vercel login
```
3. 使用指令將當前程式碼部署到vercel 生產環境
```
vercel . --prod
```


### 結構
```
.
├── api
|   └── index.go            存放api的地方
├── model                   存放api request專換的JSON
|   ├── FcmInfo.go          使用者是否開啟FCM，與FCMToken
|   ├── FirstLogin.go       /update/firstlogin時，接收JSON
|   ├── NotificationGet.go  /get/notification時接收JSON
|   ├── NotificationUpdate  /update/notification時接收JSON
|   └── TaskInfo.go         紀錄使用者未完成任務數量與任務數量總數
├── go.mod
├── go.sum
├── README.md
└── vercel.json
```
### API

/update/notification
POST方法
更新使用者是否開啟推播功能

/get/notification
POST
取得使用推播功能狀態

/update/firstlogin
POST
第一次登入時推播通能預設開啟

/notification/fcm
POST
定時任務 每天8點推播提醒使用者未完成任務或請使用者新增任務。
header需要對應環境變數的FCM_HEADER才可以發送

