package app

import (
  "fmt"
  "os"
  "path/filepath"
	"context"

  "github.com/BurntSushi/toml"
  "github.com/johnjallday/dolphin-tool-calling-agent/internal/agent"
  "github.com/johnjallday/dolphin-tool-calling-agent/internal/user"
	"github.com/johnjallday/dolphin-tool-calling-agent/internal/registry"
	"github.com/johnjallday/dolphin-tool-calling-agent/pkg/tools"
)

type AppConfig struct {
  DefaultUser string `toml:"default_user"`
}

type DefaultApp struct {
  user *user.User
	agent *agent.Agent
}

func NewApp() App { return &DefaultApp{} }

func (a *DefaultApp) Init() error {
  cfgDir := "configs"
  if err := os.MkdirAll(cfgDir, 0755); err != nil {
    return fmt.Errorf("create configs folder: %w", err)
  }

  settingPath := filepath.Join(cfgDir, "app_setting.toml")
  if _, err := os.Stat(settingPath); err != nil {
    if os.IsNotExist(err) {
      f, err := os.Create(settingPath)
      if err != nil {
        return fmt.Errorf("create app_setting.toml: %w", err)
      }
      defer f.Close()
      _, err = f.WriteString("default_user = \"\"\n")
      return err
    }
    return fmt.Errorf("stat app_setting.toml: %w", err)
  }

  var cfg AppConfig
  if _, err := toml.DecodeFile(settingPath, &cfg); err != nil {
    return fmt.Errorf("parse app_setting.toml: %w", err)
  }

  if cfg.DefaultUser == "" {
    return fmt.Errorf("default_user not set in %s", settingPath)
  }

  return a.LoadUser(cfg.DefaultUser)
}

func (a *DefaultApp) Users() []string {
  var names []string
  dir := "configs/users"
  entries, err := os.ReadDir(dir)
  if err != nil {
    return names
  }
  for _, e := range entries {
    if e.IsDir() || filepath.Ext(e.Name()) != ".toml" {
      continue
    }
    var m struct{ Name string `toml:"name"` }
    if _, err := toml.DecodeFile(filepath.Join(dir, e.Name()), &m); err == nil && m.Name != "" {
      names = append(names, m.Name)
    }
  }
  return names
}

// LoadUser loads the user TOML and then loads the default agent.
func (a *DefaultApp) LoadUser(username string) error {
  u, err := user.NewUser(username)
  if err != nil {
    return fmt.Errorf("load user %q: %w", username, err)
  }
  a.user = u
  a.agent = u.DefaultAgent
  return nil
}

func (a *DefaultApp) User() *user.User {
    if a.user == nil {
        return &user.User{
            Name:         "<none>",
            DefaultAgent: nil,
            Agents:       nil,
        }
    }
    return a.user
}

func (a *DefaultApp) Agent() *agent.Agent {
  if a.agent == nil {
    return &agent.Agent{
      Name:     "<none>",
      Model:    "<none>",
      Registry: registry.NewToolRegistry(),
    }
  }
  return a.agent

}

// LoadAgent selects one of the user’s agents by name and sets it to agent.
func (a *DefaultApp) LoadAgent(agentName string) error {
  if a.user == nil {
    return fmt.Errorf("no user loaded")
  }
  var meta *user.AgentMeta
  for i := range a.user.Agents {
    if a.user.Agents[i].Name == agentName {
      meta = &a.user.Agents[i]
      break
    }
  }
  if meta == nil {
    fmt.Println("agent %q not found for user %q", agentName, a.user.Name)
    return nil
  }

  ag, err := agent.NewAgent(meta.Name, meta.Model, meta.Plugins)
  if err != nil {
    return fmt.Errorf("init agent %q: %w", meta.Name, err)
  }

  a.agent = ag
  // also update the default in the user struct if desired
  a.user.DefaultAgent = ag
  return nil
}

func (a *DefaultApp) UnloadAgent() error {
  if a.agent == nil {
    return fmt.Errorf("no agent loaded")
  }
  a.agent = nil
  return nil
}

func (a *DefaultApp) UnloadUser() error {
  if a.user == nil {
    return fmt.Errorf("no user loaded")
  }
  a.user = nil
  a.agent = nil
  return nil
}

func (a *DefaultApp) SendMessage(ctx context.Context, msg string) error {
  if a.agent == nil {
    return fmt.Errorf("no agent loaded")
  }
  return a.agent.SendMessage(ctx, msg)
}


// Tools returns the slice of registered tools.
func (a *DefaultApp) Tools() []tools.Tool {
    return a.agent.Tools()
}
