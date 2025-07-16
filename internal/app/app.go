package app

import (
  "fmt"
  "os"
  "path/filepath"
	"context"

  "github.com/BurntSushi/toml"
  "github.com/johnjallday/dolphin-tool-calling-agent/internal/agent"
  "github.com/johnjallday/dolphin-tool-calling-agent/internal/user"
)

type AppConfig struct {
  DefaultUser string `toml:"default_user"`
}

type DefaultApp struct {
  currentUser *user.User
	currentAgent *agent.Agent
}

func NewApp() App { return &DefaultApp{} }

func (a *DefaultApp) Init() error {
  cfgDir := "configs"
  if err := os.MkdirAll(cfgDir, 0755); err != nil {
    return fmt.Errorf("create configs folder: %w", err)
  }

  settingPath := filepath.Join(cfgDir, "app_setting.toml")
  var cfg AppConfig
  if info, err := os.Stat(settingPath); err != nil {
    if os.IsNotExist(err) {
      f, err := os.Create(settingPath)
      if err != nil {
        return fmt.Errorf("create app_setting.toml: %w", err)
      }
      defer f.Close()
      _, err = f.WriteString("default_user = \"\"\n")
      fmt.Println("creating configs folder")
      fmt.Println("creating app_setting.toml")
      return err
    }
    return fmt.Errorf("stat app_setting.toml: %w", err)
  } else if info.IsDir() {
    return fmt.Errorf("app_setting.toml is a directory")
  }

  if _, err := toml.DecodeFile(settingPath, &cfg); err != nil {
    return fmt.Errorf("parse app config: %w", err)
  }

  if cfg.DefaultUser == "" {
    fmt.Errorf("default_user not set in app_setting.toml")
		return nil
  }

	if err := a.LoadUser(cfg.DefaultUser); err != nil {
    return err
  }

	a.currentAgent = a.currentUser.DefaultAgent
  return nil
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

func (a *DefaultApp) LoadUser(username string) error {
  fn := filepath.Join("configs", "users", username+".toml")
  var tmp struct {
    Name             string           `toml:"name"`
    DefaultAgentName string           `toml:"default_agent"`
    Agents           []user.AgentMeta `toml:"agents"`
  }
  if _, err := toml.DecodeFile(fn, &tmp); err != nil {
    return fmt.Errorf("load user %q: %w", username, err)
  }

  var meta *user.AgentMeta
  for i := range tmp.Agents {
    if tmp.Agents[i].Name == tmp.DefaultAgentName {
      meta = &tmp.Agents[i]
      break
    }
  }
  if meta == nil {
    return fmt.Errorf("default agent %q not found in %q", tmp.DefaultAgentName, username)
  }

  ag, err := agent.NewAgent(meta.Name, meta.Model, meta.ToolPaths)
  if err != nil {
    return fmt.Errorf("init agent %q: %w", meta.Name, err)
  }

  a.currentUser = &user.User{
    Name:         tmp.Name,
    Agents:       tmp.Agents,
    DefaultAgent: ag,
  }
  return nil
}

func (a *DefaultApp) CurrentUser() *user.User    { return a.currentUser }
func (a *DefaultApp) CurrentAgent() *agent.Agent { return a.currentAgent }

func (a *DefaultApp) SendMessage(ctx context.Context, msg string) error {
  if a.currentAgent == nil {
    return fmt.Errorf("no agent loaded")
  }
  return a.currentAgent.SendMessage(ctx, msg)
}
