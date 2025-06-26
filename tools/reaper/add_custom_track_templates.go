package reaper

import(
	"fmt"
	"path/filepath"
	"os/exec"
	"strings"
	"os"

	"github.com/openai/openai-go"
	"github.com/johnjallday/dolphin-tool-calling-agent/tools"
)
// ToolSpec defines schema and executor for CreateNewProject.
func launchRTemplate(scriptName string) error {
	homeDir := os.Getenv("HOME") 
	scriptPath := filepath.Join(homeDir, "Library", "Application Support", "REAPER", "TrackTemplates", scriptName+".RTrackTemplate") 
	cmd := exec.Command("open", "-a", "Reaper", scriptPath)
	return cmd.Run()
}

// RegisterCustomScripts scans the custom_scripts directory, builds a ToolSpec for each Lua script, 
// // and returns a slice of these specs. 
func LoadCustomTemplates() []tools.ToolSpec { 

	homeDir := os.Getenv("HOME") 
	dir := filepath.Join(homeDir, "Library", "Application Support", "REAPER", "TrackTemplates") 
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
		if !strings.HasSuffix(name, ".RTrackTemplate") {
			continue
		}
		TrackName := strings.TrimSuffix(name, ".RTrackTemplate")
		// Capture the current scriptName in a local variable to avoid closure issues.
		currentTrack := TrackName

		spec := tools.ToolSpec{
			Name:        currentTrack,
			Description: fmt.Sprintf("Add Reaper Track Template %q", currentTrack),
			Parameters:  openai.FunctionParameters{}, // No parameters in this case.
			Exec: func(args map[string]interface{}) (string, error) {
				fmt.Println("Test")
				if err := launchRTemplate(currentTrack); err != nil {
					fmt.Println("error here")
					return "", err
				}
				return fmt.Sprintf("Loaded Reaper Track Template %s", currentTrack), nil
			},
		}
		specs = append(specs, spec)
	}
	return specs
	}

var ReaperCustomTrackSpecs = LoadCustomTemplates()

