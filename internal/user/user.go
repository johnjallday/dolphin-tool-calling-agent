package user

import (
  "fmt"
  "os"
  "path/filepath"

  "github.com/BurntSushi/toml"
  "github.com/fatih/color"
  "github.com/johnjallday/dolphin-tool-calling-agent/internal/agent"
)

type AgentMeta struct {
  Name string
  Path string
}

type User struct {
  Name         string
  Agents       []AgentMeta
  DefaultAgent *agent.Agent
}

func loadAgent(meta AgentMeta) (*agent.Agent, error) {
  var ag agent.Agent
  if _, err := toml.DecodeFile(meta.Path, &ag); err != nil {
    return nil, fmt.Errorf("load agent %q: %w", meta.Name, err)
  }
  return &ag, nil
}

func NewUser(userId string) (*User, error) {
  userDir := filepath.Join("configs", userId)
  os.MkdirAll(userDir, 0755)
  cfgPath := filepath.Join(userDir, "user_setting.toml")
  if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
    f, _ := os.Create(cfgPath); defer f.Close()
    f.WriteString("default_agent = \"\"\n")
  }

  var rawCfg struct{ DefaultName string `toml:"default_agent"` }
  toml.DecodeFile(cfgPath, &rawCfg)

  agentsDir := filepath.Join(userDir, "agents")
  os.MkdirAll(agentsDir, 0755)
  entries, _ := os.ReadDir(agentsDir)

  u := &User{Name: userId}
  for _, e := range entries {
    if e.IsDir() || filepath.Ext(e.Name()) != ".toml" {
      continue
    }
    name := e.Name()[:len(e.Name())-5]
    path := filepath.Join(agentsDir, e.Name())
    meta := AgentMeta{Name: name, Path: path}
    u.Agents = append(u.Agents, meta)
    if name == rawCfg.DefaultName {
      ag, err := loadAgent(meta)
      if err != nil {
        return nil, err
      }
      u.DefaultAgent = ag
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
