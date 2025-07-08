package main

import (
	//"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/openai/openai-go"

	"github.com/johnjallday/dolphin-tool-calling-agent/registry"
	"github.com/johnjallday/dolphin-tool-calling-agent/device"
	"github.com/johnjallday/dolphin-tool-calling-agent/agents"
)

// listTools prints all available tools.
func listTools() {
	for _, ts := range registry.Specs() {
		fmt.Printf("%s: %s\n", ts.Name, ts.Description)
	}
}

func printTools() {
	// ANSI color codes:
	// "\033[1;32m" = bold green, "\033[36m" = cyan, "\033[0m" = reset
	toolsSpecs := registry.Specs()
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	// Append header row.
	t.AppendHeader(table.Row{"Name", "Description"})

	for _, ts := range toolsSpecs {
		// Using ANSI escape codes for colors: bold green for Name and cyan for Description.
		coloredName := "\033[1;32m" + ts.Name + "\033[0m"
		coloredDesc := "\033[36m" + ts.Description + "\033[0m"
		t.AppendRow(table.Row{coloredName, coloredDesc})
	}
	t.Render()
}

func readQuestion() string {
	var rl *readline.Instance
	var err error

	rl, err = readline.New("> ")
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	rl.Config.FuncFilterInputRune = func(r rune) (rune, bool) {
		if r == readline.CharCtrlL {
			fmt.Print("\033[H\033[2J") // clear screen
			rl.Refresh()
			return 0, false
		}
		return r, true
	}

	for {
		line, err := rl.Readline()
		if err != nil {
			if err == readline.ErrInterrupt {
				// Ctrl+C pressed
				fmt.Println("\nExiting...")
				os.Exit(0)
			}
			fmt.Println()
			return ""
		}
		line = strings.TrimSpace(line)
		if line != "" {
			return line
		}
	}
}


// printLogo prints a colored ASCII dolphin logo and the REPL title at startup.
func printLogo() {
	logo := `
		🐬
`
	fmt.Print("\033[36m" + logo + "\033[0m\n")
	fmt.Println("Dolphin Tool Calling REPL")
}

func main() {
	printLogo()
	device.GetCurrentAudioDevice()

	// Flags
	showTools := flag.Bool("tools", false, "list available tools")
	flag.BoolVar(showTools, "t", false, "list available tools (shorthand)")
	//configPath := flag.String("config", "./user/agents/calculator_agent.toml", "path to agent config TOML file")
	configPath := flag.String("config", "./user/agents/reaper_agent.toml", "path to agent config TOML file")
	flag.Parse()

	client := openai.NewClient()
	// Create agent from config
	agentInstance, err := agent.NewAgentFromConfig(&client, *configPath)
	if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading agent config (%s): %v\n", *configPath, err)
			os.Exit(1)
	}

	if *showTools {
		printTools()
		return
	}

	printTools()
	ctx := context.Background()
	for {
		question := readQuestion()
		switch strings.ToLower(question) {
		case "help", "tools", "-t":
			printLogo()
			printTools()
			continue
		case "exit", "quit":
			fmt.Print("Bye!\r\n")
			return
		}

		//err := chatbot.SendMessage(ctx, question)
		err := agentInstance.SendMessage(ctx, question)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
	}
}
