package gui

import (
  //"fmt"

  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/container"
  "fyne.io/fyne/v2/widget"
  //"fyne.io/fyne/v2/driver/desktop"
)

func (cw *ChatWindow) openAgentWindow() {
  // 1) If it’s already open just focus it
  if cw.agentWin != nil {
    cw.agentWin.RequestFocus()
    return
  }

  // 2) Create the agent edit window
  w := cw.app.NewWindow("Edit Agent")
  cw.agentWin = w
  w.SetOnClosed(func() { cw.agentWin = nil })

  // 3) Fetch the live agent pointer
  ag := cw.core.Agent()

  // 4) Build Entry widgets for Name & Model
  nameEntry := widget.NewEntry()
  nameEntry.SetText(ag.Name)
  modelEntry := widget.NewEntry()
  modelEntry.SetText(ag.Model)

  // 5) Put them in a Form
  form := &widget.Form{
    Items: []*widget.FormItem{
      {Text: "Name",  Widget: nameEntry},
      {Text: "Model", Widget: modelEntry},
    },
    OnSubmit: func() {
      // runs on UI thread already
      ag.Name  = nameEntry.Text
      ag.Model = modelEntry.Text
      w.Close()
    },
    OnCancel: func() {
      w.Close()
    },
  }

  w.SetContent(container.NewVBox(form))
  w.Resize(fyne.NewSize(400, 200))
  w.Show()
}
