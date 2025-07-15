package main

import (
  "bufio"
  "context"
  "fmt"
  "log"
  "os"
  //"path/filepath"
  "strconv"
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

func main() {
  tui.PrintLogo()

  var appCfg AppConfig
  if _, err := toml.DecodeFile("configs/app_setting.toml", &appCfg); err != nil {
    fmt.Fprintf(os.Stderr, "Failed to load app setting: %v\n", err)
    os.Exit(1)
  }

  usr, err := user.NewUser(appCfg.DefaultUser)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Failed to load user %q: %v\n", appCfg.DefaultUser, err)
    os.Exit(1)
  }
  usr.Print()

  client := openai.NewClient()
  var agentInstance *agent.Agent

  if usr.DefaultAgent != nil {
    agentInstance = usr.DefaultAgent
  } else {
    meta := selectAgent(usr)
    var agentCfg AgentConfig
    if _, err := toml.DecodeFile(meta.Path, &agentCfg); err != nil {
      log.Fatalf("Failed to load agent config: %v", err)
    }
    agentInstance, err = agent.NewAgent(&client, agentCfg.Name, agentCfg.Model, agentCfg.ToolPaths)
    if err != nil {
      log.Fatalf("Failed to init agent: %v", err)
    }
  }

  fig := figure.NewColorFigure(agentInstance.Name, "", "cyan", true)
  fig.Print()
  fmt.Println()
  agentInstance.PrintTools()

  ctx := context.Background()
  for {
    line := readQuestion()
    parts := strings.Fields(line)
    if len(parts) == 0 {
      continue
    }
    switch strings.ToLower(parts[0]) {
    case "help", "tools":
      tui.PrintLogo()
      agentInstance.PrintTools()
    case "exit", "quit":
      fmt.Println("Bye!")
      return
    default:
      if err := agentInstance.SendMessage(ctx, line); err != nil {
        fmt.Printf("Error: %v\n", err)
      }
    }
  }
}


func selectAgent(usr *user.User) user.AgentMeta {
  fmt.Println("Available agents:")
  for i, meta := range usr.Agents {
    fmt.Printf("[%d] %s\n", i+1, meta.Name)
  }
  fmt.Print("Select agent by number: ")
  reader := bufio.NewReader(os.Stdin)
  for {
    input, _ := reader.ReadString('\n')
    n, err := strconv.Atoi(strings.TrimSpace(input))
    if err == nil && n > 0 && n <= len(usr.Agents) {
      return usr.Agents[n-1]
    }
    fmt.Print("Invalid choice, try again: ")
  }
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
