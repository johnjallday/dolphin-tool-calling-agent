package app

import (
	"context"

	"github.com/johnjallday/dolphin-tool-calling-agent/internal/agent"
  "github.com/johnjallday/dolphin-tool-calling-agent/internal/user"
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
  SendMessage(context.Context, string) error
	LoadUser(username string) error
	LoadAgent(agentName string) error
	UnloadUser() error
	UnloadAgent() error
  //Agents() []AgentMeta
  //Tools() []ToolInfo
}
