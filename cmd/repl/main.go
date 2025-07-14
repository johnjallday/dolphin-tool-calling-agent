package main

import (
  "context"
  "fmt"
  "os"
  "strings"

  "github.com/chzyer/readline"
  "github.com/openai/openai-go"
  "github.com/common-nighthawk/go-figure"
  "github.com/BurntSushi/toml"

  "github.com/johnjallday/dolphin-tool-calling-agent/internal/agent"
  "github.com/johnjallday/dolphin-tool-calling-agent/internal/tui"
  "github.com/johnjallday/dolphin-tool-calling-agent/internal/user"
)

type AppConfig struct {
  DefaultUser string `toml:"default_user"`
}

func readQuestion() string {
  rl, _ := readline.New("> ")
  defer rl.Close()
  rl.Config.FuncFilterInputRune = func(r rune) (rune, bool) {
    if r == readline.CharCtrlL {
      fmt.Print("\033[H\033[2J")
      rl.Refresh()
      return 0, false
    }
    return r, true
  }
  for {
    line, err := rl.Readline()
    if err == readline.ErrInterrupt {
      fmt.Println("\nExiting...")
      os.Exit(0)
    }
    line = strings.TrimSpace(line)
    if line != "" {
      return line
    }
  }
}

func main() {
  tui.PrintLogo()

  var cfg AppConfig
  if _, err := toml.DecodeFile("./configs/settings.toml", &cfg); err != nil {
    fmt.Fprintf(os.Stderr, "Failed to load settings.toml: %v\n", err)
    os.Exit(1)
  }

  usr, err := user.LoadUser(cfg.DefaultUser)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Failed to load user %q: %v\n", cfg.DefaultUser, err)
    os.Exit(1)
  }
  usr.Print()

  client := openai.NewClient()
  clientPtr := &client

  agentConfigPath := "./configs/jj/agents/myagent.toml"
  fmt.Println("Loading agent from", agentConfigPath)
  fig := figure.NewColorFigure("myagent", "", "cyan", true)
  fig.Print()
  fmt.Println()

  agentInstance, err := agent.NewAgentFromConfig(clientPtr, agentConfigPath)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Failed to load agent: %v\n", err)
    os.Exit(1)
  }

  ctx := context.Background()
  for {
    question := readQuestion()
    parts := strings.Fields(question)
    if len(parts) == 0 {
      continue
    }
    switch strings.ToLower(parts[0]) {
    case "help", "tools":
      tui.PrintLogo()
			agentInstance.PrintTools()


      continue
    case "exit", "quit":
      fmt.Println("Bye!")
      return
    }

    if err := agentInstance.SendMessage(ctx, question); err != nil {
      fmt.Printf("Error: %v\n", err)
    }
  }
}
