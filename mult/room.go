package mult

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type PlayerData struct {
	ID       string  `json:"id"`
	Color    string  `json:"color"`
	Position struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
		Z float64 `json:"z"`
	} `json:"position"`
	Rotation struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
		Z float64 `json:"z"`
		W float64 `json:"w"`
	} `json:"rotation"`
	Chat string `json:"chat"`
}

var writeMutex sync.Mutex

// WebSocket のアップグレーダー
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// allowUrl := os.Getenv("ALLOW_URL")
		// if allowUrl == "" {
		// 	allowUrl = "http://localhost:3000"
		// }
		// origin := r.Header.Get("Origin")
		// return origin == allowUrl
		return true
	},
}

// ルーム管理
var roomData = make(map[*websocket.Conn]bool);
var playerDatas = make(map[*websocket.Conn]PlayerData);

// ルーム内の全員にメッセージを送信
func sendMessageAll(msg []byte) {
    writeMutex.Lock()
    defer writeMutex.Unlock()
    
    for client := range roomData {
        err := client.WriteMessage(websocket.TextMessage, msg)
        if err != nil {
            log.Println("Write Error:", err)
            client.Close()
            delete(roomData, client)
        }
    }
}


// WebSocket 接続を処理
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket Upgrade Error:", err)
		return
	}
	defer conn.Close()

	// クライアントを登録
	writeMutex.Lock()
	roomData[conn] = true
	writeMutex.Unlock()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read Error:", err)

			// 切断時にプレイヤーデータを削除
			writeMutex.Lock()
			delete(roomData, conn)
			delete(playerDatas, conn)
			writeMutex.Unlock()

			break
		}

		// JSON デコード
		var player PlayerData
		if err := json.Unmarshal(msg, &player); err != nil {
			log.Println("JSON Unmarshal Error:", err)
			continue
		}

		// `playerDatas` に格納
		writeMutex.Lock()
		playerDatas[conn] = player
		writeMutex.Unlock()
	}
}

func sendPlayerDataLoop() {
	ticker := time.NewTicker(time.Second / 60) // 60FPS で送信
	defer ticker.Stop()

	for range ticker.C {
		writeMutex.Lock()

		// `playerDatas` を `[]PlayerData` に変換
		playerDataList := make([]PlayerData, 0, len(playerDatas))
		for _, data := range playerDatas {
			playerDataList = append(playerDataList, data)
		}

		writeMutex.Unlock()

		// JSON に変換
		data, err := json.Marshal(playerDataList)
		if err != nil {
			log.Println("JSON Marshal Error:", err)
			continue
		}

		// すべてのクライアントに送信
		sendMessageAll(data)
	}
}


func init() {
	go sendPlayerDataLoop() // データ送信ループを開始
}