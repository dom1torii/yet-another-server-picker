package main

import (
	"log"

	"github.com/dom1torii/cs2-server-manager/internal/config"
	"github.com/dom1torii/cs2-server-manager/internal/tui"
	"github.com/dom1torii/cs2-server-manager/internal/platform/sudo"
)

func main() {
	sudo.CheckIfSudo()
	cfg := config.Init()

	ui := tui.New()
	tui.SetupPages(ui, cfg)

 	if err := ui.Init(); err != nil {
  	log.Fatalln(err)
  }
}
