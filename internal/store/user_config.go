package store

import (
  "fmt"
  "os"
  "path/filepath"

  "github.com/BurntSushi/toml"
  "github.com/johnjallday/dolphin-tool-calling-agent/internal/user"
)

// UserConfig mirrors the on‐disk structure of a configs/users/<name>.toml
type UserConfig struct {
  Name         string            `toml:"name"`
  DefaultAgent string            `toml:"default_agent"`
  Agents       []user.AgentMeta  `toml:"agents"`
}

// LoadUserConfig reads configs/users/<username>.toml into a UserConfig.
// Returns an error if the file does not exist or is invalid TOML.
func LoadUserConfig(username string) (*UserConfig, error) {
  userDir := filepath.Join(DefaultConfigDir, "users")
  userFile := filepath.Join(userDir, username+".toml")

  var cfg UserConfig
  if _, err := toml.DecodeFile(userFile, &cfg); err != nil {
    return nil, fmt.Errorf("store: decode %s: %w", userFile, err)
  }
  return &cfg, nil
}

// SaveUserConfig writes the given UserConfig back to
// configs/users/<cfg.Name>.toml (creating the dir if needed).
// It does an atomic‐style write via a “.tmp” file + rename.
func SaveUserConfig(cfg *UserConfig) error {
  if cfg == nil {
    return fmt.Errorf("store: cannot save nil UserConfig")
  }

  userDir := filepath.Join(DefaultConfigDir, "users")
  if err := os.MkdirAll(userDir, 0o755); err != nil {
    return fmt.Errorf("store: mkdir %s: %w", userDir, err)
  }

  target := filepath.Join(userDir, cfg.Name+".toml")
  tmpFile := target + ".tmp"

  f, err := os.OpenFile(tmpFile,
    os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644,
  )
  if err != nil {
    return fmt.Errorf("store: open temp file %s: %w", tmpFile, err)
  }
  // ensure we clean up the temp on error
  defer func() {
    f.Close()
    os.Remove(tmpFile)
  }()

  enc := toml.NewEncoder(f)
  if err := enc.Encode(cfg); err != nil {
    return fmt.Errorf("store: encode to %s: %w", tmpFile, err)
  }
  // flush to disk
  if err := f.Sync(); err != nil {
    return fmt.Errorf("store: sync %s: %w", tmpFile, err)
  }
  if err := f.Close(); err != nil {
    return fmt.Errorf("store: close %s: %w", tmpFile, err)
  }

  // atomically replace the old file
  if err := os.Rename(tmpFile, target); err != nil {
    return fmt.Errorf("store: rename %s → %s: %w", tmpFile, target, err)
  }
  return nil
}
