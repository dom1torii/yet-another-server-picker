package tui

import (
	"log"
	"sort"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dom1torii/yet-another-server-picker/internal/api"
	"github.com/dom1torii/yet-another-server-picker/internal/config"
	"github.com/dom1torii/yet-another-server-picker/internal/fs"
	"github.com/dom1torii/yet-another-server-picker/internal/ips"
	"github.com/dom1torii/yet-another-server-picker/internal/platform/firewall"
	"github.com/dom1torii/yet-another-server-picker/internal/presets"
)

func (m *model) Init() tea.Cmd {
	return tea.Batch(
		getRelays(),
		loadPresets(),
		m.updateStatus(),
		tea.SetWindowTitle("Yet Another Server Picker"),
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

func getPingToIp(index int, ip string, isBlocked bool) tea.Cmd {
	return func() tea.Msg {
		if isBlocked {
			return pingMsg{index: index, duration: -1}
		}
		duration := ips.GetPing(ip)
		return pingMsg{index: index, duration: duration}
	}
}

func (m *model) pingBatch(startIndex int) tea.Cmd {
	batchSize := 20
	var cmds []tea.Cmd

	for i := startIndex; i < startIndex+batchSize && i < len(m.Relays); i++ {
		if len(m.Relays[i].Relays) > 0 {
			ip := m.Relays[i].Relays[0].Ipv4
			isBlocked := m.BlockedMap[ip]
			cmds = append(cmds, getPingToIp(i, ip, isBlocked))
		}
	}
	return tea.Batch(cmds...)
}

func (m *model) refreshRelays() tea.Cmd {
	m.Pings = make(map[int]time.Duration)
	m.Pinged = 0
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

func (m *model) updateStatus() tea.Cmd {
	return func() tea.Msg {
		ipsCount := fs.GetFileLineCount(m.cfg.IpsPath)
		blockedMap := firewall.GetBlockedIps()
		return statusMsg{
			ipsCount:     ipsCount,
			blockedCount: len(blockedMap),
			blockedMap:   blockedMap,
		}
	}
}

func writeIps(m *model) tea.Cmd {
	return func() tea.Msg {
		ips.WriteIpsToFile(m.getUnSelectedIps(), m.cfg)
		return firewallMsg{}
	}
}
