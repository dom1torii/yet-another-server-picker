package config

import (
	"io"
	"log"
	"os"
	"strings"

	"github.com/spf13/pflag"

	"github.com/dom1torii/yet-another-server-picker/internal/fs"
)

type Config struct {
	IpsPath  string
	LogsPath string
	Logging  bool

	// cli mode
	ListRelays    bool
	SelectRelays  []string
	ListPresets   bool
	SelectPreset  string
	BlockRelays   bool
	UnBlockRelays bool
	ToBlockCount  bool
	BlockedCount  bool
}

func Init() *Config {
	cfg := &Config{}

	defaultIpsPath := fs.GetHomeDir() + "/yasp_ips.txt"
	defaultLogPath := fs.GetHomeDir() + "/yasp.log"

	ipsPath := pflag.StringP("ipspath", "i", defaultIpsPath, "Specify custom ips path path. Default path: {homedir}/yasp_ips.txt")
	logFlag := pflag.BoolP("log", "l", false, "Enable logging. Default path: {homedir}/yasp.log")
	logPath := pflag.String("logpath", "", "Specify custom log file path.")

	listRelays := pflag.Bool("listrelays", false, "List available relays")
	selectRelays := pflag.String("selectrelays", "", "Select relays from the list (separated with comma)")
	listPresets := pflag.Bool("listpresets", false, "List available presets")
	selectPreset := pflag.String("selectpreset", "", "Select a preset from the list")
	blockRelays := pflag.Bool("blockrelays", false, "Block selected relays")
	unBlockRelays := pflag.Bool("unblockrelays", false, "Unblock selected relays")
	toBlockCount := pflag.Bool("toblockcount", false, "Prints amount of relays in your ips file")
	blockedCount := pflag.Bool("blockedcount", false, "Prints amount of relays in your ips file")

	pflag.Parse()

	fs.EnsureDirectory(*ipsPath)
	cfg.IpsPath = *ipsPath

	cfg.ListRelays = *listRelays
	if *selectRelays != "" {
		cfg.SelectRelays = strings.Split(*selectRelays, ",")
	}
	cfg.BlockRelays = *blockRelays
	cfg.UnBlockRelays = *unBlockRelays
	cfg.ToBlockCount = *toBlockCount
	cfg.BlockedCount = *blockedCount
	cfg.ListPresets = *listPresets
	cfg.SelectPreset = *selectPreset

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
		log.Fatalln("Failed to open log file: ", err)
	}
	log.SetOutput(f)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
