package tools

import (
//	"encoding/json"
//	"fmt"
	"github.com/openai/openai-go"
)

// ToolSpec holds the schema and executor for a function-calling tool.
type ToolSpec struct {
	Name        string
	Description string
	Parameters  openai.FunctionParameters
	Exec        func(map[string]interface{}) (string, error)
}
