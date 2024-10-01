package directories

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
)

type FileMode = fs.FileMode

// Function to remove directory
func RemoveDirectory(path string) {
	err := os.RemoveAll(path)
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("Error removing %s directory: %v", path, err)
	}
}

// Function to ensure directory and path to directory exists and sets permissions
func CreateDirectory(path string, permissions FileMode) {
	err := os.MkdirAll(path, permissions)
	if err != nil {
		log.Fatalf("Error creating %s directory: %v", path, err)
	}
}

// Function to set owner and group of a directory
func SetOwnerAndGroup(owner, group, dir string) {
	err := exec.Command("chown", "-R", fmt.Sprintf("%s:%s", owner, group), dir).Run()
	if err != nil {
		log.Fatalf("Error setting ownership of the %s directory: %v", dir, err)
	}
}
