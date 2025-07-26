package gui

import (
  "context"
  "fmt"

  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/container"
  "fyne.io/fyne/v2/dialog"
  "fyne.io/fyne/v2/layout"
  "fyne.io/fyne/v2/widget"

  "github.com/johnjallday/dolphin-tool-calling-agent/internal/app"
)

// ChatWindow holds our Fyne window, the core App, and all tabs.
type ChatWindow struct {
  app   fyne.App
  wnd   fyne.Window
  core  app.App

  mainTabs *container.AppTabs

  // cached TabItems
  chatTab  *container.TabItem
  toolsTab *container.TabItem
  agentTab *container.TabItem
  userTab  *container.TabItem

  // chat widgets
  historyBox    *fyne.Container
  historyScroll *container.Scroll
  inputEntry    *widget.Entry

  // top bar
  statusLabel *widget.Label

  // tools widgets
  toolsList     *fyne.Container
  toolpacksList *fyne.Container

  // agent widgets
  agentList *fyne.Container

  // user widgets
  userSelect        *widget.Select
  newUserEntry      *widget.Entry
  userNameEntry     *widget.Entry
  userDefaultSelect *widget.Select
}

// NewChatWindow builds the window and all tabs exactly once.
func NewChatWindow(fy fyne.App, core app.App) *ChatWindow {
  w := fy.NewWindow("🐬 Dolphin Chat 🐬")
  cw := &ChatWindow{app: fy, wnd: w, core: core}

  // cache each TabItem
  cw.chatTab = cw.makeChatTab()
  cw.toolsTab = cw.makeToolsTab()
  cw.agentTab = cw.makeAgentTab()
  cw.userTab = cw.makeUserTab()

  // put them into AppTabs
  cw.mainTabs = container.NewAppTabs(
    cw.chatTab, cw.toolsTab, cw.agentTab, cw.userTab,
  )
  cw.mainTabs.SetTabLocation(container.TabLocationTop)

  // top bar with status
  cw.statusLabel = widget.NewLabel("")
  topBar := container.NewHBox(cw.statusLabel, layout.NewSpacer())

  // assemble
  content := container.NewBorder(topBar, nil, nil, nil, cw.mainTabs)
  w.SetContent(content)
  w.Resize(fyne.NewSize(700, 500))

  // initial fill
  cw.RefreshAll()

  return cw
}

// ShowAndRun pops up the window.
func (cw *ChatWindow) ShowAndRun() {
  cw.wnd.ShowAndRun()
}

// RefreshAll updates every pane in the window.
func (cw *ChatWindow) RefreshAll() {
  // 1) status bar
  cw.refreshUserStatus()

  // 2) chat history (in case underlying messages have changed)
  cw.historyBox.Refresh()
  cw.historyScroll.ScrollToBottom()

  // 3) tools
  cw.refreshCurrentToolsList()
  cw.refreshToolpacksList()

  // 4) agent
  cw.agentTab.Content = cw.buildAgentPane()
  cw.agentTab.Content.Refresh()

  // 5) user
  cw.userTab.Content = cw.buildUserPane()
  cw.userTab.Content.Refresh()
}

// ─────────────────────────────────────────────────────────────────────────────
// CHAT TAB
// ─────────────────────────────────────────────────────────────────────────────

func (cw *ChatWindow) makeChatTab() *container.TabItem {
  // input area
  cw.inputEntry = widget.NewEntry()
  cw.inputEntry.SetPlaceHolder("Type your message…")
  cw.inputEntry.OnSubmitted = func(_ string) { cw.sendMessage() }
  sendBtn := widget.NewButton("Send", cw.sendMessage)
  bottom := container.NewBorder(nil, nil, nil, sendBtn, cw.inputEntry)

  // history area
  cw.historyBox = container.NewVBox()
  cw.historyScroll = container.NewVScroll(cw.historyBox)

  pane := container.NewBorder(nil, bottom, nil, nil, cw.historyScroll)
  return container.NewTabItem("Chat", pane)
}

func (cw *ChatWindow) sendMessage() {
  txt := cw.inputEntry.Text
  if txt == "" {
    return
  }
  cw.appendMessage("You", txt)
  cw.inputEntry.SetText("")
  cw.wnd.Canvas().Focus(cw.inputEntry)

  go func() {
    reply, err := cw.core.SendMessage(context.Background(), txt)
    if err != nil {
      cw.appendMessage("Error", err.Error())
      return
    }
    cw.appendMessage("Agent", reply)
  }()
}

func (cw *ChatWindow) appendMessage(who, msg string) {
  cw.historyBox.Add(widget.NewLabel(fmt.Sprintf("%s: %s", who, msg)))
  cw.historyBox.Refresh()
  cw.historyScroll.ScrollToBottom()
}

// ─────────────────────────────────────────────────────────────────────────────
// TOOLS TAB
// ─────────────────────────────────────────────────────────────────────────────

