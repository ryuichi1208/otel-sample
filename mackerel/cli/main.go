package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	// GoogleのトップページのURL
	url := "http://www.google.com"

	// HTTPクライアントを作成
	client := &http.Client{}

	// HTTP GETリクエストを作成
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// リクエストを送信し、レスポンスを受け取る
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// レスポンスのボディを読み取り
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	// レスポンスのボディを表示
	fmt.Println("Response body:")
	fmt.Println(string(body))
}
