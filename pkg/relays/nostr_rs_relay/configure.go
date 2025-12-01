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
func ConfigureRelay(currentUsername, domain, pubKey, relayContact string, httpsEnabled bool) {
	spinner, _ := pterm.DefaultSpinner.Start("Configuring relay...")

	// Ensure the config directory exists and set permissions
	spinner.UpdateText("Creating config directory...")
	if currentUsername == relays.RootUser {
		directories.CreateAllDirectories(ConfigDirPath, 0755)
		directories.SetPermissions(ConfigDirPath, 0755)
	} else {
		directories.CreateAllDirectoriesUsingLinux(currentUsername, ConfigDirPath)
		directories.SetPermissionsUsingLinux(currentUsername, ConfigDirPath, "0755")
	}

	// Check for and remove existing config file
	if currentUsername == relays.RootUser {
		files.RemoveFile(ConfigFilePath)
	} else {
		files.RemoveFileUsingLinux(currentUsername, ConfigFilePath)
	}

	// Construct the sed command to change the relay url
	files.InPlaceEditUsingLinux(fmt.Sprintf(`s|relay_url = ".*"|relay_url = "%s://%s/"|`, network.WSEnabled(httpsEnabled), domain), TmpConfigFilePath)

	// Construct the sed command to change the pubkey
	files.InPlaceEditUsingLinux(fmt.Sprintf(`s|#pubkey = ".*"|pubkey = "%s"|`, pubKey), TmpConfigFilePath)

	// Construct the sed command to change the contact
	files.InPlaceEditUsingLinux(fmt.Sprintf(`s|#contact = ".*"|contact = "%s"|`, relayContact), TmpConfigFilePath)

	// Construct the sed command to change the data directory
	files.InPlaceEditUsingLinux(fmt.Sprintf(`s|#data_directory = ".*"|data_directory = "%s/%s"|`, DataDirPath, relays.DBDir), TmpConfigFilePath)

	// Construct the sed command to change the remote ip header
	files.InPlaceEditUsingLinux(fmt.Sprintf(`s|#remote_ip_header = "x-forwarded-for"|remote_ip_header = "x-forwarded-for"|`), TmpConfigFilePath)

	// Copy config file to config directory
	files.CopyFileUsingLinux(currentUsername, TmpConfigFilePath, ConfigDirPath)

	// Set permissions for the config file
	if currentUsername == relays.RootUser {
		files.SetPermissions(ConfigFilePath, 0644)
	} else {
		files.SetPermissionsUsingLinux(currentUsername, ConfigFilePath, "0644")
	}

	spinner.Success("Relay configured")
}
