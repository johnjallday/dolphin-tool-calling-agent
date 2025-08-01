package app

import (
  "fmt"
  "os"
  "path/filepath"
	"context"
	"strings"

  "github.com/BurntSushi/toml"
  "github.com/johnjallday/dolphin-tool-calling-agent/internal/agent"
  "github.com/johnjallday/dolphin-tool-calling-agent/internal/user"
	"github.com/johnjallday/dolphin-tool-calling-agent/internal/store"
	//"github.com/johnjallday/dolphin-tool-calling-agent/internal/registry"
	"github.com/johnjallday/dolphin-tool-calling-agent/pkg/tools"
)

type AppConfig struct {
  DefaultUser string `toml:"default_user"`
}

type DefaultApp struct {
  user *user.User
	agent *agent.Agent
}

// NewApp returns the concrete implementation.
func NewApp() App {
  return &DefaultApp{}
}


func (a *DefaultApp) Init() error {
	// ensure configs/ + app_setting.toml + toolpacks.toml
  if err := store.EnsureConfigDir(); err != nil {
    return fmt.Errorf("ensure config dir: %w", err)
  }

  // create subfolders if you still need them
  _ = os.MkdirAll(filepath.Join(store.DefaultConfigDir, "users"), 0755)
  _ = os.MkdirAll("plugins", 0755)

  // load the top-level app settings (which only has DefaultUser)
  settings, err := store.LoadAppSettings()
  if err != nil {
    return fmt.Errorf("read app_setting.toml: %w", err)
  }

  // if there is no default user, bail out: onboarding will kick in
  if settings.DefaultUser == "" {
    return nil
  }

  // 1) load the user from configs/users/…  
  if err := a.LoadUser(settings.DefaultUser); err != nil {
    return fmt.Errorf("load default user %q: %w", settings.DefaultUser, err)
  }

  // 2) if the user‐TOML specified a default_agent, the User.DefaultAgent
  //    field will already be set by NewUser(...).  Detect it and call
  //    your normal SetDefaultAgent logic to wire it up:
  if u := a.User(); u != nil && u.DefaultAgent != nil {
    if err := a.SetDefaultAgent(u.DefaultAgent.Name); err != nil {
      return fmt.Errorf("load user default agent %q: %w",
        u.DefaultAgent.Name, err)
    }
  }

  return nil
}


// SetDefaultUser persists & then loads that user
func (a *DefaultApp) SetDefaultUser(userName string) error {
  if err := store.SetDefaultUser(userName); err != nil {
    return fmt.Errorf("persist default user: %w", err)
  }
  if err := a.LoadUser(userName); err != nil {
    return fmt.Errorf("load user %q: %w", userName, err)
  }
  return nil
}

func (a *DefaultApp) Users() []string {
  var names []string
  dir := "configs/users"
  entries, err := os.ReadDir(dir)
  if err != nil {
    return names
  }
  for _, e := range entries {
    if e.IsDir() || filepath.Ext(e.Name()) != ".toml" {
      continue
    }
    var m struct{ Name string `toml:"name"` }
    if _, err := toml.DecodeFile(filepath.Join(dir, e.Name()), &m); err == nil && m.Name != "" {
      names = append(names, m.Name)
    }
  }
  return names
}

func (a *DefaultApp) CreateUser(userID string) error {
  u, err := user.CreateUser(userID)
  if err != nil {
    return fmt.Errorf("create user %q: %w", userID, err)
  }
  a.user = u
  a.agent = nil

  return nil
}

// LoadUser loads the user TOML and then loads the default agent.
func (a *DefaultApp) LoadUser(username string) error {
  u, err := user.NewUser(username)
  if err != nil {
    return fmt.Errorf("load user %q: %w", username, err)
  }
  a.user = u
  a.agent = u.DefaultAgent
  return nil
}

func (a *DefaultApp) User() *user.User {
    return a.user
}

func (a *DefaultApp) Agent() *agent.Agent {
  return a.agent
}

