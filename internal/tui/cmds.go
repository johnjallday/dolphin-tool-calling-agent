package tui

import (
		"fmt"
		"bufio"
    "github.com/fatih/color"
		"strings"
    "github.com/johnjallday/dolphin-tool-calling-agent/internal/app"
    "github.com/johnjallday/dolphin-tool-calling-agent/internal/user"
		"os"
  	"path/filepath"
)


// LoadUserCmd loads the named user, prints a success message, then refreshes.
func LoadUserCmd(t *TUIApp, args []string) error {
		if len(args) != 1 {
			fmt.Fprintln(t.Out, "usage: load-user <username>")
		return nil
		}
		name := args[0]


		// call into your core app to load the user
		if err := t.App.LoadUser(name); err != nil {
				return fmt.Errorf("load user %q: %w", name, err)
		}

		// on success, print a green check and refresh the UI
		color.New(color.FgGreen).
				Fprintln(t.Out, "✓ user loaded:", name)

		return t.Refresh()
}



func SwitchUserCmd(t *TUIApp, args []string) error {
    if len(args) != 1 {
        fmt.Fprintln(t.Out, "usage: switch-user <username>")
        return nil
    }
    name := args[0]

    // call your new App.SwitchUser
    if err := t.App.SwitchUser(name); err != nil {
        return fmt.Errorf("switch user: %w", err)
    }

    color.New(color.FgGreen).Fprintln(t.Out, "✓ switched to user", name)
    return t.Refresh()
}

// LoadAgentCmd loads the named agent for the current user, prints a success
// message, and then refreshes the UI.
func LoadAgentCmd(t *TUIApp, args []string) error {
    if len(args) != 1 {
        fmt.Fprintln(t.Out, "usage: load-agent <agent-name>")
        return nil
    }
    name := args[0]

    // call into your core app to load the agent
    if err := t.App.LoadAgent(name); err != nil {
        return fmt.Errorf("load agent %q: %w", name, err)
    }

    // on success, print a green check and refresh the UI
    color.New(color.FgGreen).
        Fprintln(t.Out, "✓ agent loaded:", name)

    return t.Refresh()
}

func SwitchAgentCmd(t *TUIApp, args []string) error {
    if len(args) != 1 {
        fmt.Fprintln(t.Out, "usage: switch-agent <agent-name>")
        return nil
    }
    name := args[0]

    // call into your core app to switch the agent
    if err := t.App.SwitchAgent(name); err != nil {
        return fmt.Errorf("switch agent %q: %w", name, err)
    }

    // on success, print a green check and refresh
    color.New(color.FgGreen).
        Fprintln(t.Out, "✓ switched to agent:", name)

    return t.Refresh()
}




// UnloadUserCmd unloads the current user and refreshes.
func UnloadUserCmd(t *TUIApp, _ []string) error {
    if err := t.App.UnloadUser(); err != nil {
        // print error to stderr and bail
        color.New(color.FgRed).
            Fprintln(t.Err, "error unloading user:", err)
        return err
    }
    // success!
    color.New(color.FgGreen).Fprintln(t.Out, "✓ user unloaded")
    return t.Refresh()
}

// UnloadAgentCmd unloads the current agent and refreshes.
func UnloadAgentCmd(t *TUIApp, _ []string) error {
    if err := t.App.UnloadAgent(); err != nil {
        color.New(color.FgRed).Fprintf(t.Err, "error unloading agent: %v\n", err)
        return err
    }
    color.New(color.FgGreen).Fprintln(t.Out, "✓ agent unloaded")
    return t.Refresh()
}

// CreateAgentCmd implements “create-agent <name> <model> [tool1,tool2,…]” then refreshes.
func CreateAgentCmd(t *TUIApp, args []string) error {
    if len(args) < 2 {
        fmt.Fprintln(t.Out, "usage: create-agent <name> <model> [tool1,tool2,…]")
        return nil
    }
    name, model := args[0], args[1]

    var paths []string
    if len(args) > 2 && args[2] != "" {
        paths = strings.Split(args[2], ",")
    }

    meta := app.AgentMeta{
        Name:      name,
        Model:     model,
        ToolPaths: paths,
    }
    if err := t.App.CreateAgent(meta); err != nil {
        return fmt.Errorf("create agent: %w", err)
    }

    fmt.Fprintf(t.Out,
        "✅ Agent %q (model=%q) created with %d tool(s)\n",
        name, model, len(paths),
    )
    return t.Refresh()
}


