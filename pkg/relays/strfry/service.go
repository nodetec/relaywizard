package strfry

import (
	"fmt"
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
// TODO
// Check working directory
// WorkingDirectory=/home/nostr
func SetupRelayService(domain string) {
	// Template for the systemd service file
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

	// Path for the systemd service file
	const serviceFilePath = "/etc/systemd/system/strfry.service"

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

	// Check if the service file exists and remove it if it does
	if _, err := os.Stat(serviceFilePath); err == nil {
		err = os.Remove(serviceFilePath)
		if err != nil {
			log.Fatalf("Error removing service file: %v", err)
		}
	}

	filePath := "/tmp/strfry/strfry.conf"

	// Construct the sed command to change the db path
	cmd := exec.Command("sed", "-i", fmt.Sprintf(`s|db = ".*"|db = "%s"|`, dataDir), filePath)

	// Execute the command
	if err = cmd.Run(); err != nil {
		log.Fatalf("Error changing the db path: %v", err)
	}

	// TODO
	// Determine system hard limit
	// Determine preferred nofiles value
	cmd = exec.Command("sed", "-i", `s|nofiles = .*|nofiles = 0|`, filePath)

	// Execute the command
	if err = cmd.Run(); err != nil {
		log.Fatalf("Error changing the nofiles option: %v", err)
	}

	// Check for and remove existing config file
	err = os.Remove("/etc/strfry.conf")
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("Error removing existing config file: %v", err)
	}

	// Copy config file to /etc
	err = exec.Command("cp", "/tmp/strfry/strfry.conf", "/etc").Run()
	if err != nil {
		log.Fatalf("Error copying config file: %v", err)
	}

	// Use chown command to set ownership of the config file to the nostr user
	err = exec.Command("chown", "nostr:nostr", "/etc/strfry.conf").Run()
	if err != nil {
		log.Fatalf("Error setting ownership of the config file: %v", err)
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
	err = exec.Command("systemctl", "enable", "strfry").Run()
	if err != nil {
		log.Fatalf("Error enabling Nostr relay service: %v", err)
	}

	err = exec.Command("systemctl", "start", "strfry").Run()
	if err != nil {
		log.Fatalf("Error starting Nostr relay service: %v", err)
	}

	spinner.Success("Nostr relay service configured")
}