// SetDefaultAgent persists the named agent as the default in the user's
// TOML, then loads it so that a.agent and a.user.DefaultAgent get set.
func (a *DefaultApp) SetDefaultAgent(agentName string) error {
  // 1) must have a user loaded
  if a.user == nil {
    return fmt.Errorf("no user loaded")
  }

  // 2) make sure that agentName is one of a.user.Agents
  var ok bool
  for _, m := range a.user.Agents {
    if m.Name == agentName {
      ok = true
      break
    }
  }
  if !ok {
    return fmt.Errorf("agent %q not found for user %q", agentName, a.user.Name)
  }

  // 3) load the on‐disk config, set the default_agent field, save it
  cfg, err := store.LoadUserConfig(a.user.Name)
  if err != nil {
    return fmt.Errorf("could not load user config: %w", err)
  }
  cfg.DefaultAgent = agentName
  if err := store.SaveUserConfig(cfg); err != nil {
    return fmt.Errorf("could not save user config: %w", err)
  }

  // 4) load (or reload) that agent so both a.agent and a.user.DefaultAgent update
  if err := a.LoadAgent(agentName); err != nil {
    return fmt.Errorf("could not load agent %q: %w", agentName, err)
  }

  return nil
}

// LoadAgent selects one of the user’s agents by name and sets it to agent.
func (a *DefaultApp) LoadAgent(agentName string) error {
  if a.user == nil {
    return fmt.Errorf("no user loaded")
  }
  var meta *user.AgentMeta
  for i := range a.user.Agents {
    if a.user.Agents[i].Name == agentName {
      meta = &a.user.Agents[i]
      break
    }
  }
  if meta == nil {
    fmt.Println("agent %q not found for user %q", agentName, a.user.Name)
    return nil
  }

  ag, err := agent.NewAgent(meta.Name, meta.Model, meta.Plugins)
  if err != nil {
    return fmt.Errorf("init agent %q: %w", meta.Name, err)
  }

  a.agent = ag
  // also update the default in the user struct if desired
  a.user.DefaultAgent = ag
  return nil
}

func (a *DefaultApp) SwitchUser(name string) error {
    // if there’s already a user, unload them
    if a.User != nil {
        if err := a.UnloadUser(); err != nil {
            return fmt.Errorf("could not unload existing user: %w", err)
        }
    }
    // now load the new one
    if err := a.LoadUser(name); err != nil {
        return fmt.Errorf("could not load user %q: %w", name, err)
    }
    return nil
}

func (a *DefaultApp) UnloadAgent() error {
  if a.agent == nil {
    return fmt.Errorf("no agent loaded")
  }
  a.agent = nil
  return nil
}

func (a *DefaultApp) UnloadUser() error {
  if a.user == nil {
    return fmt.Errorf("no user loaded")
  }
  a.user = nil
  a.agent = nil
  return nil
}

func (a *DefaultApp) SendMessage(ctx context.Context, msg string) (reply string, err error) {
  if a.agent == nil {
    // must return "" for reply when erroring
    return "", fmt.Errorf("no agent loaded")
  }
  // forward the two return values from your agent
  return a.agent.SendMessage(ctx, msg)
}


// Tools returns the slice of registered tools.
func (a *DefaultApp) Tools() []tools.Tool {
    // if there is no agent, just return an empty slice
    if a.agent == nil {
        return nil
    }
    return a.agent.Tools()
}

func (a *DefaultApp) Agents() []user.AgentMeta {
    if a.user == nil {
        // no user loaded → no agents
        return nil
    }
    // simply return the slice of AgentMeta from the loaded user
    return a.user.Agents
}



// CreateAgent will add a new agent entry to the currently
// loaded user’s TOML config and then refresh the in‐memory User.
func (a *DefaultApp) CreateAgent(meta AgentMeta) error {
  if a.user == nil {
    return fmt.Errorf("no user loaded")
  }

  // Path to the user’s TOML file
  userFile := filepath.Join("configs", "users", a.user.Name+".toml")

  // 1) decode existing file into a struct that mirrors your on‐disk layout
  var cfg struct {
    Name         string          `toml:"name"`
    DefaultAgent string          `toml:"default_agent"`
    Agents       []user.AgentMeta `toml:"agents"`
  }
  if _, err := toml.DecodeFile(userFile, &cfg); err != nil {
    return fmt.Errorf("decode %s: %w", userFile, err)
  }

  // 2) append the new agent meta (convert our app.AgentMeta → user.AgentMeta)
  cfg.Agents = append(cfg.Agents, user.AgentMeta{
    Name:    meta.Name,
    Model:   meta.Model,
    Plugins: meta.ToolPaths,
  })

  // 3) re‐write the TOML file (truncating)
  f, err := os.Create(userFile)
  if err != nil {
    return fmt.Errorf("rewrite %s: %w", userFile, err)
  }
  defer f.Close()
  enc := toml.NewEncoder(f)
  if err := enc.Encode(cfg); err != nil {
    return fmt.Errorf("encode %s: %w", userFile, err)
  }

  // 4) re‐load the user so that a.user.Agents is refreshed
  u, err := user.NewUser(a.user.Name)
  if err != nil {
    return fmt.Errorf("reload user %q: %w", a.user.Name, err)
  }
  a.user = u
  return nil
}

