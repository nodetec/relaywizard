package khatru29

import (
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/nodetec/rwz/pkg/utils/systemd"
	"github.com/nodetec/rwz/pkg/utils/users"
	"github.com/pterm/pterm"
)

// Function to set up the relay service
func SetupRelayService(domain, privKey string) {
	// Template for the environment file
	const envTemplate = `PORT="5577"
DOMAIN="{{.Domain}}"
RELAY_NAME="Khatru29"
RELAY_PRIVKEY="{{.PrivKey}}"
RELAY_DESCRIPTION="Khatru29 Nostr Relay"
RELAY_CONTACT="your-email@example.com"
DATABASE_PATH=/var/lib/khatru29/db
`

	// Template for the systemd service file
	const serviceTemplate = `[Unit]
Description=Khatru29 Nostr Relay Service
After=network.target

[Service]
Type=simple
User=nostr
Group=nostr
WorkingDirectory=/home/nostr
EnvironmentFile=/etc/systemd/system/khatru29.env
ExecStart=/usr/local/bin/khatru29
Restart=on-failure

[Install]
WantedBy=multi-user.target
`

	// Data directory
	const dataDir = "/var/lib/khatru29"

	// Path for the environment file
	const envFilePath = "/etc/systemd/system/khatru29.env"

	// Path for the systemd service file
	const serviceFilePath = "/etc/systemd/system/khatru29.service"

	// Relay service
	const relayService = "khatru29"

	spinner, _ := pterm.DefaultSpinner.Start("Configuring relay service...")

	// Ensure the user for the relay service exists
	if !users.UserExists("nostr") {
		spinner.UpdateText("Creating user 'nostr'...")
		users.CreateUser("nostr", true)
	} else {
		spinner.UpdateText("User 'nostr' already exists")
	}

	// Ensure the data directory exists and set ownership
	spinner.UpdateText("Creating data directory...")
	directories.CreateDirectory(dataDir, 0755)

	// Use chown command to set ownership of the data directory to the nostr user
	directories.SetOwnerAndGroup("nostr", "nostr", dataDir)

	// Check if the environment file exists and remove it if it does
	files.RemoveFile(envFilePath)

	// Check if the service file exists and remove it if it does
	files.RemoveFile(serviceFilePath)

	// Create the environment file
	spinner.UpdateText("Creating environment file...")
	envFileParams := systemd.EnvFileParams{Domain: domain, PrivKey: privKey}
	systemd.CreateEnvFile(envFilePath, envTemplate, &envFileParams)

	// Create the systemd service file
	spinner.UpdateText("Creating service file...")
	systemd.CreateServiceFile(serviceFilePath, serviceTemplate)

	// Reload systemd to apply the new service
	spinner.UpdateText("Reloading systemd daemon...")
	systemd.Reload()

	// Enable and start the Nostr relay service
	spinner.UpdateText("Enabling and starting service...")
	systemd.EnableService(relayService)
	systemd.StartService(relayService)

	spinner.Success("Nostr relay service configured")
}
