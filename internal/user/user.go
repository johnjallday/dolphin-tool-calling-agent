package user

import (
  "fmt"
  "path/filepath"
	"github.com/fatih/color"
  "github.com/BurntSushi/toml"
  "github.com/johnjallday/dolphin-tool-calling-agent/internal/agent"
  "github.com/openai/openai-go"
)

type AgentMeta struct {
  Name      string   `toml:"name"`
  Model     string   `toml:"model"`
  ToolPaths []string `toml:"tool_path"`
}

type User struct {
  Name         string
  Agents       []AgentMeta
  DefaultAgent *agent.Agent
}

func NewUser(userID string, client *openai.Client) (*User, error) {
  path := filepath.Join("configs", "users", userID+".toml")
  var raw struct {
    Name         string      `toml:"name"`
    DefaultAgent string      `toml:"default_agent"`
    Agents       []AgentMeta `toml:"agents"`
  }
  if _, err := toml.DecodeFile(path, &raw); err != nil {
    return nil, fmt.Errorf("load user %q: %w", userID, err)
  }

  u := &User{Name: raw.Name, Agents: raw.Agents}
  for _, meta := range raw.Agents {
		fmt.Println(meta.Name)
		fmt.Println(meta.ToolPaths)
    if meta.Name == raw.DefaultAgent {
      ag, err := agent.NewAgent(client, meta.Name, meta.Model, meta.ToolPaths)
      if err != nil {
        return nil, fmt.Errorf("init default agent %q: %w", meta.Name, err)
      }
      u.DefaultAgent = ag
      break
    }
  }
  return u, nil
}

func (u *User) Print() {
  cLabel := color.New(color.FgCyan, color.Bold)
  cValue := color.New(color.FgWhite)
  cList  := color.New(color.FgMagenta, color.Bold)
  cItem  := color.New(color.FgGreen)

  cLabel.Print("User: "); cValue.Println(u.Name)
  cLabel.Print("Default Agent: ")
  if u.DefaultAgent != nil {
    cValue.Println(u.DefaultAgent.Name)
  } else {
    cValue.Println("<none>")
  }
  cList.Println("Available Agents:")
  for _, meta := range u.Agents {
    fmt.Print("  "); cItem.Println("- " + meta.Name)
  }
}
