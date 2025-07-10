package main

import (
	"fmt"
	"github.com/johnjallday/dolphin-tool-calling-agent/pkg/tools"
	"github.com/openai/openai-go"
)

const (
		packName		= "Calculator"
    packVersion = "v0.0.1"
    packLink    = "https://github.com/johnjallday/dolphin-tool-calling-agent/"
)

var AddTool = tools.ToolSpec{
	Name:        "add",
	Description: "Add two numbers a and b",
	Parameters: openai.FunctionParameters{
		"type": "object",
		"properties": map[string]interface{}{
			"a": map[string]string{"type": "number"},
			"b": map[string]string{"type": "number"},
		},
		"required": []string{"a", "b"},
	},
	Exec: func(args map[string]interface{}) (string, error) {
		a, ok1 := args["a"].(float64)
		b, ok2 := args["b"].(float64)
		if !ok1 || !ok2 {
			return "", fmt.Errorf("invalid arguments, expected numbers")
		}
		return fmt.Sprintf("%v", a+b), nil
	},
}

var SubtractTool = tools.ToolSpec{
	Name:        "subtract",
	Description: "Subtract b from a",
	Parameters: openai.FunctionParameters{
		"type": "object",
		"properties": map[string]interface{}{
			"a": map[string]string{"type": "number"},
			"b": map[string]string{"type": "number"},
		},
		"required": []string{"a", "b"},
	},
	Exec: func(args map[string]interface{}) (string, error) {
			a, ok1 := args["a"].(float64)
			b, ok2 := args["b"].(float64)
			if !ok1 || !ok2 {
				return "", fmt.Errorf("invalid arguments, expected numbers")
			}
			return fmt.Sprintf("%v", a-b), nil
		},
	}

	var MultiplyTool = tools.ToolSpec{
		Name:        "multiply",
		Description: "Multiply a and b",
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]interface{}{
				"a": map[string]string{"type": "number"},
				"b": map[string]string{"type": "number"},
			},
			"required": []string{"a", "b"},
		},
		Exec: func(args map[string]interface{}) (string, error) {
			a, ok1 := args["a"].(float64)
			b, ok2 := args["b"].(float64)
			if !ok1 || !ok2 {
				return "", fmt.Errorf("invalid arguments, expected numbers")
			}
			return fmt.Sprintf("%v", a*b), nil
		},
	}

	var DivideTool = tools.ToolSpec{
		Name:        "divide",
		Description: "Divide a by b",
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]interface{}{
				"a": map[string]string{"type": "number"},
				"b": map[string]string{"type": "number"},
			},
			"required": []string{"a", "b"},
		},
		Exec: func(args map[string]interface{}) (string, error) {
			a, ok1 := args["a"].(float64)
			b, ok2 := args["b"].(float64)
			if !ok1 || !ok2 {
				return "", fmt.Errorf("invalid arguments, expected numbers")
			}
			if b == 0 {
				return "", fmt.Errorf("division by zero")
			}
			return fmt.Sprintf("%v", a/b), nil
		},
	}


// PluginSpecs is the symbol NewAgentFromConfig will look up and call.
func PluginSpecs() []tools.ToolSpec {
	return []tools.ToolSpec{
		AddTool,
		SubtractTool,
		MultiplyTool,
		DivideTool,
	}
}

func PluginPackage() tools.ToolPackage {
    return tools.ToolPackage{
				Name:		 packName,
        Version: packVersion,
        Link:    packLink,
        Specs:   []tools.ToolSpec{ 
					AddTool, 
					SubtractTool,
					MultiplyTool,
					DivideTool,
				},
    }
}
