package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	_ "embed"
	"flag"
	"github.com/fatih/color"
	"github.com/toalaah/hn/internal/app"
	"github.com/toalaah/hn/pkg/hn"
)

var (
	//go:embed LICENSE
	license string
	prog    = filepath.Base(os.Args[0])
	version = "0.0.1"
)

func main() {
	var (
		id          int
		err         error
		t           *hn.Story
		showVersion *bool
		// Certainly, there is a better way to do this. But it works...
		showVersionShort *bool
	)

	color.Unset()
	flag.Usage = func() {
		fmt.Printf("Usage: %s [-v][-h] id\n", prog)
		fmt.Printf("\nFlags:\n")
		fmt.Printf("  -version, -v  print version and exit\n")
		fmt.Printf("  -h            print this usage and exit\n")
		fmt.Printf("\n%s", license)
	}

	showVersion = flag.Bool("version", false, "Show version and exit")
	showVersionShort = flag.Bool("v", false, "Show version and exit")
	flag.Parse()
	if *showVersion || *showVersionShort {
		fmt.Printf("%s\n", version)
		os.Exit(0)
	}

	if arg := flag.Arg(0); arg == "" {
		flag.Usage()
		os.Exit(1)
	} else if id, err = strconv.Atoi(os.Args[1]); err != nil {
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
