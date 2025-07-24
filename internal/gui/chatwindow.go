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
)

type ChatWindow struct {
  app         fyne.App        // ← keep the fyne.App so we can NewWindow
  wnd         fyne.Window
  core        app.App

 	toolsWin           fyne.Window  // keep a reference so we can focus/close it
  toolsShortcutAdded bool 
 	agentWin            fyne.Window     // new window handle
  agentShortcutAdded  bool
	userWin  fyne.Window

  statusLabel *widget.Label
  agentSelect *widget.Select
	toolsBtn    *widget.Button

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
  // --- 1) top bar widgets ---
  cw.statusLabel = widget.NewLabel("")
  cw.refreshUserStatus()

  agents := cw.core.Agents()
  names := make([]string, len(agents))
  for i, m := range agents {
    names[i] = m.Name
  }
  cw.agentSelect = widget.NewSelect(names, nil)
  cw.toolsBtn = widget.NewButton("Tools", cw.openToolsWindow)
  topBar := cw.topBar()

  // --- 2) bottom bar ---
  cw.inputEntry = widget.NewEntry()
  cw.inputEntry.SetPlaceHolder("Type your message…")
  cw.inputEntry.OnSubmitted = func(_ string) { cw.sendMessage() }
  sendBtn := widget.NewButton("Send", cw.sendMessage)
  bottomBar := container.NewBorder(nil, nil, nil, sendBtn, cw.inputEntry)

  // --- 3) center pane: onboarding vs history ---
  var center fyne.CanvasObject
  if len(cw.core.Users()) == 0 {
    // On‐boarding: ask for username
    nameEntry := widget.NewEntry()
    nameEntry.SetPlaceHolder("Enter username")

    createBtn := widget.NewButton("Create User", func() {
      userID := strings.TrimSpace(nameEntry.Text)
      if userID == "" {
        dialog.ShowError(fmt.Errorf("username cannot be empty"), cw.wnd)
        return
      }
      if err := cw.core.CreateUser(userID); err != nil {
        dialog.ShowError(err, cw.wnd)
        return
      }
      // rebuild the UI now that we have at least one user
      cw.buildUI()
    })

    cw.onboardingBox = container.NewVBox(
      widget.NewLabelWithStyle("Welcome to Dolphin Chat!",
        fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
      widget.NewLabel("Please pick a username to get started:"),
      nameEntry,
      createBtn,
    )
    center = cw.onboardingBox

  } else {
    // … existing history setup …
    cw.historyBox = container.NewVBox()
    cw.historyScroll = container.NewVScroll(cw.historyBox)
    cw.historyScroll.SetMinSize(fyne.NewSize(400, 300))
    center = cw.historyScroll
  }

  content := container.NewBorder(
    cw.topBar(), bottomBar, nil, nil, center,
  )
  cw.wnd.SetContent(content)

  // --- 4) assemble and set content ---
	cw.wnd.SetContent(container.NewBorder(
		topBar,    // north
		bottomBar, // south
		nil, nil,  // west, east
		center,    // center
	))

  // --- 5) wire up agent dropdown exactly as before ---
  cw.agentSelect.OnChanged = func(name string) {
    if err := cw.core.LoadAgent(name); err != nil {
      fmt.Println("failed to load agent:", err)
      return
    }
    if cw.historyBox != nil {
      cw.historyBox.Objects = nil
      cw.historyBox.Refresh()
    }
  }
  // initial selection …
  switch {
  case cw.core.Agent() != nil:
    cw.agentSelect.SetSelected(cw.core.Agent().Name)
  case cw.core.User() != nil && cw.core.User().DefaultAgent != nil:
    cw.agentSelect.SetSelected(cw.core.User().DefaultAgent.Name)
  default:
    if len(names) > 0 {
      cw.agentSelect.SetSelected(names[0])
    }
  }
}

func (cw *ChatWindow) topBar() *fyne.Container {
  return container.NewHBox(
    cw.statusLabel,
    layout.NewSpacer(),
    cw.toolsBtn,
    cw.agentSelect,
    widget.NewButton("Edit Agent", cw.openAgentWindow),
    widget.NewButton("User", cw.openUserWindow),
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
