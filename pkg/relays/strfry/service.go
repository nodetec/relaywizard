package strfry

import (
	"fmt"
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/nodetec/rwz/pkg/utils/systemd"
	"github.com/nodetec/rwz/pkg/utils/users"
	"github.com/pterm/pterm"
)

// Function to set up the relay service
func SetupRelayService(domain string) {
	// Template for the systemd service file
	// TODO
	// Check working directory
	// WorkingDirectory=/home/nostr
	const serviceTemplate = `[Unit]
Description=strfry Nostr Relay Service
After=network.target

[Service]
Type=simple
User=nostr
Group=nostr
ExecStart=/usr/local/bin/strfry relay
Restart=on-failure
RestartSec=5
ProtectHome=yes
NoNewPrivileges=yes
ProtectSystem=full
LimitCORE=1000000000

[Install]
WantedBy=multi-user.target
`

	// Data directory
	const dataDir = "/var/lib/strfry"

	// Path for the temporary config file
	tmpConfigFilePath := "/tmp/strfry/strfry.conf"

	// Path for the config file
	configFilePath := "/etc/strfry.conf"

	// Path for the systemd service file
	const serviceFilePath = "/etc/systemd/system/strfry.service"

	// Relay service
	const relayService = "strfry"

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

	// Check if the service file exists and remove it if it does
	files.RemoveFile(serviceFilePath)

	// Construct the sed command to change the db path
	files.InPlaceEdit(fmt.Sprintf(`s|db = ".*"|db = "%s"|`, dataDir), tmpConfigFilePath)

	// Construct the sed command to change the nofiles limit
	// TODO
	// Determine system hard limit
	// Determine preferred nofiles value
	files.InPlaceEdit(`s|nofiles = .*|nofiles = 0|`, tmpConfigFilePath)

	// Check for and remove existing config file
	files.RemoveFile(configFilePath)

	// Copy config file to /etc
	files.CopyFile(tmpConfigFilePath, "/etc")

	// Use chown command to set ownership of the config file to the nostr user
	files.SetOwnerAndGroup("nostr", "nostr", configFilePath)

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
