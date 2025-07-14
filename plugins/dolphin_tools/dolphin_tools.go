package main

import (
		"github.com/johnjallday/dolphin-tool-calling-agent/pkg/tools"
		"github.com/openai/openai-go"
	)

const (
	packName		= "Dolphin Tools"
	packVersion = "v0.0.1"
	packlink		= "github.com/johnjallday/dolphin-tool-calling-agent/plugins/dolphin_tools"
)

var CreateAgent = tools.ToolSpec{
	Name:				""
}

"cli agent"
"call functions related to loading and unloading agents and tools"


"Initialize"
"load App Config"
"Load UserConfig"

"If DefaultAgent is in UserConfig"
"Load DefaultAgent"
"Load Tools"


