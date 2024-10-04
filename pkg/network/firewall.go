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
	exec.Command("ufw", "enable").Run()

	// Allow HTTP and HTTPS traffic
	err := exec.Command("ufw", "allow", "Nginx Full").Run()
	if err != nil {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Failed to allow Nginx Full: %v", err))
		os.Exit(1)
	}

	// Reload the firewall to apply the changes
	err = exec.Command("ufw", "reload").Run()
	if err != nil {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Failed to reload firewall: %v", err))
		os.Exit(1)
	}

	spinner.Success("Firewall configured successfully.")
}
