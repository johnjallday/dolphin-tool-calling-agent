package tui

import (
    "context"
    "fmt"
    "io"
		// "os"
    //"io/fs"
    // "path/filepath"
    "reflect"
    "strings"

    "github.com/fatih/color"
    "github.com/johnjallday/dolphin-tool-calling-agent/internal/app"
    "github.com/peterh/liner"
)

type TUIApp struct {
    Ctx context.Context
    App app.App
		In  io.Reader
    Out io.Writer
    Err io.Writer
    Rl  *liner.State
}

type CmdFunc func(t *TUIApp, args []string) error

func (t *TUIApp) RunInteractiveShell(
  helpKeys []string,
  commands map[string]CmdFunc,
) {
  for {
    line, err := t.Rl.Prompt("> ")
    if err != nil {
      switch err {
      case liner.ErrPromptAborted:
        fmt.Fprintln(t.Out, "\nexit on Ctrl-C")
        return
      case io.EOF:
        fmt.Fprintln(t.Out, "\nGoodbye!")
        return
      default:
        fmt.Fprintln(t.Err, "prompt error:", err)
        return
      }
    }

    line = strings.TrimSpace(line)
    if line == "" {
      if err := t.Refresh(); err != nil {
        fmt.Fprintln(t.Err, "refresh error:", err)
      }
      continue
    }

    parts := strings.Fields(line)
    rawCmd, args := parts[0], parts[1:]

    // Normalize to lower-case so "User", "USER", "uSeR" all map to "user"
    cmd := strings.ToLower(rawCmd)

    if fn, ok := commands[cmd]; ok {
      if err := fn(t, args); err != nil {
        fmt.Fprintln(t.Err, "ERROR:", err)
      }
    } else {
      // fallback → send to LLM/chat
      if err := t.App.SendMessage(t.Ctx, line); err != nil {
        fmt.Fprintln(t.Err, "ERROR:", err)
      }
    		}

    t.Rl.AppendHistory(line)
  }
}

// clearScreen emits ANSI codes to clear the terminal + move cursor home.
func (t *TUIApp) clearScreen() {
    fmt.Fprint(t.Out, "\x1b[2J\x1b[H")
}

// Refresh clears the screen, prints the logo, then status.
func (t *TUIApp) Refresh() error {
    t.clearScreen()
    t.PrintLogo()
		err := t.StatusCmd()
		if err != nil{
      fmt.Printf("error loading status bar")
		}

    return nil
}

// PrintLogo prints the dolphin logo banner.
func (t *TUIApp) PrintLogo(){
    logo := `
        🐬
`
    fmt.Print("\033[36m" + logo + "\033[0m\n")
    fmt.Println("Dolphin Tool Calling Agent Client\n")
}

// PrintToolPacks walks “./plugins” looking for .so files.


// PrintHelp uses reflection to list all methods on app.App.
func PrintHelp() {
    iface := reflect.TypeOf((*app.App)(nil)).Elem()
    fmt.Println("Available commands:")
    for i := 0; i < iface.NumMethod(); i++ {
        m := iface.Method(i)
        if m.Name == "Init" || m.Name == "SendMessage" {
            continue
        }
        sig := m.Type
        var params []string
        for j := 1; j < sig.NumIn(); j++ {
            params = append(params, sig.In(j).String())
        }
        fmt.Printf("  %s(%s)\n", m.Name, strings.Join(params, ", "))
    }
    fmt.Println()
}


// PrintStatus is used by Refresh() to show user+agent at top.
func PrintStatus(t *TUIApp) error {
    u, a := t.App.User(), t.App.Agent()
    cLabel := color.New(color.FgCyan, color.Bold)
    cValue := color.New(color.FgWhite)

    cLabel.Fprint(t.Out, "Current User:  ")
    cValue.Fprintln(t.Out, u.Name)
    cLabel.Fprint(t.Out, "Current Agent: ")
    cValue.Fprintln(t.Out, a.Name)
    return nil
}
