package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/google/go-github/v53/github"
	"github.com/k1LoW/ghfs"
)

func main() {
	ctx := context.Background()
	cli := github.NewTokenClient(ctx, "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")
	client := ghfs.Client(cli)

	fsys, err := ghfs.New("ryuichi1208", "dotfiles", client)
	if err != nil {
		log.Fatal(err)
	}
	f, err := fsys.Open("README.md")
	if err != nil {
		log.Fatal(err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", b)

	// ファイルのリストを取得
	list, err := fsys.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range list {
		fmt.Println(v.Name())
	}

}
