package main

import (
  "context"
  "fmt"
  "os"

  "github.com/johnjallday/dolphin-tool-calling-agent/internal/app"
)

func main() {
  a := app.NewApp()
  if err := a.Init("configs/app_setting.toml"); err != nil {
    fmt.Fprintf(os.Stderr, "Failed to initialize app: %v\n", err)
    os.Exit(1)
  }
  if err := a.Run(context.Background()); err != nil {
    fmt.Fprintf(os.Stderr, "Error: %v\n", err)
    os.Exit(1)
  }
}
