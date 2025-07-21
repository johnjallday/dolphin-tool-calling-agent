package main

import (
  "context"
  "fmt"
  "io"
  "os"
  "os/signal"
  "strings"
  "syscall"

  "github.com/peterh/liner"
  "github.com/johnjallday/dolphin-tool-calling-agent/internal/app"
  "github.com/johnjallday/dolphin-tool-calling-agent/internal/tui"
)

type CmdFunc func(t *tui.TUIApp, args []string) error

func main() {
  // 1) SIGINT handler
  sigCh := make(chan os.Signal, 1)
  signal.Notify(sigCh, syscall.SIGINT)
  go func() {
    <-sigCh
    fmt.Fprintln(os.Stdout, "\nreceived SIGINT, exiting")
    os.Exit(1)
  }()

  // 2) init your application
  application := app.NewApp()
  if err := application.Init(); err != nil {
    fmt.Fprintf(os.Stderr, "init error: %v\n", err)
    os.Exit(1)
  }

  // 3) init liner
  rl := liner.NewLiner()
  defer rl.Close()
  rl.SetCtrlCAborts(true)

  // 4) build your TUIApp
  t := &tui.TUIApp{
    Ctx: context.Background(),
    App: application,
    Out: os.Stdout,
    Err: os.Stderr,
    Rl:  rl,
  }

  // 5) command registry (include "clear")
  helpKeys := []string{
    "user", "users", "agent", "agents", "tools",
    "unload-user", "unload-agent",
    "help", "clear", "exit", "quit",
  }
  commands := map[string]CmdFunc{
    "user":         func(t *tui.TUIApp, _ []string) error { return tui.CurrentUser(*t) },
    "users":        func(t *tui.TUIApp, _ []string) error { return tui.Users(*t) },
    "agent":        func(t *tui.TUIApp, _ []string) error { return tui.Agent(*t) },
    "agents":       func(t *tui.TUIApp, _ []string) error { return tui.Agents(*t) },
    "tools":        func(t *tui.TUIApp, _ []string) error { return tui.Tools(*t) },
    "unload-user":  func(t *tui.TUIApp, _ []string) error { 
			if err := tui.UnloadUser(*t); err != nil {
				return err
			}
			return t.Refresh()
		},
    "unload-agent": func(t *tui.TUIApp, _ []string) error { 
			if err := tui.UnloadAgent(*t); err != nil {
				return err
			}
			return t.Refresh()
		},
    "help": func(t *tui.TUIApp, _ []string) error {
      fmt.Fprintln(t.Out, "Available commands:")
      for _, k := range helpKeys {
        fmt.Fprintf(t.Out, "  %s\n", k)
      }
      return nil
    },
    "clear": func(t *tui.TUIApp, _ []string) error {
      return t.Refresh()
    },
    "exit": func(t *tui.TUIApp, _ []string) error {
      os.Exit(0)
      return nil
    },
    "quit": func(t *tui.TUIApp, _ []string) error {
      os.Exit(0)
      return nil
    },
  }

  // initial draw
  if err := t.Refresh(); err != nil {
    fmt.Fprintln(os.Stderr, "refresh error:", err)
  }

  // 6) enter the REPL
  repl(rl, t, commands)
}

func repl(rl *liner.State, t *tui.TUIApp, commands map[string]CmdFunc) {
  for {
    line, err := rl.Prompt("> ")
    if err != nil {
      switch err {
      case liner.ErrPromptAborted:
        fmt.Fprintln(os.Stdout, "\nexit on Ctrl-C")
        os.Exit(0)
      case io.EOF:
        fmt.Println("\nGoodbye!")
        os.Exit(0)
      default:
        fmt.Fprintln(t.Err, "prompt error:", err)
        os.Exit(1)
      }
    }

    line = strings.TrimSpace(line)
    if line == "" {
      continue
    }

    parts := strings.Fields(line)
    cmd, args := parts[0], parts[1:]

    if fn, ok := commands[cmd]; ok {
      // redraw before each built-in
      if err := t.Refresh(); err != nil {
        fmt.Fprintln(t.Err, "refresh error:", err)
      }
      if err := fn(t, args); err != nil {
        fmt.Fprintln(t.Err, "ERROR:", err)
      }
    } else {
      // fallback: OpenAI
      if err := t.App.SendMessage(t.Ctx, line); err != nil {
        fmt.Fprintln(t.Err, "ERROR:", err)
      }
    }

    rl.AppendHistory(line)
  }
}
