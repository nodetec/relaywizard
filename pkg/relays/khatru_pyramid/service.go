package khatru_pyramid

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
	// Template for the environment file
	const envTemplate = `DOMAIN="{{.Domain}}"
DATABASE_PATH="/var/lib/khatru-pyramid/db"
USERDATA_PATH="/var/lib/khatru-pyramid/users.json"
MAX_INVITES_PER_PERSON="3"
RELAY_NAME="Khatru Pyramid"
RELAY_PUBKEY="{{.PubKey}}"
RELAY_DESCRIPTION="Khatru Pyramid Nostr Relay"
RELAY_CONTACT="your-email@example.com"
`
	// Template for the systemd service file
	const serviceTemplate = `[Unit]
Description=Khatru Pyramid Nostr Relay Service
After=network.target

[Service]
Type=simple
User=nostr
Group=nostr
WorkingDirectory=/home/nostr
EnvironmentFile=/etc/systemd/system/khatru-pyramid.env
ExecStart=/usr/local/bin/khatru-pyramid
Restart=on-failure

[Install]
WantedBy=multi-user.target
`

	// Data directory
	const dataDir = "/var/lib/khatru-pyramid"

	// Path for the environment file
	const envFilePath = "/etc/systemd/system/khatru-pyramid.env"

	// Path for the systemd service file
	const serviceFilePath = "/etc/systemd/system/khatru-pyramid.service"

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

	// Ensure the data directory exists and set ownership
	spinner.UpdateText("Creating data directory...")
	err := os.MkdirAll(dataDir, 0755)
	if err != nil {
		log.Fatalf("Error creating data directory: %v", err)
	}

	// Use chown command to set ownership of the data directory to the nostr user
	err = exec.Command("chown", "-R", "nostr:nostr", dataDir).Run()
	if err != nil {
		log.Fatalf("Error setting ownership of the data directory: %v", err)
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

	err = envTmpl.Execute(envFile, struct{ Domain, PubKey string }{Domain: domain, PubKey: pubKey})
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

	err = tmpl.Execute(serviceFile, struct{}{})
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
	err = exec.Command("systemctl", "enable", "khatru-pyramid").Run()
	if err != nil {
		log.Fatalf("Error enabling Nostr relay service: %v", err)
	}

	err = exec.Command("systemctl", "start", "khatru-pyramid").Run()
	if err != nil {
		log.Fatalf("Error starting Nostr relay service: %v", err)
	}

	spinner.Success("Nostr relay service configured")
}
