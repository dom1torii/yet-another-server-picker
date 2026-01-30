package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/dom1torii/yet-another-server-picker/internal/presets"
)

type presetsModel struct {
	selection int
	presetKeys []string
}

func (m *model) updatePresetSelection(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			m.state = stateStart
			return m, nil
		case "j", "down":
			if m.presets.selection < len(m.presets.presetKeys)-1 {
				m.presets.selection++
			} else {
				m.presets.selection = 0
			}
		case "k", "up":
			if m.presets.selection > 0 {
				m.presets.selection--
			} else {
				m.presets.selection = len(m.presets.presetKeys) - 1
			}
		case "enter", " ":
			selectedKey := m.presets.presetKeys[m.presets.selection]
			preset := presets.Presets[selectedKey]

			m.relays.checked = make(map[int]struct{})
			// find relays that match our selected preset
			for i, relay := range m.relays.relays {
				if _, found := preset.Pops[relay.Key]; found {
					m.relays.checked[i] = struct{}{}
				}
			}
			m.state = stateRelays
			return m, nil
		}
	}
	return m, nil
}

func (m *model) presetsView() string {
	if len(m.presets.presetKeys) == 0 {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, "Loading presets...")
	}

	var choices []string
	for i, key := range m.presets.presetKeys {
		p := presets.Presets[key]
		choices = append(choices, startItem(p.Name, m.presets.selection == i))
	}

	items := strings.Join(choices, "\n")
	styledTitle := titleStyle.Align(lipgloss.Center).Render("Presets")

	view := fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		styledTitle,
		lipgloss.NewStyle().Width(15).Render(items),
		wordwrap.String(helpStyle.Render("(↓↑: move | space/enter: select | q/esc: quit)"), m.width),
	)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, view)
}
