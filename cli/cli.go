package cli

import (
	"fmt"
	"strings"

	"lanmanvan/core"
)

// CLI manages the interactive command-line interface
type CLI struct {
	manager  *core.ModuleManager
	running  bool
	history  []string
	envMgr   *EnvironmentManager
	logger   *Logger
	builtins *BuiltinRegistry
}

// NewCLI creates a new CLI instance
func NewCLI(modulesDir string) *CLI {
	return &CLI{
		manager:  core.NewModuleManager(modulesDir),
		running:  true,
		history:  make([]string, 0),
		envMgr:   NewEnvironmentManager(),
		logger:   NewLogger(),
		builtins: NewBuiltinRegistry(),
	}
}

// Start begins the CLI loop
func (cli *CLI) Start() error {
	if err := cli.manager.DiscoverModules(); err != nil {
		return err
	}

	cli.PrintBanner()
	cli.setupSignalHandler()

	// Create readline instance with history support
	rl, err := cli.getReadlineInstance()
	if err != nil {
		return err
	}
	defer rl.Close()

	for cli.running {
		rl.SetPrompt(cli.GetPrompt())

		input, err := rl.Readline()
		if err != nil {
			if err.Error() == "Interrupt" {
				fmt.Println()
				continue
			}
			if err.Error() == "EOF" {
				break
			}
			continue
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		cli.history = append(cli.history, input)
		cli.ExecuteCommand(input)
	}

	return nil
}

// ExecuteCommand processes user commands
func (cli *CLI) ExecuteCommand(input string) {
	// Handle global environment variable syntax (key=value or key=?)
	if strings.Contains(input, "=") && !strings.Contains(input, " ") {
		parts := strings.SplitN(input, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			// Check if it's a view operation (key=?)
			if value == "?" {
				if val, exists := cli.envMgr.Get(key); exists {
					fmt.Println()
					fmt.Printf("   %s = %s\n", core.Color("cyan", key), core.Color("green", val))
					fmt.Println()
				} else {
					core.PrintWarning(fmt.Sprintf("Environment variable '%s' not set, skipping...", key))
					fmt.Println()
				}
				return
			}

			// Set environment variable
			if err := cli.envMgr.Set(key, value); err != nil {
				core.PrintError(fmt.Sprintf("Failed to set environment variable: %v, skipping...", err))
				return
			}
			fmt.Println()
			core.PrintSuccess(fmt.Sprintf("Set %s = %s", key, value))
			fmt.Println()
			return
		}
	}

	// Handle shell commands
	if strings.HasPrefix(input, "$") {
		cli.ExecuteShellCommand(input)
		return
	}

	parts := strings.Fields(input)
	if len(parts) == 0 {
		return
	}

	cmd := parts[0]
	args := parts[1:]

	switch cmd {
	case "help", "h", "?":
		cli.PrintHelp()
	case "list", "ls":
		cli.ListModules()
	case "env", "envs":
		cli.envMgr.Display()
	case "builtins":
		cli.PrintBuiltins()
	case "search":
		if len(args) > 0 {
			cli.SearchModules(strings.Join(args, " "))
		} else {
			core.PrintError("Usage: search <keyword> ... example: search network")
		}
	case "info":
		if len(args) > 0 {
			cli.ShowModuleInfo(args[0], 1)
		} else {
			core.PrintError("Usage: info <module_name> ... example: info network")
		}
	case "run":
		if len(args) > 0 {
			cli.RunModule(args[0], args[1:])
		} else {
			core.PrintError("Usage: run <module_name> [args...] ... example: run network target_network=$target_network_suffix port=80")
		}
	case "create", "new":
		if len(args) > 0 {
			cli.CreateModule(args[0], args[1:])
		} else {
			core.PrintError("Usage: create <module_name> [python|bash] ... example: create mymodule python")
		}
	case "edit":
		if len(args) > 0 {
			cli.EditModule(args[0])
		} else {
			core.PrintError("Usage: edit <module_name> ... example: edit mymodule")
		}
	case "delete", "remove", "rm":
		if len(args) > 0 {
			cli.DeleteModule(args[0])
		} else {
			core.PrintError("Usage: delete <module_name> ... example: delete mymodule")
		}
	case "history":
		cli.PrintHistory()
	case "clear", "cls":
		cli.ClearScreen()
	case "exit", "quit", "q":
		cli.running = false
		fmt.Println()
		core.PrintSuccess("Goodbye! See you next time.")
		fmt.Println()
	default:
		// Check if command ends with ! (show module info)
		if strings.HasSuffix(cmd, "!") {
			moduleName := strings.TrimSuffix(cmd, "!")
			cli.ShowModuleInfo(moduleName, 0)
		} else {
			// Try to run as a module if command is not recognized
			cli.RunModule(cmd, args)
		}
	}
}

// GetModuleManager returns the module manager instance
func (cli *CLI) GetModuleManager() *core.ModuleManager {
	return cli.manager
}

// IsRunning returns the running state
func (cli *CLI) IsRunning() bool {
	return cli.running
}

// AddHistory adds a command to history
func (cli *CLI) AddHistory(cmd string) {
	cli.history = append(cli.history, cmd)
}

// GetHistory returns the command history
func (cli *CLI) GetHistory() []string {
	return cli.history
}

// Stop stops the CLI loop
func (cli *CLI) Stop() {
	cli.running = false
}
