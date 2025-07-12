package main

import (

  "context"
  "fmt"
  "os"
  "strings"
  //"path/filepath"

  "github.com/chzyer/readline"
  "github.com/openai/openai-go"
  "github.com/BurntSushi/toml"
	"github.com/common-nighthawk/go-figure"

  "github.com/johnjallday/dolphin-tool-calling-agent/internal/agent"
  "github.com/johnjallday/dolphin-tool-calling-agent/internal/registry"
  "github.com/johnjallday/dolphin-tool-calling-agent/internal/tui"
  "github.com/johnjallday/dolphin-tool-calling-agent/internal/user"
)

type AppConfig struct {
  DefaultUser string `toml:"default_user"`
}

func readQuestion() string {
  rl, err := readline.New("> ")
  if err != nil {
    panic(err)
  }
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
    if err != nil {
      if err == readline.ErrInterrupt {
        fmt.Println("\nExiting...")
        os.Exit(0)
      }
      fmt.Println()
      return ""
    }
    line = strings.TrimSpace(line)
    if line != "" {
      return line
    }
  }
}

func loadAgent(client *openai.Client, path string) (agent.Agent, error) {
  a, err := agent.NewAgentFromConfig(client, path)
  if err != nil {
    return nil, err
  }
  fmt.Printf("Loaded agent from %s\n", path)
  return a, nil
}

func main() {
  tui.PrintLogo()

	//Initial Loading and Setting Up the App 
  var cfg AppConfig
  if _, err := toml.DecodeFile("./configs/settings.toml", &cfg); err != nil {
    fmt.Fprintf(os.Stderr, "Failed to load settings.toml: %v\n", err)
    os.Exit(1)
  }

  // Load User
  usr, err := user.LoadUser(cfg.DefaultUser)

  if err != nil {
    fmt.Fprintf(os.Stderr, "Failed to load user config for %q: %v\n", cfg.DefaultUser, err)
    os.Exit(1)
  }

	usr.Print()

	defaultAgentPath, err := usr.AgentPath(usr.DefaultAgent)
	if err != nil {
    fmt.Fprintf(os.Stderr, "%v\n", err)
    os.Exit(1)
  }

  client := openai.NewClient()
  clientPtr := &client

	fmt.Println()
	fmt.Println()
	fmt.Println("Loading:", usr.DefaultAgent)
	fig := figure.NewColorFigure(usr.DefaultAgent, "", "cyan", true)
	fig.Print()
	fmt.Println()
	fmt.Println()

  var agentInstance agent.Agent
  var agentConfigPath string


  agentInstance, err = loadAgent(clientPtr, defaultAgentPath)

  if err != nil {
    fmt.Fprintf(os.Stderr, "Error loading agent (%s): %v\n", defaultAgentPath, err)
  }

  tui.PrintTools()




	//REPL
  ctx := context.Background()

  for {
    question := readQuestion()
    parts := strings.Fields(question)
    if len(parts) == 0 {
      continue
    }

    switch strings.ToLower(parts[0]) {
    case "help", "tools", "-t":
      tui.PrintLogo()
      tui.PrintTools()
      continue
    case "exit", "quit":
      fmt.Print("Bye!\r\n")
      return
    case "list-agents", "list-agent", "list agent", "list agents":
      configs, err := agent.ListAgents("jj")
      if err != nil {
        fmt.Println("Error listing agents: Check your config folder")
        continue
      }
      for _, cfg := range configs {
        fmt.Println(cfg.Name)
        fmt.Println(cfg.Model)
        fmt.Println(cfg.ToolPaths)
      }
      continue
    case "create agent", "create-agent":
      fmt.Println("create agent")
      agent.CreateAgent()
      continue
    case "load-agent":
      if len(parts) < 2 {
        fmt.Println("Usage: load-agent <path-to-toml>")
        continue
      }
      newPath := parts[1]
      newAgent, err := loadAgent(clientPtr, newPath)
      if err != nil {
        fmt.Printf("Failed to load agent: %v\n", err)
        continue
      }
      agentInstance = newAgent
      agentConfigPath = newPath
      continue
    case "unload-agent":
      agentInstance = nil
      agentConfigPath = ""
      registry.Clear()
      fmt.Println("Agent unloaded.")
      continue
    case "current-agent":
      if agentInstance != nil && agentConfigPath != "" {
        fmt.Printf("Current agent loaded from: %s\n", agentConfigPath)
      } else {
        fmt.Println("No agent loaded.")
      }
      continue
    case "tool-pack", "tool-packs":
      tui.PrintToolPacks()
      continue
    }

    if agentInstance == nil {
      fmt.Println("No agent loaded. Use: load-agent <path-to-toml>")
      continue
    }

    if err := agentInstance.SendMessage(ctx, question); err != nil {
      fmt.Printf("Error: %v\n", err)
    }
  }
}
