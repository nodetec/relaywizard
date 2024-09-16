package strfry

import (
	"fmt"
	"github.com/pterm/pterm"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

// Function to download and make the binary executable
func InstallRelayBinary() {
	// Temporary directory for git repository
	const tempDir = "/tmp/strfry"

	// URL of the binary to download
	const downloadURL = "https://github.com/nodetec/relays/releases/download/v0.1.0/strfry-0.9.7-x86_64-linux-gnu.tar.gz"

	// Name of the binary after downloading
	const binaryName = "nostr-relay-strfry"

	// Destination directory for the binary
	const destDir = "/usr/local/bin"

	spinner, _ := pterm.DefaultSpinner.Start("Installing strfry relay...")

	// Check for and remove existing git repository
	err := os.RemoveAll(fmt.Sprintf("%s", tempDir))
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("Error removing existing repository: %v", err)
	}

	// Download git repository
	err = exec.Command("git", "clone", "-b", "0.9.7", "https://github.com/hoytech/strfry.git", fmt.Sprintf("%s", tempDir)).Run()
	if err != nil {
		log.Fatalf("Error downloading repository: %v", err)
	}

	// Install
	// Determine the file name from the URL
	tempFileName := filepath.Base(downloadURL)

	// Create the temporary file
	out, err := os.Create(fmt.Sprintf("/tmp/%s", tempFileName))
	if err != nil {
		log.Fatalf("Error creating temporary file: %v", err)
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

	// Extract binary
	err = exec.Command("tar", "-xf", fmt.Sprintf("/tmp/%s", tempFileName), "-C", fmt.Sprintf("%s", destDir)).Run()
	if err != nil {
		log.Fatalf("Error extracting binary to /usr/local/bin: %v", err)
	}

	// TODO
	// Currently, the downloaded binary is expected to have a name that matches the binaryName variable
	// Ideally, the extracted binary file should be renamed to match the binaryName variable

	// Define the final destination path
	destPath := filepath.Join(destDir, binaryName)

	// Make the file executable
	err = os.Chmod(destPath, 0755)
	if err != nil {
		log.Fatalf("Error making file executable: %v", err)
	}

	spinner.Success("strfry relay installed successfully.")
}
