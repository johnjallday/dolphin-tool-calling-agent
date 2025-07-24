package main

import (
  "fmt"
  "strings"
	"context"
	"image/color"

  fyneapp "fyne.io/fyne/v2/app"
  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/container"
  "fyne.io/fyne/v2/dialog"
  "fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/theme"

  internalapp "github.com/johnjallday/dolphin-tool-calling-agent/internal/app"
  //"github.com/johnjallday/dolphin-tool-calling-agent/internal/agent"
	//"github.com/johnjallday/dolphin-tool-calling-agent/internal/user"
)

var (
  // our core application
  application = internalapp.NewApp()

  // the Select dropdown for agents
  agentSelect *widget.Select

  // the multi-line, read-only chat history
  chatView *widget.Entry
)

type whiteTextTheme struct {
  fyne.Theme
}

func (w whiteTextTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameForeground {
    return color.White
  }
  // otherwise fall back to the embedded theme
  return w.Theme.Color(name, variant)
}

func main() {
  // initialize your application (load config, etc.)
  if err := application.Init(); err != nil {
    fmt.Println("app.Init:", err)
    return
  }

  fyApp := fyneapp.New()
  w := fyApp.NewWindow("Chat GUI")

  // 1) Agent selector + New button
  agentSelect = widget.NewSelect(nil, func(_ string) {
    onAgentSelected()
  })
  agentSelect.PlaceHolder = "Select an agent…"
  newBtn := widget.NewButton("New", func() {
    showNewAgentForm(w)
  })

  // 2) Chat history view (read-only)
  chatView = widget.NewMultiLineEntry()
  chatView.Disable()
  chatView.SetMinRowsVisible(10)

  // 3) Input + Send button
  input := widget.NewEntry()
  sendBtn := widget.NewButton("Send", func() {
    text := strings.TrimSpace(input.Text)
    if text == "" {
      return
    }
		if err := application.SendMessage(context.Background(), text); err != nil {
      dialog.ShowError(err, w)
      return
    }
    input.SetText("")
    refreshChat()
  })

  // Layout: top agentSelect+newBtn, bottom input+sendBtn, center chatView
  content := container.NewBorder(
    container.NewHBox(agentSelect, newBtn),
    container.NewHBox(input, sendBtn),
    nil, nil,
    chatView,
  )

  w.SetContent(content)
  w.Resize(fyne.NewSize(600, 400))

  // populate the agent list & show window
  refreshAgents()
  w.ShowAndRun()
}

// refreshAgents repopulates agentSelect.Options and refreshes it
func refreshAgents() {
	metas := application.Agents()  
	names := make([]string, len(metas))
	//`:writesnames := application.Agents()
  //agentSelect.Options = names
	for i, m := range metas {
    names[i] = m.Name                       // pull out the Name field
  }
	agentSelect.Options = names
  agentSelect.Refresh()
}

// onAgentSelected loads the chosen agent into the core and refreshes chat
func onAgentSelected() {
  if agentSelect.Selected == "" {
    return
  }
  if err := application.LoadAgent(agentSelect.Selected); err != nil {
    dialog.ShowError(err, nil)
    return
  }
  refreshChat()
}

// refreshChat reads the current agent’s history and writes it into chatView
func refreshChat() {
  a := application.Agent()
  if a == nil {
    chatView.SetText("")
    chatView.Refresh()
    return
  }
  var b strings.Builder
  for _, msg := range a.History() {
    switch msg.Role {
    case "system":
      b.WriteString("[system] ")
    case "user":
      b.WriteString("You: ")
    case "assistant":
      b.WriteString("Bot: ")
    }
    b.WriteString(msg.Content)
    b.WriteByte('\n')
  }
  chatView.SetText(b.String())
  chatView.Refresh()
}

// showNewAgentForm pops up a dialog to create a new agent
func showNewAgentForm(parent fyne.Window) {
  nameEntry  := widget.NewEntry()
  modelEntry := widget.NewEntry()

  form := dialog.NewForm(
    "Create Agent",
    "Create", "Cancel",
    []*widget.FormItem{
      widget.NewFormItem("Name", nameEntry),
      widget.NewFormItem("Model", modelEntry),
    },
    func(ok bool) {
      if !ok {
        return
      }
			meta := internalapp.AgentMeta{
				Name:  nameEntry.Text,
				Model: modelEntry.Text,
			}
			if err := application.CreateAgent(meta); err != nil {
        dialog.ShowError(err, parent)
        return
      }
      // refresh the selector and auto‐select the new one
      refreshAgents()
      agentSelect.SetSelected(nameEntry.Text)
      onAgentSelected()
    },
    parent,
  )
  form.Show()
}
