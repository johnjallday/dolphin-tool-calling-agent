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
  // ensure user dir
  userDir := filepath.Join("configs", userId)
  if err := os.MkdirAll(userDir, 0755); err != nil {
    return nil, fmt.Errorf("mkdir user dir: %w", err)
  }

  // ensure user_setting.toml
  cfgPath := filepath.Join(userDir, "user_setting.toml")
  if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
    // write a minimal file
    u0 := &User{Name: userId}
    f, err := os.Create(cfgPath)
    if err != nil {
      return nil, fmt.Errorf("create %q: %w", cfgPath, err)
    }
    defer f.Close()
    if err := toml.NewEncoder(f).Encode(u0); err != nil {
      return nil, fmt.Errorf("encode default user: %w", err)
    }
  } else if err != nil {
    return nil, fmt.Errorf("stat %q: %w", cfgPath, err)
  }

  // now decode it
  var u User
  if _, err := toml.DecodeFile(cfgPath, &u); err != nil {
    return nil, fmt.Errorf("decode %q: %w", cfgPath, err)
  }
  u.Name = userId

  // ensure agents dir & scan it
  agentsDir := filepath.Join(userDir, "agents")
  if err := os.MkdirAll(agentsDir, 0755); err != nil {
    return nil, fmt.Errorf("mkdir agents dir: %w", err)
  }
  entries, err := os.ReadDir(agentsDir)
  if err != nil {
    return nil, fmt.Errorf("read agents dir: %w", err)
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
