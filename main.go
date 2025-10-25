package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/fatih/color"
	"github.com/toalaah/hn/internal/app"
	"github.com/toalaah/hn/pkg/hn"
)

func main() {
	var (
		id  int
		err error
		t   *hn.Story
	)

	color.Unset()
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s [id]\n\nUnexpected number of args.\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	id, err = strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("Could not parse id: %s\n", err)
		os.Exit(1)
	}

	t, err = hn.NewThread(id)
	if err != nil {
		fmt.Printf("Error fetching thread: %s\n", err)
		os.Exit(1)
	}

	if err := app.Run(t); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
