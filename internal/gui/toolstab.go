package gui

import (
  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/container"
  "fyne.io/fyne/v2/widget"
  "fyne.io/fyne/v2/layout"
)

// buildToolsPane constructs the CanvasObject for the “Tools” tab.
func (cw *ChatWindow) buildToolsPane() fyne.CanvasObject {
  // a) current‐tools
  cw.toolsList = container.NewVBox()
  cw.refreshCurrentToolsList()
  toolsScroll := container.NewVScroll(cw.toolsList)

  // b) toolpacks
  cw.toolpacksList = container.NewVBox()
  cw.refreshToolpacksList()
  toolpacksScroll := container.NewVScroll(cw.toolpacksList)

  // c) sub‐tabs
  toolsTabs := container.NewAppTabs(
    container.NewTabItem("Current Tools", toolsScroll),
    container.NewTabItem("Toolpacks",     toolpacksScroll),
  )
  toolsTabs.SetTabLocation(container.TabLocationTop)
  return toolsTabs
}

// makeToolsTab wraps buildToolsPane in a TabItem
func (cw *ChatWindow) makeToolsTab() *container.TabItem {
  return container.NewTabItem("Tools", cw.buildToolsPane())
}

// refreshToolpacksList repopulates cw.toolpacksList from disk
func (cw *ChatWindow) refreshToolpacksList() {
  // clear out the old objects
  cw.toolpacksList.Objects = nil

  packs := cw.core.Toolpacks()
  if len(packs) == 0 {
    cw.toolpacksList.Add(
      widget.NewLabelWithStyle(
        "No toolpacks found in ./plugins",
        fyne.TextAlignCenter,
        fyne.TextStyle{Italic: true},
      ),
    )
  } else {
    for _, name := range packs {
      // btn := widget.NewButton("Load", func(n string) func() {
      //   return func() {
      //     // TODO: hook up real loading logic
      //     fmt.Println("load toolpack:", n)
      //     // after load:
      //     // cw.refreshCurrentToolsList()
      //     // cw.refreshToolpacksList()
      //   }
      // }(name))
      //
      row := container.NewHBox(
        widget.NewLabel(name),
        layout.NewSpacer(),
        // btn,
      )
      cw.toolpacksList.Add(row)
    }
  }

  cw.toolpacksList.Refresh()
}
