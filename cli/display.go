package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"lanmanvan/core"

	"github.com/fatih/color"
)

// PrintHelp prints available commands
func (cli *CLI) PrintHelp() {
	fmt.Println()
	color.New(color.FgCyan, color.Bold).Println("Available Commands:")
	fmt.Println()

	commands := []struct {
		name string
		desc string
	}{
		// ──────────────────────────────
		// Core Commands
		// ──────────────────────────────
		{"help, h, ?", "Show this help message (aliases: h, ?)"},
		{"list, ls", "List all available modules (alias: ls)"},
		{"search <keyword>", "Search modules by name/tag/description (ex: search network)"},
		{"info <module>", "Show detailed info about a module (ex: info network)"},
		{"<module>!", "Quick view module options & usage (ex: network!)"},
		{"run <module> [args...]", "Execute module with arguments (ex: run network ip=192.168.1.1)"},
		{"<module> [args...]", "Shorthand run: module arg=value (ex: network ip=192.168.1.1)"},
		{"env, envs", "Display all global environment variables (alias: envs)"},
		{"key=value", "Set persistent global environment variable (ex: timeout=30)"},
		{"key=?", "View value of a global variable (ex: timeout=?)"},
		{"create <name> [python|bash]", "Create new module (ex: create exploit python)"},
		{"edit <module>", "Edit module source code (ex: edit myexploit)"},
		{"delete, rm <module>", "Delete a module (ex: delete myexploit)"},
		{"history", "Show command history"},
		{"clear, cls", "Clear the terminal screen (alias: cls)"},
		{"refresh, reload", "Reload/refresh all modules from disk"},
		{"exit, quit, q", "Exit the framework (aliases: quit, q)"},
	}

	for _, cmd := range commands {
		fmt.Printf("  %s%-34s %s\n",
			color.GreenString(""),
			color.CyanString(cmd.name),
			cmd.desc,
		)
	}

	fmt.Println()
	color.New(color.FgWhite, color.Bold).Println("Advanced Argument Features:")
	fmt.Println()

	advancedFeatures := []struct {
		name string
		desc string
	}{
		{"Quoted Strings", "Pass multi-word arguments: arg=\"from here to there\" time=\"15:45 4/6/2025\"."},
		{"Single Quotes", "Alternative quote style: arg='value with spaces'."},
		{"Variable Expansion", "Use $var_name: run module target=$mytarget (from global env or system env)."},
		{"Builtin Functions", "Execute builtins in args: run module pwd=$(pwd) hash=$(sha256 password)."},
		{"Combined Usage", "Mix variables and builtins: run module path=$workdir sig=$(sha256 $password)."},
		{"Save Output", "Save module execution to log file: module_name arg=value save=1 ."},
		{"Threaded Execution", "Run module with multiple threads: module_name arg=value threads=5 ."},
		{"Log Location", "Output files saved to ./logs/ with timestamp: module_2006-01-02_15-04-05.log ."},
	}

	for _, feat := range advancedFeatures {
		fmt.Printf("  %s%-20s %s\n",
			color.BlueString("- "),
			color.CyanString(feat.name),
			feat.desc,
		)
	}

	fmt.Println()
	color.New(color.FgWhite, color.Bold).Println("Variable & Function Examples:")
	fmt.Println()

	varExamples := []string{
		"Set global var:      myhost=192.168.1.1",
		"View global var:     myhost=?",
		"Expand in module:    run scanner target=$myhost",
		"Builtin function:    run hasher data=$(echo \"hello world\")",
		"Combine both:        run crypto key=$(sha256 $password) iv=$(pwd)",
		"Network info:        run netmod local=$(ipaddr) host=$(hostname)",
		"Timestamp:           run logger timestamp=$(timestamp unix) save=1",
		"File content:        run reader content=$(cat /tmp/file.txt)",
	}

	for _, example := range varExamples {
		fmt.Printf("  %s\n", color.GreenString(example))
	}

	fmt.Println()
	color.New(color.FgWhite, color.Bold).Println("Shell Commands (prefix with $):")
	fmt.Println()

	shellExamples := []string{
		"$ ls -la                      (Execute shell command)",
		"$ cd /tmp && pwd              (Change directory and execute)",
		"$ ifconfig eth0               (Get interface info)",
		"$ whoami                      (Current user)",
		"$ date                        (System date)",
	}

	for _, example := range shellExamples {
		fmt.Printf("  %s\n", color.YellowString(example))
	}

	fmt.Println()
	color.New(color.FgWhite, color.Bold).Println("Quick Examples:")
	fmt.Println()

	examples := []string{
		"arp-spoofer interface=eth0 target=\"192.168.1.100\" save=1",
		"dns-resolver domain=\"google.com\" save=1",
		"port-scanner target=\"192.168.1.0/24\" threads=10 save=1",
		"html-scraper url=\"https://example.com\" depth=\"2\" save=1",
		"network-info local=$(ipaddr) hostname=$(hostname)",
	}

	for _, example := range examples {
		fmt.Printf("  - %s\n", color.CyanString(example))
	}

	fmt.Println()
}

