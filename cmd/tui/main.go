package main

import (
  "context"
  "fmt"
  "os"
  "os/signal"
  "syscall"

  "github.com/peterh/liner"

  "github.com/johnjallday/dolphin-tool-calling-agent/internal/app"
  "github.com/johnjallday/dolphin-tool-calling-agent/internal/tui"
)

func main() {
  // 1) catch SIGINT
  sigCh := make(chan os.Signal, 1)
  signal.Notify(sigCh, syscall.SIGINT)
  go func() {
    <-sigCh
    fmt.Fprintln(os.Stdout, "\nreceived SIGINT, exiting")
    os.Exit(1)
  }()

  // 2) core application
  application := app.NewApp()
  if err := application.Init(); err != nil {
    fmt.Fprintf(os.Stderr, "app.Init error: %v\n", err)
    os.Exit(1)
  }

  // 3) set up liner
  rl := liner.NewLiner()
  defer rl.Close()
  rl.SetCtrlCAborts(true)

  // 4) compose your TUI facade
  t := &tui.TUIApp{
    Ctx: context.Background(),
    App: application,
    In:  os.Stdin,
    Out: os.Stdout,
    Err: os.Stderr,
    Rl:  rl,
  }

  // ──────────────── NEW ───────────────────
  // 5) bootstrap users if none exist
  if err := tui.InitCmd(t, nil); err != nil {
    fmt.Fprintf(os.Stderr, "initialization error: %v\n", err)
    os.Exit(1)
  }
  // ────────────────────────────────────────

  // 6) initial screen draw
  if err := t.Refresh(); err != nil {
    fmt.Fprintln(os.Stderr, "refresh error:", err)
  }

  // 7) load up your REPL
  helpKeys, commands := buildCommands()
  t.RunInteractiveShell(helpKeys, commands)
}

func buildCommands() ([]string, map[string]tui.CmdFunc) {
  helpKeys := []string{
    "user", "users", "agent", "agents", "tools",
    "create-agent", "load-user", "load-agent", "unload-user", "edit-agent", "unload-agent",
    "switch-user", "switch-agent",
    "help", "clear", "exit", "quit",
  }

  cmds := map[string]tui.CmdFunc{
    "user":         tui.UserCmd,
    "users":        tui.UsersCmd,
    "agent":        tui.AgentCmd,
    "agents":       tui.AgentsCmd,
    "tools":        tui.ToolsCmd,
    "create-agent": tui.CreateAgentCmd,
		"edit-agent":		tui.EditAgentCmd,
    "load-agent":   tui.LoadAgentCmd,
    "load-user":   	tui.LoadUserCmd,
    "unload-user":  tui.UnloadUserCmd,
    "unload-agent": tui.UnloadAgentCmd,
    "switch-user":  tui.SwitchUserCmd,
    "switch-agent": tui.SwitchAgentCmd,

    "help": func(t *tui.TUIApp, _ []string) error {
			fmt.Fprintln(t.Out, "Try typing one of the available commands to get/execute the information you need.")
			fmt.Fprintln(t.Out, "If you have your agent set up then try prompting anything else.")
			fmt.Fprintln(t.Out, "")
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

  return helpKeys, cmds
}
