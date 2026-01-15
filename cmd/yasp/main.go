package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dom1torii/yet-another-server-picker/internal/cli"
	"github.com/dom1torii/yet-another-server-picker/internal/config"
	"github.com/dom1torii/yet-another-server-picker/internal/platform/sudo"
	"github.com/dom1torii/yet-another-server-picker/internal/tui"
)

func main() {
	sudo.CheckIfSudo()
	cfg := config.Init()

	if cli.IsCLIMode(cfg) {
		cli.HandleFlags(cfg)
		return
	}

	p := tea.NewProgram(tui.InitialModel(cfg), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalln("Error starting bubbletea: ", err)
		os.Exit(1)
	}
}
