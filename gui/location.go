package gui

import (
  "strings"
  "github.com/johnjallday/dolphin-tool-calling-agent/location"
)

func ShowLocation(addToHistory func(string)) {
  if ml, err := location.GetMyLocation(); err != nil {
    addToHistory("Location error: " + err.Error())
  } else {
    addToHistory("Location: " + ml.Name)
    addToHistory("Audio devices: " + strings.Join(ml.AudioDevices, ", "))
    addToHistory("Displays: " + strings.Join(ml.Displays, ", "))
    addToHistory("Networks: " + strings.Join(ml.Network, ", "))
  }
}
