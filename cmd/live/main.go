package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/toalaah/hn/internal/app"
	"github.com/toalaah/hn/pkg/hn"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No ID passed")
		os.Exit(1)
	}
	id, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("Could not parse argument to int: %s\n", os.Args[1])
		os.Exit(1)
	}
	t, err := hn.NewThread(id)
	if err != nil {
		fmt.Printf("Error fetching thread: %s\n", err)
		os.Exit(1)
	}

	if err := app.Run(t); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
