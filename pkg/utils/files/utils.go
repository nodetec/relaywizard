package files

import (
	"fmt"
	"io"
	"io/fs"
	"log"
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
			log.Fatalf("Error removing %s file: %v", path, err)
		}
	}
}

// Function to copy a file to a directory
func CopyFile(fileToCopy, destDir string) {
	err := exec.Command("cp", fileToCopy, destDir).Run()
	if err != nil {
		log.Fatalf("Error copying %s file: %v", fileToCopy, err)
	}
}

// Function to set owner and group of a file
func SetOwnerAndGroup(owner, group, file string) {
	err := exec.Command("chown", fmt.Sprintf("%s:%s", owner, group), file).Run()
	if err != nil {
		log.Fatalf("Error setting ownership of %s file: %v", file, err)
	}
}

// Function to set permissions of a file
func SetPermissions(path string, mode FileMode) {
	err := os.Chmod(path, mode)
	if err != nil {
		log.Fatalf("Error setting %s file permissions: %v", path, err)
	}
}

// Function to write content to a file
func WriteFile(path, content string, permissions FileMode) {
	err := os.WriteFile(path, []byte(content), permissions)
	if err != nil {
		log.Fatalf("Error writing content to %s file: %v", path, err)
	}
}

// Function to perform in-place editing
func InPlaceEdit(command, path string) {
	cmd := exec.Command("sed", "-i", command, path)

	// Execute the command
	if err := cmd.Run(); err != nil {
		log.Fatalf("Error editing %s in-place: %v", path, err)
	}
}

// Function to download and copy a file
func DownloadAndCopyFile(tmpFilePath, downloadURL string) {
	// Create a temporary file
	out, err := os.Create(tmpFilePath)
	if err != nil {
		log.Fatalf("Error creating %s file: %v", tmpFilePath, err)
	}
	defer out.Close()

	// Download the file
	resp, err := http.Get(downloadURL)
	if err != nil {
		log.Fatalf("Error downloading file: %v", err)
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Bad status: %s", resp.Status)
	}

	// Write the body to the temporary file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatalf("Error writing to temporary file: %v", err)
	}
}

// Function to extract a file
func ExtractFile(tmpFilePath, destDir string) {
	err := exec.Command("tar", "-xf", tmpFilePath, "-C", destDir).Run()
	if err != nil {
		log.Fatalf("Error extracting binary to %s: %v", destDir, err)
	}
}
