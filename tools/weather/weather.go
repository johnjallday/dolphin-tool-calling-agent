package weather

import (
	//"encoding/json"
	//"fmt"

	"Dolphin-Tool-Calling-Agent/tools"
	"github.com/openai/openai-go"
)

// ToolSpec defines schema and executor for get_weather.
var WeatherTool = tools.ToolSpec{
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

// getWeather is a mock for demonstration.
func getWeather(location string) string {
	return "Sunny, 25°C"
}

