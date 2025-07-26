package gui

import (
  // "fmt"

  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/container"
  "fyne.io/fyne/v2/widget"
)

// buildToolsPane constructs the CanvasObject for the “Tools” tab.
func (cw *ChatWindow) buildToolsPane() fyne.CanvasObject {
  // a) build the current-tools list
  cw.toolsList = container.NewVBox()
  cw.refreshCurrentToolsList()
  toolsScroll := container.NewVScroll(cw.toolsList)

  // b) build the toolpacks placeholder
  cw.toolpacksPane = container.NewCenter(
    widget.NewLabelWithStyle(
      "Toolpacks go here…",
      fyne.TextAlignCenter,
      fyne.TextStyle{Italic: true},
    ),
  )

  // c) wrap them in sub-tabs
  toolsTabs := container.NewAppTabs(
    container.NewTabItem("Current Tools", toolsScroll),
    container.NewTabItem("Toolpacks",     cw.toolpacksPane),
  )
  toolsTabs.SetTabLocation(container.TabLocationTop)
  return toolsTabs
}

// makeToolsTab wraps buildToolsPane in a TabItem
func (cw *ChatWindow) makeToolsTab() *container.TabItem {
  return container.NewTabItem("Tools", cw.buildToolsPane())
}
