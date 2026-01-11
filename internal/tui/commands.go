package tui

import (
	"log"
	"sort"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dom1torii/cs2-server-manager/internal/api"
	"github.com/dom1torii/cs2-server-manager/internal/config"
	"github.com/dom1torii/cs2-server-manager/internal/fs"
	"github.com/dom1torii/cs2-server-manager/internal/ips"
	"github.com/dom1torii/cs2-server-manager/internal/platform/firewall"
	"github.com/dom1torii/cs2-server-manager/internal/presets"
)

func (m *model) Init() tea.Cmd {
	return tea.Batch(
		getRelays(),
		loadPresets(),
		tea.SetWindowTitle("CS2 Server Manager"),
	)
}

func getRelays() tea.Cmd {
	return func() tea.Msg {
		resp, err := api.FetchRelays()
		if err != nil {
			log.Fatalln("Error fetching relays: ", err)
		}

		var pops []api.Pop
		for _, pop := range resp.Pops {
			pops = append(pops, pop)
		}
		return relaysMsg(pops)
	}
}

func getPingToIp(index int, ip string) tea.Cmd {
	return func() tea.Msg {
		duration := ips.GetPing(ip)
		return pingMsg{index: index, duration: duration}
	}
}

func (m *model) pingBatch(startIndex int) tea.Cmd {
	batchSize := 20
	var cmds []tea.Cmd

	for i := startIndex; i < startIndex+batchSize && i < len(m.Relays); i++ {
		// only ping if its not blocked
		if _, exists := m.Pings[i]; !exists {
			if len(m.Relays[i].Relays) > 0 {
				cmds = append(cmds, getPingToIp(i, m.Relays[i].Relays[0].Ipv4))
			}
		}
	}
	return tea.Batch(cmds...)
}

func (m *model) refreshRelays() tea.Cmd {
	m.Pings = make(map[int]time.Duration)
	m.Pinged = 0

	for i, pop := range m.Relays {
		// if the ip is blocked, add is as -1 so we can use it
		if len(pop.Relays) > 0 && firewall.IsIpBlocked(pop.Relays[0].Ipv4) {
			m.Pings[i] = -1
			m.Pinged++
		}
	}
	return m.pingBatch(0)
}

func blockIps(cfg *config.Config) tea.Cmd {
	return func() tea.Msg {
		firewall.BlockIps(cfg, nil)
		return firewallMsg{}
	}
}

func unBlockIps() tea.Cmd {
	return func() tea.Msg {
		firewall.UnBlockIps(nil)
		return firewallMsg{}
	}
}

func isFileEmpty(path string) tea.Cmd {
	return func() tea.Msg {
		empty := fs.IsFileEmpty(path)
		return isFileEmptyMsg(empty)
	}
}

func loadPresets() tea.Cmd {
	return func() tea.Msg {
		keys := make([]string, 0, len(presets.Presets))
		for k := range presets.Presets {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		return presetsMsg(keys)
	}
}
