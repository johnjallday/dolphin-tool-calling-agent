package gui


import (
  //"fmt"

  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/container"
  "fyne.io/fyne/v2/widget"

  "fyne.io/fyne/v2/dialog"

	"github.com/johnjallday/dolphin-tool-calling-agent/internal/agent"
  //"fyne.io/fyne/v2/driver/desktop"
)

func (cw *ChatWindow) openUserWindow() {
  if cw.userWin != nil {
    cw.userWin.RequestFocus()
    return
  }

  w := cw.app.NewWindow("Edit User")
  cw.userWin = w
  w.SetOnClosed(func() { cw.userWin = nil })

  usr := cw.core.User()  // *app.User

  // Name field
  nameEntry := widget.NewEntry()
  nameEntry.SetText(usr.Name)

  // Default-agent dropdown
  names := make([]string, len(usr.Agents))
  for i, m := range usr.Agents {
    names[i] = m.Name
  }
  defaultSelect := widget.NewSelect(names, nil)
  if usr.DefaultAgent != nil {
    defaultSelect.SetSelected(usr.DefaultAgent.Name)
  }

  form := &widget.Form{
    Items: []*widget.FormItem{
      {Text: "User Name",     Widget: nameEntry},
      {Text: "Default Agent", Widget: defaultSelect},
    },
    OnSubmit: func() {
      // 1) update the name
      usr.Name = nameEntry.Text

      // 2) if they chose a default, build it
      sel := defaultSelect.Selected
      if sel != "" {
        for _, m := range usr.Agents {
          if m.Name == sel {
            // pass nil for plugins since AgentMeta has none
            a, err := agent.NewAgent(m.Name, m.Model, nil)
            if err != nil {
              dialog.ShowError(err, w)
              return
            }
            usr.DefaultAgent = a
            break
          }
        }
      }

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
