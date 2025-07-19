package tui

import (
	"fmt"
	//"os"
	"io/fs"
	"path/filepath"
 	"reflect"
	"strings"

	"github.com/fatih/color"
  "github.com/johnjallday/dolphin-tool-calling-agent/internal/agent"
  "github.com/johnjallday/dolphin-tool-calling-agent/internal/user"
	"github.com/johnjallday/dolphin-tool-calling-agent/internal/app"
	//"github.com/johnjallday/dolphin-tool-calling-agent/internal/registry"
	//"github.com/jedib0t/go-pretty/v6/table"
)

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


func PrintUser(u *user.User) {
  cLabel := color.New(color.FgCyan, color.Bold)
  cValue := color.New(color.FgWhite)
  cList  := color.New(color.FgMagenta, color.Bold)

  cLabel.Print("User: ")
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

	//cList.Print("Available Tools:")
	//PrintTools(u)
}

func PrintAgent(a *agent.Agent) {
  cLabel := color.New(color.FgYellow, color.Bold)
  cValue := color.New(color.FgWhite)

  cLabel.Print("Agent: ")
  cValue.Println(a.Name)
  cLabel.Print("Model: ")
  cValue.Println(a.Model)
  // Print tools, etc. as you wish
	// fm
	//fmt.Println(a.Registry.String())
	PrintTools(a)
}


func PrintTools(a *agent.Agent) {
	cLabel := color.New(color.FgYellow, color.Bold)
	cValue := color.New(color.FgGreen)

	cLabel.Println("Tools")
	cValue.Println(a.Registry)

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