func (cw *ChatWindow) makeToolsTab() *container.TabItem {
  // a) Current Tools
  cw.toolsList = container.NewVBox()
  cw.refreshCurrentToolsList()
  curScroll := container.NewVScroll(cw.toolsList)

  // b) On‐disk Toolpacks
  cw.toolpacksList = container.NewVBox()
  cw.refreshToolpacksList()
  packScroll := container.NewVScroll(cw.toolpacksList)

  tabs := container.NewAppTabs(
    container.NewTabItem("Current Tools", curScroll),
    container.NewTabItem("Toolpacks", packScroll),
  )
  tabs.SetTabLocation(container.TabLocationTop)
  return container.NewTabItem("Tools", tabs)
}

func (cw *ChatWindow) refreshCurrentToolsList() {
  cw.toolsList.Objects = nil
  for _, t := range cw.core.Tools() {
    cw.toolsList.Add(widget.NewLabel(fmt.Sprintf(
      "%s: %s", t.Name, t.Description,
    )))
    cw.toolsList.Add(widget.NewSeparator())
  }
  if len(cw.toolsList.Objects) == 0 {
    cw.toolsList.Add(widget.NewLabel("(no tools registered)"))
  }
  cw.toolsList.Refresh()
}

func (cw *ChatWindow) refreshToolpacksList() {
  cw.toolpacksList.Objects = nil
  packs := cw.core.Toolpacks()
  if len(packs) == 0 {
    cw.toolpacksList.Add(widget.NewLabelWithStyle(
      "No toolpacks found in ./plugins",
      fyne.TextAlignCenter, fyne.TextStyle{Italic: true}),
    )
  } else {
    for _, name := range packs {
      cw.toolpacksList.Add(container.NewHBox(
        widget.NewLabel(name),
        layout.NewSpacer(),
      ))
    }
  }
  cw.toolpacksList.Refresh()
}

// ─────────────────────────────────────────────────────────────────────────────
// AGENT TAB
// ─────────────────────────────────────────────────────────────────────────────

func (cw *ChatWindow) makeAgentTab() *container.TabItem {
  return container.NewTabItem("Agent", cw.buildAgentPane())
}

func (cw *ChatWindow) buildAgentPane() fyne.CanvasObject {
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

// ─────────────────────────────────────────────────────────────────────────────
// USER TAB
// ─────────────────────────────────────────────────────────────────────────────

func (cw *ChatWindow) makeUserTab() *container.TabItem {
  return container.NewTabItem("User", cw.buildUserPane())
}

func (cw *ChatWindow) buildUserPane() fyne.CanvasObject {
  if cw.core.User() == nil {
    // no user yet
    users := cw.core.Users()
    cw.userSelect = widget.NewSelect(users, nil)
    cw.userSelect.PlaceHolder = "Choose existing…"
    load := widget.NewButton("Load", func() {
      if sel := cw.userSelect.Selected; sel != "" {
        if err := cw.core.SetDefaultUser(sel); err != nil {
          dialog.ShowError(err, cw.wnd)
          return
        }
        cw.RefreshAll()
      }
    })

    cw.newUserEntry = widget.NewEntry()
    cw.newUserEntry.SetPlaceHolder("New user ID")
    create := widget.NewButton("Create", func() {
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
      cw.RefreshAll()
    })

    return container.NewVBox(
      widget.NewLabelWithStyle("No user loaded",
        fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
      widget.NewSeparator(),
      container.NewHBox(cw.userSelect, load),
      widget.NewSeparator(),
      container.NewHBox(cw.newUserEntry, create),
    )
  }

  // editing an existing user
  usr := cw.core.User()
  cw.userNameEntry = widget.NewEntry()
  cw.userNameEntry.SetText(usr.Name)

  // select default agent
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

  form := widget.NewForm(
    &widget.FormItem{Text: "User Name",     Widget: cw.userNameEntry},
    &widget.FormItem{Text: "Default Agent", Widget: cw.userDefaultSelect},
  )
  form.OnSubmit = func() {
    usr.Name = cw.userNameEntry.Text
    if sel := cw.userDefaultSelect.Selected; sel != "" {
      if err := cw.core.SetDefaultAgent(sel); err != nil {
        dialog.ShowError(err, cw.wnd)
        return
      }
    }
    cw.RefreshAll()
    cw.mainTabs.SelectTabIndex(0)
  }
  form.OnCancel = func() {
    cw.userNameEntry.SetText(usr.Name)
    if usr.DefaultAgent != nil {
      cw.userDefaultSelect.SetSelected(usr.DefaultAgent.Name)
    }
    cw.mainTabs.SelectTabIndex(0)
  }

  return container.NewVBox(form)
}


func (cw *ChatWindow) refreshUserStatus() {
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
    fmt.Sprintf("User: %s\nCurrent Agent: %s", userPart, agentPart),
  )
}
