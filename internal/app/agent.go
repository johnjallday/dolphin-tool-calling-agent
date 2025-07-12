package app

import (
	"context"
  "github.com/johnjallday/dolphin-tool-calling-agent/internal/user"
)

// AgentManager is responsible for:
//   1) choosing which agent to run (either from the user's DefaultAgent
//      or by prompting them) and
//   2) loading/spawning that agent so REPLApp can use it.
type AgentManager interface {
    Select(ctx context.Context, u *user.User) (agentName string, err error)
    Load(ctx context.Context, u *user.User, agentName string) error
}
