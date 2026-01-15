package tui

import (
	"github.com/charmbracelet/lipgloss"
	"time"

	"github.com/dom1torii/yet-another-server-picker/internal/api"
	"github.com/dom1torii/yet-another-server-picker/internal/config"
)

var (
	selectionStyle        = lipgloss.NewStyle().Background(lipgloss.Color("8"))
	checkedStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	checkedSelectionStyle = lipgloss.NewStyle().Background(lipgloss.Color("8")).Foreground(lipgloss.Color("2")).Bold(true)
	goodPingStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	badPingStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	blockedPingStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("178"))
	timedoutPingStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
	titleStyle            = lipgloss.NewStyle().MarginLeft(2).Bold(true)
	statusStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	statusOkStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	statusWarningStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	helpStyle             = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).PaddingLeft(4).PaddingBottom(1)
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
	BlockedMap      map[string]bool
	PresetKeys      []string

	IpsCount     int
	BlockedCount int

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
		BlockedMap:      make(map[string]bool),
		Pings:           make(map[int]time.Duration),

		StartSelection:  0,
		PresetSelection: 0,

		Quitting: false,
	}
}
