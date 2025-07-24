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
  wnd         fyne.Window
  core        app.App

  statusLabel *widget.Label
  agentSelect *widget.Select

  historyBox  *fyne.Container
  history     *container.Scroll

  inputEntry  *widget.Entry
}

func NewChatWindow(fy fyne.App, core app.App) *ChatWindow {
  w := fy.NewWindow("🐬 Dolphin Chat 🐬")
  cw := &ChatWindow{wnd: w, core: core}
  cw.buildUI()
  return cw
}

func (cw *ChatWindow) buildUI() {
  // --- 1) status label (unchanged) ---
  cw.statusLabel = widget.NewLabel("")
  cw.refreshUserStatus()

  // --- 2) agent selector, but no callback yet ---
  metas := cw.core.Agents()
  names := make([]string, len(metas))
  for i, m := range metas {
    names[i] = m.Name
  }
  cw.agentSelect = widget.NewSelect(names, nil)

  // --- 3) history area (must exist before callback fires) ---
  cw.historyBox = container.NewVBox()
  cw.history    = container.NewVScroll(cw.historyBox)
  cw.history.SetMinSize(fyne.NewSize(400, 300))

  // --- 4) input + send button (unchanged) ---
  cw.inputEntry = widget.NewEntry()
  cw.inputEntry.SetPlaceHolder("Type your message…")
  cw.inputEntry.OnSubmitted = func(_ string) { cw.sendMessage() }
  sendBtn := widget.NewButton("Send", cw.sendMessage)
  bottomBar := container.NewBorder(nil, nil, nil, sendBtn, cw.inputEntry)

  // --- 5) top bar (status + spacer + selector) ---
  topBar := container.NewHBox(
    cw.statusLabel,
    layout.NewSpacer(),
    cw.agentSelect,
  )

  // --- 6) assemble everything into window ---
  content := container.NewBorder(
    topBar,       // north
    bottomBar,    // south
    nil, nil,     // west, east
    cw.history,   // center
  )
  cw.wnd.SetContent(content)

  // --- 7) now that historyBox is ready, hook up the select callback ---
  cw.agentSelect.OnChanged = func(name string) {
    if err := cw.core.LoadAgent(name); err != nil {
      fmt.Println("load agent failed:", err)
      return
    }
    // safe to clear now, since historyBox is non‐nil
    cw.historyBox.Objects = nil
    cw.historyBox.Refresh()
  }

  // --- 8) finally pick an initial agent, now that OnChanged is wired ---
  switch {
  case cw.core.Agent() != nil:
    cw.agentSelect.SetSelected(cw.core.Agent().Name)

  case cw.core.User() != nil && cw.core.User().DefaultAgent != nil:
    da := cw.core.User().DefaultAgent.Name
    cw.agentSelect.SetSelected(da)

  default:
    if len(names) > 0 {
      cw.agentSelect.SetSelected(names[0])
    }
  }
}


// refreshUserStatus updates the statusLabel from core.User()
func (cw *ChatWindow) refreshUserStatus() {
  if u := cw.core.User(); u != nil {
    defaultName := "None"
    if u.DefaultAgent != nil {
      defaultName = u.DefaultAgent.Name
    }
    cw.statusLabel.SetText(
      fmt.Sprintf("User: %s (Default Agent: %s)", u.Name, defaultName),
    )
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
