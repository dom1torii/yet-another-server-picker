//go:build linux

package sudo

import (
	"fmt"
	"os"
)

func CheckIfSudo() {
	if os.Geteuid() == 0 {
		return
	} else {
		fmt.Println("Please run with sudo")
		os.Exit(1)
	}
}
