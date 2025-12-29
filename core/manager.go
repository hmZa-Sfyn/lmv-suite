package core

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// ModuleManager handles module discovery, loading, and execution
type ModuleManager struct {
	ModulesDir string
	Modules    map[string]*ModuleConfig
}

// NewModuleManager creates a new module manager
func NewModuleManager(modulesDir string) *ModuleManager {
	return &ModuleManager{
		ModulesDir: modulesDir,
		Modules:    make(map[string]*ModuleConfig),
	}
}

// DiscoverModules scans the modules directory and loads module metadata
func (mm *ModuleManager) DiscoverModules() error {
	if err := os.MkdirAll(mm.ModulesDir, 0755); err != nil {
		return fmt.Errorf("failed to create modules directory: %w", err)
	}

	entries, err := os.ReadDir(mm.ModulesDir)
	if err != nil {
		return fmt.Errorf("failed to read modules directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			moduleDir := filepath.Join(mm.ModulesDir, entry.Name())
			mm.loadModuleFromDir(moduleDir)
		}
	}

	return nil
}

// loadModuleFromDir loads a module from a directory
func (mm *ModuleManager) loadModuleFromDir(moduleDir string) {
	moduleName := filepath.Base(moduleDir)
	moduleConfig := &ModuleConfig{
		Path:   moduleDir,
		Name:   moduleName,
		Loaded: false,
	}

	// Try to load metadata from module.yaml
	metadataPath := filepath.Join(moduleDir, "module.yaml")
	if _, err := os.Stat(metadataPath); err == nil {
		metadata, err := loadMetadata(metadataPath)
		if err != nil {
			moduleConfig.LoadError = err.Error()
			mm.Modules[moduleName] = moduleConfig
			return
		}
		moduleConfig.Metadata = metadata
		moduleConfig.Type = metadata.Type
	} else {
		// Try to infer type from available files
		moduleConfig.Type = mm.inferModuleType(moduleDir)
	}

	moduleConfig.Loaded = true
	mm.Modules[moduleName] = moduleConfig
}

// inferModuleType determines module type based on file extensions
func (mm *ModuleManager) inferModuleType(moduleDir string) string {
	entries, _ := os.ReadDir(moduleDir)
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasSuffix(name, ".py") {
			return "python"
		}
		if strings.HasSuffix(name, ".sh") {
			return "bash"
		}
		if strings.HasSuffix(name, ".go") {
			return "go"
		}
	}
	return "unknown"
}

// GetModule returns a module by name
func (mm *ModuleManager) GetModule(name string) (*ModuleConfig, error) {
	module, exists := mm.Modules[name]
	if !exists {
		return nil, fmt.Errorf("module '%s' not found, did you forget to load it?", name)
	}
	if !module.Loaded {
		return nil, fmt.Errorf("module '%s' failed to load: %s, did you forget to load it?", name, module.LoadError)
	}
	return module, nil
}

// ExecuteModule runs a module with given arguments
func (mm *ModuleManager) ExecuteModule(moduleName string, args map[string]string) (*ExecutionResult, error) {
	module, err := mm.GetModule(moduleName)
	if err != nil {
		return nil, err
	}

	switch module.Type {
	case "python":
		return executePythonModule(module, args)
	case "bash":
		return executeBashModule(module, args)
	case "go":
		return executeGoModule(module, args)
	default:
		return nil, fmt.Errorf("unsupported module type: %s, supported types are: python, bash", module.Type)
	}
}

// executePythonModule runs a Python module with real-time output
func executePythonModule(module *ModuleConfig, args map[string]string) (*ExecutionResult, error) {
	result := &ExecutionResult{
		Timestamp: time.Now(),
	}

	// Find the main Python script
	scriptPath := findMainScript(module.Path, ".py")
	if scriptPath == "" {
		result.Success = false
		result.Error = "no Python script found in module, expected .py file, e.g., main.py or run.py"
		result.ExitCode = 1
		return result, nil
	}

	// Build command
	cmd := exec.Command("python3", scriptPath)
	cmd.Dir = module.Path

	// Set environment variables for arguments
	env := os.Environ()
	for key, value := range args {
		env = append(env, fmt.Sprintf("ARG_%s=%s", strings.ToUpper(key), value))
	}
	cmd.Env = env

	// Stream output in real-time
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()

	if err != nil {
		result.Success = false
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.ExitCode = 1
		}
		result.Error = err.Error()
	} else {
		result.Success = true
		result.ExitCode = 0
	}

	return result, nil
}

// executeBashModule runs a Bash script module with real-time output
func executeBashModule(module *ModuleConfig, args map[string]string) (*ExecutionResult, error) {
	result := &ExecutionResult{
		Timestamp: time.Now(),
	}

	scriptPath := findMainScript(module.Path, ".sh")
	if scriptPath == "" {
		result.Success = false
		result.Error = "no Bash script found in module"
		result.ExitCode = 1
		return result, nil
	}

	cmd := exec.Command("bash", scriptPath)
	cmd.Dir = module.Path

	// Set environment variables
	env := os.Environ()
	for key, value := range args {
		env = append(env, fmt.Sprintf("ARG_%s=%s", strings.ToUpper(key), value))
	}
	cmd.Env = env

	// Stream output in real-time
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()

	if err != nil {
		result.Success = false
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.ExitCode = 1
		}
		result.Error = err.Error()
	} else {
		result.Success = true
		result.ExitCode = 0
	}

	return result, nil
}

// executeGoModule runs a Go module
func executeGoModule(module *ModuleConfig, args map[string]string) (*ExecutionResult, error) {
	result := &ExecutionResult{
		Timestamp: time.Now(),
		Success:   false,
		Error:     "Go module execution not yet implemented, please build and run manually, thanks!",
		ExitCode:  1,
	}
	return result, nil
}

// findMainScript finds the main script in a module directory
func findMainScript(moduleDir string, extension string) string {
	entries, _ := os.ReadDir(moduleDir)
	for _, entry := range entries {
		if !entry.IsDir() && entry.Name() == "main"+extension {
			return filepath.Join(moduleDir, entry.Name())
		}
	}
	return ""
}

// loadMetadata loads module metadata from YAML file
func loadMetadata(path string) (*ModuleMetadata, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var metadata ModuleMetadata
	if err := yaml.Unmarshal(data, &metadata); err != nil {
		return nil, err
	}

	return &metadata, nil
}

// ListModules returns all loaded modules
func (mm *ModuleManager) ListModules() []*ModuleConfig {
	var modules []*ModuleConfig
	for _, module := range mm.Modules {
		if module.Loaded {
			modules = append(modules, module)
		}
	}
	return modules
}
