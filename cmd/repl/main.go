package main

import (
	"log"
	"context"
  "github.com/johnjallday/dolphin-tool-calling-agent/internal/app"
)


func main() {
    ctx := context.Background()

    app := app.NewREPLApp(
        &app.OSUserAuth{},
        &app.TOMLUserLoader{},
        &app.REPLAgentManager{},
    )

    if err := app.Init(ctx); err != nil {
        log.Fatal(err)
    }
    if err := app.Run(ctx); err != nil {
        log.Fatal(err)
    }
}
