#!/bin/bash

go build -buildmode=plugin -o ./configs/user/tools/calculator.so ./plugins/examples/calculator_plugin/
echo "calculator build complete"
go build -buildmode=plugin -o ./configs/user/tools/mytool.so ./plugins/examples/mytool
echo "mytool build complete"
go build -buildmode=plugin -o ./configs/user/tools/reaper_tools.so ./plugins/examples/reaper_tools
echo "reaper_tool build complete"
go build -buildmode=plugin -o ./configs/user/tools/weather.so ./plugins/examples/weather
echo "weather build complete"
go build -buildmode=plugin -o ./configs/user/tools/reaper_project_manager/reaper_project_manager.so ./plugins/examples/reaper_project_manager

echo "build complete"

#go build -buildmode=plugin -o calculator.so calculator_plugin/plugin.go
