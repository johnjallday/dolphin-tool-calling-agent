package tui

import (
	"fmt"
	"os"

	"github.com/johnjallday/dolphin-tool-calling-agent/internal/registry"
	"github.com/jedib0t/go-pretty/v6/table"
)

func PrintLogo() {
	logo := `
		🐬
`
	fmt.Print("\033[36m" + logo + "\033[0m\n")
	fmt.Println("Dolphin Tool Calling REPL")
}


func PrintTools() {
	toolsSpecs := registry.Specs()
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Name", "Description"})
	for _, ts := range toolsSpecs {
		coloredName := "\033[1;32m" + ts.Name + "\033[0m"
		coloredDesc := "\033[36m" + ts.Description + "\033[0m"
		t.AppendRow(table.Row{coloredName, coloredDesc})
	}
	t.Render()
}


func PrintToolPacks() {
    const dir = "./configs/tools"
    entries, err := os.ReadDir(dir)
    if err != nil {
        fmt.Println("error reading tools:", err)
        return
    }
    for _, e := range entries {
        if e.Name() == ".DS_Store" {
            continue
        }
        fmt.Println(e.Name())
    }
}
