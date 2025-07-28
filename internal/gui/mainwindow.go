package gui

import (
  //"context"
  //"fmt"

  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/container"
  "fyne.io/fyne/v2/layout"
  "fyne.io/fyne/v2/widget"

  "github.com/johnjallday/dolphin-tool-calling-agent/internal/app"
)

// MainWindow holds our Fyne window, the core App, and all tabs.
type MainWindow struct {
  app   fyne.App
  wnd   fyne.Window
  core  app.App

  mainTabs *container.AppTabs

  // cached TabItems
  chatTab  *container.TabItem
  toolsTab *container.TabItem
  agentTab *container.TabItem
  userTab  *container.TabItem

  // chat widgets
  historyBox    *fyne.Container
  historyScroll *container.Scroll
  inputEntry    *widget.Entry

  // top bar
  statusLabel *widget.Label

  // tools widgets
  toolsList     *fyne.Container
  toolpacksList *fyne.Container

  // agent widgets
  agentList *fyne.Container

  // user widgets
  userSelect        *widget.Select
  newUserEntry      *widget.Entry
  userNameEntry     *widget.Entry
  userDefaultSelect *widget.Select
	userList *fyne.Container
}


// NewChatWindow builds the window and all tabs exactly once.
func NewMainWindow(fy fyne.App, core app.App) *MainWindow {
  w := fy.NewWindow("🐬 Dolphin Chat 🐬")
  cw := &MainWindow{app: fy, wnd: w, core: core}

  // cache each TabItem
  cw.chatTab = cw.makeChatTab()
  cw.toolsTab = cw.makeToolsTab()
  cw.agentTab = cw.makeAgentTab()
  cw.userTab = cw.makeUserTab()

  // put them into AppTabs
  cw.mainTabs = container.NewAppTabs(
    cw.chatTab, cw.toolsTab, cw.agentTab, cw.userTab,
  )
  cw.mainTabs.SetTabLocation(container.TabLocationTop)

  // top bar with status
  cw.statusLabel = widget.NewLabel("")
  topBar := container.NewHBox(cw.statusLabel, layout.NewSpacer())

  // assemble
  content := container.NewBorder(topBar, nil, nil, nil, cw.mainTabs)
  w.SetContent(content)
  w.Resize(fyne.NewSize(700, 500))

  // initial fill
  cw.RefreshAll()

  return cw
}

// ShowAndRun pops up the window.
func (cw *MainWindow) ShowAndRun() {
  cw.wnd.ShowAndRun()
}


// buildRoot returns either an onboarding pane or the normal
// tab view + status bar.
func (cw *MainWindow) buildRoot() fyne.CanvasObject {
  // No user at all?  Show the “create first user” box.
  if len(cw.core.Users()) == 0 {
    return cw.createOnboardingBox()
  }
  // We have a user, but no agents?  Show the “create first agent” box.
  if len(cw.core.Agents()) == 0 {
    return cw.createAgentOnboardingBox()
  }
  // Otherwise fall back to your normal tabs + top‐bar layout:
  topBar := container.NewHBox(cw.statusLabel, layout.NewSpacer())
  return container.NewBorder(topBar, nil, nil, nil, cw.mainTabs)
}


// RefreshAll updates every pane in the window.
func (cw *MainWindow) RefreshAll() {
  // 1) status bar
  cw.refreshUserStatus()

  // 2) chat history (in case underlying messages have changed)
  cw.historyBox.Refresh()
  cw.historyScroll.ScrollToBottom()

  // 3) tools
  cw.refreshCurrentToolsList()
  cw.refreshToolpacksList()

  // 4) agent
  cw.agentTab.Content = cw.buildAgentPane()
  cw.agentTab.Content.Refresh()

  // 5) user
  cw.userTab.Content = cw.buildUserPane()
  cw.userTab.Content.Refresh()
}

