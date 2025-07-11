package app

import (
  "context"
  "github.com/johnjallday/dolphin-tool-calling-agent/internal/agent"
  //"github.com/johnjallday/dolphin-tool-calling-agent/internal/user"
)

type App interface {
  // Initialize the app with a Settings path and a user name.
  // Loads settings.toml and user profile.
  Init(settingsPath string, userName string) error
  Run(ctx context.Context) error
  Shutdown() error
  // Agent management:
  LoadAgent(path string) error   // load a new agent from TOML
  UnloadAgent()                 // clear current agent
  CurrentAgent() agent.Agent    // inspect the currently loaded agent
  CurrentAgentConfig() string   // path to the current agent config

  // Query helpers:
  ListAgents() ([]agent.AgentConfig, error) // what configs are available?
  CreateAgent() error                       // stub out a new agent TOML interactively
}
