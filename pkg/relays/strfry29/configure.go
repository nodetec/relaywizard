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
func ConfigureRelay(domain, pubKey, relaySecretKey, relayContact string) {
	spinner, _ := pterm.DefaultSpinner.Start("Configuring relay...")

	// Ensure the config directory exists and set permissions
	spinner.UpdateText("Creating config directory...")
	directories.CreateDirectory(ConfigDirPath, 0755)

	// Check if the config file exists and remove it if it does
	files.RemoveFile(ConfigFilePath)

	// Check if the strfry29.json file exists and remove it if it does
	files.RemoveFile(PluginFilePath)

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
	files.CopyFile(TmpConfigFilePath, ConfigDirPath)

	// Set permissions for the config file
	files.SetPermissions(ConfigFilePath, 0644)

	// Create the strfry29.json file
	spinner.UpdateText("Creating plugin file...")
	pluginFileParams := plugins.PluginFileParams{Domain: domain, RelaySecretKey: relaySecretKey, ConfigFilePath: ConfigFilePath, BinaryFilePath: relays.Strfry29BinaryFilePath}
	plugins.CreatePluginFile(PluginFilePath, PluginFileTemplate, &pluginFileParams)

	// Set permissions for the strfry29.json file
	files.SetPermissions(PluginFilePath, 0600)

	spinner.Success("Relay configured")
}
