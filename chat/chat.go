
package chat

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
	"Dolphin-Tool-Calling-Agent/registry"
)


type Chatbot struct {
	client *openai.Client
	params openai.ChatCompletionNewParams
}



// NewChatbot initializes a new Chatbot with model and temperature settings.
func NewChatbot(client *openai.Client, model string) *Chatbot {
	params := openai.ChatCompletionNewParams{
		Messages:    []openai.ChatCompletionMessageParamUnion{},
		Model:       model,
		Temperature: openai.Float(0),
		Seed:        openai.Int(0),
	}
	
	registry.Register()
	registry.Initialize(&params) // register tools

	return &Chatbot{
		client: client,
		params: params,
	}
}

func (cb *Chatbot) SendMessage(ctx context.Context, userMessage string) (error) {
	// Append user message to conversation
	cb.params.Messages = append(cb.params.Messages, openai.UserMessage(userMessage))

	cmp, err := cb.client.Chat.Completions.New(ctx, cb.params)
	if err != nil {
		return err
	}

	assistantMsg := cmp.Choices[0].Message
	// Append assistant message to conversation
	cb.params.Messages = append(cb.params.Messages, assistantMsg.ToParam())

	fmt.Println(assistantMsg)

	if len(assistantMsg.ToolCalls) == 0 {
		// No tool calls, return assistant content directly
		return nil
	}

	// Dispatch tools and append their responses as messages
	dispatchTools(assistantMsg.ToolCalls, &cb.params)

	// Get final assistant response after tool executions
	finalCmp, err := cb.client.Chat.Completions.New(ctx, cb.params)
	if err != nil {
		return err
	}
	finalMsg := finalCmp.Choices[0].Message
	cb.params.Messages = append(cb.params.Messages, finalMsg.ToParam())

	//return finalMsg.Content, nil
	
	return nil
}

func SendPromptAndReceiveToolCalls(ctx context.Context, client *openai.Client, params *openai.ChatCompletionNewParams) (openai.ChatCompletionMessage, []openai.ChatCompletionMessageToolCall) {
	cmp, err := client.Chat.Completions.New(ctx, *params)
	if err != nil {
		panic(err)
	}
	msg := cmp.Choices[0].Message
	return msg, msg.ToolCalls
}



func dispatchTools(toolCalls []openai.ChatCompletionMessageToolCall, params *openai.ChatCompletionNewParams) {
	//handlers Handlers()
	handlers := registry.Handlers()
	for _, tc := range toolCalls {
		if h, ok := handlers[tc.Function.Name]; ok {
			fmt.Printf("Calling function: %s\n", tc.Function.Name)
			h(tc, params)
		}
	}
}

func runFinalChat(ctx context.Context, client *openai.Client, params *openai.ChatCompletionNewParams) {
	final, err := client.Chat.Completions.New(ctx, *params)
	if err != nil {
		panic(err)
	}
	fmt.Println(final.Choices[0].Message.Content)
}


