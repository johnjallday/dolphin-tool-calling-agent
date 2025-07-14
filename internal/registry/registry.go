package registry

import (
    "encoding/json"
    "fmt"

    "github.com/openai/openai-go"
    "github.com/johnjallday/dolphin-tool-calling-agent/pkg/tools"
)

type ToolRegistry struct {
    Tools    []tools.Tool
    handlers map[string]func(openai.ChatCompletionMessageToolCall, *openai.ChatCompletionNewParams)
}

func NewToolRegistry() *ToolRegistry {
    return &ToolRegistry{
        handlers: make(map[string]func(openai.ChatCompletionMessageToolCall, *openai.ChatCompletionNewParams)),
    }
}

// Initialize registers all tools in this registry into the chat params.
func (r *ToolRegistry) Initialize(params *openai.ChatCompletionNewParams) {
    for _, t := range r.Tools {
        params.Tools = append(params.Tools, openai.ChatCompletionToolParam{
            Function: openai.FunctionDefinitionParam{
                Name:        t.Name,
                Description: openai.String(t.Description),
                Parameters:  t.Parameters,
            },
        })
    }
}

// Register adds a tool and its execution handler.
func (r *ToolRegistry) Register(t tools.Tool) {
    r.Tools = append(r.Tools, t)
    r.handlers[t.Name] = func(call openai.ChatCompletionMessageToolCall, params *openai.ChatCompletionNewParams) {
        var args map[string]interface{}
        if err := json.Unmarshal([]byte(call.Function.Arguments), &args); err != nil {
            params.Messages = append(params.Messages, openai.ToolMessage(fmt.Sprintf("Error: %v", err), call.ID))
            return
        }
        res, err := t.Exec(args)
        if err != nil {
            params.Messages = append(params.Messages, openai.ToolMessage(fmt.Sprintf("Error: %v", err), call.ID))
            return
        }
        params.Messages = append(params.Messages, openai.ToolMessage(res, call.ID))
    }
}

// Handlers returns the map of function names to handler functions.
func (r *ToolRegistry) Handlers() map[string]func(openai.ChatCompletionMessageToolCall, *openai.ChatCompletionNewParams) {
    return r.handlers
}

// Clear resets the registry to empty.
func (r *ToolRegistry) Clear() {
    r.Tools = nil
    r.handlers = make(map[string]func(openai.ChatCompletionMessageToolCall, *openai.ChatCompletionNewParams))
}

func (r *ToolRegistry) PrintTools() {
    if len(r.Tools) == 0 {
        fmt.Println("No tools registered.")
        return
    }
    fmt.Println("Available tools:")
    for _, t := range r.Tools {
        fmt.Printf(" - %s: %s\n", t.Name, t.Description)
        // If you want to show parameters schema:
        // fmt.Printf("   parameters: %v\n", t.Parameters)
    }
}
