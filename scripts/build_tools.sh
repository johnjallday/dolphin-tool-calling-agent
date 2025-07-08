#!/bin/bash

go build -buildmode=plugin -o ./user/tools/calculator.so ./examples/calculator_plugin/
go build -buildmode=plugin -o ./user/tools/mytool.so ./examples/mytool
go build -buildmode=plugin -o ./user/tools/reaper_tools.so ./examples/reaper_tools
go build -buildmode=plugin -o ./user/tools/weather.so ./examples/weather
go build -buildmode=plugin -o ./user/tools/reaper_project_manager/reaper_project_manager.so ./examples/reaper_project_manager

#go build -buildmode=plugin -o calculator.so calculator_plugin/plugin.go
