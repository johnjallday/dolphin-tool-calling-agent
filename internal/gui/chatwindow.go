package gui

import (
  "context"
  "fmt"

  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/container"
  "fyne.io/fyne/v2/layout"
  "fyne.io/fyne/v2/widget"

  "github.com/johnjallday/dolphin-tool-calling-agent/internal/app"
  //"github.com/johnjallday/dolphin-tool-calling-agent/pkg/tools"
)

type ChatWindow struct {
  app         fyne.App        // ← keep the fyne.App so we can NewWindow
  wnd         fyne.Window
  core        app.App

 	toolsWin           fyne.Window  // keep a reference so we can focus/close it
  toolsShortcutAdded bool 

  statusLabel *widget.Label
  agentSelect *widget.Select
	toolsBtn    *widget.Button

  historyBox  *fyne.Container
  history     *container.Scroll

  inputEntry  *widget.Entry
}

func NewChatWindow(fy fyne.App, core app.App) *ChatWindow {
  w := fy.NewWindow("🐬 Dolphin Chat 🐬")
  cw := &ChatWindow{
    app:  fy,
    wnd:  w,
    core: core,
  }
  cw.buildUI()
  return cw
}

 func (cw *ChatWindow) buildUI() {
  // --- 1) Status label ---
  cw.statusLabel = widget.NewLabel("")
  cw.refreshUserStatus()

  // --- 2) Agent dropdown (no callback yet) ---
  metas := cw.core.Agents()
  names := make([]string, len(metas))
  for i, m := range metas {
    names[i] = m.Name
  }
  cw.agentSelect = widget.NewSelect(names, nil)

  // --- 3) Tools button ---
  cw.toolsBtn = widget.NewButton("Tools", cw.openToolsWindow)

  // --- 4) History area ---
  cw.historyBox = container.NewVBox()
  cw.history    = container.NewVScroll(cw.historyBox)
  cw.history.SetMinSize(fyne.NewSize(400, 300))

  // --- 5) Input + send ---
  cw.inputEntry = widget.NewEntry()
  cw.inputEntry.SetPlaceHolder("Type your message…")
  cw.inputEntry.OnSubmitted = func(_ string) { cw.sendMessage() }
  sendBtn := widget.NewButton("Send", cw.sendMessage)
  bottomBar := container.NewBorder(nil, nil, nil, sendBtn, cw.inputEntry)

  // --- 6) Top bar: status | spacer | tools | agent dropdown ---
  topBar := container.NewHBox(
    cw.statusLabel,
    layout.NewSpacer(),
    cw.toolsBtn,
    cw.agentSelect,
  )

  // --- 7) Assemble into window ---
  content := container.NewBorder(
    topBar,       // north
    bottomBar,    // south
    nil, nil,     // west, east
    cw.history,   // center
  )
  cw.wnd.SetContent(content)

  // --- 8) Wire up agentSelect ---
  cw.agentSelect.OnChanged = func(name string) {
    if err := cw.core.LoadAgent(name); err != nil {
      fmt.Println("failed to load agent:", err)
      return
    }
    cw.historyBox.Objects = nil
    cw.historyBox.Refresh()
  }

  // --- 9) Initial agent selection (fires OnChanged) ---
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
  cw.statusLabel.Refresh()
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
  lbl := widget.NewLabel(fmt.Sprintf("%s: %s", who, msg))
  cw.historyBox.Add(lbl)
  cw.historyBox.Refresh()
  cw.history.ScrollToBottom()
}

func (cw *ChatWindow) ShowAndRun() {
  cw.wnd.ShowAndRun()
}
