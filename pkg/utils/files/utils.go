package files

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/nodetec/rwz/pkg/relays"
	"github.com/pterm/pterm"
)

// TODO
// Use Go os library instead of Linux commands

type FileMode = fs.FileMode

// Function to check if a file exists
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	return !os.IsNotExist(err) && !info.IsDir()
}

// Function to check if a file exists and removes it if it does
func RemoveFile(path string) {
	if _, err := os.Stat(path); err == nil {
		err = os.Remove(path)
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to remove %s file: %v", path, err)
			os.Exit(1)
		}
	}
}

// Function to remove a file using a linux command
func RemoveFileUsingLinux(currentUsername, path string) {
	if currentUsername == relays.RootUser {
		err := exec.Command("rm", "-f", path).Run()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to remove %s file: %v", path, err)
			os.Exit(1)
		}
	} else {
		err := exec.Command("sudo", "rm", "-f", path).Run()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to remove %s file: %v", path, err)
			os.Exit(1)
		}
	}
}

// Function to copy a file to a directory
func CopyFile(currentUsername, fileToCopy, destDir string) {
	if currentUsername == relays.RootUser {
		err := exec.Command("cp", fileToCopy, destDir).Run()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to copy %s file: %v", fileToCopy, err)
			os.Exit(1)
		}
	} else {
		err := exec.Command("sudo", "cp", fileToCopy, destDir).Run()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to copy %s file: %v", fileToCopy, err)
			os.Exit(1)
		}
	}
}

// Function to move a file to a new location
func MoveFile(pathToFileBeingMoved, destFilePath string) {
	err := exec.Command("mv", pathToFileBeingMoved, destFilePath).Run()
	if err != nil {
		pterm.Println()
		pterm.Error.Printfln("Failed to move %s file: %v", pathToFileBeingMoved, err)
		os.Exit(1)
	}
}

// Function to move a file to a new location using a linux command
func MoveFileUsingLinux(currentUsername, pathToFileBeingMoved, destFilePath string) {
	if currentUsername == relays.RootUser {
		err := exec.Command("mv", pathToFileBeingMoved, destFilePath).Run()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to move %s file: %v", pathToFileBeingMoved, err)
			os.Exit(1)
		}
	} else {
		err := exec.Command("sudo", "mv", pathToFileBeingMoved, destFilePath).Run()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to move %s file: %v", pathToFileBeingMoved, err)
			os.Exit(1)
		}
	}
}

// Function to set permissions of a file
func SetPermissions(path string, mode FileMode) {
	err := os.Chmod(path, mode)
	if err != nil {
		pterm.Println()
		pterm.Error.Printfln("Failed to set %s file permissions: %v", path, err)
		os.Exit(1)
	}
}

// Function to set permissions of a file using a linux command
func SetPermissionsUsingLinux(currentUsername, path string, mode string) {
	if currentUsername == relays.RootUser {
		err := exec.Command("chmod", mode, path).Run()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to set %s file permissions: %v", path, err)
			os.Exit(1)
		}
	} else {
		err := exec.Command("sudo", "chmod", mode, path).Run()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to set %s file permissions: %v", path, err)
			os.Exit(1)
		}
	}
}

// Function to check if a file exists and set permissions of the file using a linux command
func CheckIfFileExistsAndSetPermissionsUsingLinux(currentUsername, path string, mode string) bool {
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
					pterm.Error.Printfln("Failed to set %s file permissions: %v", path, err)
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
				pterm.Error.Printfln("Failed to set %s file permissions: %v", path, err)
				os.Exit(1)
			}
		}
	}
	return true
}

// Function to set owner and group of a file
func SetOwnerAndGroup(owner, group, file string) {
	err := exec.Command("chown", fmt.Sprintf("%s:%s", owner, group), file).Run()
	if err != nil {
		pterm.Println()
		pterm.Error.Printfln("Failed to set ownership of %s file: %v", file, err)
		os.Exit(1)
	}
}

// Function to set owner and group of a file using a linux command
func SetOwnerAndGroupUsingLinux(currentUsername, owner, group, file string) {
	if currentUsername == relays.RootUser {
		err := exec.Command("chown", fmt.Sprintf("%s:%s", owner, group), file).Run()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to set ownership of %s file: %v", file, err)
			os.Exit(1)
		}
	} else {
		err := exec.Command("sudo", "chown", fmt.Sprintf("%s:%s", owner, group), file).Run()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to set ownership of %s file: %v", file, err)
			os.Exit(1)
		}
	}
}

// Function to write content to a file
func WriteFile(currentUsername, path, content string, permissions FileMode) {
	if currentUsername == relays.RootUser {
		err := os.WriteFile(path, []byte(content), permissions)
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to write content to %s file: %v", path, err)
			os.Exit(1)
		}
	} else {
		_, err := exec.Command("sudo", "touch", path).CombinedOutput()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to create %s file: %v", path, err)
			os.Exit(1)
		}

		_, err = exec.Command("sudo", "chmod", "0666", path).CombinedOutput()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to set permissions for %s file: %v", path, err)
			os.Exit(1)
		}

		err = os.WriteFile(path, []byte(content), permissions)
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to write content to %s file: %v", path, err)
			os.Exit(1)
		}

		_, err = exec.Command("sudo", "chmod", "0644", path).CombinedOutput()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to set permissions for %s file: %v", path, err)
			os.Exit(1)
		}
	}
}

