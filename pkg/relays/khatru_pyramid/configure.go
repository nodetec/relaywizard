package khatru_pyramid

import (
	"github.com/nodetec/rwz/pkg/utils/configuration"
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/pterm/pterm"
)

// Function to configure the relay
func ConfigureRelay(domain, pubKey, relayContact string) {
	spinner, _ := pterm.DefaultSpinner.Start("Configuring relay...")

	// Ensure the config directory exists and set permissions
	spinner.UpdateText("Creating config directory...")
	directories.CreateDirectory(ConfigDirPath, 0755)

	// Check if the environment file exists and remove it if it does
	files.RemoveFile(EnvFilePath)

	// Create the environment file
	spinner.UpdateText("Creating environment file...")
	envFileParams := configuration.EnvFileParams{Domain: domain, PubKey: pubKey, RelayContact: relayContact}
	configuration.CreateEnvFile(EnvFilePath, EnvFileTemplate, &envFileParams)

	// Set permissions for the environment file
	files.SetPermissions(EnvFilePath, 0644)

	spinner.Success("Relay configured")
}
