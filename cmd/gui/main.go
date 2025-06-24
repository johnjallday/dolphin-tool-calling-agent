package main

import (
	"log"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Entry Widget")

	input := widget.NewEntry()
	input.SetPlaceHolder("Enter text...")
	// Set the OnSubmitted callback so that pressing Enter prints the content
	input.OnSubmitted = func(text string) {
		log.Println("Content was:", text)
	}

	content := container.NewVBox(input, widget.NewButton("Prompt", func() {
		log.Println("Content was:", input.Text)
	}))

	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}
