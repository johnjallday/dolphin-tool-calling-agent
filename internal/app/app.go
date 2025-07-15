package app

import (
  "bufio"
  "context"
  "fmt"
  "os"
  "path/filepath"
  "strconv"
  "strings"


  "github.com/BurntSushi/toml"
  "github.com/chzyer/readline"
  "github.com/common-nighthawk/go-figure"
  "github.com/openai/openai-go"
  "github.com/johnjallday/dolphin-tool-calling-agent/internal/agent"
  "github.com/johnjallday/dolphin-tool-calling-agent/internal/tui"
  "github.com/johnjallday/dolphin-tool-calling-agent/internal/user"
)

type AppConfig struct {
  DefaultUser string `toml:"default_user"`
}

type App interface {
  Init(configPath string) error
  Run(ctx context.Context) error
}

type DefaultApp struct {
  cfg    AppConfig
  usr    *user.User
  ag     *agent.Agent
  client openai.Client
}

func NewApp() App { return &DefaultApp{} }

func (a *DefaultApp) Init(configPath string) error {
  if _, err := toml.DecodeFile(configPath, &a.cfg); err != nil {
    return fmt.Errorf("load app config: %w", err)
  }

  a.client = openai.NewClient()

  var err error
  if a.cfg.DefaultUser != "" {
    a.usr, err = user.NewUser(a.cfg.DefaultUser, &a.client)
    if err != nil {
      fmt.Fprintf(os.Stderr, "Default user %q not found.\n", a.cfg.DefaultUser)
      a.usr = a.selectUser()
    }
  } else {
    a.usr = a.selectUser()
  }

  if a.usr.DefaultAgent != nil {
    a.ag = a.usr.DefaultAgent
  } else {
    meta := a.selectAgent()
    a.ag, err = agent.NewAgent(&a.client, meta.Name, meta.Model, meta.ToolPaths)
    if err != nil {
      return fmt.Errorf("init agent: %w", err)
    }
  }
  return nil
}

func (a *DefaultApp) Run(ctx context.Context) error {
  tui.PrintLogo()
  a.usr.Print()
  fig := figure.NewColorFigure(a.ag.Name, "", "cyan", true)
  fig.Print(); fmt.Println()
  a.ag.PrintTools()

  for {
    line := a.readQuestion()
    parts := strings.Fields(line)
    if len(parts) == 0 { continue }
    switch strings.ToLower(parts[0]) {
    case "help", "tools":
      tui.PrintLogo(); a.ag.PrintTools()
    case "exit", "quit":
      fmt.Println("Bye!"); return nil
    default:
      if err := a.ag.SendMessage(ctx, line); err != nil {
        fmt.Printf("Error: %v\n", err)
      }
    }
  }
}

func (a *DefaultApp) selectUser() *user.User {
  userDir := "configs/users"
  files, err := os.ReadDir(userDir)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Failed to list users: %v\n", err)
    os.Exit(1)
  }
  var names []string
  for _, f := range files {
    if f.IsDir() || filepath.Ext(f.Name()) != ".toml" { continue }
    names = append(names, strings.TrimSuffix(f.Name(), ".toml"))
  }

  reader := bufio.NewReader(os.Stdin)
  for {
    fmt.Println("Available users:")
    for i, n := range names {
      fmt.Printf("[%d] %s\n", i+1, n)
    }
    fmt.Print("Select user by number: ")
    input, _ := reader.ReadString('\n')
    if idx, err := strconv.Atoi(strings.TrimSpace(input)); err == nil && idx > 0 && idx <= len(names) {
      usr, err := user.NewUser(names[idx-1], &a.client)
      if err != nil {
        fmt.Printf("Failed to load user: %v\n", err)
        continue
      }
      return usr
    }
    fmt.Println("Invalid choice, try again.")
  }
}

func (a *DefaultApp) selectAgent() user.AgentMeta {
  fmt.Println("Available agents:")
  for i, m := range a.usr.Agents {
    fmt.Printf("[%d] %s\n", i+1, m.Name)
  }
  reader := bufio.NewReader(os.Stdin)
  for {
    fmt.Print("Select agent by number: ")
    input, _ := reader.ReadString('\n')
    if n, err := strconv.Atoi(strings.TrimSpace(input)); err == nil && n > 0 && n <= len(a.usr.Agents) {
      return a.usr.Agents[n-1]
    }
    fmt.Println("Invalid choice, try again.")
  }
}

func (a *DefaultApp) readQuestion() string {
  rl, _ := readline.New("> ")
  defer rl.Close()
  rl.Config.FuncFilterInputRune = func(r rune) (rune, bool) {
    if r == readline.CharCtrlL {
      fmt.Print("\033[H\033[2J"); rl.Refresh()
      return 0, false
    }
    return r, true
  }
  for {
    line, err := rl.Readline()
    if err == readline.ErrInterrupt {
      fmt.Println("\nExiting..."); os.Exit(0)
    }
    if s := strings.TrimSpace(line); s != "" {
      return s
    }
  }
}
