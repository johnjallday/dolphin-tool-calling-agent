package gui

import (
  "context"
  "fmt"
	

  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/container"
  "fyne.io/fyne/v2/layout"
  "fyne.io/fyne/v2/widget"

  "github.com/johnjallday/dolphin-tool-calling-agent/internal/app"
)

type ChatWindow struct {
  app         fyne.App        // ← keep the fyne.App so we can NewWindow
  wnd         fyne.Window
  core        app.App


 	mainTabs *container.AppTabs

	// used to repopulate “Current Tools”
  toolsList     *fyne.Container
  toolpacksPane fyne.CanvasObject
	toolpacksList   *fyne.Container

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
	
  // ─── 3) PUT ’EM ALL TOGETHER ─────────────────────────────────
  cw.mainTabs = container.NewAppTabs(
    container.NewTabItem("Chat",  chatPane),
    //container.NewTabItem("Tools", toolsTabs),
		cw.makeToolsTab(),
		cw.makeAgentTab(),
		cw.makeUserTab(),
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
  // if no users at all → show user‐onboarding
  if len(cw.core.Users()) == 0 {
    return cw.createOnboardingBox()
  }

  // if we have a user, but no agent loaded yet → show agent‐onboarding
  if cw.core.Agent() == nil {
    return cw.createAgentOnboardingBox()
  }

  // otherwise we have both a user and an agent → show chat history
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
