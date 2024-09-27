package wot_relay

import (
	"github.com/pterm/pterm"
	"log"
	"os"
	"os/exec"
	"text/template"
)

// Function to check if a user exists
func userExists(username string) bool {
	cmd := exec.Command("id", "-u", username)
	err := cmd.Run()
	return err == nil
}

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

	spinner, _ := pterm.DefaultSpinner.Start("Configuring relay service...")

	// Ensure the user for the relay service exists
	if !userExists("nostr") {
		spinner.UpdateText("Creating user 'nostr'...")
		err := exec.Command("adduser", "--disabled-login", "--gecos", "", "nostr").Run()
		if err != nil {
			log.Fatalf("Error creating user: %v", err)
		}
	} else {
		spinner.UpdateText("User 'nostr' already exists")
	}

	// Ensure the templates directory exists and set ownership
	spinner.UpdateText("Creating templates directory...")
	err := os.MkdirAll(templatesDir, 0755)
	if err != nil {
		log.Fatalf("Error creating templates directory: %v", err)
	}

	// Use chown command to set ownership of the templates directory to the nostr user
	err = exec.Command("chown", "-R", "nostr:nostr", templatesDir).Run()
	if err != nil {
		log.Fatalf("Error setting ownership of the templates directory: %v", err)
	}

	// Ensure the static directory exists and set ownership
	spinner.UpdateText("Creating static directory...")
	err = os.MkdirAll(staticDir, 0755)
	if err != nil {
		log.Fatalf("Error creating static directory: %v", err)
	}

	// Use chown command to set ownership of the static directory to the nostr user
	err = exec.Command("chown", "-R", "nostr:nostr", staticDir).Run()
	if err != nil {
		log.Fatalf("Error setting ownership of the static directory: %v", err)
	}

	// Ensure the data directory exists and set ownership
	spinner.UpdateText("Creating data directory...")
	err = os.MkdirAll(dataDir, 0755)
	if err != nil {
		log.Fatalf("Error creating data directory: %v", err)
	}

	// Use chown command to set ownership of the data directory to the nostr user
	err = exec.Command("chown", "-R", "nostr:nostr", dataDir).Run()
	if err != nil {
		log.Fatalf("Error setting ownership of the data directory: %v", err)
	}

	// Check if the index.html file exists and remove it if it does
	if _, err := os.Stat(indexFilePath); err == nil {
		err = os.Remove(indexFilePath)
		if err != nil {
			log.Fatalf("Error removing index.html file: %v", err)
		}
	}

	// Check if the environment file exists and remove it if it does
	if _, err := os.Stat(envFilePath); err == nil {
		err = os.Remove(envFilePath)
		if err != nil {
			log.Fatalf("Error removing environment file: %v", err)
		}
	}

	// Check if the service file exists and remove it if it does
	if _, err := os.Stat(serviceFilePath); err == nil {
		err = os.Remove(serviceFilePath)
		if err != nil {
			log.Fatalf("Error removing service file: %v", err)
		}
	}

	// Create the index.html file
	spinner.UpdateText("Creating index.html file...")
	indexFile, err := os.Create(indexFilePath)
	if err != nil {
		log.Fatalf("Error creating index.html file: %v", err)
	}
	defer indexFile.Close()

	indexTmpl, err := template.New("index").Parse(indexTemplate)
	if err != nil {
		log.Fatalf("Error parsing index.html template: %v", err)
	}

	err = indexTmpl.Execute(indexFile, struct{ Domain, PubKey string }{Domain: domain, PubKey: pubKey})
	if err != nil {
		log.Fatalf("Error executing index.html template: %v", err)
	}

	// Create the environment file
	spinner.UpdateText("Creating environment file...")
	envFile, err := os.Create(envFilePath)
	if err != nil {
		log.Fatalf("Error creating environment file: %v", err)
	}
	defer envFile.Close()

	envTmpl, err := template.New("env").Parse(envTemplate)
	if err != nil {
		log.Fatalf("Error parsing environment template: %v", err)
	}

	err = envTmpl.Execute(envFile, struct{ PubKey, Domain string }{PubKey: pubKey, Domain: domain})
	if err != nil {
		log.Fatalf("Error executing environment template: %v", err)
	}

	// Create the systemd service file
	spinner.UpdateText("Creating service file...")
	serviceFile, err := os.Create(serviceFilePath)
	if err != nil {
		log.Fatalf("Error creating service file: %v", err)
	}
	defer serviceFile.Close()

	tmpl, err := template.New("service").Parse(serviceTemplate)
	if err != nil {
		log.Fatalf("Error parsing service template: %v", err)
	}

	err = tmpl.Execute(serviceFile, struct{ PubKey, Domain string }{PubKey: pubKey, Domain: domain})
	if err != nil {
		log.Fatalf("Error executing service template: %v", err)
	}

	// Reload systemd to apply the new service
	spinner.UpdateText("Reloading systemd daemon...")
	err = exec.Command("systemctl", "daemon-reload").Run()
	if err != nil {
		log.Fatalf("Error reloading systemd daemon: %v", err)
	}

	// Enable and start the Nostr relay service
	spinner.UpdateText("Enabling and starting service...")
	err = exec.Command("systemctl", "enable", "wot-relay").Run()
	if err != nil {
		log.Fatalf("Error enabling Nostr relay service: %v", err)
	}

	err = exec.Command("systemctl", "start", "wot-relay").Run()
	if err != nil {
		log.Fatalf("Error starting Nostr relay service: %v", err)
	}

	spinner.Success("Nostr relay service configured")
}
