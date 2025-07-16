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
	CurrentUser() *user.User
  //SelectUser(name string) error
	CurrentAgent() *agent.Agent
  SendMessage(context.Context, string) error
  //Agents() []AgentMeta
  //SelectAgent(meta AgentMeta) error
  //CurrentUser() string
  //CurrentAgent() AgentMeta
  //Tools() []ToolInfo
  //Send(ctx context.Context, prompt string) (string, error)
}
