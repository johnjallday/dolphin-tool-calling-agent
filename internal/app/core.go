package app

import (
	"context"

	"github.com/johnjallday/dolphin-tool-calling-agent/internal/agent"
  "github.com/johnjallday/dolphin-tool-calling-agent/internal/user"
	"github.com/johnjallday/dolphin-tool-calling-agent/pkg/tools"
)

type AgentMeta struct {
  Name, Model string
  ToolPaths   []string
}
type ToolInfo struct {
  Name, Description string
}

type App interface {
  Init() error
  Users() []string
	User() *user.User
	Agent() *agent.Agent
	Agents() []user.AgentMeta
	SendMessage(ctx context.Context, text string) (reply string, err error)
	CreateAgent(meta AgentMeta) error
	CreateUser(username string) error
	LoadUser(username string) error
	SwitchUser(name string) error
	LoadAgent(agentName string) error
	EditAgent(oldName string, meta AgentMeta) error
	SwitchAgent(name string) error
	UnloadUser() error
	UnloadAgent() error
	Tools() []tools.Tool
}
