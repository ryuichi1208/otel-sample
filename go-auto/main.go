package main

import (
	"fmt"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello")
}

func main() {
	http.HandleFunc("/", helloHandler)

	// ポート8080でHTTPサーバを起動
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
