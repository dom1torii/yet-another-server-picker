package tui

import (
	"sort"
	"strings"

	"github.com/dom1torii/yet-another-server-picker/internal/ips"
	"github.com/dom1torii/yet-another-server-picker/internal/presets"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case statusMsg:
		m.IpsCount = msg.ipsCount
		m.BlockedCount = msg.blockedCount
		m.BlockedMap = msg.blockedMap
		return m, m.refreshRelays()

	case presetsMsg:
		m.PresetKeys = msg
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

		if m.StartRow > m.getRows() {
			m.StartRow = 0
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
		m.Relays = msg
		return m, m.refreshRelays()

	case pingMsg:
		m.Pings[msg.index] = msg.duration
		m.Pinged++

		// ping in batches of 20
		if m.Pinged%20 == 0 && m.Pinged < len(m.Relays) {
			return m, m.pingBatch(m.Pinged)
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

func (m *model) updateStartSelection(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "1":
			m.StartSelection = 0
			m.state = stateRelays
			return m, nil
		case "2":
			m.StartSelection = 1
			m.state = statePresets
		case "3":
			m.StartSelection = 2
			return m, blockIps(m.cfg)
		case "4":
			m.StartSelection = 3
			return m, unBlockIps()
		case "q", "esc":
			m.Quitting = true
			return m, tea.Quit
		case "j", "down":
			if m.StartSelection < len(startItems)-1 {
				m.StartSelection++
			} else {
				m.StartSelection = 0
			}
		case "k", "up":
			if m.StartSelection > 0 {
				m.StartSelection--
			} else {
				m.StartSelection = len(startItems) - 1
			}
		case "enter", " ":
			if m.StartSelection == 0 {
				m.state = stateRelays
				return m, nil
			}
			if m.StartSelection == 1 {
				m.state = statePresets
				return m, nil
			}
			if m.StartSelection == 2 {
				return m, blockIps(m.cfg)
			}
			if m.StartSelection == 3 {
				return m, unBlockIps()
			}
			if m.StartSelection == 4 {
				m.Quitting = true
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m *model) updateRelaySelection(msg tea.Msg) (tea.Model, tea.Cmd) {
	rows := m.getRows()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			m.state = stateStart
			return m, nil
		case "j", "down":
			if m.RelaysSelection < len(m.Relays)-1 {
				m.RelaysSelection++
			}
		case "k", "up":
			if m.RelaysSelection > 0 {
				m.RelaysSelection--
			}
		case "h", "left":
			if m.RelaysSelection >= rows {
				m.RelaysSelection -= rows
			}
		case "l", "right":
			if m.RelaysSelection+rows < len(m.Relays) {
				m.RelaysSelection += rows
			}
		case " ":
			_, ok := m.RelaysChecked[m.RelaysSelection]
			if ok {
				delete(m.RelaysChecked, m.RelaysSelection)
			} else {
				m.RelaysChecked[m.RelaysSelection] = struct{}{}
			}
		case "enter":
			return m, isFileEmpty(m.cfg.IpsPath)
		}

	}
	return m, nil
}

func (m *model) updateConfirmSelection(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "right", "h", "l":
			m.ConfirmSelection = !m.ConfirmSelection
			return m, nil

		case "esc", "q":
			m.state = stateStart
			return m, nil

		case "enter":
			if m.ConfirmSelection {
				ips.WriteIpsToFile(m.getUnSelectedIps(), m.cfg)
				m.state = stateStart
				return m, m.updateStatus()
			}
			m.state = stateRelays
			return m, nil
		}
	}
	return m, nil
}

func (m *model) updatePresetSelection(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			m.state = stateStart
			return m, nil
		case "j", "down":
			if m.PresetSelection < len(m.PresetKeys)-1 {
				m.PresetSelection++
			} else {
				m.PresetSelection = 0
			}
		case "k", "up":
			if m.PresetSelection > 0 {
				m.PresetSelection--
			} else {
				m.PresetSelection = len(m.PresetKeys) - 1
			}
		case "enter", " ":
			selectedKey := m.PresetKeys[m.PresetSelection]
			preset := presets.Presets[selectedKey]

			m.RelaysChecked = make(map[int]struct{})
			// find relays that match our selected preset
			for i, relay := range m.Relays {
				if _, found := preset.Pops[relay.Key]; found {
					m.RelaysChecked[i] = struct{}{}
				}
			}
			m.state = stateRelays
			return m, nil
		}
	}
	return m, nil
}
