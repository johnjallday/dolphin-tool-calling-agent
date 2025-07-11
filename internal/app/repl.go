package app

import (
  "context"
  "fmt"
	"strings"
	"io"
	"os"

  "github.com/BurntSushi/toml"
  "github.com/chzyer/readline"
  openai "github.com/openai/openai-go"

  "github.com/johnjallday/dolphin-tool-calling-agent/internal/agent"
  "github.com/johnjallday/dolphin-tool-calling-agent/internal/registry"
  "github.com/johnjallday/dolphin-tool-calling-agent/internal/tui"
  "github.com/johnjallday/dolphin-tool-calling-agent/internal/user"
)

var _ App = (*REPLApp)(nil)    // compile‐time check


type nopWriteCloser struct{ io.Writer }
func (nopWriteCloser) Close() error { return nil }

type REPLApp struct {
  settingsPath     string
  settings         Settings
  usr              *user.User
  client           *openai.Client
  currentAgent     agent.Agent
  currentAgentCfg  string
  rl *readline.Instance
  in  io.Reader
  out io.Writer
}


func NewREPLAppWithIO(in io.Reader, out io.Writer) *REPLApp {
  // default to real stdin/stdout if nil
  if in == nil {
    in = os.Stdin
  }
  if out == nil {
    out = os.Stdout
  }

  // make them into ReadCloser/WriteCloser
  var rc io.ReadCloser
  if c, ok := in.(io.ReadCloser); ok {
    rc = c
  } else {
    rc = io.NopCloser(in)
  }

  var wc io.WriteCloser
  if c, ok := out.(io.WriteCloser); ok {
    wc = c
  } else {
    wc = nopWriteCloser{out}
  }

  rl, err := readline.NewEx(&readline.Config{
    Prompt: "> ",
    Stdin:  rc,
    Stdout: wc,
  })
  if err != nil {
    panic(err)
  }

  return &REPLApp{
    rl:  rl,
    in:  in,
    out: out,
  }
}



func NewREPLApp() *REPLApp {
  return NewREPLAppWithIO(nil, nil)
}

func (a *REPLApp) Init(settingsPath string, userName string) error {
  a.settingsPath = settingsPath
  // load settings
  if _, err := toml.DecodeFile(settingsPath, &a.settings); err != nil {
    return fmt.Errorf("decode settings: %w", err)
  }
  // load user
  usr, err := user.LoadUser(userName)
  if err != nil {
    return fmt.Errorf("load user %q: %w", userName, err)
  }
  a.usr    = usr
  client   := openai.NewClient()
  a.client = &client
  return nil
}
// readLine can now be:
func (a *REPLApp) readLine() (string, error) {
  line, err := a.rl.Readline()
  if err != nil {
    return "", err
  }
  return strings.TrimSpace(line), nil
}

// dispatch looks at the line, runs the right subcommand, and
// returns true if we should exit the REPL:
func (a *REPLApp) dispatch(line string, ctx context.Context) bool {
  fields := strings.Fields(line)
  if len(fields) == 0 {
    return false
  }
  cmd := fields[0]
  args := fields[1:]

  switch cmd {
  case "load-agent":
    if len(args) != 1 {
      fmt.Fprintln(a.out, "usage: load-agent <path>")
      break
    }
    if err := a.LoadAgent(args[0]); err != nil {
      fmt.Fprintln(a.out, "error:", err)
    }
  case "unload-agent":
    a.UnloadAgent()
    fmt.Fprintln(a.out, "unloaded")
	case "list-agents":
		agents, err := a.ListAgents()
		if err != nil {
			fmt.Fprintln(a.out, "error:", err)
			break
		}
		for _, cfg := range agents {
			// print the name and model (whatever fields AgentConfig actually has)
			fmt.Fprintf(a.out, "%s  (model=%s)\n", cfg.Name, cfg.Model)
		}
  case "create-agent":
    if err := a.CreateAgent(); err != nil {
      fmt.Fprintln(a.out, "error:", err)
    }
  case "exit", "quit":
    fmt.Fprintln(a.out, "Exiting.")
    return true
  default:
    fmt.Fprintln(a.out, "unknown command:", cmd)
  }
  return false
}


func (a *REPLApp) Run(ctx context.Context) error {
  tui.PrintLogo()
  a.usr.Print()

  // load the user’s default agent if any
  defaultPath, err := a.usr.AgentPath(a.usr.DefaultAgent)
  if err == nil {
    _ = a.LoadAgent(defaultPath)
    tui.PrintTools()
  }

  for {
    line, err := a.readLine()
    if err != nil { /* … interrupt handling … */ }
    if a.dispatch(line, ctx) {
      return nil
    }
  }
}

func (a *REPLApp) Shutdown() error {
  // nothing right now, but maybe close files, etc.
  return nil
}

func (a *REPLApp) LoadAgent(path string) error {
  ag, cfg, err := agent.NewAgentFromConfig(a.client, path)
  if err != nil {
    return err
  }
  a.currentAgent    = ag
  a.currentAgentCfg = path
  fmt.Printf("Loaded agent %s (model=%s)\n", cfg.Name, cfg.Model)
  return nil
}

func (a *REPLApp) UnloadAgent() {
  a.currentAgent    = nil
  a.currentAgentCfg = ""
  registry.Clear()
}

func (a *REPLApp) CurrentAgent() agent.Agent {
  return a.currentAgent
}

func (a *REPLApp) CurrentAgentConfig() string {
  return a.currentAgentCfg
}

func (a *REPLApp) ListAgents() ([]agent.AgentConfig, error) {
  return agent.ListAgents()
}

func (a *REPLApp) CreateAgent() error {

  return nil
}

// … your existing readLine() and dispatch() methods go here, plus helper sub‐commands …
