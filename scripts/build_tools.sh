#!/bin/bash

go build -buildmode=plugin -o ./user/tools/calculator.so ./examples/calculator_plugin/
echo "calculator build complete"
go build -buildmode=plugin -o ./user/tools/mytool.so ./examples/mytool
echo "mytool build complete"
go build -buildmode=plugin -o ./user/tools/reaper_tools.so ./examples/reaper_tools
echo "reaper_tool build complete"
go build -buildmode=plugin -o ./user/tools/weather.so ./examples/weather
echo "weather build complete"
go build -buildmode=plugin -o ./user/tools/reaper_project_manager/reaper_project_manager.so ./examples/reaper_project_manager

echo "build complete"

#go build -buildmode=plugin -o calculator.so calculator_plugin/plugin.go
