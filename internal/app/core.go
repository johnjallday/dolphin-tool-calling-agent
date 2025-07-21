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
  SendMessage(context.Context, string) error
	CreateAgent(meta AgentMeta) error
	LoadUser(username string) error
	LoadAgent(agentName string) error
	UnloadUser() error
	UnloadAgent() error
	Tools() []tools.Tool
	//Status()
}
