package tui

import (
	"log"
	"strconv"
	"sort"
	"strings"

	"github.com/rivo/tview"
	"github.com/gdamore/tcell/v2"

	"github.com/dom1torii/cs2-server-manager/internal/config"
	"github.com/dom1torii/cs2-server-manager/internal/api"
	"github.com/dom1torii/cs2-server-manager/internal/fs"
	"github.com/dom1torii/cs2-server-manager/internal/ips"
	"github.com/dom1torii/cs2-server-manager/internal/platform/firewall"
)

type StartListItemData struct {
  MainText  string
  Secondary func() string
  Shortcut  rune
  Action    func()
}

func InitStartPage(ui *UI, cfg *config.Config) tview.Primitive {

	// items for start menu list
	items := []StartListItemData{
	  {
	    MainText: "Select servers",
	    Secondary: nil,
			Shortcut: '1',
	    Action: func() { ui.SwitchPage("select") },
	  },
	  {
	    MainText: "Block servers you don't want",
	    Secondary: func() string {
	      ipsFile := cfg.IpsPath
	      if !fs.IsFileEmpty(ipsFile) {
	        return "[orange]" + strconv.Itoa(fs.GetFileLineCount(cfg.IpsPath)) + " IPs in " + ipsFile
	      }
	      return ipsFile + " is empty, nothing to block"
	    },
			Shortcut: '2',
	    Action: func() { firewall.BlockIps(cfg, func() { go ui.RefreshStartList() }) },
	  },
	  {
	    MainText: "Unblock all servers",
	    Secondary: func() string {
	      blockedCount, _ := firewall.GetBlockedIpCount()
				if blockedCount == 0 {
					return "[green]" + strconv.Itoa(blockedCount) + " IPs currently blocked"
				} else {
					return "[red]" + strconv.Itoa(blockedCount) + " IPs currently blocked"
				}
	    },
			Shortcut: '3',
	    Action: func() { firewall.UnBlockIps(func() { go ui.RefreshStartList() }) },
	  },
	  {
	    MainText: "Quit",
	    Secondary: nil,
			Shortcut: 'q',
	    Action: func() { ui.App.Stop() },
	  },
	}

	// create start menu list from these items
	list := tview.NewList()
	for _, item := range items {
	  secondary := ""
	  if item.Secondary != nil {
	    secondary = item.Secondary()
	  }

	  list.AddItem(item.MainText, secondary, item.Shortcut, item.Action)
	}

	// refresh secondary texts of main menu items
	refresh := func() {
	  for i, item := range items {
	    secondary := ""
	    if item.Secondary != nil {
	      secondary = item.Secondary()
	    }
	    list.SetItemText(i, item.MainText, secondary)
	  }
	}

	ui.RefreshStartList = func() {
	  ui.App.QueueUpdateDraw(refresh)
	}

	list.SetBackgroundColor(tcell.ColorDefault).SetBorder(true)

	flex := tview.NewFlex().
		AddItem(nil,0,1,false).
		AddItem(list,60,1,true).
		AddItem(nil,0,1,false)

	root := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil,0,1,false).
		AddItem(flex,10,1,true).
		AddItem(nil,0,1,false)

	return root
}

func InitSelectPage(ui *UI, cfg *config.Config) tview.Primitive {
	ipsFile := cfg.IpsPath

  response, err := api.FetchRelays()
  if err != nil {
    log.Fatalln("Failed to fetch relays:", err)
  }

  table := tview.NewTable().
    SetSelectable(true, true)

  type popItem struct {
    key     string
    desc    string
    checked bool
  }

  var items []popItem
  for key, pop := range response.Pops {
    items = append(items, popItem{
      key:  key,
      desc: pop.Desc,
    })
  }

  // sort servers alphabetically by whats inside ()
  sort.Slice(items, func(i, j int) bool {
    a := items[i].desc
    b := items[j].desc

    aStart := strings.LastIndex(a, "(")
    aEnd := strings.LastIndex(a, ")")
    bStart := strings.LastIndex(b, "(")
    bEnd := strings.LastIndex(b, ")")

    var aKey, bKey string
    if aStart != -1 && aEnd != -1 && aStart < aEnd {
      aKey = a[aStart+1 : aEnd]
    } else {
      aKey = a
    }
    if bStart != -1 && bEnd != -1 && bStart < bEnd {
      bKey = b[bStart+1 : bEnd]
    } else {
      bKey = b
    }

    return aKey < bKey
  })

  updateTable := func() {
    table.Clear()

    cols := 2
    rows := (len(items) + cols - 1) / cols

    for i, item := range items {
      row := i % rows
      col := (i / rows) * 2

      checkbox := "[ ]"
      if item.checked {
        checkbox = "[✓]"
      }

      table.SetCell(row, col,
        tview.NewTableCell(checkbox).
          SetSelectable(true))

      table.SetCell(row, col+1,
        tview.NewTableCell(item.desc).
          SetSelectable(true))
    }
  }

  updateTable()
  table.Select(0, 0)

  getUnselectedIps := func() []string {
    var ips []string
    for _, item := range items {
      if !item.checked {
        pop := response.Pops[item.key]
        for _, relay := range pop.Relays {
          ips = append(ips, relay.Ipv4)
        }
      }
    }
    return ips
  }

  // use arrow keys to navigate, space to select, enter to proceed and q/Q to go back
  table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
    row, col := table.GetSelection()
    cols := 2
    rows := (len(items) + cols - 1) / cols
    index := (col / 2) * rows + row

    switch {
    case event.Key() == tcell.KeyRune && event.Rune() == ' ':
      if index >= 0 && index < len(items) {
        items[index].checked = !items[index].checked
        updateTable()
        table.Select(row, col)
      }
      return nil

    case event.Key() == tcell.KeyEnter:
    	if !fs.IsFileEmpty(ipsFile) {
   			ui.ConfirmOverwrite(func() {
      		ips.WriteIpsToFile(getUnselectedIps(), cfg)
      	})
     	} else {
      	ips.WriteIpsToFile(getUnselectedIps(), cfg)
      	ui.SwitchPage("start")
      }

      return nil

    case event.Key() == tcell.KeyRight:
      if col%2 == 0 {
        table.Select(row, col+2)
      }
      return nil

    case event.Key() == tcell.KeyLeft:
      if col%2 == 0 && col >= 2 {
        table.Select(row, col-2)
      }
      return nil

    case event.Rune() == 'q':
   		ui.SwitchPage("start")
      return nil

    case event.Rune() == 'Q':
   		ui.SwitchPage("start")
      return nil
    }



    return event
  })

  table.SetBorder(true).
    SetTitle(" Select servers you want to play on (SPACE to select, ENTER to proceed) ")

  return table
}

func SetupPages(ui *UI, cfg *config.Config) {
  ui.Pages.AddAndSwitchToPage("start", InitStartPage(ui, cfg), true)
  ui.Pages.AddPage("select", InitSelectPage(ui, cfg), true, false)
}
