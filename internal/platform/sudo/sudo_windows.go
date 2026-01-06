//go:build windows

package sudo

import (
	"bufio"
	"fmt"
	"os"
)

func CheckIfSudo() {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		fmt.Println("Please run as administrator")
		fmt.Println("Press Enter to exit...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		os.Exit(1)
	}
}
