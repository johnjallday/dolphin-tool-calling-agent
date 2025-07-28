package gui

import (
  "fmt"

  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/container"
  "fyne.io/fyne/v2/dialog"
  "fyne.io/fyne/v2/widget"

)



// ─────────────────────────────────────────────────────────────────────────────
// USER TAB
// ─────────────────────────────────────────────────────────────────────────────

func (cw *MainWindow) makeUserTab() *container.TabItem {
  return container.NewTabItem("User", cw.buildUserPane())
}



func (cw *MainWindow) buildUserPane() fyne.CanvasObject {
  pane := container.NewVBox()

  // 1) No user loaded
  if cw.core.User() == nil {
    sel := widget.NewSelect(cw.core.Users(), nil)
    sel.PlaceHolder = "Choose existing…"
    loadBtn := widget.NewButton("Load", func() {
      if sel.Selected == "" {
        return
      }
      if err := cw.core.SetDefaultUser(sel.Selected); err != nil {
        dialog.ShowError(err, cw.wnd)
        return
      }
      cw.RefreshAll()
    })

    newEntry := widget.NewEntry()
    newEntry.SetPlaceHolder("New user ID…")
    createBtn := widget.NewButton("Create", func() {
      if newEntry.Text == "" {
        dialog.ShowInformation("Missing name", "Please enter a user ID", cw.wnd)
        return
      }
      if err := cw.core.CreateUser(newEntry.Text); err != nil {
        dialog.ShowError(err, cw.wnd)
        return
      }
      cw.RefreshAll()
    })

    pane.Add(widget.NewLabelWithStyle(
      "No user loaded",
      fyne.TextAlignCenter, fyne.TextStyle{Bold: true},
    ))
    pane.Add(widget.NewSeparator())
    pane.Add(container.NewHBox(sel, loadBtn))
    pane.Add(widget.NewSeparator())
    pane.Add(container.NewHBox(newEntry, createBtn))

  } else {
    // 2) User is loaded → edit form
    usr := cw.core.User()
    nameEntry := widget.NewEntry()
    nameEntry.SetText(usr.Name)

    metas := cw.core.Agents()
    names := make([]string, len(metas))
    for i, m := range metas {
      names[i] = m.Name
    }
    defSel := widget.NewSelect(names, nil)
    defSel.PlaceHolder = "None"
    if usr.DefaultAgent != nil {
      defSel.SetSelected(usr.DefaultAgent.Name)
    }

    form := widget.NewForm(
      &widget.FormItem{Text: "User Name",     Widget: nameEntry},
      &widget.FormItem{Text: "Default Agent", Widget: defSel},
    )
    form.OnSubmit = func() {
      usr.Name = nameEntry.Text
      if defSel.Selected != "" {
        if err := cw.core.SetDefaultAgent(defSel.Selected); err != nil {
          dialog.ShowError(err, cw.wnd)
          return
        }
      }
      cw.RefreshAll()
      cw.mainTabs.SelectTabIndex(0)
    }
    form.OnCancel = func() {
      nameEntry.SetText(usr.Name)
      if usr.DefaultAgent != nil {
        defSel.SetSelected(usr.DefaultAgent.Name)
      }
      cw.mainTabs.SelectTabIndex(0)
    }

    pane.Add(form)
  }

  // 3) Always append the agents list at the bottom
  pane.Add(widget.NewSeparator())
  pane.Add(cw.buildAgentsList())

  return pane
}

func (cw *MainWindow) refreshUserStatus() {
  var userPart, agentPart string

  // 1) Figure out user name
  if u := cw.core.User(); u != nil {
    userPart = u.Name
  } else {
    userPart = "None"
  }

  // 2) Figure out current agent name
  if a := cw.core.Agent(); a != nil {
    agentPart = a.Name
  } else {
    agentPart = "None"
  }

  // 3) Show them both
  cw.statusLabel.SetText(
    fmt.Sprintf("User: %s\n   Current Agent: %s", userPart, agentPart),
  )
}
