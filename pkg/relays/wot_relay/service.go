package wot_relay

import (
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/nodetec/rwz/pkg/utils/systemd"
	"github.com/nodetec/rwz/pkg/utils/templates"
	"github.com/nodetec/rwz/pkg/utils/users"
	"github.com/pterm/pterm"
)

// Function to set up the relay service
func SetupRelayService(domain, pubKey string) {
	// Template for index.html file
	const indexTemplate = `<!doctype html>
<html lang="en">
	<head>
		<meta charset="UTF-8" />
		<meta name="viewport" content="width=device-width, initial-scale=1.0" />
			<title>WoT Relay</title>
			<meta name="description" content="WoT Relay" />
			<link href="https://{{.Domain}}" rel="canonical" />
	</head>
	<body>
		<main>
			<div>
				<div>
					<span>WoT Relay</span>
				</div>
				<div>
					<span>Domain: {{.Domain}}</span>
				</div>
				<div>
					<span>Pubkey: {{.PubKey}}</span>
				</div>
			</div>
		</main>
	</body>
</html>
`

	// Template for the environment file
	const envTemplate = `RELAY_NAME="WoT Relay"
RELAY_DESCRIPTION="WoT Nostr Relay"
RELAY_ICON="https://pfp.nostr.build/56306a93a88d4c657d8a3dfa57b55a4ed65b709eee927b5dafaab4d5330db21f.png"
RELAY_URL="wss://{{.Domain}}"
RELAY_PUBKEY="{{.PubKey}}"
RELAY_CONTACT="{{.PubKey}}"
INDEX_PATH="/etc/wot-relay/templates/index.html"
STATIC_PATH="/etc/wot-relay/templates/static"
DB_PATH="/var/lib/wot-relay/db"
REFRESH_INTERVAL_HOURS=24
MINIMUM_FOLLOWERS=3
ARCHIVAL_SYNC="FALSE"
ARCHIVE_REACTIONS="FALSE"
`

	// Template for the systemd service file
	const serviceTemplate = `[Unit]
Description=WoT Nostr Relay Service
After=network.target

[Service]
Type=simple
User=nostr
Group=nostr
WorkingDirectory=/home/nostr
EnvironmentFile=/etc/systemd/system/wot-relay.env
ExecStart=/usr/local/bin/wot-relay
Restart=on-failure
MemoryHigh=512M
MemoryMax=1G

[Install]
WantedBy=multi-user.target
`

	// Templates directory
	const templatesDir = "/etc/wot-relay/templates"

	// Static directory
	const staticDir = "/etc/wot-relay/templates/static"

	// Data directory
	const dataDir = "/var/lib/wot-relay"

	// Path for the index.html file
	const indexFilePath = "/etc/wot-relay/templates/index.html"

	// Path for the environment file
	const envFilePath = "/etc/systemd/system/wot-relay.env"

	// Path for the systemd service file
	const serviceFilePath = "/etc/systemd/system/wot-relay.service"

	// Relay service
	const relayService = "wot-relay"

	spinner, _ := pterm.DefaultSpinner.Start("Configuring relay service...")

	// Ensure the user for the relay service exists
	if !users.UserExists("nostr") {
		spinner.UpdateText("Creating user 'nostr'...")
		users.CreateUser("nostr", true)
	} else {
		spinner.UpdateText("User 'nostr' already exists")
	}

	// Ensure the templates directory exists and set ownership
	spinner.UpdateText("Creating templates directory...")
	directories.CreateDirectory(templatesDir, 0755)

	// Use chown command to set ownership of the templates directory to the nostr user
	directories.SetOwnerAndGroup("nostr", "nostr", templatesDir)

	// Ensure the static directory exists and set ownership
	spinner.UpdateText("Creating static directory...")
	directories.CreateDirectory(staticDir, 0755)

	// Use chown command to set ownership of the static directory to the nostr user
	directories.SetOwnerAndGroup("nostr", "nostr", staticDir)

	// Ensure the data directory exists and set ownership
	spinner.UpdateText("Creating data directory...")
	directories.CreateDirectory(dataDir, 0755)

	// Use chown command to set ownership of the data directory to the nostr user
	directories.SetOwnerAndGroup("nostr", "nostr", dataDir)

	// Check if the index.html file exists and remove it if it does
	files.RemoveFile(indexFilePath)

	// Check if the environment file exists and remove it if it does
	files.RemoveFile(envFilePath)

	// Check if the service file exists and remove it if it does
	files.RemoveFile(serviceFilePath)

	// Create the index.html file
	spinner.UpdateText("Creating index.html file...")
	indexFileParams := templates.IndexFileParams{Domain: domain, PubKey: pubKey}
	templates.CreateIndexFile(indexFilePath, indexTemplate, &indexFileParams)

	// Create the environment file
	spinner.UpdateText("Creating environment file...")
	envFileParams := systemd.EnvFileParams{Domain: domain, PubKey: pubKey}
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
