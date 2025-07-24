package user

import (
  "fmt"
  "os"
  "path/filepath"

  "github.com/BurntSushi/toml"
)

// userConfig mirrors your on‐disk layout.
type userConfig struct {
  Name         string      `toml:"name"`
  DefaultAgent string      `toml:"default_agent"`
  Agents       []AgentMeta `toml:"agents"`
}

// CreateUser creates configs/users/<userID>.toml, using userID
// as the display name, and returns the loaded *User.
func CreateUser(userID string) (*User, error) {
  dir := filepath.Join("configs", "users")
  if err := os.MkdirAll(dir, 0755); err != nil {
    return nil, fmt.Errorf("mkdir %q: %w", dir, err)
  }

  userFile := filepath.Join(dir, userID+".toml")
  if _, err := os.Stat(userFile); err == nil {
    return nil, fmt.Errorf("user %q already exists", userID)
  } else if !os.IsNotExist(err) {
    return nil, fmt.Errorf("stat %q: %w", userFile, err)
  }

  // Build the on‐disk struct
  cfg := userConfig{
    Name:         userID,        // use the ID as the display name
    DefaultAgent: "",
    Agents:       []AgentMeta{},
  }

  f, err := os.Create(userFile)
  if err != nil {
    return nil, fmt.Errorf("create %q: %w", userFile, err)
  }
  defer f.Close()

  enc := toml.NewEncoder(f)
  if err := enc.Encode(cfg); err != nil {
    return nil, fmt.Errorf("encode %q: %w", userFile, err)
  }

  // Now use your existing NewUser to load the in‐memory User
  return NewUser(userID)
}