// SwitchAgent switches the current agent to one of the already‐created agents
// for the current user (loading its .so plugins under the hood).
func (a *DefaultApp) SwitchAgent(name string) error {
    // pull the loaded user back out via the method
    u := a.User()
    if u == nil {
        return fmt.Errorf("no user loaded")
    }

    // make sure the named agent exists in u.Agents
    var found bool
    for _, m := range u.Agents {
        if m.Name == name {
            found = true
            break
        }
    }
    if !found {
        return fmt.Errorf("agent %q not found for user %q", name, u.Name)
    }

    // unload the prior agent (if any)
    if err := a.UnloadAgent(); err != nil {
        return fmt.Errorf("unload existing agent: %w", err)
    }

    // load the new agent
    if err := a.LoadAgent(name); err != nil {
        return fmt.Errorf("load agent %q: %w", name, err)
    }

    return nil
}

func (a *DefaultApp) EditAgent(oldName string, meta AgentMeta) error {
    if a.user == nil {
        return fmt.Errorf("no user loaded")
    }
    userFile := filepath.Join("configs", "users", a.user.Name+".toml")

    // Decode the existing user config
    var cfg struct {
        Name         string            `toml:"name"`
        DefaultAgent string            `toml:"default_agent"`
        Agents       []user.AgentMeta  `toml:"agents"`
    }
    if _, err := toml.DecodeFile(userFile, &cfg); err != nil {
        return fmt.Errorf("decode %s: %w", userFile, err)
    }

    // Find & update the matching agent
    found := false
    for i := range cfg.Agents {
        if cfg.Agents[i].Name == oldName {
            cfg.Agents[i] = user.AgentMeta{
                Name:    meta.Name,
                Model:   meta.Model,
                Plugins: meta.ToolPaths,
            }
            // If you also want to rename the default_agent setting:
            if cfg.DefaultAgent == oldName {
                cfg.DefaultAgent = meta.Name
            }
            found = true
            break
        }
    }
    if !found {
        return fmt.Errorf("agent %q not found for user %q", oldName, a.user.Name)
    }

    // Rewrite the TOML file
    f, err := os.Create(userFile)
    if err != nil {
        return fmt.Errorf("rewrite %s: %w", userFile, err)
    }
    defer f.Close()
    enc := toml.NewEncoder(f)
    if err := enc.Encode(cfg); err != nil {
        return fmt.Errorf("encode %s: %w", userFile, err)
    }

    // Reload the in-memory user so a.user.Agents is fresh
    u, err := user.NewUser(a.user.Name)
    if err != nil {
        return fmt.Errorf("reload user %q: %w", a.user.Name, err)
    }
    a.user = u

    // If we had that agent loaded, switch to the new name
    if a.agent != nil && a.agent.Name == oldName {
        if err := a.LoadAgent(meta.Name); err != nil {
            return fmt.Errorf("reload current agent: %w", err)
        }
    }

    return nil
}



// Toolpacks returns the list of plugin “names” found under ./plugins.
// It looks for files ending in .so, and returns each filename minus the .so.
func (a *DefaultApp) Toolpacks() []string {
  const pluginDir = "plugins"
  var names []string

  entries, err := os.ReadDir(pluginDir)
  if err != nil {
    // no plugins folder or unreadable → just return empty
    return names
  }
  for _, e := range entries {
    if e.IsDir() {
      continue
    }
    if filepath.Ext(e.Name()) != ".so" {
      continue
    }
    base := strings.TrimSuffix(e.Name(), ".so")
    names = append(names, base)
  }
  return names
}

// ListRemoteToolpacks reads configs/toolpacks.toml and returns each package name.
func (a *DefaultApp) ListRemoteToolpacks() ([]string, error) {
  pkgs, err := store.LoadRemoteToolpacks()
  if err != nil {
    return nil, fmt.Errorf("load remote toolpacks: %w", err)
  }

  names := make([]string, 0, len(pkgs))
  for _, tp := range pkgs {
    names = append(names, tp.Name)
  }
  return names, nil
}
