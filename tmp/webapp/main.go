package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"
)

var httpClient *http.Client

// httpクライアントを作って返す
func newClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy:               http.ProxyFromEnvironment,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			TLSHandshakeTimeout: 5 * time.Second,
			Dial: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).Dial,
			ResponseHeaderTimeout: 5 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   100,
			IdleConnTimeout:       90,
			DisableKeepAlives:     false,
		},
	}
}

func returnHelloWorld(w http.ResponseWriter, r *http.Request) {
	res, err := httpClient.Get("https://www.google.com")
	if err != nil {
		fmt.Fprintf(w, "error: %v\n", err)
		return
	}
	defer res.Body.Close()
	// レスポンスを表示
	fmt.Fprintf(w, "status: %v\n", res.Status)
	fmt.Fprintf(w, "status: %v\n", res.Body)

	fmt.Fprintf(w, "hello world\n")
}

func middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Server-Name", "Go-Server")
		next(w, r)
	}
}

func f(ctx context.Context) {
	time.Sleep(1 * time.Second)
	fmt.Println("hello world8")
	for {
		select {
		case <-ctx.Done():
			fmt.Println("child: error: ", ctx.Err(), ctx.Value("key"))
			return
		default:
			fmt.Println("hello world9")
		}
	}
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	childCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	childCtx = context.WithValue(childCtx, "key", "value")
	for {
		select {
		case <-ctx.Done():
			cancel()
			fmt.Println("error: ", ctx.Err())
			return
		default:
			f(childCtx)
			fmt.Fprintf(w, "hello world3\n")
		}
	}
}

func main() {
	httpClient = newClient()
	// httpリクエストが来たらhello worldを返す
	http.HandleFunc("/", returnHelloWorld)

	// middlewareでhttpヘッダーを追加する
	http.HandleFunc("/middleware", middleware(returnHelloWorld))

	// httpリクエストが来たらhello worldを返す
	http.HandleFunc("/hello", helloWorld)

	// 8080ポートでサーバーを立てる
	http.ListenAndServe(":8080", nil)
}
