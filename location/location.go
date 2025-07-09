package location

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/johnjallday/dolphin-tool-calling-agent/device"
)

type Location struct {
	Name         string   `toml:"name"`
	AudioDevices []string `toml:"audio_devices"`
	Displays     []string `toml:"displays"`
	Network      []string `toml:"network"`
}

type LocationsConfig struct {
	Locations []Location `toml:"locations"`
}

// Load locations from TOML file
func LoadLocations(path string) ([]Location, error) {
	var config LocationsConfig
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, err
	}
	return config.Locations, nil
}

func GetMyLocation() {
	_, currentAudioOutput := device.GetCurrentAudioDevice()
	currentNetwork, err := device.GetMyCurrentNetwork()
	if err != nil {
		fmt.Println("Failed to get current network:", err)
		return
	}
	currentDisplay := device.GetCurrentDisplay() // returns struct

	locations, err := LoadLocations("./user/locations.toml")
	if err != nil {
		fmt.Println("Failed to load locations:", err)
		return
	}

	for _, loc := range locations {
		matchAudio := contains(loc.AudioDevices, currentAudioOutput)

		matchDisplay := false
		if currentDisplay != nil {
			matchDisplay = contains(loc.Displays, currentDisplay.Name)
		}

		matchNetwork := false
		if currentNetwork != nil {
			matchNetwork = contains(loc.Network, currentNetwork.SSID)
		}

		if matchAudio || matchDisplay || matchNetwork {
			fmt.Printf("Current location: %s\n", loc.Name)
			return
		}
	}

	fmt.Println("Location not recognized.")
}

// Helper function to check if slice contains a value
func contains(slice []string, val string) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
}
