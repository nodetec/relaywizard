package khatru29

import (
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/pterm/pterm"
)

// Function to set up the relay data directory
func SetUpRelayDataDir() {
	spinner, _ := pterm.DefaultSpinner.Start("Configuring relay data directory...")

	// Ensure the data directory exists and set permissions
	spinner.UpdateText("Creating data directory...")
	directories.CreateDirectory(DataDirPath, 0755)

	// Use chown command to set ownership of the data directory to the nostr user
	directories.SetOwnerAndGroup(relays.User, relays.User, DataDirPath)

	spinner.Success("Data directory set up")
}
