package network

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/nodetec/rwz/pkg/logs"
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/logging"
	"github.com/pterm/pterm"
)

// Function to configure the firewall
func ConfigureFirewall(currentUsername string) {
	spinner, _ := pterm.DefaultSpinner.Start("Configuring firewall...")
	if currentUsername == relays.RootUser {
		// Allow SSH connections
		err := exec.Command("ufw", "allow", "ssh").Run()
		if err != nil {
			logging.AppendRWZLogFile(currentUsername, logs.RWZLogFilePath, fmt.Sprintf("Failed to allow SSH: %v", err))
			pterm.Println()
			pterm.Error.Printfln("Failed to allow SSH: %v", err)
			os.Exit(1)
		}
		// Allow HTTP and HTTPS traffic
		err = exec.Command("ufw", "allow", "Nginx Full").Run()
		if err != nil {
			logging.AppendRWZLogFile(currentUsername, logs.RWZLogFilePath, fmt.Sprintf("Failed to allow Nginx Full: %v", err))
			pterm.Println()
			pterm.Error.Printfln("Failed to allow Nginx Full: %v", err)
			os.Exit(1)
		}

		// Disable logging
		err = exec.Command("ufw", "logging", "off").Run()
		if err != nil {
			logging.AppendRWZLogFile(currentUsername, logs.RWZLogFilePath, fmt.Sprintf("Failed to disable logging: %v", err))
			pterm.Println()
			pterm.Error.Printfln("Failed to disable logging: %v", err)
			os.Exit(1)
		}

		// Enable the firewall to apply the changes
		err = exec.Command("ufw", "--force", "enable").Run()
		if err != nil {
			logging.AppendRWZLogFile(currentUsername, logs.RWZLogFilePath, fmt.Sprintf("Failed to enable firewall: %v", err))
			pterm.Println()
			pterm.Error.Printfln("Failed to enable firewall: %v", err)
			os.Exit(1)
		}
	} else {
		// Allow SSH connections
		err := exec.Command("sudo", "ufw", "allow", "ssh").Run()
		if err != nil {
			logging.AppendRWZLogFile(currentUsername, logs.RWZLogFilePath, fmt.Sprintf("Failed to allow SSH: %v", err))
			pterm.Println()
			pterm.Error.Printfln("Failed to allow SSH: %v", err)
			os.Exit(1)
		}
		// Allow HTTP and HTTPS traffic
		err = exec.Command("sudo", "ufw", "allow", "Nginx Full").Run()
		if err != nil {
			logging.AppendRWZLogFile(currentUsername, logs.RWZLogFilePath, fmt.Sprintf("Failed to allow Nginx Full: %v", err))
			pterm.Println()
			pterm.Error.Printfln("Failed to allow Nginx Full: %v", err)
			os.Exit(1)
		}

		// Disable logging
		err = exec.Command("sudo", "ufw", "logging", "off").Run()
		if err != nil {
			logging.AppendRWZLogFile(currentUsername, logs.RWZLogFilePath, fmt.Sprintf("Failed to disable logging: %v", err))
			pterm.Println()
			pterm.Error.Printfln("Failed to disable logging: %v", err)
			os.Exit(1)
		}

		// Enable the firewall to apply the changes
		err = exec.Command("sudo", "ufw", "--force", "enable").Run()
		if err != nil {
			logging.AppendRWZLogFile(currentUsername, logs.RWZLogFilePath, fmt.Sprintf("Failed to enable firewall: %v", err))
			pterm.Println()
			pterm.Error.Printfln("Failed to enable firewall: %v", err)
			os.Exit(1)
		}
	}

	spinner.Success("Firewall configured")
}
