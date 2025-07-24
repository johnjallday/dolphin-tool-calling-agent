package gui

import (
  "image/color"

  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/canvas"
  "fyne.io/fyne/v2/widget"
)

// WhiteMultiLineEntry is a drop-in replacement for widget.NewMultiLineEntry()
// whose text (and placeholder) is forced white.
type WhiteMultiLineEntry struct{ *widget.Entry }

func NewWhiteMultiLineEntry() *WhiteMultiLineEntry {
  return &WhiteMultiLineEntry{widget.NewMultiLineEntry()}
}

func (w *WhiteMultiLineEntry) CreateRenderer() fyne.WidgetRenderer {
  rend := w.Entry.CreateRenderer()
  for _, obj := range rend.Objects() {
    if txt, ok := obj.(*canvas.Text); ok {
      txt.Color = color.White
    }
  }
  rend.Refresh()
  return rend
}