// formatModuleLine formats a single module line for display
func (cli *CLI) formatModuleLine(module *core.ModuleConfig, index int, total int) string {
	typeBadge := cli.getTypeBadge(module.Type)
	desc := ""
	tags := ""

	if module.Metadata != nil {
		if module.Metadata.Description != "" {
			desc = module.Metadata.Description
			if len(desc) > 50 {
				desc = desc[:47] + "..."
			}
		}
		if len(module.Metadata.Tags) > 0 {
			tags = strings.Join(module.Metadata.Tags[:1], "")
		}
	}

	prefix := "   ├─ "
	if index == total-1 {
		prefix = "   └─ "
	}

	return fmt.Sprintf("%s%s %s  %s %s",
		prefix,
		color.CyanString(module.Name),
		typeBadge,
		color.WhiteString(desc),
		color.MagentaString(tags),
	)
}

// ListModules displays all available modules
func (cli *CLI) ListModules() {
	modules := cli.manager.ListModules()
	if len(modules) == 0 {
		core.PrintWarning("No modules loaded. Check the modules directory or specify it with: lanmanvan -modules <path>")
		fmt.Println()
		return
	}

	fmt.Println()
	fmt.Println(core.NmapBox(fmt.Sprintf("AVAILABLE MODULES (%d)", len(modules))))

	// Sort modules by name
	sort.Slice(modules, func(i, j int) bool {
		return modules[i].Name < modules[j].Name
	})

	for i, module := range modules {
		fmt.Println(cli.formatModuleLine(module, i, len(modules)))
	}

	fmt.Println()
	core.PrintSuccess(fmt.Sprintf("Total: %d modules loaded", len(modules)))
	fmt.Println()
}

// SearchModules searches modules by keyword with highlighting
func (cli *CLI) SearchModules(keyword string) {
	modules := cli.manager.ListModules()
	keyword = strings.ToLower(keyword)

	var results []*core.ModuleConfig

	for _, module := range modules {
		name := strings.ToLower(module.Name)
		if strings.Contains(name, keyword) {
			results = append(results, module)
			continue
		}

		if module.Metadata != nil {
			desc := strings.ToLower(module.Metadata.Description)
			if strings.Contains(desc, keyword) {
				results = append(results, module)
				continue
			}

			for _, tag := range module.Metadata.Tags {
				if strings.Contains(strings.ToLower(tag), keyword) {
					results = append(results, module)
					break
				}
			}
		}
	}

	if len(results) == 0 {
		core.PrintWarning(fmt.Sprintf("No modules found for '%s', skipping...", keyword))
		return
	}

	// Sort results alphabetically by module name
	sort.Slice(results, func(i, j int) bool {
		return strings.ToLower(results[i].Name) < strings.ToLower(results[j].Name)
	})

	fmt.Println()
	fmt.Println(core.NmapBox(fmt.Sprintf("SEARCH: %s (%d results)", keyword, len(results))))

	for i, module := range results {
		fmt.Println(cli.formatModuleLineWithHighlight(module, i, len(results), keyword))
	}

	fmt.Println()
	core.PrintSuccess(fmt.Sprintf("Found %d module(s)", len(results)))
	fmt.Println()
}

