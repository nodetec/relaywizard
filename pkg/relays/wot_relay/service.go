package wot_relay

import (
	"fmt"
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/nodetec/rwz/pkg/utils/systemd"
	"github.com/pterm/pterm"
)

// Function to set up the relay service
func SetupRelayService(domain, pubKey, relayContact string, httpsEnabled bool) {
	spinner, _ := pterm.DefaultSpinner.Start("Configuring relay service...")

	// Ensure the data directory exists and set ownership
	spinner.UpdateText("Creating data directory...")
	directories.CreateDirectory(DataDirPath, 0755)
	directories.CreateDirectory(fmt.Sprintf("%s/db", DataDirPath), 0755)

	// Use chown command to set ownership of the data directory and its content to the nostr user
	directories.SetOwnerAndGroup(relays.User, relays.User, DataDirPath)

	// Ensure the config directory exists and set permissions
	spinner.UpdateText("Creating config directory...")
	directories.CreateDirectory(ConfigDirPath, 0755)

	// Ensure the templates directory exists and set permissions
	spinner.UpdateText("Creating templates directory...")
	directories.CreateDirectory(TemplatesDirPath, 0755)

	// Use chown command to set ownership of the config directory and its content to the nostr user
	directories.SetOwnerAndGroup(relays.User, relays.User, ConfigDirPath)

	// Check if the index.html file exists and remove it if it does
	files.RemoveFile(IndexFilePath)

	// Copy the index.html file to templates directory
	files.CopyFile(TmpIndexFilePath, TemplatesDirPath)

	// Use chown command to set ownership of the index.html file to the nostr user
	files.SetOwnerAndGroup(relays.User, relays.User, IndexFilePath)

	// Remove the static directory and all of its content if it exists
	spinner.UpdateText("Removing static directory...")
	directories.RemoveDirectory(StaticDirPath)

	// Copy the static directory and all of its content to the templates directory
	directories.CopyDirectory(TmpStaticDirPath, TemplatesDirPath)

	// Use chown command to set ownership of the static directory and its content to the nostr user
	directories.SetOwnerAndGroup(relays.User, relays.User, StaticDirPath)

	// Check if the environment file exists and remove it if it does
	files.RemoveFile(EnvFilePath)

	// Check if the service file exists and remove it if it does
	files.RemoveFile(ServiceFilePath)

	// Create the environment file
	spinner.UpdateText("Creating environment file...")
	envFileParams := systemd.EnvFileParams{Domain: domain, HTTPSEnabled: httpsEnabled, PubKey: pubKey, RelayContact: relayContact}
	systemd.CreateEnvFile(EnvFilePath, EnvFileTemplate, &envFileParams)

	// Create the systemd service file
	spinner.UpdateText("Creating service file...")
	serviceFileParams := systemd.ServiceFileParams{EnvFilePath: EnvFilePath, BinaryFilePath: BinaryFilePath}
	systemd.CreateServiceFile(ServiceFilePath, ServiceFileTemplate, &serviceFileParams)

	// Reload systemd to apply the new service
	spinner.UpdateText("Reloading systemd daemon...")
	systemd.Reload()

	// Enable and start the Nostr relay service
	spinner.UpdateText("Enabling and starting service...")
	systemd.EnableService(ServiceName)
	systemd.StartService(ServiceName)

	spinner.Success("Nostr relay service configured")
}
