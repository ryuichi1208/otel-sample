package main

import (
	"fmt"
	"testing"

	gomock "go.uber.org/mock/gomock"
)

func logic(db DBer) {
	fmt.Println("logic")
	red, err, i := db.Close()
	fmt.Println(red, err, i)
}

func TestClose(t *testing.T) {
	fmt.Println("TestClose")
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mock := NewMockDBer(mockCtrl)
	fmt.Println(mock)
	mock.EXPECT().Close().Return(&Redis{
		endpoint: "localhost",
		port:     6379,
	},
		nil,
		1).Times(1)
	logic(mock)

	t.Log("TestClose")
}
