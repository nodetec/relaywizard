package strfry29

import (
	"fmt"
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/nodetec/rwz/pkg/utils/git"
	"github.com/nodetec/rwz/pkg/utils/systemd"
	"github.com/nodetec/rwz/pkg/verification"
	"github.com/pterm/pterm"
	"path/filepath"
)

// Function to download and make the binary and plugin binary executable
func InstallRelayBinary() {
	downloadSpinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Downloading %s binaries...", RelayName))

	// Check for and remove existing git repository
	directories.RemoveDirectory(GitRepoTmpDirPath)

	// Download git repository
	git.Clone(GitRepoBranch, GitRepoURL, GitRepoTmpDirPath)

	directories.SetPermissions(GitRepoTmpDirPath, 0755)

	// Install
	// Determine the file name from the URL
	tmpBinaryFileName := filepath.Base(DownloadURL)

	// Temporary file path
	tmpBinaryFilePath := fmt.Sprintf("%s/%s", relays.TmpDirPath, tmpBinaryFileName)

	// Check if the temporary file exists and remove it if it does
	files.RemoveFile(tmpBinaryFilePath)

	// Download and copy the file
	files.DownloadAndCopyFile(tmpBinaryFilePath, DownloadURL)

	// Determine the file name from the URL
	tmpBinaryPluginFileName := filepath.Base(BinaryPluginDownloadURL)

	// Temporary file path
	tmpBinaryPluginFilePath := fmt.Sprintf("%s/%s", relays.TmpDirPath, tmpBinaryPluginFileName)

	// Check if the temporary file exists and remove it if it does
	files.RemoveFile(tmpBinaryPluginFilePath)

	// Download and copy the file
	files.DownloadAndCopyFile(tmpBinaryPluginFilePath, BinaryPluginDownloadURL)

	downloadSpinner.Success(fmt.Sprintf("%s binaries downloaded", RelayName))

	// Verify relay binary
	verification.VerifyRelayBinary(BinaryName, tmpBinaryFilePath)

	// Verify relay binary plugin
	verification.VerifyRelayBinary(fmt.Sprintf("%s plugin", RelayName), tmpBinaryPluginFilePath)

	installSpinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Installing %s binaries...", RelayName))

	// Check if the service file exists and disable and stop the service if it does
	if files.FileExists(ServiceFilePath) {
		// Disable and stop the Nostr relay service
		installSpinner.UpdateText("Disabling and stopping service...")
		systemd.DisableService(ServiceName)
		systemd.StopService(ServiceName)
	} else {
		installSpinner.UpdateText("Service file not found...")
	}

	// Extract relay binary
	files.ExtractFile(tmpBinaryFilePath, relays.BinaryDestDir)

	// Extract relay binary plugin
	files.ExtractFile(tmpBinaryPluginFilePath, relays.BinaryDestDir)

	// TODO
	// Currently, the downloaded binary is expected to have a name that matches the BinaryName variable
	// Ideally, the extracted binary file should be renamed to match the BinaryName variable

	// Define the final destination path
	destPath := filepath.Join(relays.BinaryDestDir, BinaryName)

	// Make the file executable
	files.SetPermissions(destPath, 0755)

	// Define the final destination path
	destPath = filepath.Join(relays.BinaryDestDir, BinaryPluginName)

	// Make the file executable
	files.SetPermissions(destPath, 0755)

	installSpinner.Success(fmt.Sprintf("%s binaries installed", RelayName))
}
