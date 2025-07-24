package main

import (

  "github.com/johnjallday/dolphin-tool-calling-agent/internal/toolmanager"
	"log"
)

func main() {
  if err := toolmanager.CheckVersion(); err != nil {
    log.Fatal(err)
  }
}
