package app

import (
	"context"
	"fmt"
	//"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/chzyer/readline"
	openai "github.com/openai/openai-go"


	"github.com/johnjallday/dolphin-tool-calling-agent/internal/agent"
	"github.com/johnjallday/dolphin-tool-calling-agent/internal/registry"
	"github.com/johnjallday/dolphin-tool-calling-agent/internal/tui"
	"github.com/johnjallday/dolphin-tool-calling-agent/internal/user"
)

// REPLApp holds all state for your interactive session.
type REPLApp struct {
	settings        Settings
	usr             *user.User
	client          *openai.Client
	currentAgent    agent.Agent
	currentAgentCfg string
}

// NewREPLApp loads settings.toml, the default user, and creates an OpenAI client.
func NewREPLApp(settingsPath string) (*REPLApp, error) {
    var s Settings
    if _, err := toml.DecodeFile(settingsPath, &s); err != nil {
        return nil, fmt.Errorf("decode settings: %w", err)
    }
    usr, err := user.LoadUser(s.DefaultUser)
    if err != nil {
        return nil, fmt.Errorf("load default user %q: %w", s.DefaultUser, err)
    }
    client := openai.NewClient()
    return &REPLApp{settings: s, usr: usr, client: &client}, nil
}
//
// Run starts the REPL loop.  It prints logo, loads default agent, prints tools, then reads commands.
func (a *REPLApp) Run(ctx context.Context) error {
    tui.PrintLogo()
    a.usr.Print()
    defaultPath, err := a.usr.AgentPath(a.usr.DefaultAgent)
    if err != nil {
        return err
    }
    if err := a.loadAgent(defaultPath); err != nil {
        return err
    }
    tui.PrintTools()

    for {
        line, err := a.readLine()
        if err != nil {
            if err == readline.ErrInterrupt {
                fmt.Println("\nExiting…")
                return nil
            }
            return err
        }
        if done := a.handle(line, ctx); done {
            return nil
        }
    }
}

// loadAgent is the former loadAgent helper.
func (a *REPLApp) loadAgent(path string) error {
	ag, cfg, err := agent.NewAgentFromConfig(a.client, path)
	if err != nil {
		return err
	}
	a.currentAgent = ag
	a.currentAgentCfg = path
	fmt.Println(cfg)
	fmt.Printf("Loaded agent from %s\n", path)
	return nil
}

// readLine encapsulates your readline setup.
func (a *REPLApp) readLine() (string, error) {
	rl, err := readline.New("> ")
	if err != nil {
		return "", err
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
			return "", err
		}
		line = strings.TrimSpace(line)
		if line != "" {
			return line, nil
		}
	}
}

// handle runs one command.  Returns true if we should exit the REPL.
func (a *REPLApp) handle(input string, ctx context.Context) bool {
	parts := strings.Fields(input)
	cmd := strings.ToLower(parts[0])

	switch cmd {
	case "exit", "quit":
		fmt.Println("Bye!")
		return true

	case "help", "tools", "-t":
		tui.PrintLogo()
		tui.PrintTools()

	case "list-agents":
		cfgs, err := agent.ListAgents()
		if err != nil {
			fmt.Println("Error listing agents:", err)
			break
		}
		for _, c := range cfgs {
			fmt.Println(c.Name, c.Model, c.ToolPaths)
		}

	case "create-agent":
		agent.CreateAgent()

	case "load-agent":
		if len(parts) < 2 {
			fmt.Println("Usage: load-agent <path-to-toml>")
			break
		}
		if err := a.loadAgent(parts[1]); err != nil {
			fmt.Println("Failed to load agent:", err)
		}

	case "unload-agent":
		a.currentAgent = nil
		a.currentAgentCfg = ""
		registry.Clear()
		fmt.Println("Agent unloaded.")

	case "current-agent":
		if a.currentAgent != nil {
			fmt.Println("Current agent config:", a.currentAgentCfg)
		} else {
			fmt.Println("No agent loaded.")
		}

	case "tool-pack", "tool-packs":
		tui.PrintToolPacks()

	default:
		// any other input is sent to the agent
		if a.currentAgent == nil {
			fmt.Println("No agent loaded. Use: load-agent <path-to-toml>")
			break
		}
		if err := a.currentAgent.SendMessage(ctx, input); err != nil {
			fmt.Println("Error:", err)
		}
	}

	return false
}
