package gui

import (
  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/container"
  "fyne.io/fyne/v2/widget"
  "fyne.io/fyne/v2/dialog"

  "github.com/johnjallday/dolphin-tool-calling-agent/internal/agent"
)

// buildUserPane builds the contents of the “User” tab.
func (cw *ChatWindow) buildUserPane() fyne.CanvasObject {
  usr := cw.core.User()
  if usr == nil {
    // no user yet: show onboarding
    return cw.createOnboardingBox()
  }

  // a) build the form widgets
  cw.userNameEntry = widget.NewEntry()
  cw.userNameEntry.SetText(usr.Name)
  names := make([]string, len(usr.Agents))
  for i, a := range usr.Agents {
    names[i] = a.Name
  }
  cw.userDefaultAgent = widget.NewSelect(names, nil)
  if usr.DefaultAgent != nil {
    cw.userDefaultAgent.SetSelected(usr.DefaultAgent.Name)
  }

  // b) assemble the form
  cw.userForm = &widget.Form{
    Items: []*widget.FormItem{
      {Text: "User Name",     Widget: cw.userNameEntry},
      {Text: "Default Agent", Widget: cw.userDefaultAgent},
    },
    OnSubmit: func() {
      usr.Name = cw.userNameEntry.Text
      if sel := cw.userDefaultAgent.Selected; sel != "" {
        for _, m := range usr.Agents {
          if m.Name == sel {
            a, err := agent.NewAgent(m.Name, m.Model, nil)
            if err != nil {
              dialog.ShowError(err, cw.wnd)
              return
            }
            usr.DefaultAgent = a
          }
        }
      }
      cw.refreshUserStatus()
      cw.mainTabs.SelectTabIndex(0) // back to Chat
    },
    OnCancel: func() {
      cw.userNameEntry.SetText(usr.Name)
      if usr.DefaultAgent != nil {
        cw.userDefaultAgent.SetSelected(usr.DefaultAgent.Name)
      }
      cw.mainTabs.SelectTabIndex(0)
    },
  }

  return container.NewVBox(cw.userForm)
}

// makeUserTab wraps it into a TabItem
func (cw *ChatWindow) makeUserTab() *container.TabItem {
  return container.NewTabItem("User", cw.buildUserPane())
}
