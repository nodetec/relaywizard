package files

import (
	"fmt"
	"github.com/pterm/pterm"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
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

// Function to copy a file to a directory
func CopyFile(fileToCopy, destDir string) {
	err := exec.Command("cp", fileToCopy, destDir).Run()
	if err != nil {
		pterm.Println()
		pterm.Error.Printfln("Failed to copy %s file: %v", fileToCopy, err)
		os.Exit(1)
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

// Function to set permissions of a file
func SetPermissions(path string, mode FileMode) {
	err := os.Chmod(path, mode)
	if err != nil {
		pterm.Println()
		pterm.Error.Printfln("Failed to set %s file permissions: %v", path, err)
		os.Exit(1)
	}
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

// Function to write content to a file
func WriteFile(path, content string, permissions FileMode) {
	err := os.WriteFile(path, []byte(content), permissions)
	if err != nil {
		pterm.Println()
		pterm.Error.Printfln("Failed to write content to %s file: %v", path, err)
		os.Exit(1)
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
func DownloadAndCopyFile(filePath, downloadURL string, permissions FileMode) {
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
func ExtractFile(tmpFilePath, destDir string) {
	err := exec.Command("tar", "-xf", tmpFilePath, "-C", destDir).Run()
	if err != nil {
		pterm.Println()
		pterm.Error.Printfln("Failed to extract binary to %s: %v", destDir, err)
		os.Exit(1)
	}
}

// Determine a file path from the base of another file path
func FilePathFromFilePathBase(filePath, newFileDirPath string) string {
	newFileName := filepath.Base(filePath)

	newFilePath := filepath.Join(newFileDirPath, newFileName)

	return newFilePath
}

// Install a compressed binary
func InstallCompressedBinary(compressedBinaryFilePath, binaryDestDir, binaryName string, permissions FileMode) {
	// Extract binary
	ExtractFile(compressedBinaryFilePath, binaryDestDir)

	// TODO
	// Currently, the downloaded binary is expected to have a name that matches the binaryName variable
	// Ideally, the extracted binary file should be renamed to match the binaryName variable

	// Define the final destination path
	destPath := filepath.Join(binaryDestDir, binaryName)

	// Make the file executable
	SetPermissions(destPath, permissions)
}
