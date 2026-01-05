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

	var exec bool
	var exec_cmd string

	var show_banner bool

	flag.StringVar(&modulesDir, "modules", "./modules", "Path to modules directory (string)")
	flag.BoolVar(&version, "version", false, "Show version (bool)")

	flag.BoolVar(&exec, "idle-exec", false, "Execute command and exit? (bool)")
	flag.StringVar(&exec_cmd, "idle-cmd", "help", "Execute command and exit (string)")

	flag.BoolVar(&show_banner, "banner", false, "Want to show the *lanmanvan* official banner? (bool)")

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
	if !exec {
		if err := cliInstance.Start(show_banner); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}

	if exec {
		if exec_cmd != "" {
			//execute command and exit!

			//init shell
			cliInstance := cli.NewCLI(absPath)

			//show banner & execute command
			if err := cliInstance.IdleStart(show_banner, exec_cmd); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			//optional: save it to file?
			//exit
			os.Exit(0)
		}
		os.Exit(0)
	}
}
