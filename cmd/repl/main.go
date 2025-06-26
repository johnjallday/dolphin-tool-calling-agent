package main

import (
	//"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/openai/openai-go"

	"github.com/johnjallday/dolphin-tool-calling-agent/registry"
	//"github.com/johnjallday/dolphin-tool-calling-agent/chat"
	"github.com/johnjallday/dolphin-tool-calling-agent/agents"
)

// listTools prints all available tools.
func listTools() {
	for _, ts := range registry.Specs() {
		fmt.Printf("%s: %s\n", ts.Name, ts.Description)
	}
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
	showTools := flag.Bool("tools", false, "list available tools")
	flag.BoolVar(showTools, "t", false, "list available tools (shorthand)")
	flag.Parse()

	client := openai.NewClient()
	//agent := chat.NewChatbot(&client, openai.ChatModelGPT4_1Nano)
	agentInstance := agent.NewAgent(&client, openai.ChatModelGPT4_1Nano)

	if *showTools {
		listTools()
		return
	}

	listTools()
	ctx := context.Background()
	for {
		question := readQuestion()
		switch strings.ToLower(question) {
		case "help", "tools", "-t":
			listTools()
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
