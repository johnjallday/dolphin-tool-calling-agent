package agent

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/johnjallday/dolphin-tool-calling-agent/registry"
)

// Agent defines the methods any agent must implement. 
type Agent interface { 
	// SendMessage sends a user message and processes the conversation. 
	SendMessage(ctx context.Context, userMessage string) error 
}

// DefaultAgent is a concrete implementation of Agent. 
type DefaultAgent struct { 
	client *openai.Client 
	params openai.ChatCompletionNewParams 
}

// NewAgent creates and returns a new DefaultAgent. 
func NewAgent(client *openai.Client, model string) Agent { 
	params := openai.ChatCompletionNewParams{ 
		Messages: []openai.ChatCompletionMessageParamUnion{}, 
		Model: model, 
		Temperature: openai.Float(0), 
		Seed: openai.Int(0), 
	}

	registry.Register()
	registry.Initialize(&params)
	//fmt.Println(params)

	return &DefaultAgent{
		client: client,
		params: params,
	}
}

// SendMessage appends the user message, processes the chat response, 
// dispatches tool calls if any, and appends the final response. 
func (a *DefaultAgent) SendMessage(ctx context.Context, userMessage string) error { 
	// Append the user’s message to the conversation. 
	a.params.Messages = append(a.params.Messages, openai.UserMessage(userMessage))

	// Get the assistant's response.
	cmp, err := a.client.Chat.Completions.New(ctx, a.params)
	if err != nil {
		return err
	}

	assistantMsg := cmp.Choices[0].Message
	a.params.Messages = append(a.params.Messages, assistantMsg.ToParam())

	// If there are no tool calls, exit early.
	if len(assistantMsg.ToolCalls) == 0 {
		return nil
	}

	// Dispatch tool calls and update the conversation with their responses.
	dispatchTools(assistantMsg.ToolCalls, &a.params)

	// Get the final assistant response after tools have been executed.
	//finalCmp, err := a.client.Chat.Completions.New(ctx, a.params)
	//if err != nil {
	//	return err
	//}
	//finalMsg := finalCmp.Choices[0].Message
	//a.params.Messages = append(a.params.Messages, finalMsg.ToParam())

	return nil
}

// SendPromptAndReceiveToolCalls sends a prompt and returns the message along with any tool calls. 
func SendPromptAndReceiveToolCalls(ctx context.Context, client *openai.Client, params *openai.ChatCompletionNewParams) (openai.ChatCompletionMessage, []openai.ChatCompletionMessageToolCall) { 
	cmp, err := client.Chat.Completions.New(ctx, *params) 
	if err != nil { 
		panic(err) 
	} 
	msg := cmp.Choices[0].Message 
	return msg, msg.ToolCalls 
}

// dispatchTools processes any tool calls by dispatching them to the registered handlers. 
func dispatchTools(toolCalls []openai.ChatCompletionMessageToolCall, params *openai.ChatCompletionNewParams) { 
	handlers := registry.Handlers() 
	for _, tc := range toolCalls { 
		fmt.Println(tc.Function)
		if h, ok := handlers[tc.Function.Name]; ok { 
			fmt.Printf("Calling function: %s\n", tc.Function.Name) 
			h(tc, params) 
			//need logic for tools with input variables
		} 
	} 
}

// runFinalChat executes the final chat call and prints the assistant’s response. 
func runFinalChat(ctx context.Context, client *openai.Client, params *openai.ChatCompletionNewParams) { 
	final, err := client.Chat.Completions.New(ctx, *params) 
	if err != nil { 
		panic(err) 
	} 
	fmt.Println(final.Choices[0].Message.Content) 
}
