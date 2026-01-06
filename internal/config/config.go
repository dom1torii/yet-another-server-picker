package config

import (
	"io"
	"log"
	"os"

	"github.com/spf13/pflag"

	"github.com/dom1torii/cs2-server-manager/internal/fs"
)

type Config struct {
	IpsPath  string
	LogsPath string
	Logging  bool
}

func Init() *Config {
	cfg := &Config{}

	defaultIpsPath := fs.GetHomeDir() + "/cs2sp_ips.txt"
	defaultLogPath := fs.GetHomeDir() + "/cs2sp.log"

	ipsPath := pflag.StringP("ipspath", "i", defaultIpsPath, "Specify custom ips path path. Default path: {homedir}/cs2sp_ips.txt")
	logFlag := pflag.BoolP("log", "l", false, "Enable logging. Default path: {homedir}/cs2sp.log")
	logPath := pflag.String("logpath", "", "Specify custom log file path.")

	pflag.Parse()

	fs.EnsureDirectory(*ipsPath)
	cfg.IpsPath = *ipsPath

	if *logFlag {
		path := *logPath
		if path == "" {
			path = defaultLogPath
		}

		fs.EnsureDirectory(path)
		cfg.LogsPath = path
		initLogger(path)
		log.Println("Started logging.")
	} else {
		log.SetOutput(io.Discard)
	}
	return cfg
}

func initLogger(loc string) {
	f, err := os.OpenFile(loc, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln("Failed to open log file:", err)
	}
	log.SetOutput(f)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
