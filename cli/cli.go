package cli

import (
	"fmt"
	"os"
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
	// Handle builtin function calls: funcname(arg1,arg2,arg3) or func(arg arg2)
	// Check if it looks like a function call: starts with identifier and has matching parentheses
	if strings.Contains(input, "(") && strings.Contains(input, ")") {
		openParen := strings.Index(input, "(")
		if openParen > 0 {
			potentialFunc := input[:openParen]
			// Check if the part before parenthesis is a valid identifier (no spaces)
			if !strings.Contains(potentialFunc, " ") && potentialFunc != "" {
				if cli.tryExecuteBuiltin(input) {
					return
				}
			}
		}
	}

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

// tryExecuteBuiltin attempts to execute a builtin function call
// Syntax: funcname(arg1, arg2, arg3) with support for:
// - Quoted strings: "hello world", 'single quotes'
// - Nested builtins: $(builtin_name args) or builtin_name()
// - Variable expansion: $varname
// - Space-separated arguments
func (cli *CLI) tryExecuteBuiltin(input string) bool {
	openParen := strings.Index(input, "(")
	if openParen == -1 {
		return false
	}

	funcName := strings.TrimSpace(input[:openParen])

	// Verify funcName is a valid identifier
	if funcName == "" {
		return false
	}
	for _, ch := range funcName {
		if !isValidVarChar(ch) {
			return false
		}
	}

	// Check if function exists
	if _, exists := cli.builtins.functions[funcName]; !exists {
		return false
	}

	// Find the matching closing parenthesis
	closeParen := cli.findMatchingParen(input, openParen+1)
	if closeParen == -1 {
		return false
	}

	// Make sure there's nothing after the closing paren (except whitespace)
	afterParen := strings.TrimSpace(input[closeParen+1:])
	if afterParen != "" {
		return false
	}

	argsStr := input[openParen+1 : closeParen]

	// Parse arguments with proper handling of quotes, builtins, and variables
	args := cli.parseAdvancedArguments(argsStr)

	// Execute builtin
	result, err := cli.builtins.Execute(funcName, args...)

	fmt.Println()
	if err != nil {
		core.PrintError(fmt.Sprintf("Error: %v", err))
	} else {
		fmt.Println(result)
	}
	fmt.Println()

	return true
}

// parseAdvancedArguments parses function arguments with support for:
// - Quoted strings (both "..." and '...')
// - Nested builtins $(builtin args) and builtin() function call syntax
// - Variable expansion $var
// - Space-separated arguments
func (cli *CLI) parseAdvancedArguments(argsStr string) []string {
	var args []string
	var currentArg strings.Builder
	i := 0

	for i < len(argsStr) {
		ch := argsStr[i]

		// Handle quoted strings
		if ch == '"' || ch == '\'' {
			quote := ch
			i++ // skip opening quote
			for i < len(argsStr) && argsStr[i] != quote {
				if argsStr[i] == '\\' && i+1 < len(argsStr) {
					// Handle escape sequences
					i++
					currentArg.WriteByte(argsStr[i])
				} else {
					currentArg.WriteByte(argsStr[i])
				}
				i++
			}
			i++ // skip closing quote
			continue
		}

		// Handle nested builtins: $(builtin args)
		if ch == '$' && i+1 < len(argsStr) && argsStr[i+1] == '(' {
			closeParen := cli.findMatchingParen(argsStr, i+2)
			if closeParen != -1 {
				nestedCall := argsStr[i : closeParen+1]
				expanded := cli.expandBuiltinCall(nestedCall)
				currentArg.WriteString(expanded)
				i = closeParen + 1
				continue
			}
		}

		// Handle builtin function calls: funcname()
		if isValidVarChar(rune(ch)) {
			// Collect identifier
			ident := cli.collectIdentifier(argsStr, &i)
			if i < len(argsStr) && argsStr[i] == '(' {
				// This is a function call
				closeParen := cli.findMatchingParen(argsStr, i+1)
				if closeParen != -1 {
					// Get the arguments inside parentheses
					innerArgs := argsStr[i+1 : closeParen]
					// Check if this is a known builtin
					if _, exists := cli.builtins.functions[ident]; exists {
						// Recursively parse the inner arguments
						expanded := cli.executeBuiltinDirectly(ident, innerArgs)
						currentArg.WriteString(expanded)
						i = closeParen + 1
						continue
					} else {
						// Not a builtin, treat as part of argument
						currentArg.WriteString(ident)
						currentArg.WriteByte('(')
						continue
					}
				}
			}
			// Not a function call, just add the identifier
			currentArg.WriteString(ident)
			continue
		}

		// Handle variable expansion: $varname
		if ch == '$' && i+1 < len(argsStr) && isValidVarChar(rune(argsStr[i+1])) {
			i++ // skip $
			var varName strings.Builder
			for i < len(argsStr) && isValidVarChar(rune(argsStr[i])) {
				varName.WriteByte(argsStr[i])
				i++
			}
			varVal := cli.expandVariable(varName.String())
			currentArg.WriteString(varVal)
			continue
		}

		// Handle comma-separated arguments
		if ch == ',' {
			arg := strings.TrimSpace(currentArg.String())
			if arg != "" {
				args = append(args, arg)
			}
			currentArg.Reset()
			i++
			continue
		}

		// Handle spaces (space-separated arguments)
		if ch == ' ' {
			arg := strings.TrimSpace(currentArg.String())
			if arg != "" {
				args = append(args, arg)
			}
			currentArg.Reset()
			i++
			continue
		}

		currentArg.WriteByte(ch)
		i++
	}

	// Add final argument
	arg := strings.TrimSpace(currentArg.String())
	if arg != "" {
		args = append(args, arg)
	}

	return args
}

// isValidVarChar checks if a rune is valid in a variable name
func isValidVarChar(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_'
}

// collectIdentifier extracts an identifier starting at position i
func (cli *CLI) collectIdentifier(s string, i *int) string {
	var ident strings.Builder
	for *i < len(s) && isValidVarChar(rune(s[*i])) {
		ident.WriteByte(s[*i])
		*i++
	}
	return ident.String()
}

// findMatchingParen finds the index of the closing parenthesis that matches
// the opening parenthesis at startIdx
func (cli *CLI) findMatchingParen(s string, startIdx int) int {
	depth := 1
	i := startIdx
	inQuote := false
	quoteChar := byte(0)

	for i < len(s) && depth > 0 {
		ch := s[i]

		// Handle quotes
		if (ch == '"' || ch == '\'') && (i == 0 || s[i-1] != '\\') {
			if !inQuote {
				inQuote = true
				quoteChar = ch
			} else if ch == quoteChar {
				inQuote = false
			}
		}

		// Handle parentheses (only outside quotes)
		if !inQuote {
			if ch == '(' {
				depth++
			} else if ch == ')' {
				depth--
				if depth == 0 {
					return i
				}
			}
		}

		i++
	}

	return -1 // Not found
}

// executeBuiltinDirectly executes a builtin function with raw argument string
func (cli *CLI) executeBuiltinDirectly(funcName string, argsStr string) string {
	// Recursively parse the inner arguments to handle nested calls
	args := cli.parseAdvancedArguments(argsStr)
	result, err := cli.builtins.Execute(funcName, args...)
	if err != nil {
		return ""
	}
	return result
}

// expandBuiltinCall expands a nested builtin call like $(sha256 abc) or $(func(arg))
func (cli *CLI) expandBuiltinCall(call string) string {
	// Handle both $(func args) and $(func(args)) syntax
	if !strings.HasPrefix(call, "$(") || !strings.HasSuffix(call, ")") {
		return call
	}

	innerCall := call[2 : len(call)-1]

	// Check if inner call is a function call syntax: func(args)
	for i, ch := range innerCall {
		if ch == '(' {
			funcName := innerCall[:i]
			if _, exists := cli.builtins.functions[funcName]; exists {
				// This is a function call like sha256(abc)
				return cli.executeBuiltinDirectly(funcName, innerCall[i+1:len(innerCall)-1])
			}
			break
		}
		if !isValidVarChar(ch) && ch != '_' {
			break
		}
	}

	// Otherwise, treat as space-separated: func arg1 arg2...
	parts := strings.Fields(innerCall)
	if len(parts) == 0 {
		return ""
	}

	funcName := parts[0]
	var funcArgs []string
	if len(parts) > 1 {
		funcArgs = parts[1:]
	}

	result, err := cli.builtins.Execute(funcName, funcArgs...)
	if err != nil {
		return ""
	}
	return result
}

// expandVariable expands a variable reference
func (cli *CLI) expandVariable(varName string) string {
	if val, exists := cli.envMgr.Get(varName); exists {
		return val
	}
	if val, exists := os.LookupEnv(varName); exists {
		return val
	}
	return "$" + varName
}
