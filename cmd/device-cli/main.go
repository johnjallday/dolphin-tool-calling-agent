package main

import (
	"fmt"
	"log"
	"github.com/johnjallday/dolphin-tool-calling-agent/device"
)

func main(){
	fmt.Println("test")
	//device.GetCurrentAudioDevice()
	myDevices, err := device.GetAudioOutputDevices()
	if err!= nil {
		log.Fatal(err)
	}
	device.PrintDevices(myDevices)
	
}
