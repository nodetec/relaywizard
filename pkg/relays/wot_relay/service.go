package wot_relay

import (
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/nodetec/rwz/pkg/utils/systemd"
	"github.com/nodetec/rwz/pkg/utils/templates"
	"github.com/pterm/pterm"
)

// Function to set up the relay service
func SetupRelayService(domain, pubKey, relayContact string, httpsEnabled bool) {
	spinner, _ := pterm.DefaultSpinner.Start("Configuring relay service...")

	// Ensure the templates directory exists and set ownership
	spinner.UpdateText("Creating templates directory...")
	directories.CreateDirectory(TemplatesDirPath, 0755)

	// Use chown command to set ownership of the templates directory to the nostr user
	directories.SetOwnerAndGroup(relays.User, relays.User, TemplatesDirPath)

	// Ensure the static directory exists and set ownership
	spinner.UpdateText("Creating static directory...")
	directories.CreateDirectory(StaticDirPath, 0755)

	// Use chown command to set ownership of the static directory to the nostr user
	directories.SetOwnerAndGroup(relays.User, relays.User, StaticDirPath)

	// Ensure the data directory exists and set ownership
	spinner.UpdateText("Creating data directory...")
	directories.CreateDirectory(DataDirPath, 0755)

	// Use chown command to set ownership of the data directory to the nostr user
	directories.SetOwnerAndGroup(relays.User, relays.User, DataDirPath)

	// Check if the index.html file exists and remove it if it does
	files.RemoveFile(IndexFilePath)

	// Check if the environment file exists and remove it if it does
	files.RemoveFile(EnvFilePath)

	// Check if the service file exists and remove it if it does
	files.RemoveFile(ServiceFilePath)

	// Create the index.html file
	spinner.UpdateText("Creating index.html file...")
	indexFileParams := templates.IndexFileParams{Domain: domain, HTTPSEnabled: httpsEnabled, PubKey: pubKey}
	templates.CreateIndexFile(IndexFilePath, IndexFileTemplate, &indexFileParams)

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
