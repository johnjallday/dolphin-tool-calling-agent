package main

import (
  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/app"
  "fyne.io/fyne/v2/container"
  "fyne.io/fyne/v2/driver/desktop"
  //"fyne.io/fyne/v2/widget"

  //"github.com/johnjallday/dolphin-tool-calling-agent/internal/agent"
  "github.com/johnjallday/dolphin-tool-calling-agent/internal/gui"
)

func main() {
  a := app.New()
  w := a.NewWindow("Dolphin Tool")

  var chatView, agentView, helpView, locationView fyne.CanvasObject

  showChat := func() {
    chatView.Show(); agentView.Hide(); helpView.Hide(); locationView.Hide()
  }
  showHelp := func() {
    helpView.Show(); chatView.Hide(); agentView.Hide(); locationView.Hide()
  }
  showLocation := func() {
    locationView.Show(); chatView.Hide(); agentView.Hide(); helpView.Hide()
  }

  // use the package‐level history adder
  addToHistory := gui.AddToHistory

  chatView, _ = gui.NewChatView(showAgent, showHelp, showLocation)

  //configs := []agent.AgentConfig{
  //  {Name: "Agent A"},
  //  {Name: "Agent B"},
  //}
  //agentView, _ = gui.NewAgentView(addToHistory, showChat)
  helpView = gui.NewHelpView(w, showChat)
	locationView = gui.NewLocationView(showChat)

  stack := container.NewMax(chatView, agentView, helpView, locationView)
  w.SetContent(stack)
  w.Resize(fyne.NewSize(600, 400))

  w.Canvas().AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyTab}, func(_ fyne.Shortcut) {
    showChat()
  })

  showChat()
  w.ShowAndRun()
}
