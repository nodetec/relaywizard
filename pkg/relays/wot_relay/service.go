package wot_relay

import (
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/nodetec/rwz/pkg/utils/systemd"
	"github.com/pterm/pterm"
)

// Function to set up the relay service
func SetUpRelayService(currentUsername, relayUser string) {
	spinner, _ := pterm.DefaultSpinner.Start("Configuring relay service...")

	// Check if the service file exists and remove it if it does
	if currentUsername == relays.RootUser {
		files.RemoveFile(ServiceFilePath)
	} else {
		files.RemoveFileUsingLinux(currentUsername, ServiceFilePath)
	}

	// Create the systemd service file
	spinner.UpdateText("Creating service file...")
	serviceFileParams := systemd.ServiceFileParams{RelayUser: relayUser, EnvFilePath: EnvFilePath, BinaryFilePath: relays.WotRelayBinaryFilePath}
	systemd.CreateServiceFile(currentUsername, ServiceFilePath, ServiceFileTemplate, "0644", &serviceFileParams)

	// Reload systemd to apply the new service
	spinner.UpdateText("Reloading systemd daemon...")
	systemd.Reload(currentUsername)

	// Enable and start the Nostr relay service
	spinner.UpdateText("Enabling and starting service...")
	systemd.EnableService(currentUsername, ServiceName)
	systemd.StartService(currentUsername, ServiceName)

	spinner.Success("Relay service enabled and started")
}
