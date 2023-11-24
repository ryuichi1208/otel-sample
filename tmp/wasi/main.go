package main

import (
	"fmt"
	"net"
)

func main() {
	_, err := net.Listen("tcp", "127.0.0.1:3000")
	if err != nil {
		fmt.Println(err)
	}

}
