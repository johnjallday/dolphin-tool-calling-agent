package gui

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
    "fyne.io/fyne/v2/driver/desktop"
)

// NewHelpView returns a CanvasObject showing help text and binding Esc to onBack.
func NewHelpView(w fyne.Window, onBack func()) fyne.CanvasObject {
    helpLabel := widget.NewLabel("Help:\n• /agent to switch to agent\n• /location to show location\n\nType /back to return to chat")
    history := container.NewVBox(helpLabel)
    scroll := container.NewVScroll(history)

    entry := widget.NewEntry()
    entry.SetPlaceHolder("/back")
    entry.OnSubmitted = func(txt string) {
        if txt == "/back" {
            entry.SetText("")
            onBack()
        }
    }

    w.Canvas().AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyEscape}, func(_ fyne.Shortcut) {
        onBack()
    })

    return container.NewBorder(nil, entry, nil, nil, scroll)
}
