package wot_relay

import (
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/configuration"
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/nodetec/rwz/pkg/utils/files"
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

	// Check if the environment file exists and remove it if it does
	if currentUsername == relays.RootUser {
		files.RemoveFile(EnvFilePath)
	} else {
		files.RemoveFileUsingLinux(currentUsername, EnvFilePath)
	}

	// Create the environment file
	spinner.UpdateText("Creating environment file...")
	envFileParams := configuration.EnvFileParams{Domain: domain, HTTPSEnabled: httpsEnabled, PubKey: pubKey, RelayContact: relayContact}
	configuration.CreateEnvFile(currentUsername, EnvFilePath, EnvFileTemplate, &envFileParams)

	// Set permissions for the environment file
	if currentUsername == relays.RootUser {
		files.SetPermissions(EnvFilePath, 0644)
	} else {
		files.SetPermissionsUsingLinux(currentUsername, EnvFilePath, "0644")
	}

	spinner.Success("Relay configured")
}
