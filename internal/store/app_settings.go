package store

import (
  "fmt"
  "os"
  "path/filepath"

  "github.com/BurntSushi/toml"
)

const DefaultConfigDir = "configs"
const SettingsFileName = "app_setting.toml"

// AppSettings mirrors the toml in configs/app_setting.toml
type AppSettings struct {
  DefaultUser string `toml:"default_user"`
}

// EnsureAppSettingsDir makes sure configs/ exists, and that
// app_setting.toml exists (initializing it if missing).
func EnsureAppSettingsDir() error {
  // 1) ensure the “configs” folder exists
  if err := os.MkdirAll(DefaultConfigDir, 0755); err != nil {
    return fmt.Errorf("mkdir %q: %w", DefaultConfigDir, err)
  }

  // 2) ensure the file exists
  path := filepath.Join(DefaultConfigDir, SettingsFileName)
  if _, err := os.Stat(path); err != nil {
    if os.IsNotExist(err) {
      f, err := os.Create(path)
      if err != nil {
        return fmt.Errorf("create %q: %w", path, err)
      }
      defer f.Close()
      // write an empty default_user
      if _, err := f.WriteString("default_user = \"\"\n"); err != nil {
        return fmt.Errorf("write %q: %w", path, err)
      }
    } else {
      return fmt.Errorf("stat %q: %w", path, err)
    }
  }
  return nil
}

// loadAppSettings reads and decodes configs/app_setting.toml
func LoadAppSettings() (*AppSettings, error) {
  path := filepath.Join(DefaultConfigDir, SettingsFileName)
  var s AppSettings
  if _, err := toml.DecodeFile(path, &s); err != nil {
    return nil, fmt.Errorf("decode %s: %w", path, err)
  }
  return &s, nil
}

// saveAppSettings truncates & writes configs/app_setting.toml
func saveAppSettings(s *AppSettings) error {
  path := filepath.Join(DefaultConfigDir, SettingsFileName)
  f, err := os.Create(path)
  if err != nil {
    return fmt.Errorf("create %s: %w", path, err)
  }
  defer f.Close()
  if err := toml.NewEncoder(f).Encode(s); err != nil {
    return fmt.Errorf("encode %s: %w", path, err)
  }
  return nil
}

// SetDefaultUser updates default_user in configs/app_setting.toml.
func SetDefaultUser(userName string) error {
  s, err := LoadAppSettings()
  if err != nil {
    return err
  }
  s.DefaultUser = userName
  return saveAppSettings(s)
}
