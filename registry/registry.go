package registry

import (

	"github.com/openai/openai-go"

	"fmt"
	"encoding/json"

	"github.com/johnjallday/dolphin-tool-calling-agent/tools"
	"github.com/johnjallday/dolphin-tool-calling-agent/tools/calculator"
	"github.com/johnjallday/dolphin-tool-calling-agent/tools/reaper"
	"github.com/johnjallday/dolphin-tool-calling-agent/tools/weather"
)

var (
	// specs accumulates all registered ToolSpecs.
	specs []tools.ToolSpec
	// handlers maps tool names to executor wrappers.
	handlers = make(map[string]func(openai.ChatCompletionMessageToolCall, *openai.ChatCompletionNewParams))
)

// Initialize registers all ToolSpecs into the ChatCompletion params.
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
func RegisterSpec(ts tools.ToolSpec) {
	specs = append(specs, ts)
	handlers[ts.Name] = func(tc openai.ChatCompletionMessageToolCall, params *openai.ChatCompletionNewParams) {
		var args map[string]interface{}
		if err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err != nil {
			panic(err)
		}
		res, err := ts.Exec(args)
		if err != nil {
			panic(err)
		}
		fmt.Println(res)
		params.Messages = append(params.Messages, openai.ToolMessage(res, tc.ID))
	}
}

func Register() {
	RegisterSpec(calculator.ToolSpec)
	RegisterSpec(calculator.MultiplySpec)
	RegisterSpec(reaper.CreateNewProjectTool)
	RegisterSpec(weather.WeatherTool)
	for _, spec := range reaper.ReaperCustomScriptsSpecs {
    RegisterSpec(spec)
	}
}

// Handlers returns the mapping of function names to handler functions.
func Handlers() map[string]func(openai.ChatCompletionMessageToolCall, *openai.ChatCompletionNewParams) {
	return handlers
}

func Specs() []tools.ToolSpec {
	return specs
}
