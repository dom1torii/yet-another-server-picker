package fs

import (
	"log"
	"os"
	"strconv"
)

func fixPermissions(path string) {
	sudoUid := os.Getenv("SUDO_UID")
	sudoGid := os.Getenv("SUDO_GID")

	if sudoUid != "" && sudoGid != "" {
		uid, _ := strconv.Atoi(sudoUid)
		gid, _ := strconv.Atoi(sudoGid)

		if err := os.Chown(path, uid, gid); err != nil {
			log.Printf("Failed to change ownership of %s: %v", path, err)
		}
	}
}
