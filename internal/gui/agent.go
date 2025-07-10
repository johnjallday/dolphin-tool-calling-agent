package gui

import (
  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/container"
  "fyne.io/fyne/v2/widget"
  "github.com/johnjallday/dolphin-tool-calling-agent/internal/agent"
)

// NewAgentView lists all agents and lets you go back with `/back`.
func NewAgentView(showChat func(), addToHistory func(string)) fyne.CanvasObject {
  listBox := container.NewVBox()
  configs, err := agent.ListAgents()
  if err != nil {
    listBox.Add(widget.NewLabel("Error listing agents: " + err.Error()))
  } else if len(configs) == 0 {
    listBox.Add(widget.NewLabel("No agents found."))
  } else {
    for _, cfg := range configs {
      btn := widget.NewButton(cfg.Name, func(name string) func() {
        return func() {
          addToHistory("Selected agent: " + name)
          showChat()
        }
      }(cfg.Name))
      listBox.Add(btn)
    }
  }

  entry := widget.NewEntry()
  entry.SetPlaceHolder("Type `/back` to return")
  entry.OnSubmitted = func(txt string) {
    entry.SetText("")
    if txt == "/back" {
      showChat()
    }
  }
  sendBtn := widget.NewButton("Send", func() {
    entry.OnSubmitted(entry.Text)
  })

  return container.NewBorder(
    nil,
    container.NewHBox(entry, sendBtn),
    nil, nil,
    container.NewScroll(listBox),
  )
}
