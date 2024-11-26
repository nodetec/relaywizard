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

// TODO
// Abstract this even more

func cloneTmpGitRepo() {
	// Check for and remove existing git repository
	directories.RemoveDirectory(GitRepoTmpDirPath)

	// Download git repository
	git.Clone(GitRepoBranch, GitRepoURL, GitRepoTmpDirPath)

	directories.SetPermissions(GitRepoTmpDirPath, 0755)
}

// Determine the temporary file name from the provided path
func tmpFilePathFromFilePath(path string) string {
	tmpFileName := filepath.Base(path)

	tmpFilePath := fmt.Sprintf("%s/%s", relays.TmpDirPath, tmpFileName)

	return tmpFilePath
}

func installRelayBinary(compressedBinaryFilePath, binaryName string) {
	// Extract relay binary
	files.ExtractFile(compressedBinaryFilePath, relays.BinaryDestDir)

	// TODO
	// Currently, the downloaded binary is expected to have a name that matches the binaryName variable
	// Ideally, the extracted binary file should be renamed to match the binaryName variable

	// Define the final destination path
	destPath := filepath.Join(relays.BinaryDestDir, binaryName)

	// Make the file executable
	files.SetPermissions(destPath, 0755)
}

// Function to download and make the binary and plugin binary executable
func InstallRelayBinaries() {
	pterm.Println()
	relayBinaryCheckSpinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Checking for existing %s binary...", BinaryName))

	cloneTmpGitRepo()

	// Check if the service file exists and disable and stop the service if it does
	systemd.DisableAndStopService(ServiceFilePath, ServiceName)

	// Check if relay binary exists
	if !files.FileExists(BinaryFilePath) {
		relayBinaryCheckSpinner.Info(fmt.Sprintf("%s binary not found", BinaryName))
		pterm.Println()

		// Determine the temporary file path
		tmpCompressedBinaryFilePath := tmpFilePathFromFilePath(DownloadURL)

		// Check if the temporary file exists and remove it if it does
		files.RemoveFile(tmpCompressedBinaryFilePath)

		// Download and copy the file
		downloadSpinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Downloading %s binary...", BinaryName))
		files.DownloadAndCopyFile(tmpCompressedBinaryFilePath, DownloadURL)
		downloadSpinner.Success(fmt.Sprintf("%s binary downloaded", BinaryName))

		// Verify relay binary
		verification.VerifyRelayBinary(BinaryName, tmpCompressedBinaryFilePath)

		installSpinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Installing %s binary...", BinaryName))
		installRelayBinary(tmpCompressedBinaryFilePath, BinaryName)
		installSpinner.Success(fmt.Sprintf("%s binary installed", BinaryName))
	} else {
		relayBinaryCheckSpinner.Info(fmt.Sprintf("%s binary found", BinaryName))
		pterm.Println()
	}

	// Determine the temporary file path
	tmpCompressedBinaryPluginFilePath := tmpFilePathFromFilePath(BinaryPluginDownloadURL)

	// Check if the temporary file exists and remove it if it does
	files.RemoveFile(tmpCompressedBinaryPluginFilePath)

	// Download and copy the file
	binaryPluginDownloadSpinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Downloading %s plugin binary...", BinaryPluginName))
	files.DownloadAndCopyFile(tmpCompressedBinaryPluginFilePath, BinaryPluginDownloadURL)
	binaryPluginDownloadSpinner.Success(fmt.Sprintf("%s plugin binary downloaded", BinaryPluginName))

	// Verify relay binary plugin
	verification.VerifyRelayBinary(fmt.Sprintf("%s plugin", BinaryPluginName), tmpCompressedBinaryPluginFilePath)

	binaryPluginInstallSpinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Installing %s plugin binary...", BinaryPluginName))
	installRelayBinary(tmpCompressedBinaryPluginFilePath, BinaryPluginName)
	binaryPluginInstallSpinner.Success(fmt.Sprintf("%s plugin binary installed", BinaryPluginName))
}