// formatModuleLineWithHighlight formats a module line with keyword highlighting
func (cli *CLI) formatModuleLineWithHighlight(module *core.ModuleConfig, index int, total int, keyword string) string {
	typeBadge := cli.getTypeBadge(module.Type)
	desc := ""

	if module.Metadata != nil {
		desc = module.Metadata.Description
	}

	prefix := "   ├─ "
	if index == total-1 {
		prefix = "   └─ "
	}

	// Highlight keyword in module name and description
	highlightedName := cli.highlightKeyword(module.Name, keyword)
	highlightedDesc := cli.highlightKeyword(desc, keyword)

	return fmt.Sprintf("%s[%s] %s - %s",
		prefix,
		typeBadge,
		highlightedName,
		highlightedDesc,
	)
}

// highlightKeyword highlights a keyword in text with purple background
func (cli *CLI) highlightKeyword(text, keyword string) string {
	if keyword == "" {
		return text
	}

	// Find all occurrences (case-insensitive)
	keywordLower := strings.ToLower(keyword)
	textLower := strings.ToLower(text)

	// Build the highlighted string
	var result strings.Builder
	lastIdx := 0

	for {
		idx := strings.Index(textLower[lastIdx:], keywordLower)
		if idx == -1 {
			result.WriteString(text[lastIdx:])
			break
		}

		idx += lastIdx
		// Add text before match
		result.WriteString(text[lastIdx:idx])

		// Add highlighted match (purple background)
		match := text[idx : idx+len(keyword)]
		result.WriteString(color.New(color.BgMagenta, color.FgWhite, color.Bold).Sprint(match))

		lastIdx = idx + len(keyword)
	}

	return result.String()
}

// ShowModuleInfo displays detailed module information
func (cli *CLI) ShowModuleInfo(moduleName string, showREADME int) {
	module, err := cli.manager.GetModule(moduleName)
	if err != nil {
		core.PrintError(fmt.Sprintf("Error: %v, skipping...", err))
		return
	}

	fmt.Println()
	fmt.Println(core.NmapBox(fmt.Sprintf("MODULE: %s", moduleName)))

	if module.Metadata != nil {
		meta := module.Metadata
		fmt.Printf("   ├─ %s %s\n", color.WhiteString("Description:"), color.WhiteString(meta.Description))
		fmt.Printf("   ├─ %s %s\n", color.WhiteString("Type:"), cli.getTypeBadge(meta.Type))
		fmt.Printf("   ├─ %s %s\n", color.WhiteString("Author:"), color.RedString(meta.Author))
		fmt.Printf("   ├─ %s %s\n", color.WhiteString("Version:"), color.MagentaString(meta.Version))

		if len(meta.Tags) > 0 {
			fmt.Printf("   ├─ %s %s\n", color.WhiteString("Tags:"), color.CyanString(strings.Join(meta.Tags, ", ")))
		}

		// Display GitHub and X URLs
		if meta.GitHubURL != "" || meta.XUrl != "" {
			if meta.GitHubURL != "" {
				fmt.Printf("   ├─ %s %s\n", color.WhiteString("GitHub:"), color.BlueString(meta.GitHubURL))
			}
			if meta.XUrl != "" {
				fmt.Printf("   ├─ %s %s\n", color.WhiteString("X/Twitter:"), color.BlueString(meta.XUrl))
			}
		}

		if len(meta.Options) > 0 {
			fmt.Printf("   └─ %s\n", color.WhiteString("Options:"))

			// Sort option names for consistent output
			optNames := make([]string, 0, len(meta.Options))
			for optName := range meta.Options {
				optNames = append(optNames, optName)
			}
			sort.Strings(optNames)

			for i, optName := range optNames {
				opt := meta.Options[optName]
				required := ""
				if opt.Required {
					required = color.RedString(" [REQUIRED]")
				}

				prefix := "       ├─ "
				childPrefix := "       │  └─ "
				if i == len(optNames)-1 {
					prefix = "       └─ "
					childPrefix = "          └─ "
				}

				fmt.Printf("%s%s %s%s\n",
					prefix,
					color.GreenString(optName),
					color.WhiteString(fmt.Sprintf("(%s)", opt.Type)),
					required,
				)
				fmt.Printf("%s%s\n", childPrefix, color.WhiteString(opt.Description))
			}
		}
	} else {
		fmt.Printf("   └─ %s %s\n", color.WhiteString("Type:"), cli.getTypeBadge(module.Type))
		fmt.Println("   (No metadata available) ")
	}

	// Display README as "About This Module"
	if showREADME == 1 {
		cli.displayReadme(moduleName, module)
	}
	fmt.Println()
}

