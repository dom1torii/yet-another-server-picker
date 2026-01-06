package tui

func (ui *UI) SwitchPage(name string) {
	ui.Pages.SwitchToPage(name)

	if name == "start" && ui.RefreshStartList != nil {
		go ui.RefreshStartList()
	}
}
