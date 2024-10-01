package strfry

import (
	"fmt"
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/nodetec/rwz/pkg/utils/git"
	"github.com/pterm/pterm"
	"path/filepath"
)

// Function to download and make the binary executable
func InstallRelayBinary() {
	// Git repository branch
	const branch = "1.0.1"

	// Git repository url
	const gitURL = "https://github.com/hoytech/strfry.git"

	// Temporary directory for git repository
	const tmpDir = "/tmp/strfry"

	// URL of the binary to download
	const downloadURL = "https://github.com/nodetec/relays/releases/download/v0.2.0/strfry-1.0.1-x86_64-linux-gnu.tar.gz"

	// Name of the binary after downloading
	const binaryName = "strfry"

	// Destination directory for the binary
	const destDir = "/usr/local/bin"

	spinner, _ := pterm.DefaultSpinner.Start("Installing strfry relay...")

	// Check for and remove existing git repository
	directories.RemoveDirectory(tmpDir)

	// Download git repository
	git.Clone(branch, gitURL, tmpDir)

	// Install
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

	spinner.Success("strfry relay installed successfully.")
}
