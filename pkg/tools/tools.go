package tools

import (
//	"encoding/json"
//	"fmt"
	"github.com/openai/openai-go"
)

// Tool holds the schema and executor for a function-calling tool.
type Tool struct {
	Name        string
	Description string
	Parameters  openai.FunctionParameters
	Exec        func(map[string]interface{}) (string, error)
}


type ToolPackage struct {
	Name		string
	Version string
	Link    string
	Tools   []Tool
}
