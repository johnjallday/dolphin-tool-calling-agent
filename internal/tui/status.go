package tui

import (
  	"github.com/fatih/color"
		"fmt"
)



func (t *TUIApp) StatusCmd() error {
    u, a := t.App.User(), t.App.Agent()

    userLoaded  := u != nil
    agentLoaded := a != nil

    var cmdList string
    switch {
    case !userLoaded && !agentLoaded:
        cmdList = "load-user | load-agent | users | create-agent"
    case userLoaded && !agentLoaded:
        cmdList = "unload-user | load-agent | switch-user | users | agents"
    default: // agentLoaded (with or without user)
        cmdList = "tools | unload-user | unload-agent | switch-agent | help"
    }

    cLabel := color.New(color.FgCyan, color.Bold)
    cValue := color.New(color.FgWhite)

    cLabel.Fprint(t.Out, "Commands:  ")
    cValue.Fprintln(t.Out, cmdList)
    fmt.Fprintln(t.Out)

    cLabel.Fprint(t.Out, "Current User:  ")
    if userLoaded {
        cValue.Fprintln(t.Out, u.Name)
    } else {
        cValue.Fprintln(t.Out, "<none>")
    }

    cLabel.Fprint(t.Out, "Current Agent: ")
    if agentLoaded {
        cValue.Fprintln(t.Out, a.Name)
    } else {
        cValue.Fprintln(t.Out, "<none>")
    }

    return nil
}


