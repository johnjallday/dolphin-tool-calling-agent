package agent

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/BurntSushi/toml"
)

// ListAgents reads all TOML files in "./user/agents" and returns their AgentConfig.
func ListAgents() ([]AgentConfig, error) {
    const dir = "./configs/user/agents"
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
