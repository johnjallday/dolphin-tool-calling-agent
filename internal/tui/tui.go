package tui

import (
	"fmt"
	//"os"
	"io/fs"
	"path/filepath"

	//"github.com/johnjallday/dolphin-tool-calling-agent/internal/registry"
	//"github.com/jedib0t/go-pretty/v6/table"
)

func PrintLogo() {
	logo := `
		🐬
`
	fmt.Print("\033[36m" + logo + "\033[0m\n")
	fmt.Println("Dolphin Tool Calling REPL")
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

