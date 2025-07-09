package main

import (
	"fmt"
	//"log"
	"github.com/johnjallday/dolphin-tool-calling-agent/device"
	"github.com/johnjallday/dolphin-tool-calling-agent/location"
)

func main(){


	location.GetMyLocation()
	display := device.GetCurrentDisplay()
	fmt.Println(display.Name)
}
