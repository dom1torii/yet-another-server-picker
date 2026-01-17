package config

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/spf13/pflag"

	"github.com/dom1torii/yet-another-server-picker/internal/fs"
)

type Config struct {
	Relays RelaysConfig `toml:"relays"`
	Ips    IpsConfig    `toml:"ips"`
	Log    LogConfig    `toml:"logging"`

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

type RelaysConfig struct {
	ShowPW bool `toml:"show_perfectworld"`
}

type IpsConfig struct {
	Path string `toml:"path"`
}

type LogConfig struct {
	Enabled bool   `toml:"enabled"`
	Path    string `toml:"path"`
}

func Init() *Config {
	cfg := &Config{}

	homeDir := fs.GetHomeDir()

	configDir := filepath.Join(homeDir, ".config", "yasp")
	configFile := filepath.Join(configDir, "config.toml")

	defaultIpsPath := filepath.Join(homeDir, "yasp_ips.txt")
	defaultLogPath := filepath.Join(homeDir, "yasp.log")

	fs.EnsureDirectory(configFile)

	info, err := os.Stat(configFile)
	if err == nil && info.Size() == 0 {
		defaultConfig(configFile, defaultIpsPath, defaultLogPath)
	}

	if _, err := toml.DecodeFile(configFile, cfg); err != nil {
		log.Fatalln("Failed to decode config: ", err)
	}

	onlyGlobal := pflag.BoolP("onlyglobal", "g", false, "Only show servers from global version of the game")
	ipsPath := pflag.StringP("ipspath", "i", getFlag(cfg.Ips.Path, defaultIpsPath), "Specify custom ips path path. Default path: {homedir}/yasp_ips.txt")
	logFlag := pflag.BoolP("log", "l", cfg.Log.Enabled, "Enable logging. Default path: {homedir}/yasp.log")
	logPath := pflag.String("logpath", getFlag(cfg.Log.Path, defaultLogPath), "Specify custom log file path.")

	listRelays := pflag.Bool("listrelays", false, "List available relays")
	selectRelays := pflag.String("selectrelays", "", "Select relays from the list (separated with comma)")
	listPresets := pflag.Bool("listpresets", false, "List available presets")
	selectPreset := pflag.String("selectpreset", "", "Select a preset from the list")
	blockRelays := pflag.Bool("blockrelays", false, "Block selected relays")
	unBlockRelays := pflag.Bool("unblockrelays", false, "Unblock selected relays")
	toBlockCount := pflag.Bool("toblockcount", false, "Prints amount of relays in your ips file")
	blockedCount := pflag.Bool("blockedcount", false, "Prints amount of blocked relays in your firewall")

	pflag.Parse()

	cfg.Ips.Path = *ipsPath
	cfg.Log.Path = *logPath

	isGlobalFlagSet := pflag.Lookup("onlyglobal").Changed
	isLogFlagSet := pflag.Lookup("log").Changed
	isPathFlagSet := pflag.Lookup("logpath").Changed

	if isGlobalFlagSet {
		cfg.Relays.ShowPW = !*onlyGlobal
	}

	fs.EnsureDirectory(cfg.Ips.Path)

	// we can just use --logpath instead of using both -l and --logpath to enable logging + change path
	if isLogFlagSet || isPathFlagSet || cfg.Log.Enabled {
		if isLogFlagSet && !*logFlag {
			cfg.Log.Enabled = false
		} else {
			cfg.Log.Enabled = true
		}
	}

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

	if cfg.Log.Enabled {
		fs.EnsureDirectory(cfg.Log.Path)
		initLogger(cfg.Log.Path)
		log.Println("Started logging at: ", cfg.Log.Path)
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

func defaultConfig(path, defaultIps, defaultLog string) {
	// fix windows paths
	defaultIps = strings.ReplaceAll(defaultIps, "\\", "\\\\")
	defaultLog = strings.ReplaceAll(defaultLog, "\\", "\\\\")

	content := []byte(strings.Join([]string{
		"[relays]",
		"show_perfectworld = true",
		"",
		"[ips]",
		"path = \"" + defaultIps + "\"",
		"",
		"[logging]",
		"enabled = false",
		"path = \"" + defaultLog + "\"",
	}, "\n"))

	if err := os.WriteFile(path, content, 0644); err != nil {
		log.Fatalln("Failed to write default config: ", err)
	}
}

func getFlag(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
