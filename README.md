<img width="1262" height="710" alt="yasp_preview" src="https://github.com/user-attachments/assets/18402ad2-9bd9-4bd8-966a-c3f0ab0ce337" />



# <img src="assets/icon-256.png" width="25"> Yet Another Server Picker

YASP - Cross-platform TUI tool that allows you to select which CS2 servers you want to play on by blocking IPs of unwanted servers in your firewall.

## How it works

YASP fetches relays from https://api.steampowered.com/ISteamApps/GetSDRConfig/v1?appid=730 and allows you to select relays you want or don't want. After you made your choice and chose "Block servers you don't want", it creates firewall rules (using iptables on linux and netsh on windows) to block unwanted servers.

## Can i get banned

No, because all it does is prevent your PC from connecting to certain IPs. It doesn't interact with the game at all and isn't a cheat.

## Installation

Releases are available for Windows and Linux on the [releases page](https://github.com/dom1torii/yet-another-server-picker/releases/).

### Linux 

#### [AUR](https://aur.archlinux.org/packages/yasp) (Arch Linux) 

```bash
paru -S yasp
```

#### Build and install

1. Install [GoLang](https://go.dev/doc/install) and [Go-Task](https://taskfile.dev/docs/installation)  
2. Clone the repository: `git clone https://github.com/dom1torii/yet-another-server-picker.git`
3. `cd` into the folder
4. Install: `sudo go-task install`

### Windows

#### Winget
```
winget install yasp
```

#### [Scoop](https://scoop.sh/)

```
scoop bucket add cs2 https://github.com/dom1torii/cs2
scoop install cs2/yasp
```

### Building from source

1. Install GoLang -> https://go.dev/doc/install
2. Clone the repository: `git clone https://github.com/dom1torii/yet-another-server-picker.git`
4. `cd` into the folder
5. Build the binary: `go build ./cmd/yasp/`

<!--## Planned features

-->
## Configuration

Config file is located at `/home/username/.config/yasp/` on linux or `C:\Users\Username\.config\yasp\` on windows and is created by default when you first launch the executable.  
It allows you to easily access and change some settings you might want to change.

Default config would look something like this:
```toml
[relays]
show_perfectworld = true

[ips]
path = "/home/username/yasp_ips.txt"

[logging]
enabled = false
path = "/home/username/yasp.log"
```

Explaination:  
`show_perfectworld` - show/hide perfect world servers in the list.  
`ips:path` - path to the file where ips you selected will be stored.  
`logging:enabled` - enable logging (for debugging purposes).   
`logging:path` - path to the file where logs will be stored.  

## Limitations

- Selecting servers is not fully accurate and sometimes you may get routed through the server you selected instead of playing directly on it.
- Sometimes server you selected may not be found. (using high `mm_dedicated_search_maxping` is recommended)

If you have any ideas on how to fix or improve something, pull requests are always welcome :)

## Libraries used

https://github.com/charmbracelet/bubbletea and others from [Charm](https://charm.land/) - TUI.  
https://github.com/spf13/pflag - CLI flags.  
https://github.com/BurntSushi/toml - TOML parser.  
https://github.com/muesli/reflow - Small library for text wrapping.  
https://github.com/prometheus-community/pro-bing - For pinging IP addresses.

## Credits

Thank you [alekun](https://x.com/akvnxii) for the logo <3
