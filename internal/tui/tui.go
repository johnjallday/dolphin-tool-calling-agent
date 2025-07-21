package tui

import (
    "context"
    "fmt"
    "io"
    "io/fs"
    "path/filepath"
    "reflect"
    "strings"

    "github.com/fatih/color"
    "github.com/johnjallday/dolphin-tool-calling-agent/internal/app"
    "github.com/peterh/liner"
)

type TUIApp struct {
    Ctx context.Context
    App app.App
    Out io.Writer
    Err io.Writer
    Rl  *liner.State
}

// clearScreen emits ANSI codes to clear the terminal + move cursor home.
func (t *TUIApp) clearScreen() {
    fmt.Fprint(t.Out, "\x1b[2J\x1b[H")
}

// Refresh clears the screen, prints the logo, then status.
func (t *TUIApp) Refresh() error {
    t.clearScreen()
    PrintLogo()
    return PrintStatus(t)
}

// PrintLogo prints the dolphin logo banner.
func PrintLogo() {
    logo := `
        🐬
`
    fmt.Print("\033[36m" + logo + "\033[0m\n")
    fmt.Println("Dolphin Tool Calling Agent Client\n")
}

// PrintToolPacks walks “./plugins” looking for .so files.
func PrintToolPacks() {
    const root = "./plugins/"
    filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            fmt.Printf("error accessing %q: %v\n", path, err)
            return nil
        }
        if d.Name() == ".DS_Store" {
            if d.IsDir() {
                return filepath.SkipDir
            }
            return nil
        }
        if !d.IsDir() && filepath.Ext(d.Name()) == ".so" {
            fmt.Println(path)
        }
        return nil
    })
}

// PrintHelp uses reflection to list all methods on app.App.
func PrintHelp() {
    iface := reflect.TypeOf((*app.App)(nil)).Elem()
    fmt.Println("Available commands:")
    for i := 0; i < iface.NumMethod(); i++ {
        m := iface.Method(i)
        if m.Name == "Init" || m.Name == "SendMessage" {
            continue
        }
        sig := m.Type
        var params []string
        for j := 1; j < sig.NumIn(); j++ {
            params = append(params, sig.In(j).String())
        }
        fmt.Printf("  %s(%s)\n", m.Name, strings.Join(params, ", "))
    }
    fmt.Println()
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
    fmt.Fprintln(t.Out, "Available Users:")
    for _, name := range t.App.Users() {
        fmt.Fprintf(t.Out, "  %s\n", name)
    }
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

// PrintStatus is used by Refresh() to show user+agent at top.
func PrintStatus(t *TUIApp) error {
    u, a := t.App.User(), t.App.Agent()
    cLabel := color.New(color.FgCyan, color.Bold)
    cValue := color.New(color.FgWhite)

    cLabel.Fprint(t.Out, "Current User:  ")
    cValue.Fprintln(t.Out, u.Name)
    cLabel.Fprint(t.Out, "Current Agent: ")
    cValue.Fprintln(t.Out, a.Name)
    return nil
}

// UnloadUserCmd unloads the current user and refreshes.
func UnloadUserCmd(t *TUIApp, _ []string) error {
    if err := t.App.UnloadUser(); err != nil {
        color.New(color.FgRed).Fprintf(t.Err, "error unloading user: %v\n", err)
        return err
    }
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
