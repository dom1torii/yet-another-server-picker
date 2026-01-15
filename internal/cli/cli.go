package cli

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/dom1torii/yet-another-server-picker/internal/api"
	"github.com/dom1torii/yet-another-server-picker/internal/config"
	"github.com/dom1torii/yet-another-server-picker/internal/fs"
	"github.com/dom1torii/yet-another-server-picker/internal/ips"
	"github.com/dom1torii/yet-another-server-picker/internal/platform/firewall"
	"github.com/dom1torii/yet-another-server-picker/internal/presets"
)

func IsCLIMode(cfg *config.Config) bool {
	return cfg.ListRelays ||
		len(cfg.SelectRelays) > 0 ||
		cfg.ListPresets ||
		cfg.SelectPreset != "" ||
		cfg.BlockRelays ||
		cfg.UnBlockRelays ||
		cfg.ToBlockCount ||
		cfg.BlockedCount
}

func HandleFlags(cfg *config.Config) {
	response, err := api.FetchRelays()
	if err != nil {
		log.Fatalln("Failed to fetch relays: ", err)
	}

	if cfg.ListRelays {
		keys := make([]string, 0, len(response.Pops))
		for key, pop := range response.Pops {
			keys = append(keys, key+" - "+pop.Desc)
		}
		sort.Strings(keys)

		fmt.Println("Available relays:\n", strings.Join(keys, "\n"))
	}
	if len(cfg.SelectRelays) > 0 {
		// check if at least one of the relays doesn't exist
		for _, name := range cfg.SelectRelays {
			if _, ok := response.Pops[name]; !ok {
				fmt.Fprintf(os.Stderr, "Relay %s doesn't exist\n", name)
				os.Exit(1)
			}
		}

		// write relays that are not selected into ips file
		var ipList []string
		selected := make(map[string]struct{})
		for _, s := range cfg.SelectRelays {
			selected[s] = struct{}{}
		}

		for popName, pop := range response.Pops {
			if _, isSelected := selected[popName]; !isSelected {
				for _, relay := range pop.Relays {
					ipList = append(ipList, relay.Ipv4)
				}
			}
		}
		ips.WriteIpsToFile(ipList, cfg)
	}
	if cfg.ListPresets {
		keys := make([]string, 0, len(presets.Presets))
		for key, item := range presets.Presets {
			keys = append(keys, key+" - "+item.Name)
		}
		sort.Strings(keys)

		fmt.Println("Available presets:\n", strings.Join(keys, "\n"))
	}
	if cfg.SelectPreset != "" {
		// write preset that is not selected to ips file
		preset, ok := presets.Presets[cfg.SelectPreset]
		if !ok {
			fmt.Fprintf(os.Stderr, "This preset doesn't exist\n")
			os.Exit(1)
		}

		var ipList []string
		for popName, pop := range response.Pops {
			if _, inPreset := preset.Pops[popName]; inPreset {
				continue
			}
			for _, relay := range pop.Relays {
				ipList = append(ipList, relay.Ipv4)
			}
		}
		ips.WriteIpsToFile(ipList, cfg)
	}
	if cfg.BlockRelays {
		firewall.BlockIps(cfg, nil)
	}
	if cfg.UnBlockRelays {
		firewall.UnBlockIps(nil)
	}
	if cfg.ToBlockCount {
		fmt.Println(fs.GetFileLineCount(cfg.IpsPath))
	}
	if cfg.BlockedCount {
		fmt.Println(len(firewall.GetBlockedIps()))
	}
}
