package gui

import (
  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/widget"
  "fyne.io/fyne/v2/dialog"
  "fyne.io/fyne/v2/container"

	"strings"
	"fmt"

  "github.com/johnjallday/dolphin-tool-calling-agent/internal/app"
)


func (cw *ChatWindow) createOnboardingBox() *fyne.Container {
  nameEntry := widget.NewEntry()
  nameEntry.SetPlaceHolder("Enter username")

  createBtn := widget.NewButton("Create User", func() {
    userID := strings.TrimSpace(nameEntry.Text)
    if userID == "" {
      dialog.ShowError(fmt.Errorf("username cannot be empty"), cw.wnd)
      return
    }
    // 1) create the user on-disk/in-memory
    if err := cw.core.CreateUser(userID); err != nil {
      dialog.ShowError(err, cw.wnd)
      return
    }
    // 2) mark them as the default user (this also calls LoadUser under the hood)
    if err := cw.core.SetDefaultUser(userID); err != nil {
      dialog.ShowError(fmt.Errorf("could not set default user: %w", err), cw.wnd)
      return
    }
    // 3) rebuild the whole UI now that we have at least one user
    cw.buildUI()
  })

  return container.NewVBox(
    widget.NewLabelWithStyle("Welcome to Dolphin Chat!",
      fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
    widget.NewLabel("Please pick a username to get started:"),
    nameEntry,
    createBtn,
  )
}


func (cw *ChatWindow) createAgentOnboardingBox() *fyne.Container {
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

    // 1) append it to the user’s TOML
    if err := cw.core.CreateAgent(meta); err != nil {
      dialog.ShowError(err, cw.wnd)
      return
    }

    // 2) persist & load as the new default agent
    if err := cw.core.SetDefaultAgent(agentName); err != nil {
      dialog.ShowError(fmt.Errorf("could not set default agent: %w", err), cw.wnd)
      return
    }

    // 3) now rebuild the UI into the real chat view
    cw.buildUI()
  })

  return container.NewVBox(
    widget.NewLabelWithStyle(
      "🎉 Welcome to Dolphin Chat! 🐬",
      fyne.TextAlignCenter, fyne.TextStyle{Bold: true},
    ),
    widget.NewLabel("Let’s create your first agent:"),
    nameEntry,
    createBtn,
  )
}
 
