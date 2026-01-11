//go:build windows

package firewall

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/dom1torii/cs2-server-manager/internal/config"
)

func BlockIps(cfg *config.Config, onDone func()) {
	UnBlockIps(onDone)
	ipsFile := cfg.IpsPath

	ruleName := "CS2_BLOCKLIST"

	file, err := os.Open(ipsFile)
	if err != nil {
		log.Fatalln("Failed to open a file containing ips: ", err)
	}
	defer file.Close()

	// read ips from a file and then join them with a comma
	var ips []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ip := strings.TrimSpace(scanner.Text())
		if ip == "" {
			continue
		}
		ips = append(ips, ip)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln("Failed to read ips file: ", err)
	}

	if len(ips) == 0 {
		log.Println("No IPs found to block.")
		return
	}

	remoteIps := strings.Join(ips, ",")

	// execute a command to block all ips at once
	cmd := exec.Command("netsh", "advfirewall", "firewall", "add", "rule",
		"name="+ruleName,
		"dir=out",
		"action=block",
		"remoteip="+remoteIps,
	)

	if err := cmd.Run(); err != nil {
		log.Fatalln("Failed to create Windows Firewall rule: ", err)
	}

	log.Println("Blocked server ips in Windows Firewall")

	if onDone != nil {
		onDone()
	}
}

func UnBlockIps(onDone func()) {
	ruleName := "CS2_BLOCKLIST"

	// just delete the rule
	exec.Command("netsh", "advfirewall", "firewall", "delete", "rule", "name="+ruleName).Run()

	log.Println("Unblocked server ips from txt")

	if onDone != nil {
		onDone()
	}
}

func CustomChainExists() bool {

	cmd := exec.Command("netsh", "advfirewall", "firewall", "show", "rule", "name=CS2_BLOCKLIST")

	return cmd.Run() == nil
}

func GetBlockedIpCount() (int, error) {
	cmd := exec.Command("netsh", "advfirewall", "firewall", "show", "rule", "name=CS2_BLOCKLIST")
	output, err := cmd.Output()
	if err != nil {
		return 0, nil
	}

	// find ips and add them together
	count := 0
	for line := range strings.SplitSeq(string(output), "\n") {
		if strings.Contains(line, "RemoteIP:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) < 2 {
				continue
			}

			ipString := strings.TrimSpace(parts[1])
			if ipString == "" || ipString == "Any" {
				return 0, nil
			}

			for ip := range strings.SplitSeq(ipString, ",") {
				if strings.TrimSpace(ip) != "" {
					count++
				}
			}
			return count, nil
		}
	}

	return 0, nil
}

func IsIpBlocked(ip string) bool {
	cmd := exec.Command("netsh", "advfirewall", "firewall", "show", "rule", "name=CS2_BLOCKLIST")
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	// find the ip and return true if its blocked and false if its not
	for line := range strings.SplitSeq(string(output), "\n") {
		if strings.Contains(line, "RemoteIP:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) < 2 {
				return false
			}

			ipList := strings.TrimSpace(parts[1])
			for entry := range strings.SplitSeq(ipList, ",") {
				cleanEntry := strings.TrimSpace(entry)
				if strings.Contains(cleanEntry, "/") {
					cleanEntry = strings.Split(cleanEntry, "/")[0]
				}

				if cleanEntry == ip {
					return true
				}
			}
		}
	}

	return false
}
