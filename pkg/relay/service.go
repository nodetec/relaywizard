package relay

import (
	"log"
	"os"
	"os/exec"
	"text/template"
	"github.com/pterm/pterm"
)


// Template for the systemd service file
const serviceTemplate = `[Unit]
Description=Nostr Relay Pyramid
After=network.target

[Service]
Type=simple
User=nostr
WorkingDirectory=/home/nostr
Environment="DOMAIN={{.Domain}}"
Environment="RELAY_NAME=nostr-relay-pyramid"
Environment="RELAY_PUBKEY={{.PubKey}}"
ExecStart=/usr/local/bin/nostr-relay-pyramid
Restart=on-failure

[Install]
WantedBy=multi-user.target
`

// Path for the systemd service file
const serviceFilePath = "/etc/systemd/system/nostr-relay-pyramid.service"

// Function to check if a user exists
func userExists(username string) bool {
	cmd := exec.Command("id", "-u", username)
	err := cmd.Run()
	return err == nil
}

// Function to set up the relay service
func SetupRelayService(domain, pubKey string) {

	spinner, _ := pterm.DefaultSpinner.Start("Configuring nginx for HTTP...")
	// Check if the service file exists and remove it if it does
	if _, err := os.Stat(serviceFilePath); err == nil {
		err = os.Remove(serviceFilePath)
		if err != nil {
			log.Fatalf("Error removing service file: %v", err)
		}
	}

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
	const dataDir = "/var/lib/nostr-relay-pyramid"
  spinner.UpdateText("Creating data directory...")
	err := os.MkdirAll(dataDir, 0755)
	if err != nil {
		log.Fatalf("Error creating data directory: %v", err)
	}

	err = os.Chown(dataDir, os.Getuid(), os.Getgid())
	if err != nil {
		log.Fatalf("Error setting ownership of the data directory: %v", err)
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

	err = tmpl.Execute(serviceFile, struct{ Domain, PubKey string }{Domain: domain, PubKey: pubKey})
	if err != nil {
		log.Fatalf("Error executing service template: %v", err)
	}

	// Reload systemd to apply the new service
  spinner.UpdateText("Reloading systemd daemon...")
	err = exec.Command("systemctl", "daemon-reload").Run()
	if err != nil {
		log.Fatalf("Error reloading systemd daemon: %v", err)
	}

	// Enable and start the nostr relay service
  spinner.UpdateText("Enabling and starting service...")
	err = exec.Command("systemctl", "enable", "nostr-relay-pyramid").Run()
	if err != nil {
		log.Fatalf("Error enabling nostr relay service: %v", err)
	}

	err = exec.Command("systemctl", "start", "nostr-relay-pyramid").Run()
	if err != nil {
		log.Fatalf("Error starting nostr relay service: %v", err)
	}

  spinner.Success("Nostr relay service configured")
}

