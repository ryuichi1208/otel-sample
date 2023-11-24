package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Printf("%d\n", os.Getpid())
	fmt.Printf("%d\n", os.Getppid())
}