// UserCmd prints the current user and their agents.
func UserCmd(t *TUIApp, _ []string) error {
    u := t.App.User()
    cLabel := color.New(color.FgCyan, color.Bold)
    cValue := color.New(color.FgWhite)
    cList := color.New(color.FgMagenta, color.Bold)

    cLabel.Print("Current User: ")
    cValue.Println(u.Name)

    cLabel.Print("Default Agent: ")
    if u.DefaultAgent != nil {
        cValue.Println(u.DefaultAgent.Name)
    } else {
        cValue.Println("<none>")
    }

    cList.Println("Available Agents:")
    for _, meta := range u.Agents {
        fmt.Fprintf(t.Out, "  - %s (plugins: %v)\n", meta.Name, meta.Plugins)
    }
    return nil
}

// UsersCmd lists all users.
func UsersCmd(t *TUIApp, _ []string) error {

    cLabel := color.New(color.FgCyan, color.Bold)
    cLabel.Fprint(t.Out, "Available Users:  ")
    //cValue.Fprintln(t.Out, u.Name)
    for _, name := range t.App.Users() {
        fmt.Fprintf(t.Out, "  %s\t", name)
    }

		fmt.Fprintln(t.Out)
    return nil
}

// AgentCmd prints the current agent (or a warning if none) and its tools.
func AgentCmd(t *TUIApp, _ []string) error {
    a := t.App.Agent()
    if a == nil {
        fmt.Fprintln(t.Out, "No agent loaded. Please load an agent.")
        return nil
    }
    cLabel := color.New(color.FgYellow, color.Bold)
    cValue := color.New(color.FgWhite)

    cLabel.Print("Agent: ")
    cValue.Println(a.Name)
    cLabel.Print("Model: ")
    cValue.Println(a.Model)

    return ToolsCmd(t, nil)
}

// AgentsCmd lists all agents for the current user.
func AgentsCmd(t *TUIApp, _ []string) error {
    fmt.Fprintln(t.Out, "Agents:")
    for _, m := range t.App.Agents() {
        fmt.Fprintf(t.Out, "  %s\t%s\n", m.Name, m.Model)
    }
    return nil
}

// ToolsCmd lists all tools on the current agent.
func ToolsCmd(t *TUIApp, _ []string) error {
    // guard against no‐agent
    if t.App.Agent() == nil {
        fmt.Fprintln(t.Out, "No agent loaded")
        return nil
    }

    cLabel := color.New(color.FgYellow, color.Bold)
    cVal   := color.New(color.FgGreen)
    cDesc  := color.New(color.FgWhite)

    cLabel.Println("Tools:")
		if t.App.Tools() == nil {
			fmt.Println("No Agent Loaded: Load an Agent")
		}

    for _, tool := range t.App.Tools() {
        cVal.Fprintf(t.Out, "  %s\t", tool.Name)
        cDesc.Fprintf(t.Out, "%s\n", tool.Description)
    }
    return nil
}

