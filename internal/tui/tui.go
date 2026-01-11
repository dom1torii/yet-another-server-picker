package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"

	"github.com/dom1torii/cs2-server-manager/internal/api"
	"github.com/dom1torii/cs2-server-manager/internal/config"
)

var (
	selectionStyle        = lipgloss.NewStyle().Background(lipgloss.Color("0"))
	checkedStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	checkedSelectionStyle = lipgloss.NewStyle().Background(lipgloss.Color("0")).Foreground(lipgloss.Color("2")).Bold(true)
	goodPingStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	badPingStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	blockedPingStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("178"))
	timedoutPingStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
	titleStyle            = lipgloss.NewStyle().MarginLeft(2).Bold(true)
	itemStyle             = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle     = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle       = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle             = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).PaddingLeft(4).PaddingBottom(1)
	quitTextStyle         = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type sessionState int

const (
	stateStart sessionState = iota
	stateRelays
	stateConfirm
	statePresets
)

type model struct {
	cfg             *config.Config
	state           sessionState
	height          int
	width           int
	StartRow        int
	Relays          []api.Pop
	RelaysSelection int
	RelaysChecked   map[int]struct{}
	Pings           map[int]time.Duration
	Pinged          int
	PresetKeys      []string

	StartSelection   int
	ConfirmSelection bool
	PresetSelection  int

	Err      error
	Quitting bool
}

func InitialModel(cfg *config.Config) *model {
	return &model{
		cfg:             cfg,
		RelaysSelection: 0,
		RelaysChecked:   make(map[int]struct{}),
		Pings:           make(map[int]time.Duration),

		StartSelection:  0,
		PresetSelection: 0,

		Quitting: false,
	}
}
