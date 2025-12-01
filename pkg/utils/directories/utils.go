package directories

import (
	"errors"
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
	dirInfo, err := os.Stat(dirPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false
		}
		pterm.Error.Printfln("Failed to check if %s directory exists: %v", dirPath, err)
		os.Exit(1)
	}

	return dirInfo.IsDir()
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

// Function to ensure directory and path to directory exists and sets permissions for all created directories, if directories aren't created then the permissions aren't set
// No data is overwritten and no error is thrown if any of the directories already exist
func CreateAllDirectories(path string, permissions FileMode) {
	err := os.MkdirAll(path, permissions)
	if err != nil {
		pterm.Println()
		pterm.Error.Printfln("Failed to create %s directory: %v", path, err)
		os.Exit(1)
	}
}

// Function to ensure directory and path to directory exists using a linux command
// No data is overwritten and no error is thrown if any of the directories already exist
func CreateAllDirectoriesUsingLinux(currentUsername, path string) {
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

// Function to copy a directory and all of its content using a linux command
func CopyDirectoryUsingLinux(currentUsername, dirToCopyPath, destDirPath string) {
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
func SetPermissionsUsingLinux(currentUsername, path, mode string) {
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
func CheckIfDirectoryExistsAndSetPermissionsUsingLinux(currentUsername, path, mode string) bool {
	if currentUsername == relays.RootUser {
		err := exec.Command("chmod", mode, path).Run()
		if err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				errorCode := exitError.ExitCode()
				// Directory not found
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
			// Directory not found
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

// Function to set owner and group of a directory and all of the directories and files within the specified directory using a linux command
func SetOwnerAndGroupForAllContentUsingLinux(currentUsername, owner, group, dirPath string) {
	ownerGroupArgument := fmt.Sprintf("%s:%s", owner, group)

	if currentUsername == relays.RootUser {
		err := exec.Command("chown", "-R", ownerGroupArgument, dirPath).Run()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to set ownership of the %s directory: %v", dirPath, err)
			os.Exit(1)
		}
	} else {
		err := exec.Command("sudo", "chown", "-R", ownerGroupArgument, dirPath).Run()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to set ownership of the %s directory: %v", dirPath, err)
			os.Exit(1)
		}
	}
}
