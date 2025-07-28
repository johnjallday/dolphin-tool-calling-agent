package gui

import (
  "context"
  "fmt"

  "fyne.io/fyne/v2/container"
  //"fyne.io/fyne/v2/layout"
  "fyne.io/fyne/v2/widget"

  //"github.com/johnjallday/dolphin-tool-calling-agent/internal/app"
)


// ─────────────────────────────────────────────────────────────────────────────
// CHAT TAB
// ─────────────────────────────────────────────────────────────────────────────

func (cw *MainWindow) makeChatTab() *container.TabItem {
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

func (cw *MainWindow) sendMessage() {
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

func (cw *MainWindow) appendMessage(who, msg string) {
  cw.historyBox.Add(widget.NewLabel(fmt.Sprintf("%s: %s", who, msg)))
  cw.historyBox.Refresh()
  cw.historyScroll.ScrollToBottom()
}
