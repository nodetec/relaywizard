package wot_relay

import (
	"fmt"
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/pterm/pterm"
)

// Function to set up the relay data directory
func SetUpRelayDataDir() {
	spinner, _ := pterm.DefaultSpinner.Start("Configuring relay data directory...")

	// TODO
	// Look into how to back up WoT Relay databases
	spinner.UpdateText("Checking for existing data directory...")
	if directories.DirExists(DataDirPath) {
		spinner.UpdateText("Removing existing data directory...")
		directories.RemoveDirectory(DataDirPath)
	}

	// Ensure the data directory exists and set permissions
	spinner.UpdateText("Creating data directory...")
	directories.CreateDirectory(DataDirPath, 0755)
	directories.CreateDirectory(fmt.Sprintf("%s/%s", DataDirPath, relays.DBDir), 0755)

	spinner.Success("Data directory set up")
}
