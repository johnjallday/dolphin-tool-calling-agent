package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/openai/openai-go"

	"Dolphin-Tool-Calling-Agent/chat"
	"Dolphin-Tool-Calling-Agent/tools"
)

func main() {
	showTools := flag.Bool("tools", false, "list available tools")
	flag.Parse()

	if *showTools {
		for _, ts := range tools.Specs() {
			fmt.Printf("%s: %s\n", ts.Name, ts.Description)
		}
		return
	}

	client := openai.NewClient()
	chatClient := chat.NewChatClient(&client)

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Dolphin CLI Chatbot. Type 'exit' or Ctrl+C to quit.")
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break // EOF or error
		}
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			continue
		}
		if strings.EqualFold(text, "exit") || strings.EqualFold(text, "quit") {
			fmt.Println("Bye!")
			break
		}
		if strings.EqualFold(text, "tools") {
			for _, ts := range tools.Specs() {
				fmt.Printf("%s: %s\n", ts.Name, ts.Description)
			}
			continue
		}

		resp, err := chatClient.HandleQuestion(context.Background(), text)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		fmt.Println(resp)
	}
}
