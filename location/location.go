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

func LoadLocations(path string) ([]Location, error) {
    var config LocationsConfig
    if _, err := toml.DecodeFile(path, &config); err != nil {
        return nil, err
    }
    return config.Locations, nil
}

// GetMyLocation returns the matching Location or an error if none found.
func GetMyLocation() (*Location, error) {
    _, currentAudio := device.GetCurrentAudioDevice()

    currentNetwork, err := device.GetMyCurrentNetwork()
    if err != nil {
        return nil, fmt.Errorf("get network: %w", err)
    }

    currentDisplay := device.GetCurrentDisplay()

    locations, err := LoadLocations("./user/locations.toml")
    if err != nil {
        return nil, fmt.Errorf("load locations: %w", err)
    }

    for _, loc := range locations {
        if contains(loc.AudioDevices, currentAudio) ||
            (currentDisplay != nil && contains(loc.Displays, currentDisplay.Name)) ||
            (currentNetwork != nil && contains(loc.Network, currentNetwork.SSID)) {
            return &loc, nil
        }
    }

    return nil, fmt.Errorf("location not recognized")
}

func contains(slice []string, val string) bool {
    for _, s := range slice {
        if s == val {
            return true
        }
    }
    return false
}
