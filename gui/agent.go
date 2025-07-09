package gui

import (
    "strings"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
    "github.com/johnjallday/dolphin-tool-calling-agent/agent"
)

func NewAgentView(
    w fyne.Window,
    configs []agent.AgentConfig,
    addToHistory func(string),
    showChat func(),
) fyne.CanvasObject {
    localHist := container.NewVBox()
    localScroll := container.NewScroll(localHist)
    addLocal := func(txt string) {
        if txt == "" {
            return
        }
        localHist.Add(widget.NewLabel(txt))
        localHist.Refresh()
        localScroll.ScrollToBottom()
    }

    list := container.NewVBox()
    for _, cfg := range configs {
        cfg := cfg
        list.Add(widget.NewButton(cfg.Name, func() {
            msg := "Loaded agent: " + cfg.Name
            addLocal(msg)
            addToHistory(msg)
        }))
    }

    backBtn := widget.NewButton("◀ Back to Chat", showChat)
    input := widget.NewEntry()
    input.SetPlaceHolder("Chat or type `/back`…")
    sendLocal := func(text string) {
        input.SetText("")
        cmd := strings.TrimSpace(text)
        if cmd == "/back" {
            showChat()
            return
        }
        if cmd == "" {
            return
        }
        addLocal("You: " + cmd)
    }
    sendBtn := widget.NewButton("Send", func() { sendLocal(input.Text) })
    input.OnSubmitted = sendLocal

    bottom := container.NewVBox(input, sendBtn)
    content := container.NewVBox(widget.NewLabel("Agent Chat"), list, localScroll)
    return container.NewBorder(backBtn, bottom, nil, nil, content)
}
