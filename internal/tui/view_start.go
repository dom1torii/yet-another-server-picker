package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	tea "github.com/charmbracelet/bubbletea"
)

type startModel struct {
	selection    int
	blockedMap   map[string]bool
	ipsCount     int
	blockedCount int
}

func (m *model) updateStartSelection(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "1":
			m.start.selection = 0
			m.state = stateRelays
			return m, nil
		case "2":
			m.start.selection = 1
			m.state = statePresets
		case "3":
			m.start.selection = 2
			return m, blockIps(m.cfg)
		case "4":
			m.start.selection = 3
			return m, unBlockIps()
		case "q", "esc":
			m.Quitting = true
			return m, tea.Quit
		case "j", "down":
			if m.start.selection < len(startItems)-1 {
				m.start.selection++
			} else {
				m.start.selection = 0
			}
		case "k", "up":
			if m.start.selection > 0 {
				m.start.selection--
			} else {
				m.start.selection = len(startItems) - 1
			}
		case "enter", " ":
			if m.start.selection == 0 {
				m.state = stateRelays
				return m, nil
			}
			if m.start.selection == 1 {
				m.state = statePresets
				return m, nil
			}
			if m.start.selection == 2 {
				return m, blockIps(m.cfg)
			}
			if m.start.selection == 3 {
				return m, unBlockIps()
			}
			if m.start.selection == 4 {
				m.Quitting = true
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m *model) startView() string {
	var startChoices []string
	for i, label := range startItems {
		startChoices = append(startChoices, startItem(label, m.start.selection == i))
		// add status lines
		if i == 2 {
			status := fmt.Sprintf("    %d IP(s) to block", m.start.ipsCount)
			startChoices = append(startChoices, statusStyle.Render(status))
		}
		if i == 3 {
			status := fmt.Sprintf("    %d IP(s) currently blocked", m.start.blockedCount)
			if m.start.blockedCount == 0 {
				startChoices = append(startChoices, statusOkStyle.Render(status))
			} else {
				startChoices = append(startChoices, statusWarningStyle.Render(status))
			}
		}
	}

	items := strings.Join(startChoices, "\n")
	view := fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		wordwrap.String(titleStyle.Render("Yet Another Server Picker"), m.width),
		lipgloss.NewStyle().Width(35).Render(items),
		wordwrap.String(helpStyle.Render("(↓↑: move | space/enter: select | q/esc: quit)"), m.width),
	)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		view,
	)
}

var startItems = []string{
	"(1) Select servers",
	"(2) Use a preset",
	"(3) Block unwanted servers",
	"(4) Unblock all servers",
	"(q) Quit",
}

func startItem(label string, isSelected bool) string {
	if isSelected {
		return selectionStyle.Render(label)
	}
	return label
}
