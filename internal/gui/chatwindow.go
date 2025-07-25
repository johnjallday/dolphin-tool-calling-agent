package gui

import (
  "context"
  "fmt"

  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/container"
  "fyne.io/fyne/v2/layout"
  "fyne.io/fyne/v2/widget"
  "fyne.io/fyne/v2/dialog"
	"strings"

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

  // ─── 4) AGENT TAB ────────────────────────────────────────────
  ag := cw.core.Agent()

  // a) build the Select with NO callback first
  agents := cw.core.Agents()
  names  := make([]string, len(agents))
  for i, a := range agents {
    names[i] = a.Name
  }
  cw.agentSelect = widget.NewSelect(names, nil)

  // b) pre-select the current agent (won’t trigger OnChanged)
  if ag != nil {
    cw.agentSelect.SetSelected(ag.Name)
  }

  // c) now wire up OnChanged
  cw.agentSelect.OnChanged = func(sel string) {
    // your “switch agent” logic here, for example
    cw.core.LoadAgent(sel)
    // rebuild the UI so everything picks up the new agent
    cw.buildUI()
  }

  var agentPane fyne.CanvasObject
  if ag == nil {
    // no agent at all yet: show selector + onboarding
    agentPane = container.NewVBox(
      cw.agentSelect,
      widget.NewSeparator(),
      cw.createAgentOnboardingBox(),
    )
  } else {
    // agent exists: show selector + edit form
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
        // redraw the UI so dropdown & other tabs update
        cw.buildUI()
        cw.mainTabs.SelectTabIndex(0) // back to Chat
      },
      OnCancel: func() {
        cw.agentNameEntry.SetText(ag.Name)
        cw.agentModelEntry.SetText(ag.Model)
        cw.mainTabs.SelectTabIndex(0)
      },
    }

    agentPane = container.NewVBox(
      cw.agentSelect,
      widget.NewSeparator(),
      cw.agentForm,
    )
  }

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
    container.NewTabItem("Agent", agentPane),
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
    if err := cw.core.CreateAgent(meta); err != nil {
      dialog.ShowError(err, cw.wnd)
      return
    }
    // now that we have an agent, rebuild into the real chat UI
    cw.buildUI()
  })

  return container.NewVBox(
    widget.NewLabelWithStyle(
      "🎉 Welcome to Dolphin Chat! 🐬", fyne.TextAlignCenter, fyne.TextStyle{Bold: true},
    ),
    widget.NewLabel("Let’s create your first agent:"),
    nameEntry,
    createBtn,
  )
}
