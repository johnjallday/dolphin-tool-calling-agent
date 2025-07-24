package main

import (
    "fmt"

    fyneapp "fyne.io/fyne/v2/app"
		"fyne.io/fyne/v2/theme"
    "github.com/johnjallday/dolphin-tool-calling-agent/internal/app"
    "github.com/johnjallday/dolphin-tool-calling-agent/internal/gui"
)

func main() {
    core := app.NewApp()       // core is an app.App interface
    if err := core.Init(); err != nil {
        fmt.Println("init error:", err)
        return
    }

    fy := fyneapp.New()
    //fy.Settings().SetTheme(gui.NewWhiteTextTheme())
		fy.Settings().SetTheme(gui.NewGreyedTextTheme(theme.DarkTheme()))

    w := gui.NewChatWindow(fy, core)
    w.ShowAndRun()
}
