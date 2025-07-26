package gui

import (
  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/container"
  "fyne.io/fyne/v2/dialog"
  "fyne.io/fyne/v2/layout"
  "fyne.io/fyne/v2/widget"

  "github.com/johnjallday/dolphin-tool-calling-agent/internal/app"
)

// buildUserPane builds the contents of the “User” tab.
func (cw *ChatWindow) buildUserPane() fyne.CanvasObject {
  // if no user loaded, offer to select or create
  if cw.core.User() == nil {
    users := cw.core.Users()
    cw.userSelect = widget.NewSelect(users, nil)
    cw.userSelect.PlaceHolder = "Choose existing…"

    loadBtn := widget.NewButton("Load", func() {
      if name := cw.userSelect.Selected; name != "" {
        if err := cw.core.SetDefaultUser(name); err != nil {
          dialog.ShowError(err, cw.wnd)
          return
        }
        cw.refreshUserTab()
      }
    })

    cw.newUserEntry = widget.NewEntry()
    cw.newUserEntry.SetPlaceHolder("New user ID (filename)")
    createBtn := widget.NewButton("Create", func() {
      id := cw.newUserEntry.Text
      if id == "" {
        dialog.ShowInformation("Missing name",
          "Please enter a user ID", cw.wnd)
        return
      }
      if err := cw.core.CreateUser(id); err != nil {
        dialog.ShowError(err, cw.wnd)
        return
      }
      cw.refreshUserTab()
    })

    return container.NewVBox(
      widget.NewLabelWithStyle("No user loaded",
        fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
      widget.NewSeparator(),
      widget.NewLabel("Select an existing user:"),
      container.NewHBox(cw.userSelect, loadBtn),
      widget.NewSeparator(),
      widget.NewLabel("…or create a new one:"),
      container.NewHBox(cw.newUserEntry, createBtn),
    )
  }

  // --- user is loaded, edit name & default agent ---
  usr := cw.core.User()

  // a) build Name + DefaultAgent form
  cw.userNameEntry = widget.NewEntry()
  cw.userNameEntry.SetText(usr.Name)

  metas := cw.core.Agents()
  names := make([]string, len(metas))
  for i, m := range metas {
    names[i] = m.Name
  }
  cw.userDefaultSelect = widget.NewSelect(names, nil)
  cw.userDefaultSelect.PlaceHolder = "None"
  if usr.DefaultAgent != nil {
    cw.userDefaultSelect.SetSelected(usr.DefaultAgent.Name)
  }

  cw.userForm = &widget.Form{
    Items: []*widget.FormItem{
      {Text: "User Name",     Widget: cw.userNameEntry},
      {Text: "Default Agent", Widget: cw.userDefaultSelect},
    },
    OnSubmit: func() {
      usr.Name = cw.userNameEntry.Text
      if sel := cw.userDefaultSelect.Selected; sel != "" {
        if err := cw.core.SetDefaultAgent(sel); err != nil {
          dialog.ShowError(err, cw.wnd)
          return
        }
      }
      cw.refreshUserTab()
      cw.mainTabs.SelectTabIndex(0)
    },
    OnCancel: func() {
      cw.userNameEntry.SetText(usr.Name)
      if usr.DefaultAgent != nil {
        cw.userDefaultSelect.SetSelected(usr.DefaultAgent.Name)
      } else {
        cw.userDefaultSelect.SetSelected("")
      }
      cw.mainTabs.SelectTabIndex(0)
    },
  }

  // b) Available Agents list
  cw.availAgentsBox = container.NewVBox(
    widget.NewLabelWithStyle("Available Agents",
      fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
  )
  if len(metas) == 0 {
    cw.availAgentsBox.Add(widget.NewLabel("— none —"))
  } else {
    for _, m := range metas {
      cw.availAgentsBox.Add(container.NewHBox(
        widget.NewLabel(m.Name),
        layout.NewSpacer(),
        widget.NewButton("Switch to", func(name string) func() {
          return func() {
            if err := cw.core.SwitchAgent(name); err != nil {
              dialog.ShowError(err, cw.wnd)
              return
            }
            cw.refreshUserTab()
          }
        }(m.Name)),
      ))
    }
  }

  // c) “Add Agent” sub‐form using our custom widget
  addAgentForm := NewAddAgentForm(
    cw.core.Toolpacks(),
    func(name, model string, tools []string) {
      if name == "" || model == "" {
        dialog.ShowInformation("Missing fields",
          "Please fill Agent Name and Model", cw.wnd)
        return
      }
      meta := app.AgentMeta{
        Name:      name,
        Model:     model,
        ToolPaths: tools,
      }
      if err := cw.core.CreateAgent(meta); err != nil {
        dialog.ShowError(err, cw.wnd)
        return
      }
      cw.refreshUserTab()
    },
  )

  // d) assemble everything
  return container.NewVBox(
    cw.userForm,
    widget.NewSeparator(),
    cw.availAgentsBox,
    widget.NewSeparator(),
    addAgentForm,
  )
}

func (cw *ChatWindow) makeUserTab() *container.TabItem {
  return container.NewTabItem("User", cw.buildUserPane())
}

func (cw *ChatWindow) refreshUserTab() {
  const idx = 1
  if idx >= len(cw.mainTabs.Items) {
    return
  }
  cw.mainTabs.Items[idx].Content = cw.buildUserPane()
  cw.mainTabs.Refresh()
  cw.mainTabs.SelectTabIndex(idx)
}
