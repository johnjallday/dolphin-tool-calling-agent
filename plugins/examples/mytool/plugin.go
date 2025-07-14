package main

import (
		"github.com/johnjallday/dolphin-tool-calling-agent/pkg/tools"
		"github.com/openai/openai-go"
	)


const (
		packName		= "Greeting"
    packVersion = "v0.0.1"
    packLink    = "https://github.com/johnjallday/dolphin-tool-calling-agent/"
)

	// Define your Tool(s)
var HelloTool = tools.Tool{
	Name:        "say_hello",
	Description: "Returns a greeting",
	Parameters:  openai.FunctionParameters{ "type":"object",
"properties":map[string]interface{}{} },
	Exec: func(args map[string]interface{}) (string,error) {
		return "👋 Hello from plugin!", nil
	},
}

// Export a function called PluginSpecs
func PluginSpecs() []tools.Tool {
	return []tools.Tool{ HelloTool }
}

func PluginPackage() tools.ToolPackage {
    return tools.ToolPackage{
				Name:		 packName,
        Version: packVersion,
        Link:    packLink,
        Tools:   []tools.Tool{ HelloTool },
    }
}
