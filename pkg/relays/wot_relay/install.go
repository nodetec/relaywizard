package wot_relay

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
	// URL of the binary to download
	const downloadURL = "https://github.com/nodetec/relays/releases/download/v0.2.0/wot-relay-0.1.12-x86_64-linux-gnu.tar.gz"

	// Name of the binary after downloading
	const binaryName = "wot-relay"

	// Destination directory for the binary
	const destDir = "/usr/local/bin"

	// Templates directory for the relay
	const templatesDir = "/etc/wot-relay/templates"

	// Static directory for the relay
	const staticDir = "/etc/wot-relay/templates/static"

	// Data directory for the relay
	const dataDir = "/var/lib/wot-relay"

	spinner, _ := pterm.DefaultSpinner.Start("Installing WoT Relay...")

	// Ensure the templates directory exists
	err := os.MkdirAll(templatesDir, 0755)
	if err != nil {
		log.Fatalf("Error creating templates directory: %v", err)
	}

	// Ensure the static directory exists
	err = os.MkdirAll(staticDir, 0755)
	if err != nil {
		log.Fatalf("Error creating static directory: %v", err)
	}

	// Ensure the data directory exists
	err = os.MkdirAll(dataDir, 0755)
	if err != nil {
		log.Fatalf("Error creating data directory: %v", err)
	}

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

	spinner.Success("WoT Relay installed successfully.")
}
