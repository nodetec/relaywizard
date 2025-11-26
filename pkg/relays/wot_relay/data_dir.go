package wot_relay

import (
	"fmt"
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/pterm/pterm"
)

// Function to set up the relay data directory
func SetUpRelayDataDir(currentUsername string) {
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

	// Ensure the data directory exists and set permissions
	spinner.UpdateText("Creating data directory...")
	if currentUsername == relays.RootUser {
		directories.CreateDirectory(DataDirPath, 0755)
		directories.CreateDirectory(fmt.Sprintf("%s/%s", DataDirPath, relays.DBDir), 0755)
	} else {
		directories.CreateDirectoryUsingLinux(currentUsername, DataDirPath)
		directories.SetPermissionsUsingLinux(currentUsername, DataDirPath, "0755")
		directories.CreateDirectoryUsingLinux(currentUsername, fmt.Sprintf("%s/%s", DataDirPath, relays.DBDir))
		directories.SetPermissionsUsingLinux(currentUsername, fmt.Sprintf("%s/%s", DataDirPath, relays.DBDir), "0755")
	}

	spinner.Success("Data directory set up")
}
