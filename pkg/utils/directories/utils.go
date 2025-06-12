package directories

import (
	"fmt"
	"github.com/pterm/pterm"
	"io/fs"
	"os"
	"os/exec"
)

type FileMode = fs.FileMode

// Function to check if a directory exists
func DirExists(dirPath string) bool {
	info, err := os.Stat(dirPath)
	return !os.IsNotExist(err) && info.IsDir()
}

// Function to remove directory and its content
func RemoveDirectory(path string) {
	err := os.RemoveAll(path)
	if err != nil && !os.IsNotExist(err) {
		pterm.Println()
		pterm.Error.Printfln("Failed to remove %s directory: %v", path, err)
		os.Exit(1)
	}
}

// Function to ensure directory and path to directory exists and sets permissions if created
func CreateDirectory(path string, permissions FileMode) {
	err := os.MkdirAll(path, permissions)
	if err != nil {
		pterm.Println()
		pterm.Error.Printfln("Failed to create %s directory: %v", path, err)
		os.Exit(1)
	}
}

// Function to copy a directory and all of its content
func CopyDirectory(dirToCopyPath, destDirPath string) {
	err := exec.Command("cp", "-R", dirToCopyPath, destDirPath).Run()
	if err != nil {
		pterm.Println()
		pterm.Error.Printfln("Failed to copy the %s directory to the %s directory: %v", dirToCopyPath, destDirPath, err)
		os.Exit(1)
	}
}

// Function to set permissions of a directory
func SetPermissions(path string, mode FileMode) {
	err := os.Chmod(path, mode)
	if err != nil {
		pterm.Println()
		pterm.Error.Printfln("Failed to set %s directory permissions: %v", path, err)
		os.Exit(1)
	}
}

// Function to set owner and group of a directory
func SetOwnerAndGroup(owner, group, dir string) {
	err := exec.Command("chown", "-R", fmt.Sprintf("%s:%s", owner, group), dir).Run()
	if err != nil {
		pterm.Println()
		pterm.Error.Printfln("Failed to set ownership of the %s directory: %v", dir, err)
		os.Exit(1)
	}
}
