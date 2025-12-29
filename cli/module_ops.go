package cli

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"lanmanvan/core"
)

// RunModule executes a module with provided arguments
func (cli *CLI) RunModule(moduleName string, args []string) {
	// First, get the module to check for required args
	module, err := cli.manager.GetModule(moduleName)
	if err != nil {
		core.PrintError(fmt.Sprintf("%v", err))
		return
	}

	// Parse arguments into key=value pairs (supports both arg=value and arg = value)
	moduleArgs := make(map[string]string)
	threads := 1
	i := 0

	for i < len(args) {
		arg := args[i]

		// Check if it's a key=value pair
		if strings.Contains(arg, "=") {
			parts := strings.SplitN(arg, "=", 2)
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			if key == "threads" {
				fmt.Sscanf(value, "%d", &threads)
			} else {
				moduleArgs[key] = value
			}
		} else if i+2 < len(args) && args[i+1] == "=" {
			// Handle "key = value" format
			key := strings.TrimSpace(arg)
			value := strings.TrimSpace(args[i+2])

			if key == "threads" {
				fmt.Sscanf(value, "%d", &threads)
			} else {
				moduleArgs[key] = value
			}
			i += 2 // Skip the = and value
		} else {
			// Positional argument
			moduleArgs[fmt.Sprintf("arg%d", i)] = arg
		}

		i++
	}

	// Merge global environment variables (command-line args take precedence)
	for key, value := range cli.envMgr.GetAll() {
		if _, exists := moduleArgs[key]; !exists {
			moduleArgs[key] = value
		}
	}

	// Check for required arguments
	if module.Metadata != nil && len(module.Metadata.Required) > 0 {
		missingArgs := []string{}
		for _, required := range module.Metadata.Required {
			if _, exists := moduleArgs[required]; !exists {
				missingArgs = append(missingArgs, required)
			}
		}

		// If required arguments are missing, show usage info
		if len(missingArgs) > 0 {
			fmt.Println()
			core.PrintWarning(fmt.Sprintf("Module '%s' requires arguments, skipping...", moduleName))
			fmt.Println()
			fmt.Println(core.NmapBox(fmt.Sprintf("MODULE: %s - USAGE", moduleName)))
			fmt.Printf("   Description: %s\n\n", module.Metadata.Description)
			fmt.Println("   Required Arguments:")
			for _, opt := range missingArgs {
				if optMeta, exists := module.Metadata.Options[opt]; exists {
					fmt.Printf("      * %s (%s) - %s\n", opt, optMeta.Type, optMeta.Description)
				}
			}

			if len(module.Metadata.Options) > 0 {
				fmt.Println("\n   Optional Arguments:")
				for optName, optMeta := range module.Metadata.Options {
					isRequired := false
					for _, req := range module.Metadata.Required {
						if req == optName {
							isRequired = true
							break
						}
					}
					if !isRequired {
						fmt.Printf("      • %s (%s) - %s\n", optName, optMeta.Type, optMeta.Description)
					}
				}
			}

			fmt.Printf("\n   Example Usage:\n")
			fmt.Printf("      %s %s=%s\n\n", moduleName, missingArgs[0], "value")
			return
		}
	}

	startTime := time.Now()

	fmt.Println()
	if threads > 1 {
		core.PrintInfo(fmt.Sprintf("Executing module '%s' with %d threads...", core.Color("cyan", moduleName), threads))
	} else {
		core.PrintInfo(fmt.Sprintf("Executing module '%s'...", core.Color("cyan", moduleName)))
	}
	fmt.Println()

	// Mark module as running
	cli.startModuleExecution()
	defer cli.stopModuleExecution()

	var result *core.ExecutionResult

	if threads > 1 {
		result, err = cli.runModuleThreaded(moduleName, moduleArgs, threads)
	} else {
		result, err = cli.manager.ExecuteModule(moduleName, moduleArgs)
	}

	if err != nil {
		core.PrintError(fmt.Sprintf("%v", err))
		return
	}

	duration := time.Since(startTime)

	// Print output section
	if result.Output != "" {
		fmt.Println(core.NmapBox("Output"))
		for _, line := range strings.Split(strings.TrimSpace(result.Output), "\n") {
			if line != "" {
				fmt.Println(core.NmapSubBox(line))
			}
		}
		fmt.Println()
	}

	// Print error section if exists
	if result.Error != "" {
		core.PrintError("Error Output:")
		for _, line := range strings.Split(result.Error, "\n") {
			if line != "" {
				fmt.Printf("  %s\n", core.Color("red", line))
			}
		}
		fmt.Println()
	}

	// Print simple result log
	fmt.Println()
	if result.Success {
		core.PrintSuccess(fmt.Sprintf("Completed in %s [exit: %d], success!", duration.String(), result.ExitCode))
	} else {
		core.PrintError(fmt.Sprintf("Failed in %s [exit: %d], skipping...", duration.String(), result.ExitCode))
	}
	fmt.Println()
}

