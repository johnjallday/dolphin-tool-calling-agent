package reaper

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/openai/openai-go"
	//"github.com/pelletier/go-toml/v2"

	"github.com/johnjallday/dolphin-tool-calling-agent/tools"
)

type ReaperConfig struct {
	DefaultTemplate string `toml:"default_template"`
	ScriptPath      string `toml:"script_path"`
}

// Declare a package-level variable to hold the config
var reaperConfig ReaperConfig

// ToolSpec defines schema and executor for CreateNewProject.
var CreateNewProjectTool = tools.ToolSpec{
	Name:        "create_new_project",
	Description: "Create a new Reaper project with a name and bpm",
	Parameters: openai.FunctionParameters{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]string{"type": "string"},
			"bpm":  map[string]string{"type": "integer"},
		},
		"required": []string{"name"},
	},
	Exec: func(args map[string]interface{}) (string, error) {
		name := args["name"].(string)
		var bpm int
		if v, ok := args["bpm"].(float64); ok {
			bpm = int(v)
		}
		return CreateNewProject(name, bpm)
	},
}

func CreateNewProject(name string, bpm int) (string, error) {
	if reaperConfig.DefaultTemplate == "" {
		return "", fmt.Errorf("default template not configured")
	}
	projectDir := name
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return "", err
	}
	dest := filepath.Join(projectDir, name+".RPP")
	data, err := os.ReadFile(reaperConfig.DefaultTemplate)
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(dest, data, 0644); err != nil {
		return "", err
	}
	if bpm > 0 {
		contentBytes, err := os.ReadFile(dest)
		if err != nil {
			return "", err
		}
		lines := strings.Split(string(contentBytes), "\n")
		for i, line := range lines {
			trimmed := strings.TrimLeft(line, " \t")
			if strings.HasPrefix(trimmed, "TEMPO ") {
				indent := line[:len(line)-len(trimmed)]
				parts := strings.Fields(trimmed)
				if len(parts) >= 2 {
					parts[1] = strconv.Itoa(bpm)
					lines[i] = indent + strings.Join(parts, " ")
				}
				break
			}
		}
		replaced := strings.Join(lines, "\n")
		if err := os.WriteFile(dest, []byte(replaced), 0644); err != nil {
			return "", err
		}
	}
	cmd := exec.Command("open", "-a", "Reaper", dest)
	fmt.Printf("Executing command: %s\n", strings.Join(cmd.Args, " "))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	msg := fmt.Sprintf("Created and launched project: %s", dest)
	if bpm > 0 {
		msg += fmt.Sprintf(" (BPM %d)", bpm)
	}
	return msg, nil
}


