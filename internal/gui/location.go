package gui

import (
  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/container"
  "fyne.io/fyne/v2/widget"
  "github.com/johnjallday/dolphin-tool-calling-agent/internal/location"
)

func ShowLocation(addToHistory func(string)) {
	if ml, err := location.GetMyLocation(); err != nil {
	addToHistory("Location error: " + err.Error())
	} else {
	addToHistory("Location: " + ml.Name)
	//addToHistory("Audio devices: " + strings.Join(ml.AudioDevices, ", "))
	//addToHistory("Displays: " + strings.Join(ml.Displays, ", "))
	//addToHistory("Networks: " + strings.Join(ml.Network, ", "))
	}
}
// NewLocationView shows location info and lets the user type `/back` to return.
func NewLocationView(showChat func()) fyne.CanvasObject {
  label := widget.NewLabel("Loading location...")
  // populate label via ShowLocation, appending lines
  ShowLocation(func(line string) {
    if label.Text == "" || label.Text == "Loading location..." {
      label.SetText(line)
    } else {
      label.SetText(label.Text + "\n" + line)
    }
  })

  entry := widget.NewEntry()
  entry.SetPlaceHolder("Type `/back` to return")
  entry.OnSubmitted = func(txt string) {
    entry.SetText("")
    if txt == "/back" || txt == "/b" {
      showChat()
    }
  }
  sendBtn := widget.NewButton("Send", func() {
    entry.OnSubmitted(entry.Text)
  })

  return container.NewBorder(
    nil,
    container.NewVBox(entry, sendBtn),
    nil, nil,
    container.NewCenter(label),
  )
}
