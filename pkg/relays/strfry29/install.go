package strfry29

import (
	"fmt"
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/nodetec/rwz/pkg/utils/git"
	"github.com/pterm/pterm"
	"path/filepath"
)

// Function to download and make the binary executable
func InstallRelayBinary() {
	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Installing %s relay...", RelayName))

	// Check for and remove existing git repository
	directories.RemoveDirectory(GitRepoTmpDir)

	// Download git repository
	git.Clone(GitRepoBranch, GitRepoURL, GitRepoTmpDir)

	// Install
	// Determine the file name from the URL
	tmpFileName := filepath.Base(DownloadURL)

	// Temporary file path
	tmpFilePath := fmt.Sprintf("%s/%s", relays.TmpDirPath, tmpFileName)

	// Check if the temporary file exists and remove it if it does
	files.RemoveFile(tmpFilePath)

	// Download and copy the file
	files.DownloadAndCopyFile(tmpFilePath, DownloadURL)

	// Extract binary
	files.ExtractFile(tmpFilePath, relays.BinaryDestDir)

	// Determine the file name from the URL
	tmpFileName = filepath.Base(BinaryPluginDownloadURL)

	// Temporary file path
	tmpFilePath = fmt.Sprintf("%s/%s", relays.TmpDirPath, tmpFileName)

	// Check if the temporary file exists and remove it if it does
	files.RemoveFile(tmpFilePath)

	// Download and copy the file
	files.DownloadAndCopyFile(tmpFilePath, BinaryPluginDownloadURL)

	// Extract binary
	files.ExtractFile(tmpFilePath, relays.BinaryDestDir)

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

	spinner.Success(fmt.Sprintf("%s relay installed successfully.", RelayName))
}
