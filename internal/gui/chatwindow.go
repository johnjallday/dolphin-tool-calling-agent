package gui

import (
  "context"
  "fmt"

  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/container"
  "fyne.io/fyne/v2/layout"
  "fyne.io/fyne/v2/widget"
  "fyne.io/fyne/v2/dialog"

  "github.com/johnjallday/dolphin-tool-calling-agent/internal/app"
  "github.com/johnjallday/dolphin-tool-calling-agent/internal/agent"
)

type ChatWindow struct {
  app         fyne.App        // ← keep the fyne.App so we can NewWindow
  wnd         fyne.Window
  core        app.App


 	mainTabs *container.AppTabs

	// used to repopulate “Current Tools”
  toolsList     *fyne.Container
  toolpacksPane fyne.CanvasObject

	userNameEntry    *widget.Entry
  userDefaultAgent *widget.Select
  userForm         *widget.Form


  statusLabel *widget.Label
	agentSelect     *widget.Select
	agentNameEntry  *widget.Entry
	agentModelEntry *widget.Entry
	agentForm       *widget.Form

  historyBox 		 	*fyne.Container
  historyScroll		*container.Scroll

  onboardingBox *fyne.Container

  inputEntry  *widget.Entry
}

func NewChatWindow(fy fyne.App, core app.App) *ChatWindow {
  w := fy.NewWindow("🐬 Dolphin Chat 🐬")
  cw := &ChatWindow{app: fy, wnd: w, core: core}
  cw.buildUI()
  return cw
}


func (cw *ChatWindow) buildUI() {
  // ─── 1) TOP BAR ──────────────────────────────────────────────
  cw.statusLabel = widget.NewLabel("")
  cw.refreshUserStatus()

  topBar := container.NewHBox(
    cw.statusLabel,
    layout.NewSpacer(),
  )

  // ─── 2) CHAT TAB ─────────────────────────────────────────────
  cw.inputEntry = widget.NewEntry()
  cw.inputEntry.SetPlaceHolder("Type your message…")
  cw.inputEntry.OnSubmitted = func(_ string) { cw.sendMessage() }
  sendBtn := widget.NewButton("Send", cw.sendMessage)

  chatBottom := container.NewBorder(nil, nil, nil, sendBtn, cw.inputEntry)
  chatCenter := cw.buildCenter()
  chatPane   := container.NewBorder(nil, chatBottom, nil, nil, chatCenter)

  // ─── 3) TOOLS TAB ────────────────────────────────────────────
  cw.toolsList     = container.NewVBox()
  cw.refreshCurrentToolsList()
  toolsScroll      := container.NewVScroll(cw.toolsList)
  cw.toolpacksPane = container.NewCenter(
    widget.NewLabelWithStyle("Toolpacks go here…",
      fyne.TextAlignCenter, fyne.TextStyle{Italic: true}),
  )
  toolsTabs := container.NewAppTabs(
    container.NewTabItem("Current Tools", toolsScroll),
    container.NewTabItem("Toolpacks",     cw.toolpacksPane),
  )
  toolsTabs.SetTabLocation(container.TabLocationTop)


  // ─── 5) USER TAB ──────────────────────────────────────────────
  var userPane fyne.CanvasObject
  usr := cw.core.User()
  if usr == nil {
    userPane = cw.createOnboardingBox()
  } else {
    cw.userNameEntry = widget.NewEntry()
    cw.userNameEntry.SetText(usr.Name)

    names2 := make([]string, len(usr.Agents))
    for i, a := range usr.Agents {
      names2[i] = a.Name
    }
    cw.userDefaultAgent = widget.NewSelect(names2, nil)
    if usr.DefaultAgent != nil {
      cw.userDefaultAgent.SetSelected(usr.DefaultAgent.Name)
    }
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
        cw.mainTabs.SelectTabIndex(0)
      },
      OnCancel: func() {
        cw.userNameEntry.SetText(usr.Name)
        if usr.DefaultAgent != nil {
          cw.userDefaultAgent.SetSelected(usr.DefaultAgent.Name)
        }
        cw.mainTabs.SelectTabIndex(0)
      },
    }
    userPane = container.NewVBox(cw.userForm)
  }

  // ─── 6) PUT ’EM ALL TOGETHER ─────────────────────────────────
  cw.mainTabs = container.NewAppTabs(
    container.NewTabItem("Chat",  chatPane),
    container.NewTabItem("Tools", toolsTabs),
		cw.makeAgentTab(),
    container.NewTabItem("User",  userPane),
  )
  cw.mainTabs.SetTabLocation(container.TabLocationTop)

  content := container.NewBorder(
    topBar, nil, nil, nil,
    cw.mainTabs,
  )
  cw.wnd.SetContent(content)
  cw.wnd.Resize(fyne.NewSize(600, 480))
}



func (cw *ChatWindow) buildCenter() fyne.CanvasObject {
  // … your existing logic to choose between
  // onboardingBox, agentOnboardingBox, or cw.historyScroll …
  // e.g.:
  if len(cw.core.Users()) == 0 {
    return cw.createOnboardingBox()
  }
  // … etc …
  // default:
  cw.historyBox = container.NewVBox()
  cw.historyScroll = container.NewVScroll(cw.historyBox)
  return cw.historyScroll
}


// refreshCurrentToolsList repopulates the list of .so plugins
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

func (cw *ChatWindow) topBar() *fyne.Container {
  return container.NewHBox(
    cw.statusLabel,
    layout.NewSpacer(),
  )
}

func (cw *ChatWindow) refreshUserStatus() {
  if u := cw.core.User(); u != nil {
    def := "None"
    if u.DefaultAgent != nil {
      def = u.DefaultAgent.Name
    }
    cw.statusLabel.SetText(fmt.Sprintf("User: %s (Default Agent: %s)", u.Name, def))
  } else {
    cw.statusLabel.SetText("User: None")
  }
}

// helper to build the history pane (so we can re-use it)
func (cw *ChatWindow) buildCenterHistory(bottomBar *fyne.Container) {
  cw.historyBox = container.NewVBox()
  cw.historyScroll = container.NewVScroll(cw.historyBox)
  cw.historyScroll.SetMinSize(fyne.NewSize(400, 300))
}


func (cw *ChatWindow) sendMessage() {
  txt := cw.inputEntry.Text
  if txt == "" {
    return
  }
  cw.appendMessage("You", txt)
  cw.inputEntry.SetText("")
  cw.wnd.Canvas().Focus(cw.inputEntry)

  fyne.Do(func() {
    reply, err := cw.core.SendMessage(context.Background(), txt)
    if err != nil {
      cw.appendMessage("Error", err.Error())
      return
    }
    cw.appendMessage("Agent", reply)
  })
}

func (cw *ChatWindow) appendMessage(who, msg string) {
  lbl := widget.NewLabel(fmt.Sprintf("%s: %s", who, msg))
  cw.historyBox.Add(lbl)
  cw.historyBox.Refresh()
  cw.historyScroll.ScrollToBottom()
}

func (cw *ChatWindow) ShowAndRun() {
  cw.wnd.ShowAndRun()
}








