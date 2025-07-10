package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/openai/openai-go"
	"github.com/johnjallday/dolphin-tool-calling-agent/pkg/tools"
	"github.com/BurntSushi/toml"
)

const (
		packName		= "Reaper Reascript Launcher"
    packVersion = "v0.0.1"
    packLink    = "https://github.com/johnjallday/dolphin-tool-calling-agent/"
)


type ReaperConfig struct {
	DefaultTemplate string `toml:"default_template"`
	ScriptPath      string `toml:"script_path"`
}

var reaperConfig ReaperConfig

func launchTool(scriptName string) error {
	scriptPath := filepath.Join("tools", "reaper", "custom_scripts", scriptName+".lua")
	cmd := exec.Command("open", "-a", "Reaper", scriptPath)
	return cmd.Run()
}

// RegisterCustomScripts scans the custom_scripts directory, builds a ToolSpec for each Lua script, 
// // and returns a slice of these specs. 
func LoadCustomScripts() []tools.ToolSpec { 
	var specs []tools.ToolSpec

	registerConfig()
	if reaperConfig.ScriptPath == "" {
		fmt.Errorf("ScriptPath Missing")
		return specs
	}

	//dir := filepath.Join("tools", "reaper", "custom_scripts") 
	entries, err := os.ReadDir(reaperConfig.ScriptPath) 
	if err != nil { 
		panic(fmt.Errorf("failed to read custom scripts directory: %w", err))
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".lua") {
			continue
		}
		scriptName := strings.TrimSuffix(name, ".lua")
		// Capture the current scriptName in a local variable to avoid closure issues.
		currentScript := scriptName

		spec := tools.ToolSpec{
			Name:        currentScript,
			Description: fmt.Sprintf("Run custom Reaper script %q", currentScript),
			Parameters:  openai.FunctionParameters{}, // No parameters in this case.
			Exec: func(args map[string]interface{}) (string, error) {
				if err := launchTool(currentScript); err != nil {
					return "", err
				}
				return fmt.Sprintf("Launched Reaper script %s", currentScript), nil
			},
		}
		specs = append(specs, spec)
	}
	return specs
}


func registerConfig() error {
	configPath := "./user/tools/reaper_project_manager/settings.toml"
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return fmt.Errorf("could not resolve settings.toml path: %w", err)
	}

	// Check if settings.toml exists, if not, create it with default values
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		fmt.Println("settings.toml not found, creating default config...")

		defaultConfig := ReaperConfig{
			DefaultTemplate: "./tools/reaper/Default.RPP", // You can change this placeholder
			ScriptPath:      "./tools/reaper/scripts",              // You can change this placeholder
		}
		// Serialize and write default config
		f, ferr := os.Create(absPath)
		if ferr != nil {
			return fmt.Errorf("could not create settings.toml: %w", ferr)
		}
		defer f.Close()
		enc := toml.NewEncoder(f)
		if err := enc.Encode(defaultConfig); err != nil {
			return fmt.Errorf("could not encode default config: %w", err)
		}
		fmt.Printf("Created %s with default values. Please edit as needed.\n", absPath)
		reaperConfig = defaultConfig
		return nil
	}

	// Load settings.toml
	if _, err := toml.DecodeFile(absPath, &reaperConfig); err != nil {
		return fmt.Errorf("failed to decode settings.toml: %w", err)
	}
	return nil
}



var ReaperCustomScriptsSpecs = LoadCustomScripts()
// PluginSpecs is looked up by the host application (via NewAgentFromConfig).
// It returns the slice of ToolSpec that the plugin makes available.
// If you need to merge these with other tools (e.g. arithmetic tools), you could combine the slices here.
func PluginSpecs() []tools.ToolSpec {
	return ReaperCustomScriptsSpecs
}


func PluginPackage() tools.ToolPackage {
    return tools.ToolPackage{
				Name:		 packName,
        Version: packVersion,
        Link:    packLink,
        Specs:   LoadCustomScripts(),
    }
}

