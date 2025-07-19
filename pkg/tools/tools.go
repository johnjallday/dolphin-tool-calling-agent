package tools

import (
	"fmt"
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

func (tp ToolPackage) String() string {
	return fmt.Sprintf("ToolPackage: %s (v%s)\nLink: %s\nTools:\n%s",
		tp.Name, tp.Version, tp.Link, tp.listTools())
}

func (tp ToolPackage) listTools() string {
	if len(tp.Tools) == 0 {
		return "  (none)"
	}
	s := ""
	for _, t := range tp.Tools {
		s += fmt.Sprintf("  - %s: %s\n", t.Name, t.Description)
	}
	return s
}
