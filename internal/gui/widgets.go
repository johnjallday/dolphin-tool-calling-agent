package gui

import (
  "image/color"

  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/canvas"
  "fyne.io/fyne/v2/widget"
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
