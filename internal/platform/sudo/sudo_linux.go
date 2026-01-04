//go:build linux

package sudo

import (
	"os"
	"fmt"
)

func CheckIfSudo() {
	if os.Geteuid() == 0 {
		return
	} else {
		fmt.Println("Please run with sudo")
		os.Exit(1)
	}
}
