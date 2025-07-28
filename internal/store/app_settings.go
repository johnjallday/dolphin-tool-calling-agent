package store

import (
  "fmt"
  "os"
  "path/filepath"

  "github.com/BurntSushi/toml"
  "github.com/johnjallday/dolphin-tool-calling-agent/pkg/tools"
)

const (
  DefaultConfigDir     = "configs"
  SettingsFileName     = "app_setting.toml"
  ToolpacksFileName    = "toolpacks.toml"
)

// AppSettings mirrors configs/app_setting.toml
type AppSettings struct {
  DefaultUser string `toml:"default_user"`
}

// toolpacksConfig mirrors the [[toolpacks]] table in configs/toolpacks.toml
type toolpacksConfig struct {
  Toolpacks []tools.ToolPackage `toml:"toolpack"`
}

// EnsureConfigDir makes sure configs/ exists and that both
// app_setting.toml and toolpacks.toml exist, initializing them if missing.
func EnsureConfigDir() error {
  // 1) make configs folder
  if err := os.MkdirAll(DefaultConfigDir, 0755); err != nil {
    return fmt.Errorf("mkdir %q: %w", DefaultConfigDir, err)
  }

  // 2) ensure app_setting.toml
  if err := ensureFileWithDefault(
    filepath.Join(DefaultConfigDir, SettingsFileName),
    `default_user = ""`+"\n",
  ); err != nil {
    return err
  }

  // 3) ensure toolpacks.toml
  if err := ensureFileWithDefault(
    filepath.Join(DefaultConfigDir, ToolpacksFileName),
    `# List your remote tool-packages here (will unmarshal into []tools.ToolPackage)
# [[toolpacks]]
# name = "reaper_project_manager"
# version = "0.1.0"
# link = "https://example.com/rpm.so"
# description = "Manage Reaper projects in your DAW"
` ,
  ); err != nil {
    return err
  }

  return nil
}

// ensureFileWithDefault creates the file at path if it does not exist,
// writing defaultContents into it.
func ensureFileWithDefault(path, defaultContents string) error {
  if _, err := os.Stat(path); err != nil {
    if os.IsNotExist(err) {
      f, err := os.Create(path)
      if err != nil {
        return fmt.Errorf("create %q: %w", path, err)
      }
      defer f.Close()
      if _, err := f.WriteString(defaultContents); err != nil {
        return fmt.Errorf("write default to %q: %w", path, err)
      }
    } else {
      return fmt.Errorf("stat %q: %w", path, err)
    }
  }
  return nil
}

// LoadAppSettings reads and decodes configs/app_setting.toml
func LoadAppSettings() (*AppSettings, error) {
  path := filepath.Join(DefaultConfigDir, SettingsFileName)
  var s AppSettings
  if _, err := toml.DecodeFile(path, &s); err != nil {
    return nil, fmt.Errorf("decode %s: %w", path, err)
  }
  return &s, nil
}

// SetDefaultUser updates default_user in configs/app_setting.toml
func SetDefaultUser(userName string) error {
  s, err := LoadAppSettings()
  if err != nil {
    return err
  }
  s.DefaultUser = userName
  return saveToml(filepath.Join(DefaultConfigDir, SettingsFileName), s)
}

// LoadRemoteToolpacks reads and decodes configs/toolpacks.toml
// returning the slice of ToolPackage declared there.
func LoadRemoteToolpacks() ([]tools.ToolPackage, error) {
  path := filepath.Join(DefaultConfigDir, ToolpacksFileName)
  var cfg toolpacksConfig
  if _, err := toml.DecodeFile(path, &cfg); err != nil {
    return nil, fmt.Errorf("decode %s: %w", path, err)
  }
  return cfg.Toolpacks, nil
}

// SaveRemoteToolpacks overwrites configs/toolpacks.toml with the given list.
// Useful if you fetch from a central registry and want to snapshot locally.
func SaveRemoteToolpacks(packs []tools.ToolPackage) error {
  wrapper := toolpacksConfig{Toolpacks: packs}
 	return saveToml(filepath.Join(DefaultConfigDir, ToolpacksFileName), wrapper)
}

// saveToml is a small helper to encode a struct to TOML in path.
func saveToml(path string, v interface{}) error {
  f, err := os.Create(path)
  if err != nil {
    return fmt.Errorf("create %s: %w", path, err)
  }
  defer f.Close()

  enc := toml.NewEncoder(f)
  if err := enc.Encode(v); err != nil {
    return fmt.Errorf("encode %s: %w", path, err)
  }
  return nil
}
