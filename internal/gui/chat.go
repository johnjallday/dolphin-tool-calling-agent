package gui

import (
  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/container"
  "fyne.io/fyne/v2/widget"
)

var (
  historyBox   *fyne.Container
  historyScroll *container.Scroll
)

// AddToHistory appends a line to the chat and scrolls to bottom.
func AddToHistory(text string) {
  if historyBox == nil || historyScroll == nil {
    return
  }
  historyBox.Add(widget.NewLabel(text))
  historyScroll.ScrollToBottom()
}

func NewChatView(
  showAgent, showHelp, showLocation func(),
) (fyne.CanvasObject, *widget.Entry) {

  historyBox = container.NewVBox()
  historyScroll = container.NewScroll(historyBox)

  input := widget.NewEntry()
  input.SetPlaceHolder("Type message or `/agent` or `/help`…")

  send := func(txt string) {
    input.SetText("")
    switch txt {
    case "/agent":
      showAgent()
    case "/help":
      showHelp()
    case "/location":
      showLocation()
    default:
      if txt != "" {
        AddToHistory("You: " + txt)
      }
    }
  }
  input.OnSubmitted = send
  sendBtn := widget.NewButton("Send", func() { send(input.Text) })

  view := container.NewBorder(nil, container.NewVBox(input, sendBtn), nil, nil, historyScroll)
  return view, input
}
