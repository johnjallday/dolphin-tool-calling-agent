package gui

import (
  "image/color"
	"fmt"

  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/canvas"
  "fyne.io/fyne/v2/widget"
  "fyne.io/fyne/v2/container"
)

// WhiteMultiLineEntry is a decorator around a normal multi‐line Entry
// that forces the text (and placeholder) to be white.
type WhiteMultiLineEntry struct {
  widget.Entry
}

func NewWhiteMultiLineEntry() *WhiteMultiLineEntry {
  // Note: widget.NewMultiLineEntry returns *widget.Entry
  e := &WhiteMultiLineEntry{Entry: *widget.NewMultiLineEntry()}
  e.ExtendBaseWidget(e)
  return e
}

func (w *WhiteMultiLineEntry) CreateRenderer() fyne.WidgetRenderer {
  // get the normal Entry renderer
  rend := w.Entry.CreateRenderer()
  // tweak every canvas.Text in it
  for _, obj := range rend.Objects() {
    if txt, ok := obj.(*canvas.Text); ok {
      txt.Color = color.White
    }
  }
  // no need to call rend.Refresh() here; Fyne will do that for you after CreateRenderer
  return rend
}


// buildAgentsList returns a VBox of all AgentMeta from core.Agents()
func (cw *ChatWindow) buildAgentsList() *fyne.Container {
  if cw.agentList == nil {
    cw.agentList = container.NewVBox()
  }
  // clear it out
  cw.agentList.Objects = nil

  metas := cw.core.Agents()
  cw.agentList.Add(widget.NewLabelWithStyle(
    "Agents", fyne.TextAlignLeading, fyne.TextStyle{Bold: true},
  ))

  if len(metas) == 0 {
    cw.agentList.Add(widget.NewLabel("— no agents for this user —"))
  } else {
    for _, m := range metas {
      cw.agentList.Add(widget.NewLabel(
        fmt.Sprintf("%s (model=%s) plugins=%v",
          m.Name, m.Model, m.Plugins),
      ))
    }
  }

  cw.agentList.Refresh()
  return cw.agentList
}
