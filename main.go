package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Semi-koron/c3-tyoukankaku-2025-backend/mult"
	"github.com/gorilla/mux"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // ローカル環境用デフォルトポート
	}
	
	r := mux.NewRouter()

	// ルーム作成用の API
	r.HandleFunc("/", mult.HandleConnection)

	serverAddr := ":8080"
	fmt.Println("WebSocket server started on", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, r))
}