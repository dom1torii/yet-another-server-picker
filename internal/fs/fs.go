package fs

import (
	"bufio"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

func GetHomeDir() string {
	if sudoUser := os.Getenv("SUDO_USER"); sudoUser != "" {
		if usr, err := user.Lookup(sudoUser); err == nil {
			return usr.HomeDir
		}
	}

	if home := os.Getenv("HOME"); home != "" {
		return home
	}

	// should work for windows idk
	usr, err := user.Current()
	if err != nil {
		log.Fatalln("Failed to get user:", err)
	}
	return usr.HomeDir
}

func EnsureDirectory(path string) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		log.Fatalln("Failed to create directory: ", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		f, err := os.Create(path)
		if err != nil {
			log.Fatalln("Failed to create file: ", err)
		}
		f.Close()
	}
}

func IsFileEmpty(filePath string) bool {
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return true
	} else if err != nil {
		log.Fatalln("Failed to get file info:", err)
	}

	return fileInfo.Size() == 0
}

func GetFileLineCount(filePath string) int {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalln("Failed to open file:", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Fatalln("Failed to close file:", err)
		}
	}()

	scanner := bufio.NewScanner(file)

	lines := 0
	for scanner.Scan() {
		lines++
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln("Failed to read file:", err)
	}

	return lines
}
