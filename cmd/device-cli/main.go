package main

import (
	"fmt"
	"log"
	"github.com/johnjallday/dolphin-tool-calling-agent/device"
)

func main(){
	//device.GetCurrentAudioDevice()
	myDevices, err := device.GetAudioOutputDevices()
	if err!= nil {
		log.Fatal(err)
	}
	device.PrintDevices(myDevices)
	
	network, err := device.GetMyCurrentNetwork()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("SSID: %s\nBSSID: %s\n", network.SSID, network.BSSID)
}
