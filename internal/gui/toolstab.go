package gui

import (
	"fmt"

  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/layout"
  "fyne.io/fyne/v2/container"
  "fyne.io/fyne/v2/widget"
)


// ─────────────────────────────────────────────────────────────────────────────
// TOOLS TAB
// ─────────────────────────────────────────────────────────────────────────────

func (cw *MainWindow) makeToolsTab() *container.TabItem {
  // a) Current Tools
  cw.toolsList = container.NewVBox()
  cw.refreshCurrentToolsList()
  curScroll := container.NewVScroll(cw.toolsList)

  // b) On-disk Toolpacks
  cw.toolpacksList = container.NewVBox()
  cw.refreshToolpacksList()
  localScroll := container.NewVScroll(cw.toolpacksList)

  // c) Remote Toolpacks
  cw.remotetoolpacksList = container.NewVBox()
  cw.refreshRemoteToolpacksList()
  remoteScroll := container.NewVScroll(cw.remotetoolpacksList)

  tabs := container.NewAppTabs(
    container.NewTabItem("Current Tools", curScroll),
    container.NewTabItem("Local Toolpacks", localScroll),
    container.NewTabItem("Remote Toolpacks", remoteScroll),
  )
  tabs.SetTabLocation(container.TabLocationTop)
  return container.NewTabItem("Tools", tabs)
}

// refreshRemoteToolpacksList loads and displays the registry entries
func (cw *MainWindow) refreshRemoteToolpacksList() {
  cw.remotetoolpacksList.Objects = nil
  packs, err := cw.core.ListRemoteToolpacks()
  if err != nil {
    cw.remotetoolpacksList.Add(widget.NewLabel(fmt.Sprintf(
      "Error loading remote toolpacks: %v", err)))
  } else if len(packs) == 0 {
    cw.remotetoolpacksList.Add(widget.NewLabelWithStyle(
      "No remote toolpacks found", fyne.TextAlignCenter,
      fyne.TextStyle{Italic: true}))
  } else {
    for _, name := range packs {
      cw.remotetoolpacksList.Add(container.NewHBox(
        widget.NewLabel(name),
        layout.NewSpacer(),
      ))
    }
  }
  cw.remotetoolpacksList.Refresh()
}





func (cw *MainWindow) refreshCurrentToolsList() {
  cw.toolsList.Objects = nil
  for _, t := range cw.core.Tools() {
    cw.toolsList.Add(widget.NewLabel(fmt.Sprintf(
      "%s: %s", t.Name, t.Description,
    )))
    cw.toolsList.Add(widget.NewSeparator())
  }
  if len(cw.toolsList.Objects) == 0 {
    cw.toolsList.Add(widget.NewLabel("(no tools registered)"))
  }
  cw.toolsList.Refresh()
}


func (cw *MainWindow) refreshToolpacksList() {
  cw.toolpacksList.Objects = nil
  packs := cw.core.Toolpacks()
  if len(packs) == 0 {
    cw.toolpacksList.Add(widget.NewLabelWithStyle(
      "No toolpacks found in ./plugins",
      fyne.TextAlignCenter, fyne.TextStyle{Italic: true}),
    )
  } else {
    for _, name := range packs {
      cw.toolpacksList.Add(container.NewHBox(
        widget.NewLabel(name),
        layout.NewSpacer(),
      ))
    }
  }

}


