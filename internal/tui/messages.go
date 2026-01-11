package tui

import (
	"time"

	"github.com/dom1torii/cs2-server-manager/internal/api"
)

type firewallMsg struct{}

type relaysMsg []api.Pop

type presetsMsg []string

type isFileEmptyMsg bool

type pingMsg struct {
	index    int
	duration time.Duration
}
