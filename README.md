<img width="1500" height="844" alt="aboba123" src="https://github.com/user-attachments/assets/b99ae40c-49c9-4eda-9847-299f38ef8302" />


# CS2 Server Manager
CS2 Server Manager is a TUI tool that allows you to choose CS2 servers you want to play on with ease.

## How it works

It takes relays from https://api.steampowered.com/ISteamApps/GetSDRConfig/v1?appid=730 and allows you to select relays you want. After you made your choice and chose "Block relays you don't want", it uses firewall rules (using iptables on linux and netsh on windows) to block servers you didn't choose.

## Can i get banned

No, because all it does is prevent your PC from connecting to certain IPs. It doesn't interact with the game at all and isn't a cheat.

## Installation

Releases are available for Windows and Linux on the [releases page](https://github.com/dom1torii/cs2-server-manager/releases/).

### Linux 

#### AUR (Arch Linux) 

```bash
Soon
```

### Windows

#### [Scoop](https://scoop.sh/) (recommended)

```
scoop bucket add cs2 https://github.com/dom1torii/cs2
scoop install cs2/cs2-server-manager
```

### Building from source

1. Install GoLang -> https://go.dev/doc/install
2. Clone the repo `git clone https://github.com/dom1torii/cs2-server-manager.git`
3. `cd` into the folder
4. Run `go build ./cmd/cs2-server-manager/`

## Planned features

- Support linux firewall utils other than iptables
- Toggle between blocking/allowing servers instead of just allowing
- Global/China version switch
- Settings and config file

## Notes

On linux, only `iptables` is currently supported. Working on making `ufw` and `nftables` work too.

The tool is not fully accurate and sometimes will connect you to server that are **routed** through the server you chose. 
Its also possible that you wont find the server you selected. 
If you have any ideas how to fix that, pull requests are welcome :)

## Libraries used

https://github.com/charmbracelet/bubbletea and others from [Charm](https://charm.land/) - TUI. 
https://github.com/spf13/pflag - CLI flags. 
https://github.com/muesli/reflow - Small library for text wrapping. 
https://github.com/prometheus-community/pro-bing - For pinging IP addresses. 
