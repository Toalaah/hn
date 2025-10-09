package main

import (
	"fmt"
	"io"
	"os"

	"github.com/toalaah/hn/internal/app"
	"github.com/toalaah/hn/pkg/hn"
)

func main() {
	f, err := os.Open("./comments.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	body, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	t, err := hn.NewThreadFromData(body)
	if err != nil {
		panic(err)
	}

	if err := app.Run(t); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
