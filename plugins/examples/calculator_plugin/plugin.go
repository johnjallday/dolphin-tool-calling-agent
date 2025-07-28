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

var AddTool = tools.Tool{
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

var SubtractTool = tools.Tool{
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

	var MultiplyTool = tools.Tool{
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

	var DivideTool = tools.Tool{
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




func PluginPackage() tools.ToolPackage {
    return tools.ToolPackage{
				Name:		 packName,
        Version: packVersion,
        Link:    packLink,
				Description: "Sample Calculator plugin",
        Tools:   []tools.Tool{ 
					AddTool, 
					SubtractTool,
					MultiplyTool,
					DivideTool,
				},
    }
}
