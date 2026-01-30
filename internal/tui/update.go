package tui

import (
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case statusMsg:
		m.start.ipsCount = msg.ipsCount
		m.start.blockedCount = msg.blockedCount
		m.start.blockedMap = msg.blockedMap
		return m, m.refreshRelays()

	case presetsMsg:
		m.presets.presetKeys = msg
		return m, nil

	case isFileEmptyMsg:
		if msg {
			m.state = stateStart
			return m, tea.Sequence(
				writeIps(m),
				m.updateStatus(),
			)
		} else {
			m.state = stateConfirm
			return m, nil
		}

	case firewallMsg:
		return m, m.updateStatus()

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		if m.relays.startRow > m.getRows() {
			m.relays.startRow = 0
		}

	case relaysMsg:
		sort.SliceStable(msg, func(i, j int) bool {
			keyI := getStringToSort(msg[i].Desc)
			keyJ := getStringToSort(msg[j].Desc)
			if keyI != keyJ {
				return keyI < keyJ
			}
			return strings.ToLower(msg[i].Desc) < strings.ToLower(msg[j].Desc)
		})
		m.relays.relays = msg
		return m, m.refreshRelays()

	case pingMsg:
		m.relays.pings[msg.index] = msg.duration
		m.relays.pinged++

		// ping in batches of 20
		if m.relays.pinged%20 == 0 && m.relays.pinged < len(m.relays.relays) {
			return m, m.pingBatch(m.relays.pinged)
		}
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			m.Quitting = true
			return m, tea.Quit
		}

		switch m.state {
		case stateStart:
			return m.updateStartSelection(msg)
		case stateRelays:
			return m.updateRelaySelection(msg)
		case stateConfirm:
			return m.updateConfirmSelection(msg)
		case statePresets:
			return m.updatePresetSelection(msg)
		}
	}

	return m, nil
}
