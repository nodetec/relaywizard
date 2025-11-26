package strfry29

import (
	"fmt"
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/nodetec/rwz/pkg/utils/plugins"
	"github.com/pterm/pterm"
)

// Function to configure the relay
func ConfigureRelay(currentUsername, domain, pubKey, relaySecretKey, relayContact string) {
	spinner, _ := pterm.DefaultSpinner.Start("Configuring relay...")

	// Ensure the config directory exists and set permissions
	spinner.UpdateText("Creating config directory...")
	if currentUsername == relays.RootUser {
		directories.CreateDirectory(ConfigDirPath, 0755)
	} else {
		directories.CreateDirectoryUsingLinux(currentUsername, ConfigDirPath)
		directories.SetPermissionsUsingLinux(currentUsername, ConfigDirPath, "0755")
	}

	// Check if the config file exists and remove it if it does
	if currentUsername == relays.RootUser {
		files.RemoveFile(ConfigFilePath)
	} else {
		files.RemoveFileUsingLinux(currentUsername, ConfigFilePath)
	}

	// Check if the strfry29.json file exists and remove it if it does
	if currentUsername == relays.RootUser {
		files.RemoveFile(PluginFilePath)
	} else {
		files.RemoveFileUsingLinux(currentUsername, PluginFilePath)
	}

	// Construct the sed command to change the db path
	files.InPlaceEdit(fmt.Sprintf(`s|db = ".*"|db = "%s/%s"|`, DataDirPath, relays.DBDir), TmpConfigFilePath)

	// Construct the sed command to change the realIpHeader
	files.InPlaceEdit(`s|realIpHeader = .*|realIpHeader = "x-forwarded-for"|`, TmpConfigFilePath)

	// Construct the sed command to change the info description
	files.InPlaceEdit(fmt.Sprintf(`s|description = ".*"|description = "%s"|`, ConfigFileInfoDescription), TmpConfigFilePath)

	// Construct the sed command to change the pubkey
	files.InPlaceEdit(fmt.Sprintf(`s|pubkey = .*|pubkey = "%s"|`, pubKey), TmpConfigFilePath)

	// Construct the sed command to change the contact
	files.InPlaceEdit(fmt.Sprintf(`s|contact = ".*"|contact = "%s"|`, relayContact), TmpConfigFilePath)

	// Construct the sed command to change the plugin path
	files.InPlaceEdit(fmt.Sprintf(`s|plugin = ".*"|plugin = "%s"|`, BinaryPluginFilePath), TmpConfigFilePath)

	// Copy config file to /etc/strfry29
	files.CopyFile(currentUsername, TmpConfigFilePath, ConfigDirPath)

	// Set permissions for the config file
	if currentUsername == relays.RootUser {
		files.SetPermissions(ConfigFilePath, 0644)
	} else {
		files.SetPermissionsUsingLinux(currentUsername, ConfigFilePath, "0644")
	}

	// Create the strfry29.json file
	spinner.UpdateText("Creating plugin file...")
	pluginFileParams := plugins.PluginFileParams{Domain: domain, RelaySecretKey: relaySecretKey, ConfigFilePath: ConfigFilePath, BinaryFilePath: relays.Strfry29BinaryFilePath}
	plugins.CreatePluginFile(currentUsername, PluginFilePath, PluginFileTemplate, &pluginFileParams)

	// Set permissions for the strfry29.json file
	if currentUsername == relays.RootUser {
		files.SetPermissions(PluginFilePath, 0600)
	} else {
		files.SetPermissionsUsingLinux(currentUsername, ConfigFilePath, "0600")
	}

	spinner.Success("Relay configured")
}
