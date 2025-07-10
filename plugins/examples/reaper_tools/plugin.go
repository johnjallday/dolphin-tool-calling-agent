package main

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strings"

    "github.com/johnjallday/dolphin-tool-calling-agent/pkg/tools"
    "github.com/jedib0t/go-pretty/v6/table"
    "github.com/BurntSushi/toml"
    "github.com/openai/openai-go"
)

type ReaperConfig struct {
    DefaultTemplate string `toml:"default_template"`
    ScriptPath      string `toml:"script_path"`
}

var reaperConfig ReaperConfig

func registerConfig() error {
    // switch to the reaper folder
    configPath := "./configs/user/tools/reaper/settings.toml"
    absPath, err := filepath.Abs(configPath)
    if err != nil {
        return fmt.Errorf("could not resolve %s: %w", configPath, err)
    }
    if _, err := os.Stat(absPath); os.IsNotExist(err) {
        // default points at custom_scripts
        defaultConfig := ReaperConfig{
            DefaultTemplate: "./configs/user/tools/reaper/Default.RPP",
            ScriptPath:      "./configs/user/tools/reaper/custom_scripts",
        }
        f, ferr := os.Create(absPath)
        if ferr != nil {
            return fmt.Errorf("could not create %s: %w", absPath, ferr)
        }
        defer f.Close()
        if err := toml.NewEncoder(f).Encode(defaultConfig); err != nil {
            return fmt.Errorf("could not encode default config: %w", err)
        }
        reaperConfig = defaultConfig
        fmt.Println("Created default reaper settings at", absPath)
        return nil
    }
    if _, err := toml.DecodeFile(absPath, &reaperConfig); err != nil {
        return fmt.Errorf("failed to decode %s: %w", absPath, err)
    }
    return nil
}

func launchTool(scriptName string) error {
    if err := registerConfig(); err != nil {
        return err
    }
    if reaperConfig.ScriptPath == "" {
        return fmt.Errorf("script_path missing in config")
    }
    scriptPath := filepath.Join(reaperConfig.ScriptPath, scriptName+".lua")
    fmt.Println("Launching REAPER with:", scriptPath)
    cmd := exec.Command("open", "-a", "Reaper", scriptPath)
    return cmd.Run()
}

func LoadCustomScripts() []tools.ToolSpec {
    if err := registerConfig(); err != nil {
        panic(err)
    }
    entries, err := os.ReadDir(reaperConfig.ScriptPath)
    if err != nil {
        panic(fmt.Errorf("failed to read %s: %w", reaperConfig.ScriptPath, err))
    }

    var specs []tools.ToolSpec
    for _, e := range entries {
        if e.IsDir() || !strings.HasSuffix(e.Name(), ".lua") {
            continue
        }
        name := strings.TrimSuffix(e.Name(), ".lua")
        specs = append(specs, tools.ToolSpec{
            Name:        name,
            Description: fmt.Sprintf("Run custom Reaper script %q", name),
            Parameters:  openai.FunctionParameters{},
            Exec: func(name string) func(map[string]interface{}) (string, error) {
                return func(_ map[string]interface{}) (string, error) {
                    if err := launchTool(name); err != nil {
                        return "", err
                    }
                    return fmt.Sprintf("Launched %s", name), nil
                }
            }(name),
        })
    }
    return specs
}
