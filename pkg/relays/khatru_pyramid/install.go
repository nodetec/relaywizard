package khatru_pyramid

import (
	"fmt"
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/pterm/pterm"
	"path/filepath"
)

// Function to download and make the binary executable
func InstallRelayBinary() {
	// URL of the binary to download
	const downloadURL = "https://github.com/nodetec/relays/releases/download/v0.2.0/khatru-pyramid-0.1.0-x86_64-linux-gnu.tar.gz"

	// Name of the binary after downloading
	const binaryName = "khatru-pyramid"

	// Destination directory for the binary
	const destDir = "/usr/local/bin"

	// Data directory for the relay
	const dataDir = "/var/lib/khatru-pyramid"

	spinner, _ := pterm.DefaultSpinner.Start("Installing Khatru Pyramid relay...")

	// Ensure the data directory exists
	directories.CreateDirectory(dataDir, 0755)

	// Determine the file name from the URL
	tmpFileName := filepath.Base(downloadURL)

	// Temporary file path
	tmpFilePath := fmt.Sprintf("/tmp/%s", tmpFileName)

	// Check if the temporary file exists and remove it if it does
	files.RemoveFile(tmpFilePath)

	// Download and copy the file
	files.DownloadAndCopyFile(tmpFilePath, downloadURL)

	// Extract binary
	files.ExtractFile(tmpFilePath, destDir)

	// TODO
	// Currently, the downloaded binary is expected to have a name that matches the binaryName variable
	// Ideally, the extracted binary file should be renamed to match the binaryName variable

	// Define the final destination path
	destPath := filepath.Join(destDir, binaryName)

	// Make the file executable
	files.SetPermissions(destPath, 0755)

	spinner.Success("Khatru Pyramid relay installed successfully.")
}
