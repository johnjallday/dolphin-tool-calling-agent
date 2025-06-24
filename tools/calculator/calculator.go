package calculator

import (
	"fmt"

	"github.com/johnjallday/dolphin-tool-calling-agent/tools"

	"github.com/openai/openai-go"
)

// ToolSpec defines schema and executor for addNumbers.
var ToolSpec = tools.ToolSpec{
	Name:        "addNumbers",
	Description: "Add two numbers",
	Parameters: openai.FunctionParameters{
		"type": "object",
		"properties": map[string]interface{}{
			"a": map[string]string{"type": "number"},
			"b": map[string]string{"type": "number"},
		},
		"required": []string{"a", "b"},
	},
	Exec: func(args map[string]interface{}) (string, error) {
		a := args["a"].(float64)
		b := args["b"].(float64)
		return fmt.Sprintf("%v", a+b), nil
	},
}

// MultiplySpec defines schema and executor for multiply.
var MultiplySpec = tools.ToolSpec{
	Name:        "multiply",
	Description: "Multiply two numbers",
	Parameters: openai.FunctionParameters{
		"type": "object",
		"properties": map[string]interface{}{
			"a": map[string]string{"type": "number"},
			"b": map[string]string{"type": "number"},
		},
		"required": []string{"a", "b"},
	},
	Exec: func(args map[string]interface{}) (string, error) {
		a := args["a"].(float64)
		b := args["b"].(float64)
		return fmt.Sprintf("%v", a*b), nil
	},
}