// displayReadme reads and displays the README.md as "About This Module"
func (cli *CLI) displayReadme(moduleName string, module *core.ModuleConfig) {
	readmePath := filepath.Join(cli.manager.ModulesDir, moduleName, "README.md")

	// Check if README exists
	if _, err := os.Stat(readmePath); err != nil {
		return
	}

	// Read README content
	content, err := ioutil.ReadFile(readmePath)
	if err != nil {
		return
	}

	readmeText := string(content)
	if readmeText == "" {
		return
	}

	fmt.Println()
	fmt.Println(core.NmapBox("ABOUT THIS MODULE"))

	// Create markdown renderer
	renderer := NewMarkdownRenderer()

	// Display README content with markdown rendering
	lines := strings.Split(strings.TrimSpace(readmeText), "\n")
	inCodeBlock := false
	var codeBlockContent []string
	var codeBlockLang string

	for _, line := range lines {
		// Handle code blocks
		if strings.HasPrefix(line, "```") {
			if !inCodeBlock {
				inCodeBlock = true
				codeBlockLang = strings.TrimPrefix(line, "```")
				if codeBlockLang == "" {
					codeBlockLang = "code"
				}
				codeBlockContent = []string{}
			} else {
				// End code block
				inCodeBlock = false
				if len(codeBlockContent) > 0 {
					fmt.Println()
					fmt.Println(renderer.RenderCodeBlock(strings.Join(codeBlockContent, "\n"), codeBlockLang))
					fmt.Println()
				}
				codeBlockContent = []string{}
			}
			continue
		}

		if inCodeBlock {
			codeBlockContent = append(codeBlockContent, line)
		} else {
			if line == "" {
				fmt.Println()
			} else {
				renderedLine := renderer.renderLine(line)
				fmt.Printf("   %s\n", renderedLine)
			}
		}
	}

	fmt.Println()
}

// PrintHistory shows command history
func (cli *CLI) PrintHistory() {
	if len(cli.history) == 0 {
		core.PrintWarning("No command history, skipping...")
		return
	}

	fmt.Println()
	fmt.Println(core.NmapBox(fmt.Sprintf("COMMAND HISTORY (%d)", len(cli.history))))

	for i, cmd := range cli.history {
		prefix := "   ├─ "
		if i == len(cli.history)-1 {
			prefix = "   └─ "
		}

		fmt.Printf("%s%s %s\n",
			prefix,
			color.GreenString(fmt.Sprintf("[%d]", i+1)),
			color.WhiteString(cmd),
		)
	}

	fmt.Println()
}

func HighlightPurple(text string, keyword string) string {

	if keyword == "" {
		return text
	}

	// Find all occurrences (case-insensitive)
	keywordLower := strings.ToLower(keyword)
	textLower := strings.ToLower(text)

	// Build the highlighted string
	var result strings.Builder
	lastIdx := 0

	for {
		idx := strings.Index(textLower[lastIdx:], keywordLower)
		if idx == -1 {
			result.WriteString(text[lastIdx:])
			break
		}

		idx += lastIdx
		// Add text before match
		result.WriteString(text[lastIdx:idx])

		// Add highlighted match (purple background)
		match := text[idx : idx+len(keyword)]
		result.WriteString(color.New(color.BgMagenta, color.FgWhite, color.Bold).Sprint(match))

		lastIdx = idx + len(keyword)
	}

	return result.String()
}

// HighlightPurple highlights all occurrences of keyword in text with purple background.
// If you don't have this function, use this simple version:

// getTypeBadge returns a colored badge for module type
func (cli *CLI) getTypeBadge(moduleType string) string {
	switch moduleType {
	case "python":
		return color.BlueString("[PY]")
	case "bash":
		return color.CyanString("[SH]")
	case "go":
		return color.MagentaString("[GO]")
	default:
		return color.WhiteString("[??]")
	}
}
