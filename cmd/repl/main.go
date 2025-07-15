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

type AgentConfig struct {
  Name      string   `toml:"name"`
  Model     string   `toml:"model"`
  ToolPaths []string `toml:"tool_path"`
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
    if s := strings.TrimSpace(line); s != "" {
      return s
    }
  }
}

func main() {
  tui.PrintLogo()

  var appCfg AppConfig
  if _, err := toml.DecodeFile("configs/settings.toml", &appCfg); err != nil {
    fmt.Fprintf(os.Stderr, "Failed to load settings.toml: %v\n", err)
    os.Exit(1)
  }

  usr, err := user.LoadUser(appCfg.DefaultUser)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Failed to load user %q: %v\n", appCfg.DefaultUser, err)
    os.Exit(1)
  }
  usr.Print()

  client := openai.NewClient()

  agentConfigPath := "configs/jj/agents/myagent.toml"
  fmt.Println("Loading agent from", agentConfigPath)

  var agentCfg AgentConfig
  if _, err := toml.DecodeFile(agentConfigPath, &agentCfg); err != nil {
    fmt.Fprintf(os.Stderr, "Failed to load agent config: %v\n", err)
    os.Exit(1)
  }

  fig := figure.NewColorFigure(agentCfg.Name, "", "cyan", true)
  fig.Print()
  fmt.Println()

  agentInstance, err := agent.NewAgent(&client, agentCfg.Name, agentCfg.Model, agentCfg.ToolPaths)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Failed to init agent: %v\n", err)
    os.Exit(1)
  }

  agentInstance.PrintTools()

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
