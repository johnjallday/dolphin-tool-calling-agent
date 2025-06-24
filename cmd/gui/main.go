package main

import (
    "fmt"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"

    "Dolphin-Tool-Calling-Agent/tools" // adjust import path

)

func main() {
    a := app.New()
    w := a.NewWindow("Dolphin LLM Chatbot GUI")

    input := widget.NewEntry()
    input.SetPlaceHolder("Type your question...")

    output := widget.NewLabel("")
    scroll := container.NewScroll(output)
    scroll.SetMinSize(fyne.NewSize(400, 200))

    input.OnSubmitted = func(text string) {
        response, err := tools.HandleQuestion(text)
        if err != nil {
            output.SetText(fmt.Sprintf("Error: %v", err))
            return
        }
        output.SetText(response)
        input.SetText("")
    }

    w.SetContent(container.NewVBox(
        scroll,
        input,
    ))

    w.ShowAndRun()
}
