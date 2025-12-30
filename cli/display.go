package cli

import (
	"fmt"
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

// SearchModules searches modules by keyword
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

	fmt.Println()
	fmt.Println(core.NmapBox(fmt.Sprintf("SEARCH: %s (%d results)", keyword, len(results))))

	for i, module := range results {
		fmt.Println(cli.formatModuleLine(module, i, len(results)))
	}

	fmt.Println()
	core.PrintSuccess(fmt.Sprintf("Found %d module(s)", len(results)))
	fmt.Println()
}

// ShowModuleInfo displays detailed module information
func (cli *CLI) ShowModuleInfo(moduleName string) {
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
