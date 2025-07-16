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

  tui.PrintLogo()
  rl, err := readline.New("> ")
  if err != nil {
    fmt.Fprintln(os.Stderr, "readline error:", err)
    os.Exit(1)
  }
  defer rl.Close()

  a.CurrentUser().Print()
  for {
    line, err := rl.Readline()
    if err != nil {
      break
    }
    raw := strings.TrimSpace(line)
    lower := strings.ToLower(raw)
    fields := strings.Fields(raw)


    switch {

		case len(fields) >= 2 && strings.ToLower(fields[0]) == "unload" && strings.ToLower(fields[1]) == "agent":
			if err := a.UnloadAgent(); err != nil {
				fmt.Fprintln(os.Stderr, "unload agent error:", err)
			} else {
				fmt.Println("Agent unloaded")
			}

		case len(fields) >= 2 && strings.ToLower(fields[0]) == "unload" && strings.ToLower(fields[1]) == "user":
			if err := a.UnloadUser(); err != nil {
				fmt.Fprintln(os.Stderr, "unload user error:", err)
			} else {
				fmt.Println("User and agent unloaded")
			}
    case len(fields) >= 2 && strings.ToLower(fields[0]) == "load" && strings.ToLower(fields[1]) == "user":
      if len(fields) < 3 {
        fmt.Println("Usage: load user <username>")
        continue
      }
      username := fields[2]
      if err := a.LoadUser(username); err != nil {
        fmt.Fprintln(os.Stderr, "load user error:", err)
      } else {
        fmt.Println("Loaded user:", username)
				a.CurrentAgent().Print()
      }

		// inside your main() REPL switch, above the "current agent" case
		case len(fields) >= 2 && strings.ToLower(fields[0]) == "load" && strings.ToLower(fields[1]) == "agent":
			if len(fields) < 3 {
				fmt.Println("Usage: load agent <agentName>")
				continue
			}
			agentName := fields[2]
			if err := a.LoadAgent(agentName); err != nil {
				fmt.Fprintln(os.Stderr, "load agent error:", err)
			} else {
				fmt.Println("Loaded agent:", agentName)
				a.CurrentAgent().Print()
			}
    case lower == "current agent" || lower == "current_agent":
      a.CurrentAgent().Print()
    case lower == "current user" || lower == "current_user" || lower == "agents":
      a.CurrentUser().Print()
    case lower == "users":
      fmt.Println(a.Users())
    case lower == "help":
      tui.PrintLogo()
    case lower == "exit" || lower == "quit":
      fmt.Println("Bye!")
      return
    default:
      if raw != "" {
        if err := a.SendMessage(ctx, raw); err != nil {
          fmt.Fprintln(os.Stderr, "send error:", err)
        }
      }
    }
  }
}
