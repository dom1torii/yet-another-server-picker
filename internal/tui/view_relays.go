package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/dom1torii/yet-another-server-picker/internal/api"
)

type relaysModel struct {
	selection int
	checked   map[int]struct{}
	mode      string
	pings     map[int]time.Duration
	pinged    int
	relays    []api.Pop
	startRow  int
}

func (m *model) updateRelaySelection(msg tea.Msg) (tea.Model, tea.Cmd) {
	rows := m.getRows()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			m.state = stateStart
			return m, nil
		case "t":
			if m.relays.mode == "allow" {
				m.relays.mode = "block"
			} else {
				m.relays.mode = "allow"
			}
			return m, nil
		case "j", "down":
			if m.relays.selection < len(m.relays.relays)-1 {
				m.relays.selection++
			}
		case "k", "up":
			if m.relays.selection > 0 {
				m.relays.selection--
			}
		case "h", "left":
			if m.relays.selection >= rows {
				m.relays.selection -= rows
			}
		case "l", "right":
			if m.relays.selection+rows < len(m.relays.relays) {
				m.relays.selection += rows
			}
		case " ":
			_, ok := m.relays.checked[m.relays.selection]
			if ok {
				delete(m.relays.checked, m.relays.selection)
			} else {
				m.relays.checked[m.relays.selection] = struct{}{}
			}
		case "enter":
			return m, isFileEmpty(m.cfg.Ips.Path)
		}

	}
	return m, nil
}

func (m *model) relaysView() string {
	if len(m.relays.relays) == 0 {
		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			"Loading servers data...",
		)
	}

	currentMode := ""

	if m.relays.mode == "allow" {
		currentMode = modeAllowStyle.Render("allow")
	} else {
		currentMode = modeBlockStyle.Render("block")
	}

	// create checkboxes
	var checkboxes []string
	for i, pop := range m.relays.relays {
		_, checked := m.relays.checked[i]

		label := pop.Desc
		maxLabelWidth := (m.width / 2) - 20
		if len(label) > maxLabelWidth && maxLabelWidth > 3 {
			label = label[:maxLabelWidth-3] + "..."
		}

		pingDur, found := m.relays.pings[i]
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

		checkboxes = append(checkboxes, checkbox(label, pingDisplay, m.relays.selection == i, checked, m.relays.mode))
	}

	// create 2 columns
	rows := m.getRows()
	maxVisibleRows := max(m.height-8, 1)
	currentRow := m.relays.selection % rows
	if currentRow < m.relays.startRow {
		m.relays.startRow = currentRow
	} else if currentRow >= m.relays.startRow+maxVisibleRows {
		m.relays.startRow = currentRow - maxVisibleRows + 1
	}

	if rows > maxVisibleRows && m.relays.startRow > rows-maxVisibleRows {
		m.relays.startRow = rows - maxVisibleRows
	}

	if m.relays.startRow < 0 || rows <= maxVisibleRows {
		m.relays.startRow = 0
	}

	mid := rows
	leftEnd := min(m.relays.startRow+maxVisibleRows, mid)
	leftVisible := checkboxes[m.relays.startRow:leftEnd]

	var rightVisible []string
	if mid < len(checkboxes) {
		rightTotal := checkboxes[mid:]
		rStart := m.relays.startRow
		rEnd := m.relays.startRow + maxVisibleRows

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

	if m.relays.startRow > 0 {
		topIndicator = "↑ more"
	}
	if m.relays.startRow+maxVisibleRows < rows {
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
