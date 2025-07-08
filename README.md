# 🐬Dolphin Tool Calling Agent
This is a simple tool calling agent written in Go.
If you want something compiled, then this is a way to go.
Often times we don't really need a complicated AI that does everything for us.
You want an AI agent that can help you with specific tasks.

## Requirement
Go 1.24 or later
OpenAI api key

## Install
```bash
git clone https://github.com/johnjallday/dolphin-tool-calling-agent.git
```

## Build
To Build an REPL version
```bash
go build -o dolphin_repl ./cmd/repl
./dolphin_repl
```

Or simply
```bash
go run ./cmd/repl
```

## Usage
For Reaper users, I created simple tools that can read and launch your custom Lua scripts. 
Everyone has a different workflow, so I can’t provide a one-size-fits-all solution. 
However, if you are comfortable writing your own custom scripts, this tool might be very useful.

## Plugins

go build -buildmode=plugin -o calculator.so calculator_plugin/plugin.go
create an agent with plugins or tools

## Roadmap
-[] GUI version using Fyne
-[] Agent Building UX
-[] Python Support


## Support
coff.ee/johnjallday

