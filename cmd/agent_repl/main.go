package main

import(
	"fmt"
	"github.com/johnjallday/dolphin-tool-calling-agent/agents"
)

func main() {
	fmt.Println("Hello")
	agentInstance := agent.NewAgent(&client, openai.ChatModelGPT4_1Nano)
}
