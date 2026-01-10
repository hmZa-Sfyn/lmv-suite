package cli

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"

	"lanmanvan/core"
)

// CLI manages the interactive command-line interface
type CLI struct {
	manager *core.ModuleManager
	running bool
	history []string
	envMgr  *EnvironmentManager
	logger  *Logger

	//v1.5 #macros
	macros        map[string]string
	macroParams   map[string][]string
	macroRequired map[string]map[string]bool
	builtinMacros map[string]bool
}

// NewCLI creates a new CLI instance
func NewCLI(modulesDir string) *CLI {
	return &CLI{
		manager: core.NewModuleManager(modulesDir),
		running: true,
		history: make([]string, 0),
		envMgr:  NewEnvironmentManager(),
		logger:  NewLogger(),

		//v1.5
		macros:        make(map[string]string),
		macroParams:   make(map[string][]string),
		macroRequired: make(map[string]map[string]bool),
		builtinMacros: make(map[string]bool),
	}
}

// Start begins the CLI loop
func (cli *CLI) Start(banner__ bool) error {
	if err := cli.manager.DiscoverModules(); err != nil {
		return err
	}

	if banner__ { // why is this showing when i told it not to? this one too!
		cli.PrintBanner()
	}
	cli.setupSignalHandler()

	///////////////////////////////////
	// v1.5

	// In your CLI initialization (NewCLI or similar):
	cli.macros = make(map[string]string)
	cli.macroParams = make(map[string][]string)
	cli.macroRequired = make(map[string]map[string]bool)
	cli.builtinMacros = map[string]bool{
		"echo":   true,
		"if":     true,
		"else":   true,
		"define": true,
		"def":    true,
	}

	// END v1.5

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

// Idle start
func (cli *CLI) IdleStart(banner__ bool, command__ string) error {
	if err := cli.manager.DiscoverModules(); err != nil {
		return err
	}

	if banner__ { // why is this showing when i told it not to?
		cli.PrintBanner()
	}
	cli.setupSignalHandler()

	// Create readline instance with history support
	rl, err := cli.getReadlineInstance()
	if err != nil {
		return err
	}
	defer rl.Close()

	for cli.running {
		//rl.SetPrompt(cli.GetPrompt())

		input := command__

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		cli.history = append(cli.history, input)
		cli.ExecuteCommand(input)

		break
	}

	return nil
}

// ExecuteCommand processes user commands

// handleBuiltinMacro returns true if the macro was handled (built-in), false otherwise
func (cli *CLI) ExecuteCommand(input string) {
	input = strings.TrimSpace(input)
	if input == "" {
		return
	}

	// 1. First: structured syntax commands that contain -> (for loops, etc.)
	if strings.HasPrefix(input, "for ") &&
		strings.Contains(input, " in ") &&
		strings.Contains(input, " -> ") {

		cli.executeForLoop(input)
		return
	}

	// 2. Special prefixes: #proxychains and #sudo → run original command via idle executor with prefix
	if strings.HasPrefix(input, "#proxychains ") || strings.HasPrefix(input, "#sudo ") {
		prefix := ""
		cmdPart := ""

		if strings.HasPrefix(input, "#proxychains ") {
			prefix = "proxychains "
			cmdPart = strings.TrimSpace(input[len("#proxychains "):])
		} else if strings.HasPrefix(input, "#sudo ") {
			prefix = "sudo "
			cmdPart = strings.TrimSpace(input[len("#sudo "):])
		}

		if cmdPart == "" {
			core.PrintError("Missing command after #" + prefix[:len(prefix)-1])
			return
		}

		// The REAL command we want to run in background/idle mode
		innerCommand := cmdPart

		// Build the wrapper that uses idle-exec
		wrapper := fmt.Sprintf(
			`%s~/bin/lanmanvan -modules ~/lanmanvan/modules -idle-exec -idle-cmd %q`,
			prefix,
			innerCommand,
		)

		core.PrintInfo("Executing in background/idle mode with prefix:")
		fmt.Printf("  → %s\n\n", wrapper)

		// Run the wrapper via shell (this should now work)
		cli.ExecuteShellCommand(wrapper)
		return
	}

	// 3. Output redirection > and >>  (only after special syntaxes!)
	// We check for space before > or >> to reduce false positives
	if strings.Contains(input, " > ") || strings.Contains(input, " >> ") ||
		(strings.HasSuffix(input, ">") && !strings.HasSuffix(input, "->")) ||
		(strings.HasSuffix(input, ">>")) {

		// Find the LAST occurrence of > or >>
		greaterPos := strings.LastIndex(input, ">>")
		if greaterPos == -1 {
			greaterPos = strings.LastIndex(input, ">")
		}

		if greaterPos > 0 {
			cmd := strings.TrimSpace(input[:greaterPos])
			redirectPart := strings.TrimSpace(input[greaterPos:])

			fields := strings.Fields(redirectPart)
			if len(fields) < 2 {
				core.PrintError("Redirection syntax: command > file  or  command >> file")
				return
			}

			op := fields[0] // > or >>
			filename := strings.Join(fields[1:], " ")
			filename = strings.Trim(filename, "\"'")

			wrapper := fmt.Sprintf(
				`lmv -idle-exec -idle-cmd %q %s %q`,
				cmd,
				op,
				filename,
			)

			core.PrintInfo("Redirecting output via idle/background executor...")
			fmt.Printf("  → %s\n\n", wrapper)

			cli.ExecuteShellCommand(wrapper)
			return
		}
	}

	// env var set / view
	if strings.Contains(input, "=") && !strings.Contains(input, " ") {
		parts := strings.SplitN(input, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			if value == "?" {
				if val, exists := cli.envMgr.Get(key); exists {
					fmt.Println()
					fmt.Printf("   %s = %s\n", core.Color("cyan", key), core.Color("green", val))
					fmt.Println()
				} else {
					core.PrintWarning(fmt.Sprintf("Variable '%s' not set", key))
					fmt.Println()
				}
				return
			}

			if err := cli.envMgr.Set(key, value); err != nil {
				core.PrintError(fmt.Sprintf("Failed to set variable: %v", err))
				return
			}

			fmt.Println()
			core.PrintSuccess(fmt.Sprintf("Set %s = %s", key, value))
			fmt.Println()
			return
		}
	}

	// 6. Shell command ($ prefix)
	if strings.HasPrefix(input, "$") {
		cli.ExecuteShellCommand(input)
		return
	}

	// 7. Regular commands / modules
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
	case "search":
		if len(args) > 0 {
			cli.SearchModules(strings.Join(args, " "))
		} else {
			core.PrintError("Usage: search <keyword>")
		}
	case "info":
		if len(args) > 0 {
			cli.ShowModuleInfo(args[0], 1)
		} else {
			core.PrintError("Usage: info <module>")
		}
	case "run":
		if len(args) > 0 {
			cli.RunModule(args[0], args[1:])
		} else {
			core.PrintError("Usage: run <module> [args...]")
		}
	case "create", "new":
		if len(args) > 0 {
			cli.CreateModule(args[0], args[1:])
		} else {
			core.PrintError("Usage: create <name> [python|bash]")
		}
	case "edit":
		if len(args) > 0 {
			cli.EditModule(args[0])
		} else {
			core.PrintError("Usage: edit <module>")
		}
	case "delete", "remove", "rm":
		if len(args) > 0 {
			cli.DeleteModule(args[0])
		} else {
			core.PrintError("Usage: delete <module>")
		}
	case "history":
		cli.PrintHistory()
	case "clear", "cls":
		cli.ClearScreen()
	case "refresh", "reload":
		cli.RefreshModules()
	case "exit", "quit", "q":
		cli.running = false
		fmt.Println()
		core.PrintSuccess("Goodbye! See you next time.")
		fmt.Println()
		return

	default:
		// Quick module info: module!
		if strings.HasSuffix(cmd, "!") {
			moduleName := strings.TrimSuffix(cmd, "!")
			cli.ShowModuleInfo(moduleName, 0)
		} else {
			// Try to run as module
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

// RefreshModules refreshes and reloads all modules from the modules directory
func (cli *CLI) RefreshModules() {
	fmt.Println()
	core.PrintInfo("Refreshing modules...")
	fmt.Println()

	// Clear and reinitialize the module manager with the same directory
	modulesDirPath := cli.manager.ModulesDir
	cli.manager = core.NewModuleManager(modulesDirPath)

	// Discover modules again
	if err := cli.manager.DiscoverModules(); err != nil {
		core.PrintError(fmt.Sprintf("Failed to refresh modules: %v", err))
		fmt.Println()
		return
	}

	// Count loaded modules
	modules := cli.manager.ListModules()
	moduleCount := len(modules)

	fmt.Println()
	core.PrintSuccess(fmt.Sprintf("Modules refreshed successfully! Loaded %d module(s)", moduleCount))
	fmt.Println()

	// Display summary of loaded modules
	if moduleCount > 0 {
		fmt.Println(core.NmapBox("Loaded Modules"))
		for i, module := range modules {
			status := ""
			fmt.Printf("   [%d] %s %s\n", i+1, status, core.Color("cyan", module.Name))
		}
		fmt.Println()
	}
}

// Iterator represents something that can produce values one by one
type Iterator interface {
	Next() (string, bool) // value, ok
	Len() int             // total expected items (for progress)
	Close() error         // optional cleanup
}

func (cli *CLI) executeForLoop(input string) {
	// Supported syntaxes:
	// for $x in 1..100 -> command
	// for x in a..z -> command
	// for ip in 192.168.1.1..192.168.1.50 -> ping $ip
	// for c in a..z+A..Z+0..9 -> echo $c
	// for user in admin|root|guest -> hydra -l $user ...

	input = strings.TrimSpace(input)

	// Flexible regex - supports both $var and var
	re := regexp.MustCompile(`(?i)^for\s+(?:\$?(\w+))\s+(?:in\s+)?(.+?)\s*[-=]{1,2}>\s*(.+)$`)
	matches := re.FindStringSubmatch(input)
	if len(matches) != 4 {
		core.PrintError("Invalid for-loop syntax.\nExamples:\n  for $x in 1..100 -> echo $x\n  for ip in 192.168.1.1..50 -> ping $ip\n  for c in a..z+A..Z -> echo $c")
		return
	}

	varName := matches[1]
	source := strings.TrimSpace(matches[2])
	command := strings.TrimSpace(matches[3])

	iter, err := parseRangeSource(source)
	if err != nil {
		E_msg := "Cannot parse range: " + err.Error() + "\nSource was: " + source + ""
		core.PrintError(E_msg)
		return
	}
	defer iter.Close()

	total := iter.Len()
	if total == 0 {
		core.PrintWarning("Empty range - nothing to do")
		return
	}

	fmt.Println()
	core.PrintInfo(fmt.Sprintf("Loop: %s ∈ %s  (%d items)", varName, source, total))
	fmt.Println()

	results := []string{}
	count := 0

	for {
		value, ok := iter.Next()
		if !ok {
			break
		}
		count++

		// Support both $var and ${var}
		expanded := regexp.MustCompile(`\$\{`+regexp.QuoteMeta(varName)+`\}|\$`+regexp.QuoteMeta(varName)).
			ReplaceAllString(command, value)

		fmt.Printf("  [%3d/%3d] → %s\n", count, total, expanded)

		var result string
		if strings.Contains(expanded, "|>") {
			result = cli.executePipedCommandsForLoop(expanded)
			if result != "" {
				results = append(results, result)
			}
		} else {
			cli.ExecuteCommand(expanded)
		}
	}

	if len(results) > 0 {
		fmt.Println()
		core.PrintSuccess("Collected results (" + string(len(results)) + "):")
		for i, res := range results {
			fmt.Printf("  [%2d] %s\n", i+1, strings.TrimSpace(res))
		}
		fmt.Println()
	}
}

// parseRangeSource returns an iterator for different kinds of ranges
func parseRangeSource(s string) (Iterator, error) {
	s = strings.TrimSpace(s)

	// 1. List style: item1|item2|item3
	if strings.Contains(s, "|") {
		items := strings.Split(s, "|")
		cleanItems := make([]string, 0, len(items))
		for _, item := range items {
			trimmed := strings.TrimSpace(item)
			if trimmed != "" {
				cleanItems = append(cleanItems, trimmed)
			}
		}
		return &listIterator{items: cleanItems}, nil
	}

	// 2. Multiple ranges with + : a..z+A..Z+0..9
	if strings.Contains(s, "+") {
		parts := strings.Split(s, "+")
		iterators := make([]Iterator, 0, len(parts))
		for _, part := range parts {
			it, err := parseSingleRange(strings.TrimSpace(part))
			if err != nil {
				return nil, fmt.Errorf("invalid part %q: %v", part, err)
			}
			iterators = append(iterators, it)
		}
		return newChainIterator(iterators...), nil
	}

	// 3. Single range
	return parseSingleRange(s)
}

// ────────────────────────────────────────────────────────────────────────────────
// Chain Iterator (for a..z + 0..9 + !@# style)
// ────────────────────────────────────────────────────────────────────────────────

type chainIterator struct {
	iterators []Iterator
	current   int
}

func newChainIterator(iters ...Iterator) Iterator {
	return &chainIterator{
		iterators: iters,
		current:   0,
	}
}

func (it *chainIterator) Next() (string, bool) {
	for it.current < len(it.iterators) {
		val, ok := it.iterators[it.current].Next()
		if ok {
			return val, true
		}
		it.current++
	}
	return "", false
}

func (it *chainIterator) Len() int {
	total := 0
	for _, i := range it.iterators {
		total += i.Len()
	}
	return total
}

func (it *chainIterator) Close() error {
	for _, i := range it.iterators {
		_ = i.Close() // best effort
	}
	return nil
}

// ────────────────────────────────────────────────────────────────────────────────
// IP Range Iterator (full IPs: 192.168.1.1 .. 192.168.1.50)
// ────────────────────────────────────────────────────────────────────────────────

type ipRangeIterator struct {
	start net.IP
	end   net.IP
	curr  net.IP
}

func newIPRangeIterator(start, end net.IP) Iterator {
	// Make copies because net.IP is slice
	curr := make(net.IP, len(start))
	copy(curr, start)

	return &ipRangeIterator{
		start: start,
		end:   end,
		curr:  curr,
	}
}

func (it *ipRangeIterator) Next() (string, bool) {
	if bytes.Compare(it.curr, it.end) > 0 {
		return "", false
	}

	result := it.curr.String()

	// Increment IP
	for i := len(it.curr) - 1; i >= 0; i-- {
		it.curr[i]++
		if it.curr[i] > 0 {
			break
		}
		// carry over
		it.curr[i] = 0
	}

	return result, true
}

func (it *ipRangeIterator) Len() int {
	// Very rough estimate - good enough for progress bar
	diff := ipToInt(it.end) - ipToInt(it.start)
	if diff < 0 {
		return 0
	}
	return int(diff) + 1
}

func (it *ipRangeIterator) Close() error { return nil }

// Helper: IPv4 only!
func ipToInt(ip net.IP) int64 {
	ip = ip.To4()
	if ip == nil {
		return 0
	}
	return int64(ip[0])<<24 | int64(ip[1])<<16 | int64(ip[2])<<8 | int64(ip[3])
}

// ────────────────────────────────────────────────────────────────────────────────
// Partial IP Range (last octet only)  e.g. 192.168.1.10..192.168.1.50
// ────────────────────────────────────────────────────────────────────────────────

type partialIPRangeIterator struct {
	prefix string
	start  int
	end    int
	curr   int
}

func newPartialIPRangeIterator(startStr, endStr string) (Iterator, error) {
	// Assume format like: 192.168.1.   or  10.0.0.
	// We expect startStr and endStr to be last octet numbers

	// Very naive version - you can make it more robust later
	start, err1 := strconv.Atoi(startStr)
	if err1 != nil {
		return nil, fmt.Errorf("invalid start octet: %v", err1)
	}
	end, err2 := strconv.Atoi(endStr)
	if err2 != nil {
		return nil, fmt.Errorf("invalid end octet: %v", err2)
	}

	// Guess prefix from startStr if possible
	prefix := ""
	if dot := strings.LastIndex(startStr, "."); dot != -1 {
		prefix = startStr[:dot+1]
	}

	return &partialIPRangeIterator{
		prefix: prefix,
		start:  start,
		end:    end,
		curr:   start,
	}, nil
}

func (it *partialIPRangeIterator) Next() (string, bool) {
	if it.curr > it.end {
		return "", false
	}
	ip := fmt.Sprintf("%s%d", it.prefix, it.curr)
	it.curr++
	return ip, true
}

func (it *partialIPRangeIterator) Len() int {
	return it.end - it.curr + 1
}

func (it *partialIPRangeIterator) Close() error { return nil }

// ────────────────────────────────────────────────────────────────────────────────
// Character Range Iterator   a..z    or   0..9
// ────────────────────────────────────────────────────────────────────────────────

type charRangeIterator struct {
	current byte
	end     byte
}

func newCharRangeIterator(start, end byte) Iterator {
	return &charRangeIterator{
		current: start,
		end:     end,
	}
}

func (it *charRangeIterator) Next() (string, bool) {
	if it.current > it.end {
		return "", false
	}
	val := string(it.current)
	it.current++
	return val, true
}

func (it *charRangeIterator) Len() int {
	return int(it.end - it.current + 1)
}

func (it *charRangeIterator) Close() error { return nil }

func parseSingleRange(s string) (Iterator, error) {
	if strings.Contains(s, "..") {
		parts := strings.SplitN(s, "..", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid .. range format")
		}
		startStr, endStr := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])

		// Try IP range first
		startIP := net.ParseIP(startStr)
		if startIP != nil {
			endIP := net.ParseIP(endStr)
			if endIP != nil {
				return newIPRangeIterator(startIP, endIP), nil
			}
			// Maybe last octet only?
			if strings.HasPrefix(startStr, "192.168.") || strings.Contains(startStr, ".") {
				return newPartialIPRangeIterator(startStr, endStr)
			}
		}

		// Numeric range
		start, err1 := strconv.Atoi(startStr)
		end, err2 := strconv.Atoi(endStr)
		if err1 == nil && err2 == nil {
			return newNumericRangeIterator(start, end), nil
		}

		// Character range
		if len(startStr) == 1 && len(endStr) == 1 {
			return newCharRangeIterator(startStr[0], endStr[0]), nil
		}
	}

	return nil, fmt.Errorf("unsupported range format: %s", s)
}

// ────────────────────────────────────────────────────────────────────────────────
// Simple iterators implementations (you can put them in separate file)
// ────────────────────────────────────────────────────────────────────────────────

type listIterator struct {
	items []string
	idx   int
}

func (it *listIterator) Next() (string, bool) {
	if it.idx >= len(it.items) {
		return "", false
	}
	v := it.items[it.idx]
	it.idx++
	return v, true
}
func (it *listIterator) Len() int     { return len(it.items) }
func (it *listIterator) Close() error { return nil }

type numericRangeIterator struct {
	current, end int
}

func newNumericRangeIterator(start, end int) *numericRangeIterator {
	return &numericRangeIterator{current: start, end: end}
}
func (it *numericRangeIterator) Next() (string, bool) {
	if it.current > it.end {
		return "", false
	}
	v := fmt.Sprintf("%d", it.current)
	it.current++
	return v, true
}
func (it *numericRangeIterator) Len() int     { return it.end - it.current + 1 }
func (it *numericRangeIterator) Close() error { return nil }

// Add these yourself (similar style):
// - newIPRangeIterator
// - newPartialIPRangeIterator
// - newCharRangeIterator
// - newChainIterator (for + separated ranges)

// executePipedCommandsForLoop handles pipes and returns output instead of printing
func (cli *CLI) executePipedCommandsForLoop(input string) string {
	parts := strings.Split(input, "|>")
	if len(parts) < 2 {
		return ""
	}

	var result string
	var err error

	// Execute first command
	firstCmd := strings.TrimSpace(parts[0])
	result, err = cli.executePipedCommand(firstCmd, "")
	if err != nil {
		return ""
	}

	// Execute remaining commands, passing output as input
	for i := 1; i < len(parts); i++ {
		nextCmd := strings.TrimSpace(parts[i])
		result, err = cli.executePipedCommand(nextCmd, result)
		if err != nil {
			return ""
		}
	}

	return result
}

// executePipedCommands handles piped commands with |> syntax
// Example: whoami() |> sha256() or cat(file.txt) |> base64()
func (cli *CLI) executePipedCommands(input string) {
	parts := strings.Split(input, "|>")
	if len(parts) < 2 {
		return
	}

	var result string
	var err error

	// Execute first command
	firstCmd := strings.TrimSpace(parts[0])
	result, err = cli.executePipedCommand(firstCmd, "")
	if err != nil {
		core.PrintError(fmt.Sprintf("Pipe error in first command: %v", err))
		return
	}

	// Execute remaining commands, passing output as input
	for i := 1; i < len(parts); i++ {
		nextCmd := strings.TrimSpace(parts[i])
		result, err = cli.executePipedCommand(nextCmd, result)
		if err != nil {
			core.PrintError(fmt.Sprintf("Pipe error at step %d: %v", i+1, err))
			return
		}
	}

	// Only print result if the last command is not file() - file() handles its own output
	lastCmd := strings.TrimSpace(parts[len(parts)-1])
	if !strings.HasPrefix(lastCmd, "file(") {
		fmt.Println()
		fmt.Println(result)
		fmt.Println()
	}
}

// executePipedCommand executes a single command in a pipe chain
// Supports: builtin(args), module, module arg=value
func (cli *CLI) executePipedCommand(cmd string, input string) (string, error) {
	cmd = strings.TrimSpace(cmd)

	// Handle string literals in pipes: "\n", "\t", "text", etc.
	if (strings.HasPrefix(cmd, "\"") && strings.HasSuffix(cmd, "\"")) ||
		(strings.HasPrefix(cmd, "'") && strings.HasSuffix(cmd, "'")) {
		// Remove quotes and process escape sequences
		literal := cmd[1 : len(cmd)-1]

		// Process escape sequences
		literal = strings.ReplaceAll(literal, "\\n", "\n")
		literal = strings.ReplaceAll(literal, "\\t", "\t")
		literal = strings.ReplaceAll(literal, "\\r", "\r")
		literal = strings.ReplaceAll(literal, "\\\\", "\\")

		// String literals just pass through, potentially appending to input
		return input + literal, nil
	}

	// If input from previous command, inject it appropriately
	if input != "" {
		// If command is a builtin function call
		if strings.Contains(cmd, "(") && strings.Contains(cmd, ")") {
			openParen := strings.Index(cmd, "(")
			closeParen := strings.LastIndex(cmd, ")")
			if openParen > 0 && closeParen > openParen {
				funcName := cmd[:openParen]
				args := cmd[openParen+1 : closeParen]
				if args != "" {
					args += ", \"" + input + "\""
				} else {
					args = "\"" + input + "\""
				}
				cmd = funcName + "(" + args + ")"
			}
		} else {
			// It's a module call with potential arguments
			// Check if there's an argument pattern like: modulename ip=$somevar
			if strings.Contains(cmd, "=") {
				// Module with specific arguments - find what argument to inject into
				// If pattern is "module arg=$var", replace $var with input
				if strings.Contains(cmd, "$") {
					// Find the variable and replace it
					parts := strings.Split(cmd, "=")
					if len(parts) >= 2 {
						// Replace the variable value with piped input
						lastPart := parts[len(parts)-1]
						if strings.HasPrefix(lastPart, "$") {
							// Replace the variable
							varName := strings.TrimSpace(lastPart)
							cmd = strings.Replace(cmd, varName, "\""+input+"\"", 1)
						} else {
							// Append input as new argument
							cmd = cmd + " input=\"" + input + "\""
						}
					}
				} else {
					// Append input as new argument
					cmd = cmd + " input=\"" + input + "\""
				}
			} else {
				// No arguments, just module name
				cmd = cmd + " input=\"" + input + "\""
			}
		}
	}

	// Try to execute as module
	parts := strings.Fields(cmd)
	if len(parts) > 0 {
		moduleName := parts[0]
		args := parts[1:]

		// Check if module exists
		if _, err := cli.manager.GetModule(moduleName); err == nil {
			// Execute module and capture output
			return cli.executeModuleForPipe(moduleName, args)
		}
	}

	return "", fmt.Errorf("invalid pipe command: %s", cmd)
}

// executeModuleForPipe executes a module and returns its output
func (cli *CLI) executeModuleForPipe(moduleName string, args []string) (string, error) {
	_, err := cli.manager.GetModule(moduleName)
	if err != nil {
		return "", err
	}

	// Parse arguments with support for variable expansion
	moduleArgs := make(map[string]string)
	parsedArgs := cli.parseArguments(args)

	for key, value := range parsedArgs {
		switch key {
		case "threads", "save":
			// Skip these
		default:
			moduleArgs[key] = value
		}
	}

	// Merge global environment variables
	for key, value := range cli.envMgr.GetAll() {
		if _, exists := moduleArgs[key]; !exists {
			moduleArgs[key] = value
		}
	}

	// Save original stdout to restore later
	saveOut := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		// Fallback: execute without capturing stdout
		result, execErr := cli.manager.ExecuteModule(moduleName, moduleArgs)
		if execErr != nil {
			return "", execErr
		}
		return strings.TrimSpace(result.Output), nil
	}

	// Redirect stdout to our pipe
	os.Stdout = writer

	// Execute module
	result, err := cli.manager.ExecuteModule(moduleName, moduleArgs)

	// Restore stdout
	writer.Close()
	os.Stdout = saveOut

	// Read captured stdout (not used but needed to drain the pipe)
	_ = reader
	reader.Close()

	if err != nil {
		return "", err
	}

	// Return the module output directly, which contains the structured output
	return strings.TrimSpace(result.Output), nil
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
