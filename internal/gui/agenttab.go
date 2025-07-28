package gui

import (
  "fmt"

  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/container"
  "fyne.io/fyne/v2/dialog"
  "fyne.io/fyne/v2/layout"
  "fyne.io/fyne/v2/widget"

  "github.com/johnjallday/dolphin-tool-calling-agent/internal/app"
)


// ─────────────────────────────────────────────────────────────────────────────
// AGENT TAB
// ─────────────────────────────────────────────────────────────────────────────

func (cw *MainWindow) makeAgentTab() *container.TabItem {
  return container.NewTabItem("Agent", cw.buildAgentPane())
}

func (cw *MainWindow) buildAgentPane() fyne.CanvasObject {
  metas := cw.core.Agents()
  cw.agentList = container.NewVBox()
  cw.agentList.Add(widget.NewLabelWithStyle(
    "Existing Agents", fyne.TextAlignLeading, fyne.TextStyle{Bold: true},
  ))

  if len(metas) == 0 {
    cw.agentList.Add(widget.NewLabel("— none —"))
  } else {
    for _, m := range metas {
      btn := widget.NewButton("Switch to", func(name string) func() {
        return func() {
          if err := cw.core.SwitchAgent(name); err != nil {
            dialog.ShowError(err, cw.wnd)
            return
          }
          cw.refreshUserStatus()
        }
      }(m.Name))
      cw.agentList.Add(container.NewHBox(
        widget.NewLabel(fmt.Sprintf("%s (%s)", m.Name, m.Model)),
        layout.NewSpacer(),
        btn,
      ))
    }
  }

  // AddAgentForm (your existing form)
  form := NewAddAgentForm(
    cw.core.Toolpacks(),
    func(name, model string, tools []string) {
      if name == "" || model == "" {
        dialog.ShowInformation("Missing fields",
          "Please fill Agent Name and Model", cw.wnd)
        return
      }
      meta := app.AgentMeta{Name: name, Model: model, ToolPaths: tools}
      if err := cw.core.CreateAgent(meta); err != nil {
        dialog.ShowError(err, cw.wnd)
        return
      }
      // fully refresh the window
      cw.RefreshAll()
    },
  )

  return container.NewVBox(cw.agentList, widget.NewSeparator(), form)
}