// Function to create a file if it doesn't exist, open the file in write only mode, and append content to the end of the file
func AppendContentToFile(path, content string, permissions FileMode) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, permissions)
	if err != nil {
		pterm.Println()
		pterm.Error.Printfln("Failed to open %s file: %v", path, err)
		os.Exit(1)
	}

	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf("%s\n", content))
	if err != nil {
		pterm.Println()
		pterm.Error.Printfln("Failed to write to %s file: %v", path, err)
		os.Exit(1)
	}
}

// Function to perform in-place editing
func InPlaceEdit(command, path string) {
	cmd := exec.Command("sed", "-i", command, path)

	// Execute the command
	if err := cmd.Run(); err != nil {
		pterm.Println()
		pterm.Error.Printfln("Failed to edit %s file in-place: %v", path, err)
		os.Exit(1)
	}
}

// Function to check if a line exists
func LineExists(pattern, path string) bool {
	cmd := exec.Command("grep", "-q", pattern, path)
	err := cmd.Run()

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode := exitError.ExitCode()
			if exitCode == 1 {
				return false
			} else {
				pterm.Println()
				pterm.Error.Printfln("Failed to search %s for %s: %v", path, pattern, err)
				os.Exit(1)
			}
		}
	}
	return true
}

// Function to download and copy a file
func DownloadAndCopyFile(currentUsername, filePath, downloadURL string, permissions FileMode) {
	// Create the file from the provided file path
	out, err := os.Create(filePath)
	if err != nil {
		pterm.Println()
		pterm.Error.Printfln("Failed to create %s file: %v", filePath, err)
		os.Exit(1)
	}
	defer out.Close()

	// Set the permissions for the created file
	SetPermissions(filePath, permissions)

	// Set owner and group for the created file to be the current user
	// if the current user isn't the root user
	// to ensure the current user has write access to the file
	if currentUsername != relays.RootUser {
		SetOwnerAndGroup(currentUsername, currentUsername, filePath)
	}

	// Download the file to copy
	resp, err := http.Get(downloadURL)
	if err != nil {
		pterm.Println()
		pterm.Error.Printfln("Failed to download file: %v", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		pterm.Println()
		pterm.Error.Printfln("Bad repsonse status code: %s", resp.Status)
		os.Exit(1)
	}

	// Write the response body to the created file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		pterm.Println()
		pterm.Error.Printfln("Failed to write to temporary file: %v", err)
		os.Exit(1)
	}
}

// Function to extract a file
func ExtractFile(currentUsername, tmpFilePath, destDir string) {
	if currentUsername == relays.RootUser {
		err := exec.Command("tar", "-xf", tmpFilePath, "-C", destDir).Run()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to extract binary to %s: %v", destDir, err)
			os.Exit(1)
		}
	} else {
		err := exec.Command("sudo", "tar", "-xf", tmpFilePath, "-C", destDir).Run()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to extract binary to %s: %v", destDir, err)
			os.Exit(1)
		}
	}
}

// Determine a file path from the base of another file path
func FilePathFromFilePathBase(filePath, newFileDirPath string) string {
	newFileName := filepath.Base(filePath)

	newFilePath := filepath.Join(newFileDirPath, newFileName)

	return newFilePath
}

// Install a compressed binary
func InstallCompressedBinary(currentUsername, compressedBinaryFilePath, binaryDestDir, binaryName string, permissions FileMode) {
	// Extract binary
	ExtractFile(currentUsername, compressedBinaryFilePath, binaryDestDir)

	// TODO
	// Currently, the downloaded binary is expected to have a name that matches the binaryName variable
	// Ideally, the extracted binary file should be renamed to match the binaryName variable

	// Define the final destination path
	destPath := filepath.Join(binaryDestDir, binaryName)

	// Make the file executable
	if currentUsername == relays.RootUser {
		SetPermissions(destPath, permissions)
	} else {
		SetPermissionsUsingLinux(currentUsername, destPath, "0755")
	}
}

// TODO
// Improve backup process by creating a unique and descriptive backup file name
// E.g., <database-file-name>-<pubkey-of-relay-runner>-<utc-timestamp-of-backup>-<unique-identifier>-bak.<database-file-extension>
// E.g., <users-file-name>-<pubkey-of-main-user>-<utc-timestamp-of-backup>-<unique-identifier>-bak.<users-file-extension>
// Then check if the file exists and create the file if it doesn't or try to create a new unique file name if it already exists
// Function to create a unique backup file name
func CreateUniqueBackupFileName(backupsDirPath, backupFileNameBase string) string {
	backupFileNumber := 0
	uniqueBackupFileName := fmt.Sprintf("%s-%s", backupFileNameBase, strconv.Itoa((backupFileNumber)))

	for FileExists(fmt.Sprintf("%s/%s", backupsDirPath, uniqueBackupFileName)) {
		backupFileNumber++
		uniqueBackupFileName = fmt.Sprintf("%s-%s", backupFileNameBase, strconv.Itoa(backupFileNumber))
	}

	return uniqueBackupFileName
}
