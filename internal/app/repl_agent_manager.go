package app

import (
    //"bufio"
    "context"
    "fmt"
    //"os"
  	"github.com/johnjallday/dolphin-tool-calling-agent/internal/user"
  	"github.com/johnjallday/dolphin-tool-calling-agent/internal/agent"
)

type REPLAgentManager struct{}

func (m *REPLAgentManager) Select(ctx context.Context, u *user.User) (string, error) {
    // if they've already got a default agent, and it's in the list, use it
    for _, a := range u.Agents {
        if a == u.DefaultAgent {
            return a, nil
        }
    }

    // otherwise prompt them to pick one
    fmt.Println("Pick an agent:")
    for i, a := range u.Agents {
        fmt.Printf("  %d) %s\n", i+1, a)
    }
    fmt.Print("> ")
    var choice int
    if _, err := fmt.Scanln(&choice); err != nil {
        return "", err
    }
    if choice < 1 || choice > len(u.Agents) {
        return "", fmt.Errorf("invalid choice")
    }
    return u.Agents[choice-1], nil
}

func (m *REPLAgentManager) Load(ctx context.Context, u *user.User, agentName string) error {
    path, err := u.AgentPath(agentName)
    if err != nil {
        return err
    }
		fmt.Println(path)
		fmt.Println("REPLAgentManager")

		cfg, err := agent.LoadConfig(path)
    // TODO: read the agent's TOML, construct your tool/agent client, etc.
    // For example:
    //   cfg, err := agent.LoadConfig(path)
    //   a.client = agent.New(cfg)
    return nil
}
