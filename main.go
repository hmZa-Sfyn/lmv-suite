package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"lanmanvan/cli"
)

func main() {
	var modulesDir string
	var version bool

	flag.StringVar(&modulesDir, "modules", "./modules", "Path to modules directory")
	flag.BoolVar(&version, "version", false, "Show version")
	flag.Parse()

	if version {
		fmt.Println("LanManVan %s - Advanced Modular Framework in Go", version)
		os.Exit(0)
	}

	// Expand home directory if needed
	if modulesDir == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: could not determine home directory: %v\n", err)
			os.Exit(1)
		}
		modulesDir = home
	}

	// Make absolute path
	absPath, err := filepath.Abs(modulesDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: invalid modules path: %v\n", err)
		os.Exit(1)
	}

	// Create and start CLI
	cliInstance := cli.NewCLI(absPath)
	if err := cliInstance.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
