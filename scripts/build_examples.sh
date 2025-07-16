#!/bin/bash

go build -buildmode=plugin -o ./plugins/calculator.so ./plugins/examples/calculator_plugin/
echo "calculator build complete"
go build -buildmode=plugin -o ./plugins/mytool.so ./plugins/examples/mytool
echo "mytool build complete"
go build -buildmode=plugin -o ./plugins/weather.so ./plugins/examples/weather
echo "weather build complete"
go build -buildmode=plugin -o ./plugins/reaper_project_manager/reaper_project_manager.so ./plugins/examples/reaper_project_manager

echo "build complete"

#go build -buildmode=plugin -o calculator.so calculator_plugin/plugin.go
