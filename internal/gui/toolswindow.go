package gui

import (
  "fmt"

  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/container"
  "fyne.io/fyne/v2/widget"
  "fyne.io/fyne/v2/driver/desktop"
)

func (cw *ChatWindow) openToolsWindow() {
  // If it already exists, just focus it
  if cw.toolsWin != nil {
    cw.toolsWin.RequestFocus()
    return
  }

  w := cw.app.NewWindow("Tools")
  cw.toolsWin = w      // remember for later focus/close

  // --- build the two tabs: Current Tools / Toolpacks placeholder ---
  // (same as your previous code)
  currentList := container.NewVBox()
  for _, t := range cw.core.Tools() {
    currentList.Add(widget.NewLabel(fmt.Sprintf("%s: %s", t.Name, t.Description)))
    currentList.Add(widget.NewSeparator())
  }
  if len(currentList.Objects) == 0 {
    currentList.Add(widget.NewLabel("(no tools registered)"))
  }
  currentScroll := container.NewVScroll(currentList)

  toolpacksPlaceholder := container.NewCenter(
    widget.NewLabelWithStyle("Toolpacks will go here", fyne.TextAlignCenter, fyne.TextStyle{Italic: true}),
  )

  tabs := container.NewAppTabs(
    container.NewTabItem("Current Tools", currentScroll),
    container.NewTabItem("Toolpacks", toolpacksPlaceholder),
  )
  tabs.SetTabLocation(container.TabLocationTop)

  w.SetContent(tabs)
  w.Resize(fyne.NewSize(360, 480))

  // --- define the shortcuts we’ll use ---
  ctrlTab := &desktop.CustomShortcut{
    KeyName:  fyne.KeyTab,
    Modifier: fyne.KeyModifierControl,
  }
  ctrlShiftTab := &desktop.CustomShortcut{
    KeyName:  fyne.KeyTab,
    Modifier: fyne.KeyModifierControl | fyne.KeyModifierShift,
  }
  esc := &desktop.CustomShortcut{KeyName: fyne.KeyEscape}

  // --- register them on the Tools window’s Canvas ---
  //  Esc → close the tools window
  w.Canvas().AddShortcut(esc, func(_ fyne.Shortcut) {
    w.Close()
  })

  //  Ctrl+Shift+Tab → go back to the chat window
  w.Canvas().AddShortcut(ctrlShiftTab, func(_ fyne.Shortcut) {
    cw.wnd.RequestFocus()
  })

  // --- register Ctrl+Tab on the Chat window (once only) ---
  if !cw.toolsShortcutAdded {
    cw.wnd.Canvas().AddShortcut(ctrlTab, func(_ fyne.Shortcut) {
      if cw.toolsWin != nil {
        cw.toolsWin.RequestFocus()
      } else {
        cw.openToolsWindow()
      }
    })
    cw.toolsShortcutAdded = true
  }

  w.Show()
}
