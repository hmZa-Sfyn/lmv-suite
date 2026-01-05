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
		{"help, h, ?", "Show this help message, aliases: h, ?"},
		{"list, ls", "List all modules, aliases: ls"},
		{"search <keyword>", "Search modules by name/tag, example: search network"},
		{"info <module>", "Show detailed module information, example: info network"},
		{"<module>!", "Quick show module options and usage, example: network!"},
		{"run <module> [args]", "Execute a module with arguments, example: run network ip="},
		{"<module> [args]", "Shorthand: <module> arg_key=value, example: network ip=192.168.1.1"},
		{"<module> arg_key = value", "Format with spaces (alternative), example: network ip = 192.168.1.1"},
		{"env, envs", "Show all global environment variables, aliases: envs"},
		{"builtins", "Show all 30+ builtin functions with examples"},
		{"builtins <keyword>", "Show all builtin functions matching to keyword with examples"},
		{"key=value", "Set global environment variable (persistent), example: timeout=10"},
		{"key=?", "View global environment variable value, example: timeout=?"},
		{"create <name> [type]", "Create a new module (python/bash), example: create mymodule python"},
		{"edit <module>", "Edit module files, example: edit mymodule"},
		{"delete <module>", "Delete a module, example: delete mymodule"},
		{"history", "Show command history"},
		{"clear", "Clear screen, aliases: cls"},
		{"exit, quit, q", "Exit framework, aliases: quit, q"},
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

