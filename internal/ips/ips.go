package ips

import (
	"bufio"
	"log"
	"os"
	"strings"
	"time"

	"github.com/dom1torii/yet-another-server-picker/internal/config"
	"github.com/prometheus-community/pro-bing"
)

func WriteIpsToFile(ips []string, cfg *config.Config) {
	ipsFile := cfg.IpsPath
	file, err := os.Create(ipsFile)
	if err != nil {
		log.Fatalln("Failed to create file: ", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Fatalln("Failed to close file: ", err)
		}
	}()

	// 1 ip on every line
	writer := bufio.NewWriter(file)
	for _, ip := range ips {
		_, err := writer.WriteString(ip + "\n")
		if err != nil {
			log.Fatalln("Failed to write ips to a file: ", err)
			return
		}
	}

	if err := writer.Flush(); err != nil {
		log.Fatalln("Failed to flush writer: ", err)
	}

	log.Println("Wrote ips to " + ipsFile)
}

func GetPing(ip string) time.Duration {
	pinger, err := probing.NewPinger(ip)
	if err != nil {
		log.Fatalln("Failed to create pinger: ", err)
	}

	pinger.SetPrivileged(true)

	pinger.Count = 3
	pinger.Timeout = time.Millisecond * 500

	err = pinger.Run()
	if err != nil {
		// don't do anything if ip is blocked with firewall
		if strings.Contains(err.Error(), "operation not permitted") {
			return -1
		}
		log.Fatalln("Failed to run pinger: ", err)
	}

	stats := pinger.Statistics()
	return stats.AvgRtt
}
