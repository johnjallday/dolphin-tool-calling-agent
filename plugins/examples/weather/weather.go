package main

import (
		"github.com/johnjallday/dolphin-tool-calling-agent/pkg/tools"
		"github.com/openai/openai-go"
	)

const (
		packName		= "Weather"
    packVersion = "v0.0.1"
    packLink    = "https://github.com/johnjallday/dolphin-tool-calling-agent/"
)

var WeatherTool = tools.Tool{
	Name:        "get_weather",
	Description: "Get weather at the given location",
	Parameters: openai.FunctionParameters{
		"type": "object",
		"properties": map[string]interface{}{
			"location": map[string]string{"type": "string"},
		},
		"required": []string{"location"},
	},
	Exec: func(args map[string]interface{}) (string, error) {
		loc := args["location"].(string)
		return getWeather(loc), nil
	},
}


func getWeather(location string) string {
	return "Sunny, 25°C"
}

//exposes the plugin to dolphin-tool-calling-agent
func PluginPackage() tools.ToolPackage {
    return tools.ToolPackage{
				Name:		 packName,
        Version: packVersion,
        Link:    packLink,
        Tools:   []tools.Tool{ WeatherTool },
    }
}
