package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	// リクエストヘッダーを取得する
	fmt.Printf("Request Header: %v\n", r.Header)
	// リクエストボディを取得する
	fmt.Printf("Request Body: %v\n", r.Body)
	// リクエストパラメータを取得する
	fmt.Printf("Request Form: %v\n", r.Form)

	// httpヘッダーを書き込む
	w.Header().Set("Content-Type", "text/plain")
	// ステータスコードを書き込む
	w.WriteHeader(http.StatusOK)
	// レスポンスの内容を書き込む
	fmt.Fprintf(w, "Hello, World!")
}

func main() {
	// ハンドラを登録
	http.HandleFunc("/", handler)

	// サーバを起動
	port := 8080
	fmt.Printf("Starting server on port %d...\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
