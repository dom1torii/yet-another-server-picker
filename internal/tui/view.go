package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"

	"github.com/dom1torii/yet-another-server-picker/internal/presets"
)

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

func (m *model) startView() string {
	var startChoices []string
	for i, label := range startItems {
		startChoices = append(startChoices, startItem(label, m.StartSelection == i))
		// add status lines
		if i == 2 {
			status := fmt.Sprintf("    %d IP(s) to block", m.IpsCount)
			startChoices = append(startChoices, statusStyle.Render(status))
		}
		if i == 3 {
			status := fmt.Sprintf("    %d IP(s) currently blocked", m.BlockedCount)
			if m.BlockedCount == 0 {
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

func checkbox(label string, pingMs string, isSelected bool, isChecked bool, mode string) string {
	var indicator string
	var indicatorStyle, labelStyle lipgloss.Style

	if mode == "block" {
		if isChecked {
			indicator = "[x]"
			indicatorStyle = crossedStyle
		} else {
			indicator = "[ ]"
		}
		if isSelected {
			indicatorStyle = crossedSelectionStyle
			labelStyle = selectionStyle
		}
	} else {
		if isChecked {
			indicator = "[✓]"
			indicatorStyle = checkedStyle
		} else {
			indicator = "[ ]"
		}
		if isSelected {
			indicatorStyle = checkedSelectionStyle
			labelStyle = selectionStyle
		}
	}

	if isSelected && !isChecked {
		indicatorStyle = selectionStyle
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Center,
		indicatorStyle.Render(indicator),
		labelStyle.Render(" "+label),
		" ",
		pingMs,
	)
}

func (m *model) relaysView() string {
	if len(m.Relays) == 0 {
		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			"Loading servers data...",
		)
	}

	currentMode := ""

	if m.Mode == "allow" {
		currentMode = modeAllowStyle.Render("allow")
	} else {
		currentMode = modeBlockStyle.Render("block")
	}

	// create checkboxes
	var checkboxes []string
	for i, pop := range m.Relays {
		_, checked := m.RelaysChecked[i]

		label := pop.Desc
		maxLabelWidth := (m.width / 2) - 20
		if len(label) > maxLabelWidth && maxLabelWidth > 3 {
			label = label[:maxLabelWidth-3] + "..."
		}

		pingDur, found := m.Pings[i]
		pingDisplay := "(...)" // display ... if its still loading

		if found {
			if pingDur > 0 {
				ms := pingDur.Milliseconds()
				rawPing := fmt.Sprintf("(%dms)", ms)

				if ms < 100 {
					pingDisplay = goodPingStyle.Render(rawPing)
				} else {
					pingDisplay = badPingStyle.Render(rawPing)
				}
			} else if pingDur == -1 {
				pingDisplay = blockedPingStyle.Render("(blocked)")
			} else {
				pingDisplay = timedoutPingStyle.Render("(timed out)")
			}
		}

		checkboxes = append(checkboxes, checkbox(label, pingDisplay, m.RelaysSelection == i, checked, m.Mode))
	}

	// create 2 columns
	rows := m.getRows()
	maxVisibleRows := max(m.height-8, 1)
	currentRow := m.RelaysSelection % rows
	if currentRow < m.StartRow {
		m.StartRow = currentRow
	} else if currentRow >= m.StartRow+maxVisibleRows {
		m.StartRow = currentRow - maxVisibleRows + 1
	}

	if rows > maxVisibleRows && m.StartRow > rows-maxVisibleRows {
		m.StartRow = rows - maxVisibleRows
	}

	if m.StartRow < 0 || rows <= maxVisibleRows {
		m.StartRow = 0
	}

	mid := rows
	leftEnd := min(m.StartRow+maxVisibleRows, mid)
	leftVisible := checkboxes[m.StartRow:leftEnd]

	var rightVisible []string
	if mid < len(checkboxes) {
		rightTotal := checkboxes[mid:]
		rStart := m.StartRow
		rEnd := m.StartRow + maxVisibleRows

		if rStart < len(rightTotal) {
			if rEnd > len(rightTotal) {
				rEnd = len(rightTotal)
			}
			rightVisible = rightTotal[rStart:rEnd]
		}
	}

	// join columns
	leftCol := lipgloss.NewStyle().PaddingRight(4).Render(lipgloss.JoinVertical(lipgloss.Left, leftVisible...))
	rightCol := lipgloss.JoinVertical(lipgloss.Left, rightVisible...)
	columns := lipgloss.JoinHorizontal(lipgloss.Top, leftCol, rightCol)

	// indicator if you can scroll up or down
	topIndicator := " "
	bottomIndicator := " "

	if m.StartRow > 0 {
		topIndicator = "↑ more"
	}
	if m.StartRow+maxVisibleRows < rows {
		bottomIndicator = "↓ more"
	}

	topIndicator = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render(topIndicator)
	bottomIndicator = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render(bottomIndicator)

	view := fmt.Sprintf(
		"%s\n%s\n\n%s\n\n%s\n%s",
		titleStyle.Render("Select servers to ")+currentMode,
		topIndicator,
		columns,
		bottomIndicator,
		wordwrap.String(helpStyle.Render("(←↓↑→: move | space: select | enter: apply | t: toggle block/allow | q/esc: back)"), m.width),
	)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		view,
	)
}

func (m *model) confirmView() string {
	var yes, no string

	if m.ConfirmSelection {
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

func (m *model) presetsView() string {
	if len(m.PresetKeys) == 0 {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, "Loading presets...")
	}

	var choices []string
	for i, key := range m.PresetKeys {
		p := presets.Presets[key]
		choices = append(choices, startItem(p.Name, m.PresetSelection == i))
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
