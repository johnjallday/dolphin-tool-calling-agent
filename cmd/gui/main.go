package main

import (
    "fmt"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
    "github.com/johnjallday/dolphin-tool-calling-agent/agent"
    "github.com/johnjallday/dolphin-tool-calling-agent/gui"
)

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
        if txt == "" {
            return
        }
        globalHist.Add(widget.NewLabel(txt))
        globalHist.Refresh()
        globalScroll.ScrollToBottom()
    }

    commands := map[string]string{
        "/agent":    "Switch to agent selection",
        "/help":     "List available commands",
        "/location": "Show current location & devices",
    }

    showChat := func() { w.SetContent(chatView) }
    agentView = gui.NewAgentView(w, configs, addToHistory, showChat)

    handle := func(txt string) {
        input.SetText("")
        cmd := txt
        switch cmd {
        case "/agent":
            w.SetContent(agentView)
        case "/help":
            for k, d := range commands {
                addToHistory(fmt.Sprintf("%-9s %s", k, d))
            }
        case "/location":
            gui.ShowLocation(addToHistory)
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

    gui.ShowLocation(addToHistory)
    w.ShowAndRun()
}
