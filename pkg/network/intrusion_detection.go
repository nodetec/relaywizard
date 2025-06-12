package network

import (
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/nodetec/rwz/pkg/utils/network"
	"github.com/pterm/pterm"
	"os"
	"os/exec"
)

// Function to configure the intrusion detection system
func ConfigureIntrusionDetection() {
	spinner, _ := pterm.DefaultSpinner.Start("Configuring intrusion detection system...")

	// Check if the SSH jail file exists and remove it if it does
	files.RemoveFile(SSHJailFilePath)

	// Create the SSH jail file
	spinner.UpdateText("Creating SSH jail file...")
	network.CreateJailFile(SSHJailFilePath, SSHJailFileTemplate)

	// Restart the intrusion detection system to apply the changes
	err := exec.Command("systemctl", "restart", "fail2ban").Run()
	if err != nil {
		pterm.Println()
		pterm.Error.Printfln("Failed to restart intrusion detection system: %v", err)
		os.Exit(1)
	}

	spinner.Success("Intrusion detection system configured")
}
