package main

import (
    "fmt"
    "strings"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
    "github.com/johnjallday/dolphin-tool-calling-agent/agent"
    "github.com/johnjallday/dolphin-tool-calling-agent/location"
)

func newAgentView(
    w fyne.Window,
    configs []agent.AgentConfig,
    addToHistory func(string),
    showChat func(),
) fyne.CanvasObject {
    localHist := container.NewVBox()
    localScroll := container.NewScroll(localHist)
    addLocal := func(txt string) {
        if txt == "" { return }
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

func main() {
    a := app.New()
    w := a.NewWindow("Chat with Agents")

    configs, err := agent.ListAgents()
    if err != nil {
        fmt.Println("load agents:", err)
    }

    var chatView, agentView fyne.CanvasObject
    globalHist := container.NewVBox()
    globalScroll := container.NewScroll(globalHist)
    input := widget.NewEntry()
    input.SetPlaceHolder("Type message or `/agent` or `/help`…")

    addToHistory := func(txt string) {
        if txt == "" { return }
        globalHist.Add(widget.NewLabel(txt))
        globalHist.Refresh()
        globalScroll.ScrollToBottom()
    }

    commands := map[string]string{
        "/agent":  "Switch to agent selection",
        "/help":   "List available commands",
        "/device": "List connected devices",
    }

    showChat := func() { w.SetContent(chatView) }
    agentView = newAgentView(w, configs, addToHistory, showChat)

    handle := func(txt string) {
        input.SetText("")
        cmd := strings.TrimSpace(txt)
        switch cmd {
        case "/agent":
            w.SetContent(agentView)
        case "/help":
            for k, d := range commands {
                addToHistory(fmt.Sprintf("%-7s %s", k, d))
            }
        case "/device":
						if ml, err := location.GetMyLocation(); err != nil {
									addToHistory("Location error: " + err.Error())
							} else {
									addToHistory("Location: " + ml.Name)
									addToHistory("Audio devices: " + strings.Join(ml.AudioDevices, ", "))
									addToHistory("Displays: " + strings.Join(ml.Displays, ", "))
									addToHistory("Networks: " + strings.Join(ml.Network, ", "))
							}
        default:
            if cmd != "" {
                addToHistory("You: " + cmd)
            }
        }
    }

    sendBtn := widget.NewButton("Send", func() { handle(input.Text) })
    bottom := container.NewVBox(input, sendBtn)
    chatView = container.NewBorder(nil, bottom, nil, nil, globalScroll)
    input.OnSubmitted = handle

    w.SetContent(chatView)
    w.Resize(fyne.NewSize(500, 400))

    // print location and its devices
    if ml, err := location.GetMyLocation(); err != nil {
        addToHistory("Location error: " + err.Error())
    } else {
        addToHistory("Location: " + ml.Name)
        addToHistory("Audio devices: " + strings.Join(ml.AudioDevices, ", "))
        addToHistory("Displays: " + strings.Join(ml.Displays, ", "))
        addToHistory("Networks: " + strings.Join(ml.Network, ", "))
    }

    w.ShowAndRun()
}
