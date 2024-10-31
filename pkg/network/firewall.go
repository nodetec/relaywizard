package network

import (
	"fmt"
	"github.com/pterm/pterm"
	"os"
	"os/exec"
)

// Function to configure the firewall
func ConfigureFirewall() {
	spinner, _ := pterm.DefaultSpinner.Start("Configuring firewall...")
	// Allow SSH connections
	err := exec.Command("ufw", "allow", "ssh").Run()
	if err != nil {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Failed to allow SSH: %v", err))
		os.Exit(1)
	}

	// Allow HTTP and HTTPS traffic
	err = exec.Command("ufw", "allow", "Nginx Full").Run()
	if err != nil {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Failed to allow Nginx Full: %v", err))
		os.Exit(1)
	}

	// Disable logging
	err = exec.Command("ufw", "logging", "off").Run()
	if err != nil {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Failed to disable logging: %v", err))
		os.Exit(1)
	}

	// Enable the firewall to apply the changes
	err = exec.Command("ufw", "--force", "enable").Run()
	if err != nil {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Failed to enable firewall: %v", err))
		os.Exit(1)
	}

	spinner.Success("Firewall configured successfully.")
}
