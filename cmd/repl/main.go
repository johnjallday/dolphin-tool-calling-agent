package main

import (
	"context"
	//"flag"
	"fmt"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/openai/openai-go"

	"github.com/johnjallday/dolphin-tool-calling-agent/internal/registry"
	"github.com/johnjallday/dolphin-tool-calling-agent/internal/agent"
	"github.com/johnjallday/dolphin-tool-calling-agent/internal/tui"
)

func readQuestion() string {
	rl, err := readline.New("> ")
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	rl.Config.FuncFilterInputRune = func(r rune) (rune, bool) {
		if r == readline.CharCtrlL {
			fmt.Print("\033[H\033[2J")
			rl.Refresh()
			return 0, false
		}
		return r, true
	}

	for {
		line, err := rl.Readline()
		if err != nil {
			if err == readline.ErrInterrupt {
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


// Load an agent from a config file
func loadAgent(client *openai.Client, path string) (agent.Agent, error) {
	a, err := agent.NewAgentFromConfig(client, path)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Loaded agent from %s\n", path)
	return a, nil
}

func main() {

	client := openai.NewClient()
	clientPtr := &client
	fmt.Println("OpenAI Client starting")

	DefaultAgentPath := "./configs/user/agents/reaper_agent.toml"

	tui.PrintLogo()

	var agentInstance agent.Agent
	var agentConfigPath string

	// Try to load default agent
	agentInstance, err := loadAgent(clientPtr, DefaultAgentPath)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading agent config (%s): %v\n", DefaultAgentPath, err)
	} else {
		//agentConfigPath = *defaultConfigPath
		fmt.Println("test")
	}

	tui.PrintTools()
	ctx := context.Background()
	for {
		question := readQuestion()
		parts := strings.Fields(question)
		if len(parts) == 0 {
			continue
		}

		switch strings.ToLower(parts[0]) {
		case "help", "tools", "-t":
			tui.PrintLogo()
			tui.PrintTools()
			continue
		case "exit", "quit":
			fmt.Print("Bye!\r\n")
			return
		case "load-agent":
			if len(parts) < 2 {
				fmt.Println("Usage: load-agent <path-to-toml>")
				continue
			}
			newPath := parts[1]
			//newAgent, err := loadAgent(client, newPath)
			newAgent, err := loadAgent(clientPtr, newPath)
			if err != nil {
				fmt.Printf("Failed to load agent: %v\n", err)
				continue
			}
			agentInstance = newAgent
			agentConfigPath = newPath
			continue
		case "unload-agent":
			agentInstance = nil
			agentConfigPath = ""
			registry.Clear() 
			fmt.Println("Agent unloaded.")
			continue
		case "current-agent":
			if agentInstance != nil && agentConfigPath != "" {
				fmt.Printf("Current agent loaded from: %s\n", agentConfigPath)
			} else {
				fmt.Println("No agent loaded.")
			}
			continue
		}

		if agentInstance == nil {
			fmt.Println("No agent loaded. Use: load-agent <path-to-toml>")
			continue
		}

		err := agentInstance.SendMessage(ctx, question)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}
