package registry

import (
    "encoding/json"
    "fmt"
    "sort"

    "github.com/openai/openai-go"
    "github.com/johnjallday/dolphin-tool-calling-agent/pkg/tools"
)

type ToolRegistry struct {
    // tools maps the tool‐name to its definition
    tools    map[string]tools.Tool
    // handlers maps the tool‐name to the code that executes it
    handlers map[string]func(openai.ChatCompletionMessageToolCall, *openai.ChatCompletionNewParams)
}

func NewToolRegistry() *ToolRegistry {
    return &ToolRegistry{
        tools:    make(map[string]tools.Tool),
        handlers: make(map[string]func(openai.ChatCompletionMessageToolCall, *openai.ChatCompletionNewParams)),
    }
}

// Initialize registers all tools in this registry into the chat params.
func (r *ToolRegistry) Initialize(params *openai.ChatCompletionNewParams) {
    for _, t := range r.Tools() {
        params.Tools = append(params.Tools, openai.ChatCompletionToolParam{
            Function: openai.FunctionDefinitionParam{
                Name:        t.Name,
                Description: openai.String(t.Description),
                Parameters:  t.Parameters,
            },
        })
    }
}

// Register adds or updates a tool and wires up its handler.
func (r *ToolRegistry) Register(t tools.Tool) {
    r.tools[t.Name] = t

    // overwrite any existing handler for this name
    r.handlers[t.Name] = func(call openai.ChatCompletionMessageToolCall, params *openai.ChatCompletionNewParams) {
        var args map[string]interface{}
        if err := json.Unmarshal([]byte(call.Function.Arguments), &args); err != nil {
            params.Messages = append(params.Messages,
                openai.ToolMessage(fmt.Sprintf("Error parsing arguments: %v", err), call.ID))
            return
        }
        res, err := t.Exec(args)
        if err != nil {
            params.Messages = append(params.Messages,
                openai.ToolMessage(fmt.Sprintf("Error running %s: %v", t.Name, err), call.ID))
            return
        }
        params.Messages = append(params.Messages, openai.ToolMessage(res, call.ID))
    }
}

// Handlers returns the map of function names to handler functions.
func (r *ToolRegistry) Handlers() map[string]func(openai.ChatCompletionMessageToolCall, *openai.ChatCompletionNewParams) {
    return r.handlers
}

// Tools returns a sorted slice of all registered tools.
func (r *ToolRegistry) Tools() []tools.Tool {
    names := make([]string, 0, len(r.tools))
    for name := range r.tools {
        names = append(names, name)
    }
    sort.Strings(names)

    out := make([]tools.Tool, 0, len(names))
    for _, name := range names {
        out = append(out, r.tools[name])
    }
    return out
}

// ListToolNames returns just the names, sorted.
func (r *ToolRegistry) ListToolNames() []string {
    names := make([]string, 0, len(r.tools))
    for name := range r.tools {
        names = append(names, name)
    }
    sort.Strings(names)
    return names
}

// Clear resets the registry to empty.
func (r *ToolRegistry) Clear() {
    r.tools = make(map[string]tools.Tool)
    r.handlers = make(map[string]func(openai.ChatCompletionMessageToolCall, *openai.ChatCompletionNewParams))
}

// String prints a human‐readable list of tools.
func (r *ToolRegistry) String() string {
    ts := r.Tools()
    if len(ts) == 0 {
        return "No tools registered.\n"
    }
    out := "Available tools:\n"
    for _, t := range ts {
        out += fmt.Sprintf(" - %-20s  %s\n", t.Name, t.Description)
    }
    return out
}
