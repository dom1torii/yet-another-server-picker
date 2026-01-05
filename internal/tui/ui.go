package tui

import (
	"github.com/rivo/tview"
	"github.com/dom1torii/cs2-server-manager/internal/app"
)

type UI struct {
  App   *tview.Application
  Pages *tview.Pages
  State *app.AppState

  RefreshStartList func()
  RefreshSelectPage func()
}

func New() *UI {
	tviewapp := tview.NewApplication()
	pages := tview.NewPages()

	return &UI{
		App: tviewapp,
		Pages: pages,
		State: &app.AppState{},
	}
}

func (ui *UI) Init() error {
  return ui.App.SetRoot(ui.Pages, true).Run()
}
