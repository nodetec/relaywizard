package nostr_rs_relay

import (
	"fmt"
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/nodetec/rwz/pkg/utils/network"
	"github.com/pterm/pterm"
)

// Function to configure the relay
func ConfigureRelay(domain, pubKey, relayContact string, httpsEnabled bool) {
	spinner, _ := pterm.DefaultSpinner.Start("Configuring relay...")

	// Ensure the config directory exists and set permissions
	spinner.UpdateText("Creating config directory...")
	directories.CreateDirectory(ConfigDirPath, 0755)

	// Check for and remove existing config file
	files.RemoveFile(ConfigFilePath)

	// Construct the sed command to change the relay url
	files.InPlaceEdit(fmt.Sprintf(`s|relay_url = ".*"|relay_url = "%s://%s/"|`, network.WSEnabled(httpsEnabled), domain), TmpConfigFilePath)

	// Construct the sed command to change the pubkey
	files.InPlaceEdit(fmt.Sprintf(`s|#pubkey = ".*"|pubkey = "%s"|`, pubKey), TmpConfigFilePath)

	// Construct the sed command to change the contact
	files.InPlaceEdit(fmt.Sprintf(`s|#contact = ".*"|contact = "%s"|`, relayContact), TmpConfigFilePath)

	// Construct the sed command to change the data directory
	files.InPlaceEdit(fmt.Sprintf(`s|#data_directory = ".*"|data_directory = "%s/%s"|`, DataDirPath, relays.DBDir), TmpConfigFilePath)

	// Construct the sed command to change the remote ip header
	files.InPlaceEdit(fmt.Sprintf(`s|#remote_ip_header = "x-forwarded-for"|remote_ip_header = "x-forwarded-for"|`), TmpConfigFilePath)

	// Copy config file to config directory
	files.CopyFile(TmpConfigFilePath, ConfigDirPath)

	// Set permissions for the config file
	files.SetPermissions(ConfigFilePath, 0644)

	spinner.Success("Relay configured")
}
