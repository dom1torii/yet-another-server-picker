package ips

import (
	"bufio"
	"log"
	"os"

	"github.com/dom1torii/cs2-server-manager/internal/config"
)

func WriteIpsToFile(ips []string, cfg *config.Config) {
	ipsFile := cfg.IpsPath
	file, err := os.Create(ipsFile)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	// 1 ip on every line
	writer := bufio.NewWriter(file)
	for _, ip := range ips {
		_, err := writer.WriteString(ip + "\n")
		if err != nil {
			log.Fatalln("Failed to write ips to a file:", err)
			return
		}
	}

	writer.Flush()
	log.Println("Wrote ips to " + ipsFile)
}
