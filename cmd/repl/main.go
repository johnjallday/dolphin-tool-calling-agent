package main

import (
	"context"
	"log"

  "github.com/johnjallday/dolphin-tool-calling-agent/internal/app"
)

func main() {
	repl, err := app.NewREPLApp("./configs/settings.toml")
	if err != nil {
		log.Fatalf("failed to initialize REPL: %v", err)
	}

	if err := repl.Run(context.Background()); err != nil {
		log.Fatalf("REPL error: %v", err)
	}
}
