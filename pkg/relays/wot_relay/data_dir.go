package wot_relay

import (
	"fmt"

	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/pterm/pterm"
)

// Function to set up the relay data directory
func SetUpRelayDataDir(currentUsername, relayUser string) {
	spinner, _ := pterm.DefaultSpinner.Start("Configuring relay data directory...")

	// TODO
	// Look into how to back up WoT Relay databases
	spinner.UpdateText("Checking for existing data directory...")
	if directories.DirExists(DataDirPath) {
		spinner.UpdateText("Removing existing data directory...")
		if currentUsername == relays.RootUser {
			directories.RemoveDirectory(DataDirPath)
		} else {
			directories.RemoveDirectoryUsingLinux(currentUsername, DataDirPath)
		}
	}

	dataDBDirPath := fmt.Sprintf("%s/%s", DataDirPath, relays.DBDir)

	// Ensure the data directory exists and set permissions
	spinner.UpdateText("Creating data directory...")
	if currentUsername == relays.RootUser {
		directories.CreateAllDirectories(dataDBDirPath, 0755)
		directories.SetPermissions(DataDirPath, 0755)
		directories.SetPermissions(dataDBDirPath, 0755)
	} else {
		directories.CreateAllDirectoriesUsingLinux(currentUsername, dataDBDirPath)
		directories.SetPermissionsUsingLinux(currentUsername, DataDirPath, "0755")
		directories.SetPermissionsUsingLinux(currentUsername, dataDBDirPath, "0755")
		directories.SetOwnerAndGroupForAllContentUsingLinux(currentUsername, relayUser, relayUser, DataDirPath)
	}

	spinner.Success("Data directory set up")
}
