<img width="1920" height="1080" alt="image" src="https://github.com/user-attachments/assets/b199a561-7703-4cd9-a244-6327f930c423" />

# CS2 Server Manager
CS2 Server Manager is a TUI tool that allows you to choose CS2 servers you want to play on with ease.

## How it works

It takes relays from https://api.steampowered.com/ISteamApps/GetSDRConfig/v1?appid=730 and allows you to select relays you want. After you chose relays, it uses firewall rules (using iptables on linux and netsh on windows) to block server you didn't choose.

## Can i get banned

No, because it just doesn't allow your PC to connect to some IPs. It doesn't interact with the game at all and isn't a cheat.

## Installation
### Releases

Later

### Build from source

1. Install GoLang -> https://go.dev/doc/install
2. Clone the repo `git clone https://github.com/dom1torii/cs2-server-manager.git`
3. `cd` into the folder
4. Run `go build ./cmd/cs2-server-manager/`

## Notes

The script is not fully accurate and sometimes will connect you to server that are **routed** through the server you chose. 
Its also possible that you wont find the server you selected. 
If you have any ideas how to fix that, pull requests are welcome :)

## Credits 

https://github.com/rivo/tview - Very cool TUI library.
