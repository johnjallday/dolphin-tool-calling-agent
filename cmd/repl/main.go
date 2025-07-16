package main

import (
  "context"
  "fmt"
  "os"
  "strings"

  "github.com/chzyer/readline"
  "github.com/johnjallday/dolphin-tool-calling-agent/internal/app"
  "github.com/johnjallday/dolphin-tool-calling-agent/internal/tui"
)

func main() {
  ctx := context.Background()
  a := app.NewApp()
  if err := a.Init(); err != nil {
    fmt.Fprintln(os.Stderr, "Init error:", err)
    os.Exit(1)
  }

  fmt.Println("Loading Default User")
  fmt.Println(a.Users())
  tui.PrintLogo()

  rl, err := readline.New("> ")
  if err != nil {
    fmt.Fprintln(os.Stderr, "readline error:", err)
    os.Exit(1)
  }
  defer rl.Close()

  for {
    line, err := rl.Readline()
    if err != nil {
      break
    }
    cmd := strings.TrimSpace(line)
    switch strings.ToLower(cmd) {
    case "current agent", "current_agent":
      a.CurrentAgent().Print()
    case "current user", "current_user":
      a.CurrentUser().Print()
    case "users":
      fmt.Println(a.Users())
    case "help":
      tui.PrintLogo()
    case "exit", "quit":
      fmt.Println("Bye!")
      return
    default:
      if cmd != "" {
        if err := a.SendMessage(ctx, cmd); err != nil {
          fmt.Fprintln(os.Stderr, "send error:", err)
        }
      }
    }
  }
}
