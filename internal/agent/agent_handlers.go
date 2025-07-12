package agent

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/BurntSushi/toml"
   	"github.com/johnjallday/dolphin-tool-calling-agent/pkg/tools"
)

// ListAgents reads all the .toml files in ./configs/<userName>/agents
func ListAgents(userName string) ([]AgentConfig, error) {
    dir := filepath.Join("configs", userName, "agents")

    entries, err := os.ReadDir(dir)
    if err != nil {
        return nil, fmt.Errorf("read dir %q: %w", dir, err)
    }

    var configs []AgentConfig
    for _, e := range entries {
        if e.IsDir() || filepath.Ext(e.Name()) != ".toml" {
            continue
        }
        path := filepath.Join(dir, e.Name())
        var cfg AgentConfig
        if _, err := toml.DecodeFile(path, &cfg); err != nil {
            return nil, fmt.Errorf("decode %q: %w", e.Name(), err)
        }
        configs = append(configs, cfg)
    }
    return configs, nil
}


func CreateAgent(){
	//CreateAgent
	//
	fmt.Println("Creating Agent")
	//ask for name:
	files := tools.GetAvailableToolPacks()
	fmt.Println(files)
}

