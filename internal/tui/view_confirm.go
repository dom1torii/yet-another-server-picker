package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/dom1torii/yet-another-server-picker/internal/ips"
)

type confirmModel struct {
	selection bool
}

func (m *model) updateConfirmSelection(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "right", "h", "l":
			m.confirm.selection = !m.confirm.selection
			return m, nil

		case "esc", "q":
			m.state = stateStart
			return m, nil

		case "enter":
			if m.confirm.selection {
				if m.relays.mode == "allow" {
					ips.WriteIpsToFile(m.getUnSelectedIps(), m.cfg)
				} else {
					ips.WriteIpsToFile(m.getSelectedIps(), m.cfg)
				}
				m.state = stateStart
				return m, m.updateStatus()
			}
			m.state = stateRelays
			return m, nil
		}
	}
	return m, nil
}

func (m *model) confirmView() string {
	var yes, no string

	if m.confirm.selection {
		yes = selectionStyle.Render(" YES ")
		no = " NO "
	} else {
		yes = " YES "
		no = selectionStyle.Render(" NO ")
	}

	buttons := lipgloss.JoinHorizontal(lipgloss.Top, yes, "    ", no)

	styledTitle := titleStyle.Align(lipgloss.Center).Render("Your IPs file already contains some IPs.\nAre you sure you wanna proceed?")
	styledHelp := helpStyle.Align(lipgloss.Center).Render("(←→: select | enter: confirm | q/esc: back)")

	view := fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		styledTitle,
		buttons,
		wordwrap.String(styledHelp, m.width),
	)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, view)
}
