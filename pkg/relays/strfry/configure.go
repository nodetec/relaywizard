package strfry

import (
	"fmt"
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/pterm/pterm"
)

// Function to configure the relay
func ConfigureRelay(pubKey, relayContact string) {
	spinner, _ := pterm.DefaultSpinner.Start("Configuring relay...")

	// Ensure the config directory exists and set permissions
	spinner.UpdateText("Creating config directory...")
	directories.CreateDirectory(ConfigDirPath, 0755)

	// Check for and remove existing config file
	files.RemoveFile(ConfigFilePath)

	// Construct the sed command to change the db path
	files.InPlaceEdit(fmt.Sprintf(`s|db = ".*"|db = "%s/%s"|`, DataDirPath, relays.DBDir), TmpConfigFilePath)

	// Construct the sed command to change the nofiles limit
	files.InPlaceEdit(`s|nofiles = .*|nofiles = 0|`, TmpConfigFilePath)

	// Construct the sed command to change the realIpHeader
	files.InPlaceEdit(`s|realIpHeader = .*|realIpHeader = "x-forwarded-for"|`, TmpConfigFilePath)

	// Construct the sed command to change the pubkey
	files.InPlaceEdit(fmt.Sprintf(`s|pubkey = .*|pubkey = "%s"|`, pubKey), TmpConfigFilePath)

	// Construct the sed command to change the contact
	files.InPlaceEdit(fmt.Sprintf(`s|contact = ".*"|contact = "%s"|`, relayContact), TmpConfigFilePath)

	// Copy config file to config directory
	files.CopyFile(TmpConfigFilePath, ConfigDirPath)

	// Set permissions for the config file
	files.SetPermissions(ConfigFilePath, 0644)

	spinner.Success("Relay configured")
}
