package main

import (
  "fmt"
  "os"

  "github.com/johnjallday/dolphin-tool-calling-agent/internal/app"
  "github.com/johnjallday/dolphin-tool-calling-agent/internal/bubbletui"
)

func main() {
  // 1) initialize your core application
  core := app.NewApp()
  if err := core.Init(); err != nil {
    fmt.Fprintf(os.Stderr, "init error: %v\n", err)
    os.Exit(1)
  }

  // 2) launch the Bubble Tea TUI
  if err := bubbletui.RunStatusTUI(core); err != nil {
    fmt.Fprintf(os.Stderr, "fatal: %v\n", err)
    os.Exit(1)
  }
}
