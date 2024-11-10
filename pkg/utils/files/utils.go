package files

import (
	"fmt"
	"github.com/pterm/pterm"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
)

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
			pterm.Error.Println(fmt.Sprintf("Failed to remove %s file: %v", path, err))
			os.Exit(1)
		}
	}
}

// Function to copy a file to a directory
func CopyFile(fileToCopy, destDir string) {
	err := exec.Command("cp", fileToCopy, destDir).Run()
	if err != nil {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Failed to copy %s file: %v", fileToCopy, err))
		os.Exit(1)
	}
}

// Function to set permissions of a file
func SetPermissions(path string, mode FileMode) {
	err := os.Chmod(path, mode)
	if err != nil {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Failed to set %s file permissions: %v", path, err))
		os.Exit(1)
	}
}

// Function to set owner and group of a file
func SetOwnerAndGroup(owner, group, file string) {
	err := exec.Command("chown", fmt.Sprintf("%s:%s", owner, group), file).Run()
	if err != nil {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Failed to set ownership of %s file: %v", file, err))
		os.Exit(1)
	}
}

// Function to write content to a file
func WriteFile(path, content string, permissions FileMode) {
	err := os.WriteFile(path, []byte(content), permissions)
	if err != nil {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Failed to write content to %s file: %v", path, err))
		os.Exit(1)
	}
}

// Function to perform in-place editing
func InPlaceEdit(command, path string) {
	cmd := exec.Command("sed", "-i", command, path)

	// Execute the command
	if err := cmd.Run(); err != nil {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Failed to edit %s file in-place: %v", path, err))
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
				pterm.Error.Println(fmt.Sprintf("Failed to search %s for %s: %v", path, pattern, err))
				os.Exit(1)
			}
		}
	}
	return true
}

// Function to download and copy a file
func DownloadAndCopyFile(tmpFilePath, downloadURL string) {
	// Create a temporary file
	out, err := os.Create(tmpFilePath)
	if err != nil {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Failed to create %s file: %v", tmpFilePath, err))
		os.Exit(1)
	}
	defer out.Close()

	SetPermissions(tmpFilePath, 0644)

	// Download the file
	resp, err := http.Get(downloadURL)
	if err != nil {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Failed to download file: %v", err))
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Bad repsonse status code: %s", resp.Status))
		os.Exit(1)
	}

	// Write the body to the temporary file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Failed to write to temporary file: %v", err))
		os.Exit(1)
	}
}

// Function to extract a file
func ExtractFile(tmpFilePath, destDir string) {
	err := exec.Command("tar", "-xf", tmpFilePath, "-C", destDir).Run()
	if err != nil {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Failed to extract binary to %s: %v", destDir, err))
		os.Exit(1)
	}
}
