package user

import (
  "fmt"
  "path/filepath"

  "github.com/BurntSushi/toml"
  "github.com/johnjallday/dolphin-tool-calling-agent/internal/agent"
)

type AgentMeta struct {
  Name    string   `toml:"name"`
  Model   string   `toml:"model"`
  Plugins []string `toml:"plugins"`
}

type User struct {
  Name         string
  Agents       []AgentMeta
  DefaultAgent *agent.Agent
}

func NewUser(userID string) (*User, error) {
  path := filepath.Join("configs", "users", userID+".toml")
  fmt.Println("Loading user config:", path)

  var raw struct {
    Name         string      `toml:"name"`
    DefaultAgent string      `toml:"default_agent"`
    Agents       []AgentMeta `toml:"agents"`
  }

  if _, err := toml.DecodeFile(path, &raw); err != nil {
    return nil, fmt.Errorf("decode %s: %w", path, err)
  }

  u := &User{Name: raw.Name, Agents: raw.Agents}
  for _, meta := range raw.Agents {
    if meta.Name == raw.DefaultAgent {
			ag, err := agent.NewAgent(meta.Name, meta.Model, meta.Plugins)
      if err != nil {
        return nil, fmt.Errorf("init default agent %q: %w", meta.Name, err)
      }
      u.DefaultAgent = ag
      break
    }
  }
  return u, nil
}


func (u *User) String() string {
  s := fmt.Sprintf("User: %s\n", u.Name)
  if u.DefaultAgent != nil {
    s += fmt.Sprintf("Default Agent: %s\n", u.DefaultAgent.Name)
  } else {
    s += "Default Agent: <none>\n"
  }
  s += "Available Agents:\n"
  for _, meta := range u.Agents {
    s += fmt.Sprintf("  - %s (plugins: %v)\n", meta.Name, meta.Plugins)
  }
  return s
}


