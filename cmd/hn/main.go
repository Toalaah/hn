package main

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/toalaah/hn/internal/app"
	"github.com/toalaah/hn/pkg/hn"
)

func main() {
	var (
		id  int
		err error
		t   *hn.Story
	)
	if len(os.Args) > 1 {
		id, err = strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Printf("Error opening file: %s\n", err)
			os.Exit(1)
		}
	}

	switch id {
	case 0:
		f, err := os.Open("./hn.json")
		if err != nil {
			fmt.Printf("Error opening file: %s\n", err)
			os.Exit(1)
		}
		defer f.Close()
		body, err := io.ReadAll(f)
		if err != nil {
			fmt.Printf("Error reading file: %s\n", err)
			os.Exit(1)
		}
		t, err = hn.NewThreadFromData(body)
		if err != nil {
			fmt.Printf("Error fetching thread: %s\n", err)
			os.Exit(1)
		}
	default:
		t, err = hn.NewThread(id)
		if err != nil {
			fmt.Printf("Error fetching thread: %s\n", err)
			os.Exit(1)
		}
	}

	if err := app.Run(t); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
