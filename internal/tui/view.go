package tui

func (m *model) View() string {
	switch m.state {
	case stateStart:
		return m.startView()
	case stateRelays:
		return m.relaysView()
	case stateConfirm:
		return m.confirmView()
	case statePresets:
		return m.presetsView()
	default:
		return ""
	}
}
