package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHello(t *testing.T) {
	t.Log("Hello World")
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	w := httptest.NewRecorder()
	helloWorld(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Response code is %v", w.Code)
	}

	if w.Body.String() != "hello world\n" {
		t.Errorf("Response body is %v", w.Body.String())
	}
}

// サーバを立ててテストする
func TestHelloServer(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(helloWorld))
	t.Cleanup(func() {
		testServer.Close()
	})

	t.Log(testServer.URL)
	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/hello", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
}

type S3er interface {
	Upload() error
}

type S3 struct {
}

func (s *S3) Upload() error {
	fmt.Println("upload")
	return nil
}

type S3Mock struct {
}

func (s *S3Mock) Upload() error {
	fmt.Println("upload mock")
	return nil
}

func upload(s3 S3er) {
	s3.Upload()
}

func TestUpload(t *testing.T) {
	s3 := &S3Mock{}
	upload(s3)
}
