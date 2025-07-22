package tui

import (
		"fmt"
    "github.com/fatih/color"
		"strings"
    "github.com/johnjallday/dolphin-tool-calling-agent/internal/app"
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
