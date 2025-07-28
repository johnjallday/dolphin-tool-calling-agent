
package gui

import (
  "context"
  "fmt"

  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/container"
  "fyne.io/fyne/v2/widget"
)


func (cw *MainWindow) makeChatTab() *container.TabItem {
  cw.inputEntry = widget.NewEntry()
  cw.inputEntry.SetPlaceHolder("Type your message…")
  cw.inputEntry.OnSubmitted = func(_ string) { cw.sendMessage() }

  sendBtn := widget.NewButton("Send", cw.sendMessage)
  bottom := container.NewBorder(nil, nil, nil, sendBtn, cw.inputEntry)

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

  // This is already on Fyne’s UI thread:
  cw.appendMessage("You", txt)
  cw.inputEntry.SetText("")
  cw.wnd.Canvas().Focus(cw.inputEntry)

  // Do the network/agent call in a goroutine
  go func(userText string) {
    reply, err := cw.core.SendMessage(context.Background(), userText)
    if err != nil {
      // schedule error on the UI thread
      fyne.Do(func() {
        cw.appendMessage("Error", err.Error())
      })
      return
    }
    // schedule agent reply on the UI thread
    fyne.Do(func() {
      cw.appendMessage("Agent", reply)
    })
  }(txt)
}

// appendMessage _must_ run on the UI thread.
func (cw *MainWindow) appendMessage(who, msg string) {
  lbl := widget.NewLabel(fmt.Sprintf("%s: %s", who, msg))
  cw.historyBox.Add(lbl)
  cw.historyBox.Refresh()
  cw.historyScroll.ScrollToBottom()
}
