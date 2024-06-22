package network

import (
	"log"
	"os/exec"

	"github.com/pterm/pterm"
)

// Function to configure the firewall
func ConfigureFirewall() {
	spinner, _ := pterm.DefaultSpinner.Start("Configuring firewall...")
	exec.Command("ufw", "enable").Run()

	// Allow HTTP and HTTPS traffic
	err := exec.Command("ufw", "allow", "Nginx Full").Run()
	if err != nil {
		log.Fatalf("Error allowing Nginx Full: %v", err)
	}

	// Reload the firewall to apply the changes
	err = exec.Command("ufw", "reload").Run()
	if err != nil {
		log.Fatalf("Error reloading firewall: %v", err)
	}

	spinner.Success("Firewall configured successfully.")
}
