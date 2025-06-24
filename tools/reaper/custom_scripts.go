package reaper

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/openai/openai-go"

	"github.com/johnjallday/dolphin-tool-calling-agent/tools"
)

func launchTool(scriptName string) error {
	scriptPath := filepath.Join("tools", "reaper", "custom_scripts", scriptName+".lua")
	cmd := exec.Command("open", "-a", "Reaper", scriptPath)
	return cmd.Run()
}

// RegisterCustomScripts scans the custom_scripts directory, builds a ToolSpec for each Lua script, 
// // and returns a slice of these specs. 
func LoadCustomScripts() []tools.ToolSpec { 
	dir := filepath.Join("tools", "reaper", "custom_scripts") 
	entries, err := os.ReadDir(dir) 
	if err != nil { 
		panic(fmt.Errorf("failed to read custom scripts directory: %w", err))
	}

	var specs []tools.ToolSpec
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

var ReaperCustomScriptsSpecs = LoadCustomScripts()
