package gui

import (
  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/container"
  "fyne.io/fyne/v2/widget"
  // "fyne.io/fyne/v2/dialog"

  // "github.com/johnjallday/dolphin-tool-calling-agent/internal/agent"
)

// buildAgentPane constructs just the CanvasObject for the Agent tab.
func (cw *ChatWindow) buildAgentPane() fyne.CanvasObject {
  ag := cw.core.Agent()

  // a) build the Select
  agents := cw.core.Agents()
  names := make([]string, len(agents))
  for i, a := range agents {
    names[i] = a.Name
  }
  cw.agentSelect = widget.NewSelect(names, nil)
  if ag != nil {
    cw.agentSelect.SetSelected(ag.Name)
  }
  cw.agentSelect.OnChanged = func(sel string) {
    cw.core.LoadAgent(sel)
    cw.buildUI()
  }

  // b) if no agent yet, show onboarding
  if ag == nil {
    return container.NewVBox(
      cw.agentSelect,
      widget.NewSeparator(),
      cw.createAgentOnboardingBox(),
    )
  }

  // c) else show the edit form
  cw.agentNameEntry  = widget.NewEntry()
  cw.agentNameEntry.SetText(ag.Name)
  cw.agentModelEntry = widget.NewEntry()
  cw.agentModelEntry.SetText(ag.Model)

  cw.agentForm = &widget.Form{
    Items: []*widget.FormItem{
      {Text: "Name",  Widget: cw.agentNameEntry},
      {Text: "Model", Widget: cw.agentModelEntry},
    },
    OnSubmit: func() {
      ag.Name  = cw.agentNameEntry.Text
      ag.Model = cw.agentModelEntry.Text
      cw.buildUI()                    // redraw tabs
      cw.mainTabs.SelectTabIndex(0)   // switch back to Chat
    },
    OnCancel: func() {
      cw.agentNameEntry.SetText(ag.Name)
      cw.agentModelEntry.SetText(ag.Model)
      cw.mainTabs.SelectTabIndex(0)
    },
  }

  return container.NewVBox(
    cw.agentSelect,
    widget.NewSeparator(),
    cw.agentForm,
  )
}

// makeAgentTab wraps the pane in a TabItem
func (cw *ChatWindow) makeAgentTab() *container.TabItem {
  return container.NewTabItem("Agent", cw.buildAgentPane())
}