// EditAgentCmd prompts the user to update an existing agent’s name, model,
// and tool paths. Usage: edit-agent <old-name>
func EditAgentCmd(t *TUIApp, args []string) error {
    if len(args) != 1 {
        fmt.Fprintln(t.Out, "usage: edit-agent <agent-name>")
        return nil
    }
    oldName := args[0]

    // find the existing agent
    var current user.AgentMeta
    found := false
    for _, m := range t.App.Agents() {
        if m.Name == oldName {
            current = m
            found = true
            break
        }
    }
    if !found {
        return fmt.Errorf("agent %q not found", oldName)
    }

    // Print current info
    cyan := color.New(color.FgCyan, color.Bold)
    white := color.New(color.FgWhite)
    cyan.Fprint(t.Out, "Editing Agent: ")
    white.Fprintln(t.Out, current.Name)
    cyan.Fprint(t.Out, "  Model: ")
    white.Fprintln(t.Out, current.Model)
    cyan.Fprint(t.Out, "  Tools: ")
    white.Fprintln(t.Out, strings.Join(current.Plugins, ", "))

    reader := bufio.NewReader(t.In) // ← now t.In is defined

    // prompt for new name
    fmt.Fprintf(t.Out, "New name [%s]: ", current.Name)
    line, err := reader.ReadString('\n')
    if err != nil {
        return err
    }
    line = strings.TrimSpace(line)
    newName := current.Name
    if line != "" {
        newName = line
    }

    // prompt for new model
    fmt.Fprintf(t.Out, "New model [%s]: ", current.Model)
    line, err = reader.ReadString('\n')
    if err != nil {
        return err
    }
    line = strings.TrimSpace(line)
    newModel := current.Model
    if line != "" {
        newModel = line
    }

    // prompt for new tools
    fmt.Fprintf(t.Out,
        "New tool paths (comma-separated) [%s]: ",
        strings.Join(current.Plugins, ","),
    )
    line, err = reader.ReadString('\n')
    if err != nil {
        return err
    }
    line = strings.TrimSpace(line)

    var newTools []string
    if line == "" {
        newTools = current.Plugins
    } else {
        for _, p := range strings.Split(line, ",") {
            if t := strings.TrimSpace(p); t != "" {
                newTools = append(newTools, t)
            }
        }
    }

    // build the app.AgentMeta and call EditAgent
    updated := app.AgentMeta{
        Name:      newName,
        Model:     newModel,
        ToolPaths: newTools,
    }
    if err := t.App.EditAgent(oldName, updated); err != nil {
        return fmt.Errorf("edit agent: %w", err)
    }

    color.New(color.FgGreen).
        Fprintln(t.Out, "✓ agent updated:", oldName, "→", newName)

    return t.Refresh()
}

// InitCmd ensures configs/users exists, and if there are no .toml users,
// immediately calls CreateUserCmd to bootstrap the first user.
func InitCmd(t *TUIApp, args []string) error {
  usersDir := filepath.Join("configs", "users")
  if err := os.MkdirAll(usersDir, 0755); err != nil {
    return fmt.Errorf("mkdir %q: %w", usersDir, err)
  }

  entries, err := os.ReadDir(usersDir)
  if err != nil {
    return fmt.Errorf("read dir %q: %w", usersDir, err)
  }
  // look for any .toml files
  found := false
  for _, e := range entries {
    if !e.IsDir() && filepath.Ext(e.Name()) == ".toml" {
      found = true
      break
    }
  }
  if !found {
    // no users yet ⇒ prompt to create one
    color.New(color.FgYellow).Fprintln(t.Out,
      "No users found; let's create your first user.")
    if err := CreateUserCmd(t, nil); err != nil {
      return fmt.Errorf("bootstrap user: %w", err)
    }
  }
  return nil
}

// CreateUserCmd prompts the user for a userID, calls App.CreateUser,
// then optionally SwitchUser to make it the default.
func CreateUserCmd(t *TUIApp, args []string) error {
  reader := bufio.NewReader(t.In)

  // 1) ask for user ID
  fmt.Fprint(t.Out, "Enter new user ID: ")
  line, err := reader.ReadString('\n')
  if err != nil {
    return err
  }
  userID := strings.TrimSpace(line)
  if userID == "" {
    fmt.Fprintln(t.Out, "aborted: empty user ID")
    return nil
  }

  // 2) create it
  if err := t.App.CreateUser(userID); err != nil {
    return fmt.Errorf("CreateUser %q: %w", userID, err)
  }
  color.New(color.FgGreen).Fprintln(t.Out,
    "✓ user created:", userID)

  // 3) ask whether to make it default
  fmt.Fprint(t.Out, "Make this the default user? [Y/n]: ")
  yn, err := reader.ReadString('\n')
  if err != nil {
    return err
  }
  yn = strings.TrimSpace(yn)
  if yn == "" || strings.ToLower(yn[:1]) == "y" {
    if err := t.App.SwitchUser(userID); err != nil {
      return fmt.Errorf("SwitchUser %q: %w", userID, err)
    }
    color.New(color.FgGreen).Fprintln(t.Out,
      "✓ set as default user:", userID)
  } else {
    color.New(color.FgCyan).Fprintln(t.Out,
      "Default user left unchanged.")
  }

  return t.Refresh()
}
