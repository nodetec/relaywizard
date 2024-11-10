package strfry

import (
	"fmt"
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/nodetec/rwz/pkg/utils/git"
	"github.com/nodetec/rwz/pkg/utils/systemd"
	"github.com/pterm/pterm"
	"path/filepath"
)

// Function to download and make the binary executable
func InstallRelayBinary() {
	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Installing %s relay...", RelayName))

	// Check for and remove existing git repository
	directories.RemoveDirectory(GitRepoTmpDirPath)

	// Download git repository
	git.Clone(GitRepoBranch, GitRepoURL, GitRepoTmpDirPath)

	directories.SetPermissions(GitRepoTmpDirPath, 0755)

	// Install
	// Determine the file name from the URL
	tmpFileName := filepath.Base(DownloadURL)

	// Temporary file path
	tmpFilePath := fmt.Sprintf("%s/%s", relays.TmpDirPath, tmpFileName)

	// Check if the temporary file exists and remove it if it does
	files.RemoveFile(tmpFilePath)

	// Download and copy the file
	files.DownloadAndCopyFile(tmpFilePath, DownloadURL)

	// Check if the service file exists and disable and stop the service if it does
	if files.FileExists(ServiceFilePath) {
		// Disable and stop the Nostr relay service
		spinner.UpdateText("Disabling and stopping service...")
		systemd.DisableService(ServiceName)
		systemd.StopService(ServiceName)
	} else {
		spinner.UpdateText("Service file not found...")
	}

	// Extract binary
	files.ExtractFile(tmpFilePath, relays.BinaryDestDir)

	// TODO
	// Currently, the downloaded binary is expected to have a name that matches the BinaryName variable
	// Ideally, the extracted binary file should be renamed to match the BinaryName variable

	// Define the final destination path
	destPath := filepath.Join(relays.BinaryDestDir, BinaryName)

	// Make the file executable
	files.SetPermissions(destPath, 0755)

	spinner.Success(fmt.Sprintf("%s relay installed successfully.", RelayName))
}
