# 🐬 Dolphin Tool Calling Agent 🐬
This is a simple tool calling agent written in Go.
Often times, we don't really need a complicated AI that does everything for us.
Not all AI fits all. Why not build your own tool?


## Requirement
Go 1.24 or later\n
OpenAI api key

## Install
```bash
git clone https://github.com/johnjallday/dolphin-tool-calling-agent.git
```

## Build
To Build an REPL version
```bash
go build -o dolphin_repl ./cmd/gui
./dolphin_repl
```

Or simply
```bash
go run ./cmd/gui
```
Or just type
```bash
make
```

## Note
This app is at a very early stage. More userfriendly updates will come soon.
In order to setup properly, checkout all the .toml files in the project.

## Usage
For Reaper users, I created simple tools that can read and launch your custom Lua scripts. 
Everyone has a different workflow, so I can’t provide a one-size-fits-all solution. 
However, if you are comfortable writing your own custom scripts, this tool might be very useful.

## Plugins

To checkout how some plugins were built:
github.com/johnjallday/dolphin-tool-calling-agent/examples

```bash
go build -buildmode=plugin -o calculator.so calculator_plugin/plugin.go
create an agent with plugins or tools
```


## Roadmap
-[] GUI version using Fyne
-[] Agent Builder
-[] Python Support
-[] Windows DLL Support
-[] Web3 Integration


## Potential Ideas
-[] Tool Market


## Download Tools
Here are actively managed tools by me:
Dolphin-Reaper-Tools

## Support
coff.ee/johnjallday

