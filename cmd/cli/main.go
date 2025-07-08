package main

import (

	"github.com/johnjallday/dolphin-tool-calling-agent/device"
	"github.com/johnjallday/dolphin-tool-calling-agent/tools"

)


func main() {
	tools.BuildPlugin()
	device.GetAudioDevice()
}
