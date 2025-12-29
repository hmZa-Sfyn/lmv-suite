package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"lanmanvan/core"
)

// EnvironmentManager handles global environment variables
type EnvironmentManager struct {
	vars     map[string]string
	filePath string
}

// NewEnvironmentManager creates a new environment manager
func NewEnvironmentManager() *EnvironmentManager {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".lanmanvan")
	os.MkdirAll(configDir, 0700)

	em := &EnvironmentManager{
		vars:     make(map[string]string),
		filePath: filepath.Join(configDir, "env.json"),
	}

	em.Load()
	return em
}

// Set sets a global environment variable
func (em *EnvironmentManager) Set(key, value string) error {
	em.vars[key] = value
	return em.Save()
}

// Get retrieves a global environment variable
func (em *EnvironmentManager) Get(key string) (string, bool) {
	val, exists := em.vars[key]
	return val, exists
}

// GetAll returns all environment variables
func (em *EnvironmentManager) GetAll() map[string]string {
	return em.vars
}

// Delete removes an environment variable
func (em *EnvironmentManager) Delete(key string) error {
	delete(em.vars, key)
	return em.Save()
}

// Save persists environment variables to JSON file
func (em *EnvironmentManager) Save() error {
	data, err := json.MarshalIndent(em.vars, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(em.filePath, data, 0600)
}

// Load reads environment variables from JSON file
func (em *EnvironmentManager) Load() error {
	data, err := ioutil.ReadFile(em.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // File doesn't exist yet, that's okay
		}
		return err
	}

	return json.Unmarshal(data, &em.vars)
}

// Clear removes all environment variables
func (em *EnvironmentManager) Clear() error {
	em.vars = make(map[string]string)
	return em.Save()
}

// Display shows all environment variables
func (em *EnvironmentManager) Display() {
	if len(em.vars) == 0 {
		core.PrintWarning("No global environment variables set, use '<key> = <value>' to add some, or type '<key>=?' to view its value")
		fmt.Println()
		return
	}

	fmt.Println()
	fmt.Println(core.NmapBox("GLOBAL ENVIRONMENT VARIABLES"))

	i := 0
	for key, value := range em.vars {
		i++
		prefix := "   ├─ "
		if i == len(em.vars) {
			prefix = "   └─ "
		}
		fmt.Printf("%s%s = %s\n", prefix, core.Color("cyan", key), core.Color("green", value))
	}
	fmt.Println()
}
