package main


import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"


	"github.com/openai/openai-go"
	"github.com/BurntSushi/toml"


	"github.com/johnjallday/dolphin-tool-calling-agent/pkg/tools"
)

const (
		packName		= "Reaper Project Manager"
    packVersion = "v0.0.1"
    packLink    = "https://github.com/johnjallday/dolphin-tool-calling-agent/"
)

type ReaperConfig struct {
	DefaultTemplate string `toml:"default_template"`
	ScriptPath      string `toml:"script_path"`
}

var reaperConfig ReaperConfig

// Tool defines schema and executor for CreateNewProject.
var CreateNewProjectTool = tools.Tool{
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
	registerConfig()
	if reaperConfig.DefaultTemplate == "" {
		return "", fmt.Errorf("default template not configured")
	}
	projectDir := name
	fmt.Println(projectDir)
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


func PluginSpecs() []tools.Tool {

	return []tools.Tool{ CreateNewProjectTool }
}

func PluginPackage() tools.ToolPackage {
    return tools.ToolPackage{
				Name:		 packName,
        Version: packVersion,
        Link:    packLink,
        Tools:   []tools.Tool{ CreateNewProjectTool },
    }
}
