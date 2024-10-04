package directories

import (
	"fmt"
	"github.com/pterm/pterm"
	"io/fs"
	"os"
	"os/exec"
)

type FileMode = fs.FileMode

// Function to remove directory
func RemoveDirectory(path string) {
	err := os.RemoveAll(path)
	if err != nil && !os.IsNotExist(err) {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Failed to remove %s directory: %v", path, err))
		os.Exit(1)
	}
}

// Function to ensure directory and path to directory exists and sets permissions
func CreateDirectory(path string, permissions FileMode) {
	err := os.MkdirAll(path, permissions)
	if err != nil {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Failed to create %s directory: %v", path, err))
		os.Exit(1)
	}
}

// Function to set owner and group of a directory
func SetOwnerAndGroup(owner, group, dir string) {
	err := exec.Command("chown", "-R", fmt.Sprintf("%s:%s", owner, group), dir).Run()
	if err != nil {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Failed to set ownership of the %s directory: %v", dir, err))
		os.Exit(1)
	}
}