func (cli *CLI) PrintBuiltins(keyword string) {
	fmt.Println()
	fmt.Println(core.NmapBox("BUILTIN FUNCTIONS (150+) - DETAILED REFERENCE"))
	fmt.Println()

	builtins := cli.builtins.GetAll()

	// Categories without spaces in names
	categories := map[string][]*BuiltinFunction{
		"FileSystem":        {},
		"FileOperations":    {},
		"System":            {},
		"Hashing":           {},
		"AdvancedHashing":   {},
		"Encoding":          {},
		"AdvancedEncoding":  {},
		"Strings":           {},
		"StringProcessing":  {},
		"Converters":        {},
		"NetworkValidation": {},
		"Network":           {},
		"NetworkExtended":   {},
		"Validation":        {},
		"Math":              {},
		"MathExtended":      {},
		"Cryptography":      {},
		"DateTime":          {},
		"JSONData":          {},
		"Utilities":         {},
	}

	// Categorize all builtins
	for _, fn := range builtins {
		switch fn.Name {
		// FileSystem core + extended
		case "pwd", "cd", "ls", "mkdir", "rm", "cp", "mv", "cat", "exists", "filesize",
			"find", "tail", "head", "touch", "chmod", "stat", "isdir", "isfile", "file",
			"readfile", "writefile":
			categories["FileSystem"] = append(categories["FileSystem"], fn)

		// System
		case "whoami", "hostname", "date", "uname", "arch", "ostype", "uptime", "ps",
			"getenv", "which":
			categories["System"] = append(categories["System"], fn)

		// Hashing basic
		case "md5", "sha1", "sha256", "hash", "checksum", "crc32":
			categories["Hashing"] = append(categories["Hashing"], fn)

		// AdvancedHashing
		case "sha512", "blake2b", "blake2s", "hmac_sha256", "hmac_sha512",
			"murmur3", "xxhash", "fnv1", "fnv1a", "djb2":
			categories["AdvancedHashing"] = append(categories["AdvancedHashing"], fn)

		// Encoding basic
		case "base64", "hex", "url", "json", "csv", "xml", "ascii", "unicode":
			categories["Encoding"] = append(categories["Encoding"], fn)

		// AdvancedEncoding
		case "base32", "base58", "base85", "punycode", "morse", "binary", "octal",
			"quoted_printable", "percent_encode", "htmlescape":
			categories["AdvancedEncoding"] = append(categories["AdvancedEncoding"], fn)

		// Strings core
		case "strlen", "toupper", "tolower", "reverse", "trim", "substr", "replace",
			"split", "startswith", "endswith", "contains", "repeat":
			categories["Strings"] = append(categories["Strings"], fn)

		// StringProcessing
		case "camelcase", "snakecase", "kebabcase", "capitalize", "lowercase", "uppercase",
			"swapcase", "ltrim", "rtrim", "center", "ljust", "rjust", "indent",
			"dedent", "wordcount":
			categories["StringProcessing"] = append(categories["StringProcessing"], fn)

		// Converters
		case "btoa", "atob", "bin2hex", "hex2bin", "bin2dec", "dec2bin", "hex2dec",
			"dec2hex", "oct2dec", "dec2oct", "rot13", "rot47", "caesar", "reverse_bytes",
			"toascii":
			categories["Converters"] = append(categories["Converters"], fn)

		// NetworkValidation
		case "isipv4", "isipv6", "isemail", "isurl", "ismac", "isdomain", "ispath",
			"isport", "iscdr", "getcidr", "getiprange", "ip2int", "int2ip", "reverseip",
			"parseurl", "isvalidip":
			categories["NetworkValidation"] = append(categories["NetworkValidation"], fn)

		// Network core
		case "ping", "nslookup", "ipaddr", "gethostbyname", "getipversion", "iplookup",
			"getport", "getmac", "gateway", "getdns":
			categories["Network"] = append(categories["Network"], fn)

		// NetworkExtended
		case "cidrmatch", "hostmask", "broadcast", "ipversion", "isprivate",
			"isloopback", "ismulticast", "getmactable":
			categories["NetworkExtended"] = append(categories["NetworkExtended"], fn)

		// Validation (data types)
		case "isint", "isfloat", "isalpha", "isalnum", "isnumeric", "isspace",
			"isbinary", "ishex", "isuuid", "isbase64", "ismd5", "issha1", "issha256",
			"isjson":
			categories["Validation"] = append(categories["Validation"], fn)

		// Math core
		case "calc", "abs", "min", "max", "sum", "avg", "random", "randomstr",
			"randint", "rand":
			categories["Math"] = append(categories["Math"], fn)

		// MathExtended
		case "pow", "sqrt", "cbrt", "log", "exp", "factorial":
			categories["MathExtended"] = append(categories["MathExtended"], fn)

		// Cryptography
		case "uuid", "timestamp", "genpass", "seed":
			categories["Cryptography"] = append(categories["Cryptography"], fn)

		// DateTime
		case "now", "epoch", "iso8601", "strtotime", "timeformat", "strftime",
			"timezone", "dayofweek":
			categories["DateTime"] = append(categories["DateTime"], fn)

		// JSONData
		case "jsonpretty", "jsonminify", "jsonformat":
			categories["JSONData"] = append(categories["JSONData"], fn)

		// Utilities
		case "sleep", "echo", "list":
			categories["Utilities"] = append(categories["Utilities"], fn)

		default:
			categories["Utilities"] = append(categories["Utilities"], fn)
		}
	}

	normalizedKeyword := strings.ToLower(strings.TrimSpace(keyword))
	isCategorySearch := false
	targetCategory := ""

	if strings.HasPrefix(keyword, "#") {
		targetCategory = strings.TrimSpace(strings.TrimPrefix(keyword, "#"))
		if _, ok := categories[targetCategory]; ok {
			isCategorySearch = true
		} else {
			// --- Invalid category handling ---
			fmt.Println(color.RedString("Unknown category: %s", targetCategory))
			fmt.Println()

			// Suggest possible matches (case-insensitive substring)
			var suggestions []string
			lcTarget := strings.ToLower(targetCategory)
			for cat := range categories {
				lcCat := strings.ToLower(cat)
				if strings.Contains(lcCat, lcTarget) || strings.Contains(lcTarget, lcCat) {
					suggestions = append(suggestions, cat)
				}
			}

			if len(suggestions) > 0 {
				fmt.Println(color.YellowString("Did you mean one of these?"))
				sort.Strings(suggestions)
				for _, s := range suggestions {
					fmt.Printf("  #%-20s\n", s)
				}
				fmt.Println()
			}

			// Show all valid categories
			fmt.Println(color.CyanString("Valid categories:"))
			order := []string{
				"FileSystem", "FileOperations", "System",
				"Hashing", "AdvancedHashing",
				"Encoding", "AdvancedEncoding",
				"Strings", "StringProcessing",
				"Converters",
				"NetworkValidation", "Network", "NetworkExtended",
				"Validation",
				"Math", "MathExtended",
				"Cryptography",
				"DateTime",
				"JSONData",
				"Utilities",
			}
			for _, cat := range order {
				count := len(categories[cat])
				if count > 0 {
					fmt.Printf("  #%-20s (%d)\n", cat, count)
				}
			}
			fmt.Println()
			printQuickReference()
			return
		}
	}

	type FunctionWithCategory struct {
		Func     *BuiltinFunction
		Category string
	}

	var results []FunctionWithCategory

	if isCategorySearch {
		for _, fn := range categories[targetCategory] {
			results = append(results, FunctionWithCategory{Func: fn, Category: targetCategory})
		}
	} else if normalizedKeyword == "" {
		// Full categorized view
		order := []string{
			"FileSystem", "FileOperations", "System",
			"Hashing", "AdvancedHashing",
			"Encoding", "AdvancedEncoding",
			"Strings", "StringProcessing",
			"Converters",
			"NetworkValidation", "Network", "NetworkExtended",
			"Validation",
			"Math", "MathExtended",
			"Cryptography",
			"DateTime",
			"JSONData",
			"Utilities",
		}

		for _, cat := range order {
			fns := categories[cat]
			if len(fns) == 0 {
				continue
			}
			fmt.Printf(" %s\n", color.CyanString("═ %s (%d) ═", cat, len(fns)))
			sort.Slice(fns, func(i, j int) bool { return fns[i].Name < fns[j].Name })
			for i, fn := range fns {
				isLast := i == len(fns)-1
				prefix := " ├─ "
				if isLast {
					prefix = " └─ "
				}
				catTag := color.CyanString("[%s]", cat)
				name := color.WhiteString(fmt.Sprintf("%-18s", fn.Name))
				fmt.Printf("%s%s %s %s\n", prefix, catTag, name, fn.Description)

				if fn.DetailedDesc != "" {
					detailPrefix := " │ "
					if isLast {
						detailPrefix = "   "
					}
					fmt.Printf("%s%s\n", detailPrefix, color.WhiteString(fn.DetailedDesc))

					if len(fn.Examples) > 0 {
						fmt.Printf("%s%s\n", detailPrefix, color.YellowString("Examples:"))
						for j, ex := range fn.Examples {
							exPrefix := " │ ├─ "
							if j == len(fn.Examples)-1 {
								exPrefix = " │ └─ "
								if isLast {
									exPrefix = "   └─ "
								}
							} else if isLast {
								exPrefix = "   ├─ "
							}
							fmt.Printf("%s%s\n", exPrefix, color.MagentaString(ex))
						}
					}
				}
				if isLast {
					fmt.Println()
				}
			}
		}
		printQuickReference()
		return
	} else {
		// Keyword search across name, description, detailed desc, examples
		for cat, fns := range categories {
			for _, fn := range fns {
				lowerName := strings.ToLower(fn.Name)
				lowerDesc := strings.ToLower(fn.Description)
				lowerDetailed := strings.ToLower(fn.DetailedDesc)
				lowerExamples := ""
				for _, ex := range fn.Examples {
					lowerExamples += strings.ToLower(ex) + " "
				}

				if strings.Contains(lowerName, normalizedKeyword) ||
					strings.Contains(lowerDesc, normalizedKeyword) ||
					strings.Contains(lowerDetailed, normalizedKeyword) ||
					strings.Contains(lowerExamples, normalizedKeyword) {
					results = append(results, FunctionWithCategory{Func: fn, Category: cat})
				}
			}
		}
	}

	if len(results) == 0 {
		fmt.Println(color.RedString(" No matching builtins found for: %s", keyword))
		fmt.Println()
		return
	}

	// Sort results by name
	sort.Slice(results, func(i, j int) bool {
		return results[i].Func.Name < results[j].Func.Name
	})

	header := "MATCHING BUILTINS"
	if isCategorySearch {
		header = fmt.Sprintf("CATEGORY: %s", targetCategory)
	}
	fmt.Printf(" %s\n", color.CyanString("═ %s (%d) ═", header, len(results)))

	for i, item := range results {
		fn := item.Func
		cat := item.Category
		isLast := i == len(results)-1
		prefix := " ├─ "
		if isLast {
			prefix = " └─ "
		}

		hlName := HighlightPurple(fn.Name, normalizedKeyword)
		hlDesc := HighlightPurple(fn.Description, normalizedKeyword)

		catTag := color.New(color.FgCyan, color.Bold).Sprint(fmt.Sprintf("[%s]", cat))
		nameFormatted := fmt.Sprintf("%-18s", hlName)
		nameColored := color.WhiteString(nameFormatted)

		fmt.Printf("%s%s %s %s\n", prefix, catTag, nameColored, hlDesc)

		if fn.DetailedDesc != "" {
			detailPrefix := " │ "
			if isLast {
				detailPrefix = "   "
			}
			hlDetailed := HighlightPurple(fn.DetailedDesc, normalizedKeyword)
			fmt.Printf("%s%s\n", detailPrefix, color.WhiteString(hlDetailed))

			if len(fn.Examples) > 0 {
				fmt.Printf("%s%s\n", detailPrefix, color.YellowString("Examples:"))
				for j, ex := range fn.Examples {
					exPrefix := " │ ├─ "
					if j == len(fn.Examples)-1 {
						exPrefix = " │ └─ "
						if isLast {
							exPrefix = "   └─ "
						}
					} else if isLast {
						exPrefix = "   ├─ "
					}
					hlEx := HighlightPurple(ex, normalizedKeyword)
					fmt.Printf("%s%s\n", exPrefix, color.MagentaString(hlEx))
				}
			}
		}
		if isLast {
			fmt.Println()
		}
	}

	if !isCategorySearch && normalizedKeyword != "" {
		fmt.Println(color.CyanString("\n Tip: Use #CategoryName (e.g., #Strings or #NetworkValidation) to show full category"))
		fmt.Println()
	}

	printQuickReference()
}
func printQuickReference() {
	fmt.Println("\n   " + color.CyanString("═ Quick Reference ═"))
	fmt.Println("   ├─ List all:                " + color.YellowString("list") + "  or just press Tab twice")
	fmt.Println("   ├─ Show full category:     " + color.YellowString("#CategoryName") + "  (e.g., #FileSystem, #NetworkValidation)")
	fmt.Println("   ├─ Search functions:       " + color.YellowString("keyword") + "  (case-insensitive, matches name/description/examples)")
	fmt.Println("   ├─ Call syntax:            " + color.YellowString("funcname(arg1, arg2, ...)"))
	fmt.Println("   ├─ Nested calls:           " + color.YellowString("echo($(sha256($(randomstr 16))))"))
	fmt.Println("   ├─ In module args:         " + color.YellowString("run module target=$(hostname) port=$(random 1024-65535)"))
	fmt.Println("   ├─ Variables:              " + color.YellowString("set myip $(ipaddr) ; echo $myip"))
	fmt.Println("   ├─ Quoted arguments:       " + color.YellowString("echo(\"hello world with spaces\")"))
	fmt.Println("   ├─ Multi-line input:       Use \\ at end of line to continue")
	fmt.Println("   └─ Pipe output (future):   Coming soon")
	fmt.Println()
}

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
