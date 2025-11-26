package directories

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"

	"github.com/nodetec/rwz/pkg/relays"
	"github.com/pterm/pterm"
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

// Function to remove directory and its content using a linux command
func RemoveDirectoryUsingLinux(currentUsername, path string) {
	if currentUsername == relays.RootUser {
		err := exec.Command("rm", "-rf", path).Run()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to remove %s directory: %v", path, err)
			os.Exit(1)
		}
	} else {
		err := exec.Command("sudo", "rm", "-rf", path).Run()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to remove %s directory: %v", path, err)
			os.Exit(1)
		}
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

// Function to ensure directory and path to directory exists using a linux command
func CreateDirectoryUsingLinux(currentUsername, path string) {
	if !DirExists(path) {
		if currentUsername == relays.RootUser {
			err := exec.Command("mkdir", "-p", path).Run()
			if err != nil {
				pterm.Println()
				pterm.Error.Printfln("Failed to create %s directory: %v", path, err)
				os.Exit(1)
			}
		} else {
			err := exec.Command("sudo", "mkdir", "-p", path).Run()
			if err != nil {
				pterm.Println()
				pterm.Error.Printfln("Failed to create %s directory: %v", path, err)
				os.Exit(1)
			}
		}
	}
}

// Function to copy a directory and all of its content
func CopyDirectory(currentUsername, dirToCopyPath, destDirPath string) {
	if currentUsername == relays.RootUser {
		err := exec.Command("cp", "-R", dirToCopyPath, destDirPath).Run()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to copy the %s directory to the %s directory: %v", dirToCopyPath, destDirPath, err)
			os.Exit(1)
		}
	} else {
		err := exec.Command("sudo", "cp", "-R", dirToCopyPath, destDirPath).Run()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to copy the %s directory to the %s directory: %v", dirToCopyPath, destDirPath, err)
			os.Exit(1)
		}
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

// Function to set permissions of a directory using a linux command
func SetPermissionsUsingLinux(currentUsername, path string, mode string) {
	if currentUsername == relays.RootUser {
		err := exec.Command("chmod", mode, path).Run()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to set %s directory permissions: %v", path, err)
			os.Exit(1)
		}
	} else {
		err := exec.Command("sudo", "chmod", mode, path).Run()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to set %s directory permissions: %v", path, err)
			os.Exit(1)
		}
	}
}

// Function to check if a directory exists and set permissions of the directory using a linux command
func CheckIfDirectoryExistsAndSetPermissionsUsingLinux(currentUsername, path string, mode string) bool {
	if currentUsername == relays.RootUser {
		err := exec.Command("chmod", mode, path).Run()
		if err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				errorCode := exitError.ExitCode()
				// File not found
				if errorCode == 1 {
					return false
				} else {
					pterm.Println()
					pterm.Error.Printfln("Failed to set %s directory permissions: %v", path, err)
					os.Exit(1)
				}
			}
		}
	} else {
		err := exec.Command("sudo", "chmod", mode, path).Run()
		if exitError, ok := err.(*exec.ExitError); ok {
			errorCode := exitError.ExitCode()
			// File not found
			if errorCode == 1 {
				return false
			} else {
				pterm.Println()
				pterm.Error.Printfln("Failed to set %s directory permissions: %v", path, err)
				os.Exit(1)
			}
		}
	}
	return true
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

// Function to set owner and group of a directory using a linux command
func SetOwnerAndGroupUsingLinux(currentUsername, owner, group, dir string) {
	if currentUsername == relays.RootUser {
		err := exec.Command("chown", "-R", fmt.Sprintf("%s:%s", owner, group), dir).Run()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to set ownership of the %s directory: %v", dir, err)
			os.Exit(1)
		}
	} else {
		err := exec.Command("sudo", "chown", "-R", fmt.Sprintf("%s:%s", owner, group), dir).Run()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to set ownership of the %s directory: %v", dir, err)
			os.Exit(1)
		}
	}
}