// runModuleThreaded executes a module with multiple threads
func (cli *CLI) runModuleThreaded(moduleName string, args map[string]string, threads int) (*core.ExecutionResult, error) {
	_, err := cli.manager.GetModule(moduleName)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	results := make(chan *core.ExecutionResult, threads)
	var outputs []string
	var mu sync.Mutex

	wg.Add(threads)

	for i := 0; i < threads; i++ {
		go func(threadID int) {
			defer wg.Done()
			result, _ := cli.manager.ExecuteModule(moduleName, args)
			if result != nil {
				mu.Lock()
				outputs = append(outputs, fmt.Sprintf("[Thread : %d] %s", threadID, result.Output))
				mu.Unlock()
			}
			results <- result
		}(i + 1)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect all results
	finalResult := &core.ExecutionResult{
		Success:   true,
		Timestamp: time.Now(),
		ExitCode:  0,
	}

	for result := range results {
		if result != nil && !result.Success {
			finalResult.Success = false
		}
	}

	mu.Lock()
	finalResult.Output = strings.Join(outputs, "\n")
	mu.Unlock()

	return finalResult, nil
}

// CreateModule creates a new module
func (cli *CLI) CreateModule(moduleName string, args []string) {
	moduleType := "python"
	if len(args) > 0 {
		moduleType = strings.ToLower(args[0])
	}

	if moduleType != "python" && moduleType != "bash" {
		core.PrintError("Invalid type. Use 'python' or 'bash', default is 'python'")
		return
	}

	moduleDir := filepath.Join(cli.manager.ModulesDir, moduleName)

	// Check if already exists
	if _, err := os.Stat(moduleDir); err == nil {
		core.PrintError(fmt.Sprintf("Module '%s' already exists, skipping...", moduleName))
		return
	}

	// Create directory
	if err := os.MkdirAll(moduleDir, 0755); err != nil {
		core.PrintError(fmt.Sprintf("Failed to create module directory: %v, skipping...", err))
		return
	}

	// Create module.yaml
	yamlContent := fmt.Sprintf(`name: %s
description: "Description of your module"
type: %s
author: Your Name
version: 1.0.0
tags:
  - custom
options:
  target:
    type: string
    description: Target parameter
    required: true
required:
  - target
`, moduleName, moduleType)

	yamlPath := filepath.Join(moduleDir, "module.yaml")
	if err := ioutil.WriteFile(yamlPath, []byte(yamlContent), 0644); err != nil {
		core.PrintError(fmt.Sprintf("Failed to create module.yaml: %v", err))
		return
	}

	// Create main script
	var scriptContent, scriptName string

	if moduleType == "python" {
		scriptName = "main.py"
		scriptContent = `#!/usr/bin/env python3
"""
Module: ` + moduleName + `
Description: Your module description
"""

import os
import sys

def main():
    # Get arguments from environment variables
    target = os.getenv('ARG_TARGET') or 'localhost'
    
    print(f"[*] Module executing on {target}")
    
    try:
        # Your code here
        print("[+] Module completed successfully!")
    except Exception as e:
        print(f"[!] Error: {e}")
        sys.exit(1)

if __name__ == '__main__':
    main()
`
	} else {
		scriptName = "main.sh"
		scriptContent = `#!/bin/bash
# Module: ` + moduleName + `
# Description: Your module description

TARGET="${ARG_TARGET:-localhost}"

echo "[*] Module executing on $TARGET"

# Your code here

echo "[+] Module completed successfully!"
`
	}

	scriptPath := filepath.Join(moduleDir, scriptName)
	if err := ioutil.WriteFile(scriptPath, []byte(scriptContent), 0755); err != nil {
		core.PrintError(fmt.Sprintf("Failed to create main script: %v", err))
		return
	}

	core.PrintSuccess(fmt.Sprintf("Module '%s' created successfully", moduleName))
	core.PrintInfo(fmt.Sprintf("Location: %s", moduleDir))
	fmt.Println()
}

// EditModule allows editing module files
func (cli *CLI) EditModule(moduleName string) {
	module, err := cli.manager.GetModule(moduleName)
	if err != nil {
		core.PrintError(fmt.Sprintf("Module not found: %v, try: 'search %v'", err, moduleName))
		return
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nano"
	}

	core.PrintInfo(fmt.Sprintf("Edit module '%s' (directory: %s)", moduleName, module.Path))
	fmt.Println()

	files, err := ioutil.ReadDir(module.Path)
	if err != nil {
		core.PrintError(fmt.Sprintf("Failed to read module directory: %v", err))
		return
	}

	fmt.Println("Files in module:")
	for i, file := range files {
		fmt.Printf("  %d. %s\n", i+1, core.Color("cyan", file.Name()))
	}
	fmt.Println()

	core.PrintInfo("Tip: Use 'run' command to test your changes")
	fmt.Println()
}

// DeleteModule removes a module
func (cli *CLI) DeleteModule(moduleName string) {
	module, err := cli.manager.GetModule(moduleName)
	if err != nil {
		core.PrintError(fmt.Sprintf("Module not found: %v", err))
		return
	}

	fmt.Println()
	core.PrintWarning(fmt.Sprintf("About to delete module: %s", moduleName))
	fmt.Printf("Are you sure? (yes/no): ")

	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(response)

	if strings.ToLower(response) != "yes" {
		core.PrintInfo("Cancelled")
		return
	}

	if err := os.RemoveAll(module.Path); err != nil {
		core.PrintError(fmt.Sprintf("Failed to delete module: %v", err))
		return
	}

	core.PrintSuccess(fmt.Sprintf("Module '%s' deleted successfully", moduleName))
	fmt.Println()
}
