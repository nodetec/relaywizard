package khatru_pyramid

import (
	"fmt"
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/nodetec/rwz/pkg/utils/systemd"
	"github.com/nodetec/rwz/pkg/verification"
	"github.com/pterm/pterm"
	"path/filepath"
)

// Function to download and make the binary executable
func InstallRelayBinary(pubKey string) {
	downloadSpinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Downloading %s binary...", RelayName))

	// Determine the file name from the URL
	tmpFileName := filepath.Base(DownloadURL)

	// Temporary file path
	tmpFilePath := fmt.Sprintf("%s/%s", relays.TmpDirPath, tmpFileName)

	// Check if the temporary file exists and remove it if it does
	files.RemoveFile(tmpFilePath)

	// Download and copy the file
	files.DownloadAndCopyFile(tmpFilePath, DownloadURL)

	downloadSpinner.Success(fmt.Sprintf("%s binary downloaded", RelayName))

	// Verify relay binary
	verification.VerifyRelayBinary(RelayName, tmpFilePath)

	installSpinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Installing %s binary...", RelayName))

	// Check if the service file exists and disable and stop the service if it does
	if files.FileExists(ServiceFilePath) {
		// Disable and stop the Nostr relay service
		installSpinner.UpdateText("Disabling and stopping service...")
		systemd.DisableService(ServiceName)
		systemd.StopService(ServiceName)
	} else {
		installSpinner.UpdateText("Service file not found...")
	}

	// Check if users.json file exists
	if files.FileExists(UsersFilePath) {
		// Check if the pubKey exists in the users.json file
		installSpinner.UpdateText("Checking for public key in users.json file...")
		lineExists := files.LineExists(fmt.Sprintf(`"%s":""`, pubKey), UsersFilePath)

		// If false remove data directory
		if !lineExists {
			installSpinner.UpdateText("Public key not found, removing data directory...")
			directories.RemoveDirectory(DataDirPath)
		} else {
			installSpinner.UpdateText("Public key found, keeping data directory.")
		}
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

	installSpinner.Success(fmt.Sprintf("%s binary installed", RelayName))
}
