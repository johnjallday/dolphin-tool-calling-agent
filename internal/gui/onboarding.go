package gui

import (
  "fmt"
  "strings"

  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/container"
  "fyne.io/fyne/v2/dialog"
  "fyne.io/fyne/v2/widget"

  "github.com/johnjallday/dolphin-tool-calling-agent/internal/app"
)

// createOnboardingBox shows when there are zero users.
func (cw *MainWindow) createOnboardingBox() fyne.CanvasObject {
  nameEntry := widget.NewEntry()
  nameEntry.SetPlaceHolder("Enter username")

  createBtn := widget.NewButton("Create User", func() {
    userID := strings.TrimSpace(nameEntry.Text)
    if userID == "" {
      dialog.ShowError(fmt.Errorf("username cannot be empty"), cw.wnd)
      return
    }
    // 1) create the user
    if err := cw.core.CreateUser(userID); err != nil {
      dialog.ShowError(err, cw.wnd)
      return
    }
    // 2) set as default (loads it, too)
    if err := cw.core.SetDefaultUser(userID); err != nil {
      dialog.ShowError(fmt.Errorf("could not set default user: %w", err), cw.wnd)
      return
    }
    // 3) now re‐draw everything
    cw.RefreshAll()
    // go back to “Chat” tab
    cw.mainTabs.SelectTabIndex(0)
  })

  return container.NewVBox(
    widget.NewLabelWithStyle("Welcome to Dolphin Chat!",
      fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
    widget.NewLabel("Please pick a username to get started:"),
    nameEntry,
    createBtn,
  )
}

// createAgentOnboardingBox shows when we have a user but no agent. 
func (cw *MainWindow) createAgentOnboardingBox() fyne.CanvasObject {
  nameEntry := widget.NewEntry()
  nameEntry.SetPlaceHolder("Enter agent name")

  createBtn := widget.NewButton("Create Agent", func() {
    agentName := strings.TrimSpace(nameEntry.Text)
    if agentName == "" {
      dialog.ShowError(fmt.Errorf("agent name cannot be empty"), cw.wnd)
      return
    }

    meta := app.AgentMeta{
      Name:      agentName,
      Model:     "gpt-4.1-nano",
      ToolPaths: nil,
    }

    // 1) create the agent
    if err := cw.core.CreateAgent(meta); err != nil {
      dialog.ShowError(err, cw.wnd)
      return
    }
    // 2) set as default & load
    if err := cw.core.SetDefaultAgent(agentName); err != nil {
      dialog.ShowError(fmt.Errorf("could not set default agent: %w", err), cw.wnd)
      return
    }
    // 3) re‐draw everything
    cw.RefreshAll()
    // back to “Chat” tab
    cw.mainTabs.SelectTabIndex(0)
  })

  return container.NewVBox(
    widget.NewLabelWithStyle("🎉 Welcome to Dolphin Chat! 🐬",
      fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
    widget.NewLabel("Let’s create your first agent:"),
    nameEntry,
    createBtn,
  )
}
