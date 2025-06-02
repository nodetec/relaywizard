package strfry

import (
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/nodetec/rwz/pkg/utils/systemd"
	"github.com/pterm/pterm"
)

// Function to set up the relay service
func SetUpRelayService(relayUser string) {
	spinner, _ := pterm.DefaultSpinner.Start("Configuring relay service...")

	// Check if the service file exists and remove it if it does
	files.RemoveFile(ServiceFilePath)

	// Create the systemd service file
	spinner.UpdateText("Creating service file...")
	serviceFileParams := systemd.ServiceFileParams{RelayUser: relayUser, BinaryFilePath: BinaryFilePath, ConfigFilePath: ConfigFilePath}
	systemd.CreateServiceFile(ServiceFilePath, ServiceFileTemplate, &serviceFileParams)

	// Reload systemd to apply the new service
	spinner.UpdateText("Reloading systemd daemon...")
	systemd.Reload()

	// Enable and start the Nostr relay service
	spinner.UpdateText("Enabling and starting service...")
	systemd.EnableService(ServiceName)
	systemd.StartService(ServiceName)

	spinner.Success("Relay service enabled and started")
}
