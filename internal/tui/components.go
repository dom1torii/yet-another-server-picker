package tui

import (
	"github.com/rivo/tview"
)

func (ui *UI) ConfirmOverwrite(callback func()) {
	modal := tview.NewModal().
		SetText("Your ips file isn't empty.\nDo you want to overwrite it?").
		AddButtons([]string{"Yes", "No"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" {
				callback()
			}
			ui.SwitchPage("start")
		})
	ui.Pages.AddAndSwitchToPage("confirm", modal, true)
}
