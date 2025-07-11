package main

import (
	"log"
	"context"
  "github.com/johnjallday/dolphin-tool-calling-agent/internal/app"
)

func main() {

  app := app.NewREPLApp()
  if err := app.Init("./configs/settings.toml", "jj"); err != nil {
    log.Fatal(err)
  }
  ctx := context.Background()
  if err := app.Run(ctx); err != nil {
    log.Fatal(err)
  }
  _ = app.Shutdown()
}
