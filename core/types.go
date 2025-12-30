package core

import "time"

// ModuleMetadata holds information about a module
type ModuleMetadata struct {
	Name        string                `yaml:"name"`
	Description string                `yaml:"description"`
	Type        string                `yaml:"type"` // python, bash, go
	Author      string                `yaml:"author"`
	Version     string                `yaml:"version"`
	Options     map[string]OptionMeta `yaml:"options"`
	Required    []string              `yaml:"required"`
	Tags        []string              `yaml:"tags"`
	GitHubURL   string                `yaml:"github_url"`
	XUrl        string                `yaml:"x_url"`
}

// OptionMeta describes a module option
type OptionMeta struct {
	Type        string `yaml:"type"` // string, int, bool, file
	Description string `yaml:"description"`
	Default     string `yaml:"default"`
	Required    bool   `yaml:"required"`
}

// ExecutionRequest represents a module execution request
type ExecutionRequest struct {
	ModuleName string
	Arguments  map[string]string
	Timestamp  time.Time
}

// ExecutionResult represents module execution output
type ExecutionResult struct {
	Success   bool
	Output    string
	Error     string
	ExitCode  int
	Timestamp time.Time
}

// ModuleConfig represents runtime configuration
type ModuleConfig struct {
	Path      string
	Name      string
	Type      string
	Metadata  *ModuleMetadata
	Loaded    bool
	LoadError string
}
