package user

import (
  "fmt"
  "os"
  "path/filepath"

  "github.com/BurntSushi/toml"
	"github.com/fatih/color"
)

type User struct {
  Name         string   // set from userId
  DefaultAgent string   `toml:"default_agent"`
  Agents       []string // populated from configs/{user}/agents/*.toml
}

func LoadUser(userId string) (*User, error) {
  // 1. load default_agent
  userConfigPath := filepath.Join("configs", userId, "user_setting.toml")
  var u User
  if _, err := toml.DecodeFile(userConfigPath, &u); err != nil {
    return nil, fmt.Errorf("failed to load user config %q: %w", userConfigPath, err)
  }
  u.Name = userId

  // 2. scan agents directory
  agentsDir := filepath.Join("configs", userId, "agents")
  entries, err := os.ReadDir(agentsDir)
  if err != nil {
    return nil, fmt.Errorf("failed to read agents dir %q: %w", agentsDir, err)
  }
  for _, e := range entries {
    if e.IsDir() || filepath.Ext(e.Name()) != ".toml" {
      continue
    }
    name := e.Name()[:len(e.Name())-len(".toml")]
    u.Agents = append(u.Agents, name)
  }
  return &u, nil
}

func (u *User) AgentPath(name string) (string, error) {
  for _, a := range u.Agents {
    if a == name {
      return filepath.Join("configs", u.Name, "agents", name+".toml"), nil
    }
  }
  return "", fmt.Errorf("agent %q not found for user %q", name, u.Name)
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
