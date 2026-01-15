//go:build linux

package firewall

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/dom1torii/yet-another-server-picker/internal/config"
)

func BlockIps(cfg *config.Config, onDone func()) {
	UnBlockIps(onDone)
	ipsFile := cfg.IpsPath

	chain := "OUTPUT"
	customChain := "CS2_BLOCKLIST"

	// create custom chain for cs2
	createChain := exec.Command("iptables", "-N", customChain)
	if err := createChain.Run(); err != nil {
		log.Fatalln("Failed to create custom chain: ")
	}

	// tie our chain to OUTPUT
	if err := exec.Command("iptables", "-C", chain, "-j", customChain).Run(); err != nil {
		if err := exec.Command("iptables", "-A", chain, "-j", customChain).Run(); err != nil {
			log.Fatalln("Failed to tie chains: ", err)
		}
	}

	file, err := os.Open(ipsFile)
	if err != nil {
		log.Fatalln("Failed to open a file containing ips: ", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatalln("Failed to close file: ")
		}
	}()

	// read ips from a file and add them to commands 1 by 1
	var cmds []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ip := strings.TrimSpace(scanner.Text())
		if ip == "" {
			continue
		}

		// if rule doesn't exist, add it to cmds
		checkCmd := exec.Command("iptables", "-C", customChain, "-d", ip, "-j", "DROP")
		if err := checkCmd.Run(); err != nil {
			cmds = append(cmds, "iptables -A "+customChain+" -d "+ip+" -j DROP")
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalln("Failed to open a file containing ips: ", err)
	}

	// if ips are the same, exit
	if len(cmds) == 0 {
		log.Println("No new IPs to block.")
		return
	}

	// join commands together and run at once
	batch := strings.Join(cmds, " && ")
	command := exec.Command("bash", "-c", batch)
	if err := command.Run(); err != nil {
		log.Fatalln("Failed to run iptables commands: ", err)
	}

	log.Println("Blocked server ips from txt")

	if onDone != nil {
		onDone()
	}
}

func UnBlockIps(onDone func()) {
	chain := "OUTPUT"
	customChain := "CS2_BLOCKLIST"

	// untie and delete custom chain
	exec.Command("iptables", "-D", chain, "-j", customChain).Run()
	exec.Command("iptables", "-F", customChain).Run()
	exec.Command("iptables", "-X", customChain).Run()

	log.Println("Unblocked server ips from txt")

	if onDone != nil {
		onDone()
	}
}

func CustomChainExists() bool {
	cmd := exec.Command("iptables", "-L", "CS2_BLOCKLIST", "-n")
	return cmd.Run() == nil
}

func GetBlockedIps() map[string]bool {
	blocked := make(map[string]bool)

	cmd := exec.Command("iptables", "-S", "CS2_BLOCKLIST")
	output, err := cmd.Output()
	if err != nil {
		return blocked
	}

	lines := strings.SplitSeq(string(output), "\n")
	for line := range lines {
		if strings.Contains(line, "-d ") {
			fields := strings.Fields(line)
			for i, f := range fields {
				if f == "-d" && i+1 < len(fields) {
					ip := fields[i+1]
					cleanIP := strings.Split(ip, "/")[0]
					blocked[cleanIP] = true
				}
			}
		}
	}

	return blocked
}
