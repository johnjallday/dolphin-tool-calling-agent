package user

import (
  "fmt"
  "os"
  "path/filepath"

  "github.com/BurntSushi/toml"
	"github.com/fatih/color"
  "github.com/johnjallday/dolphin-tool-calling-agent/internal/agent"
)

// User holds per-user settings, available agents, and the default agent instance.
type User struct {
  Name         string        // set from userId
  Agents       []string      // filenames (no .toml)
  DefaultAgent *agent.Agent  // loaded instance, or nil if unset
}

func NewUser(userId string) (*User, error) {
  userDir := filepath.Join("configs", userId)
  if err := os.MkdirAll(userDir, 0755); err != nil {
    return nil, fmt.Errorf("mkdir user dir: %w", err)
  }

  cfgPath := filepath.Join(userDir, "user_setting.toml")
  if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
    // write minimal config with no default_agent
    f, err := os.Create(cfgPath)
    if err != nil {
      return nil, fmt.Errorf("create %q: %w", cfgPath, err)
    }
    defer f.Close()
    // only default_agent field, empty by default
    _, err = f.WriteString("default_agent = \"\"\n")
    if err != nil {
      return nil, fmt.Errorf("init config: %w", err)
    }
  } else if err != nil {
    return nil, fmt.Errorf("stat %q: %w", cfgPath, err)
  }

  // decode only default_agent name
  var rawCfg struct {
    DefaultName string `toml:"default_agent"`
  }
  if _, err := toml.DecodeFile(cfgPath, &rawCfg); err != nil {
    return nil, fmt.Errorf("decode %q: %w", cfgPath, err)
  }

  // scan agents directory
  agentsDir := filepath.Join(userDir, "agents")
  if err := os.MkdirAll(agentsDir, 0755); err != nil {
    return nil, fmt.Errorf("mkdir agents dir: %w", err)
  }
  entries, err := os.ReadDir(agentsDir)
  if err != nil {
    return nil, fmt.Errorf("read agents dir: %w", err)
  }
  u := &User{Name: userId}
  for _, e := range entries {
    if !e.IsDir() && filepath.Ext(e.Name()) == ".toml" {
      name := e.Name()[:len(e.Name())-5]
      u.Agents = append(u.Agents, name)
    }
  }

  // load default agent instance or print notice
  if rawCfg.DefaultName == "" {
    fmt.Println("no default agent set")
  } else {
    if !contains(u.Agents, rawCfg.DefaultName) {
      return nil, fmt.Errorf("configured default %q not found", rawCfg.DefaultName)
    }
    path := filepath.Join(agentsDir, rawCfg.DefaultName+".toml")
    var ag agent.Agent
    if _, err := toml.DecodeFile(path, &ag); err != nil {
      return nil, fmt.Errorf("load default agent %q: %w", rawCfg.DefaultName, err)
    }
    u.DefaultAgent = &ag
  }

  return u, nil
}

func contains(list []string, v string) bool {
  for _, x := range list {
    if x == v {
      return true
    }
  }
  return false
}



// Print writes the user’s details in colorized, indented format.
func (u *User) Print() {
  // Define colors
  labelColor := color.New(color.FgCyan, color.Bold)
  valueColor := color.New(color.FgWhite)
  listLabel := color.New(color.FgMagenta, color.Bold)
  itemColor := color.New(color.FgGreen)

  // Print header
  labelColor.Print("User: ")
  valueColor.Println(u.Name)

  // Print default agent
  labelColor.Print("Default Agent: ")
  valueColor.Println(u.DefaultAgent)

  // Print agent list
  listLabel.Println("Available Agents:")
  for _, a := range u.Agents {
    fmt.Print("  ")        // indent
    itemColor.Println("- " + a)
  }
}
