package registry

import (
    "encoding/json"
    "fmt"

    "github.com/openai/openai-go"
    "github.com/johnjallday/dolphin-tool-calling-agent/pkg/tools"
)

var (
    // registeredTools accumulates all registered Tool.
    registeredTools []tools.Tool
    // handlers maps tool names to executor wrappers.
    handlers = make(map[string]func(openai.ChatCompletionMessageToolCall, *openai.ChatCompletionNewParams))
)

// Initialize registers all Tool into the ChatCompletion params.
func Initialize(params *openai.ChatCompletionNewParams) {
    for _, t := range registeredTools {
        params.Tools = append(params.Tools, openai.ChatCompletionToolParam{
            Function: openai.FunctionDefinitionParam{
                Name:        t.Name,
                Description: openai.String(t.Description),
                Parameters:  t.Parameters,
            },
        })
    }
}

// RegisterTool adds a Tool and its handler.
func RegisterTool(t tools.Tool) {
    registeredTools = append(registeredTools, t)
    handlers[t.Name] = func(call openai.ChatCompletionMessageToolCall, params *openai.ChatCompletionNewParams) {
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

// Handlers returns the mapping of function names to handler functions.
func Handlers() map[string]func(openai.ChatCompletionMessageToolCall, *openai.ChatCompletionNewParams) {
    return handlers
}

// Clear resets the registry.
func Clear() {
    registeredTools = nil
    handlers = make(map[string]func(openai.ChatCompletionMessageToolCall, *openai.ChatCompletionNewParams))
}
