package tui

import (
	"fmt"
	//"os"
	"io/fs"
	"path/filepath"
 	"reflect"
	"strings"
	"io"
	"context"
	

	"github.com/fatih/color"
  //"github.com/johnjallday/dolphin-tool-calling-agent/internal/agent"
  //"github.com/johnjallday/dolphin-tool-calling-agent/internal/user"
	"github.com/johnjallday/dolphin-tool-calling-agent/internal/app"
	"github.com/peterh/liner"
	//"github.com/johnjallday/dolphin-tool-calling-agent/internal/registry"
	//"github.com/jedib0t/go-pretty/v6/table"
)

type TUIApp struct {
  Ctx   context.Context
  App   app.App                // your domain interface, see below
  Out   io.Writer          // usually os.Stdout
  Err   io.Writer          // usually os.Stderr
  Rl   *liner.State      // your REPL state if you need it
}


func PrintLogo() {
	logo := `
		🐬
`
	fmt.Print("\033[36m" + logo + "\033[0m\n")
	fmt.Println("Dolphin Tool Calling REPL\n")
}


func PrintToolPacks() {
    const root = "./plugins/"

    err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            // something went wrong reading this entry; just report and keep going
            fmt.Printf("error accessing %q: %v\n", path, err)
            return nil
        }

        // skip macOS metadata files
        if d.Name() == ".DS_Store" {
            if d.IsDir() {
                return filepath.SkipDir
            }
            return nil
        }

        // if this is a file and has .so extension, print it
        if !d.IsDir() && filepath.Ext(d.Name()) == ".so" {
            fmt.Println(path)
        }
        return nil
    })

    if err != nil {
        fmt.Println("error walking plugins directory:", err)
    }
}


// Help prints all App interface methods except Init, without return types.
func PrintHelp() {
    iface := reflect.TypeOf((*app.App)(nil)).Elem()
    fmt.Println("Available commands:")
    for i := 0; i < iface.NumMethod(); i++ {
        m := iface.Method(i)
        // skip the ones you don't want
        if m.Name == "Init" || m.Name == "SendMessage" {
            continue
        }
        sig := m.Type
        // collect parameter types, skipping the receiver at In(0)
        var params []string
        for j := 1; j < sig.NumIn(); j++ {
            params = append(params, sig.In(j).String())
        }
        fmt.Printf("  %s(%s)\n", m.Name, strings.Join(params, ", "))
    }
    fmt.Println()
}

func CurrentUser(t TUIApp) error {
  u := t.App.User()
  cLabel := color.New(color.FgCyan, color.Bold)
  cValue := color.New(color.FgWhite)
  cList  := color.New(color.FgMagenta, color.Bold)

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
    fmt.Printf("  - %s (plugins: %v)\n", meta.Name, meta.Plugins)
  }
  return nil
}

func Users(t TUIApp) error {
  list := t.App.Users()
  fmt.Fprintln(t.Out, "Available Users:")
  for _, name := range list {
    fmt.Fprintf(t.Out, "  %s\n", name)
  }
  return nil
}

func Agent(t TUIApp) error {
  a := t.App.Agent()

	if a == nil {
		fmt.Println(t.Out, "No agent loaded. Please load an agent.")
		return nil
	}
  //fmt.Fprintf(t.Out, "Current agent: %s (%s)\n", a.Name, a.Model)
  cLabel := color.New(color.FgYellow, color.Bold)
  cValue := color.New(color.FgWhite)

  cLabel.Print("Agent: ")
  cValue.Println(a.Name)
  cLabel.Print("Model: ")
  cValue.Println(a.Model)
	Tools(t)
  return nil
}

func Agents(t TUIApp) error {
  metas := t.App.Agents()
  fmt.Fprintln(t.Out, "Agents:")
  for _, m := range metas {
    fmt.Fprintf(t.Out, "  %s\t%s\n", m.Name, m.Model)
  }
  return nil
}

func Tools(t TUIApp) error {
  tools := t.App.Tools()
  //fmt.Fprintln(t.Out, "Tools:")

	cLabel := color.New(color.FgYellow, color.Bold)
	cValue := color.New(color.FgGreen)
	cValueDesc := color.New(color.FgWhite)


	cLabel.Println("Tools")
  for _, tool := range tools {
    cValue.Fprintf(t.Out, "  %s\t", tool.Name)
    cValueDesc.Fprintf(t.Out, "  %s\t\n", tool.Description)
  }

  return nil
}


// Status prints the current user and agent names.
func PrintStatus(t TUIApp) error {
  //s := t.App
  u, a := t.App.User(), t.App.Agent()

  cLabel := color.New(color.FgCyan, color.Bold)
  cValue := color.New(color.FgWhite)

  // if you want to keep using t.Out rather than passing it around:
  if _, err := cLabel.Fprint(t.Out, "Current User:  "); err != nil {
    return err
  }
  if _, err := cValue.Fprintln(t.Out, u.Name); err != nil {
    return err
  }
  if _, err := cLabel.Fprint(t.Out, "Current Agent: "); err != nil {
    return err
  }
  if _, err := cValue.Fprintln(t.Out, a.Name); err != nil {
    return err
  }
  return nil
}

// UnloadUser calls the core UnloadUser and reports success/failure.
func UnloadUser(t TUIApp) error {
  err := t.App.UnloadUser()
  if err != nil {
    // print in red
    color.New(color.FgRed).Fprintf(t.Err, "error unloading user: %v\n", err)
    return err
  }
  color.New(color.FgGreen).Fprintln(t.Out, "✓ user unloaded")
  return nil
}

// UnloadAgent calls the core UnloadAgent and reports success/failure.
func UnloadAgent(t TUIApp) error {
  err := t.App.UnloadAgent()
  if err != nil {
    color.New(color.FgRed).Fprintf(t.Err, "error unloading agent: %v\n", err)
    return err
  }
  color.New(color.FgGreen).Fprintln(t.Out, "✓ agent unloaded")
  return nil
}



// clearScreen emits the ANSI codes to clear the terminal + home the cursor
func (t TUIApp) clearScreen() {
   fmt.Fprint(t.Out, "\x1b[2J\x1b[H")
}

// Refresh clears the screen, reprints the logo, and then the status panel
func (t TUIApp) Refresh() error {
   t.clearScreen()
   PrintLogo()
   return PrintStatus(t)
}
