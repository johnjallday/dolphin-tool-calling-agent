package registry

import (

	"github.com/openai/openai-go"

	"fmt"
	"encoding/json"

	"github.com/johnjallday/dolphin-tool-calling-agent/pkg/tools"
)

var (
	// specs accumulates all registered Tool.
	specs []tools.Tool
	// handlers maps tool names to executor wrappers.
	handlers = make(map[string]func(openai.ChatCompletionMessageToolCall, *openai.ChatCompletionNewParams))
)

// Initialize registers all Tool into the ChatCompletion params.
func Initialize(params *openai.ChatCompletionNewParams) {
	for _, ts := range specs {
		params.Tools = append(params.Tools, openai.ChatCompletionToolParam{
			Function: openai.FunctionDefinitionParam{
				Name:        ts.Name,
				Description: openai.String(ts.Description),
				Parameters:  ts.Parameters,
			},
		})
	}
}

func RegisterSpec(ts tools.Tool) {
	specs = append(specs, ts)
	handlers[ts.Name] = func(tc openai.ChatCompletionMessageToolCall, params *openai.ChatCompletionNewParams) {
		var args map[string]interface{}
		if err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err != nil {
			params.Messages = append(params.Messages, openai.ToolMessage(fmt.Sprintf("Error: %v", err), tc.ID))
			return
		}
		res, err := ts.Exec(args)
		if err != nil {
			params.Messages = append(params.Messages, openai.ToolMessage(fmt.Sprintf("Error: %v", err), tc.ID))
			return
		}

		params.Messages = append(params.Messages, openai.ToolMessage(res, tc.ID))
	}
}

// Handlers returns the mapping of function names to handler functions.
func Handlers() map[string]func(openai.ChatCompletionMessageToolCall, *openai.ChatCompletionNewParams) {
	return handlers
}


func Clear() {
	specs = nil
	handlers = make(map[string]func(openai.ChatCompletionMessageToolCall, *openai.ChatCompletionNewParams))
}
