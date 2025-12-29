package network

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/nodetec/rwz/pkg/logs"
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/nodetec/rwz/pkg/utils/logging"
	"github.com/nodetec/rwz/pkg/utils/network"
	"github.com/pterm/pterm"
)

// Function to configure the intrusion detection system
func ConfigureIntrusionDetection(currentUsername string) {
	spinner, _ := pterm.DefaultSpinner.Start("Configuring intrusion detection system...")

	if currentUsername == relays.RootUser {
		// Check if the SSH jail file exists and remove it if it does
		files.RemoveFile(SSHJailFilePath)

		// Create the SSH jail file
		spinner.UpdateText("Creating SSH jail file...")
		network.CreateJailFile(currentUsername, SSHJailFilePath, SSHJailFileTemplate)

		// Restart the intrusion detection system to apply the changes
		err := exec.Command("systemctl", "restart", "fail2ban").Run()
		if err != nil {
			logging.AppendRWZLogFile(currentUsername, logs.RWZLogFilePath, fmt.Sprintf("Failed to restart intrusion detection system: %v", err))
			pterm.Println()
			pterm.Error.Printfln("Failed to restart intrusion detection system: %v", err)
			os.Exit(1)
		}
	} else {
		// Check if the SSH jail file exists and remove it if it does
		files.RemoveFileUsingLinux(currentUsername, SSHJailFilePath)

		// Create the SSH jail file
		spinner.UpdateText("Creating SSH jail file...")
		network.CreateJailFile(currentUsername, SSHJailFilePath, SSHJailFileTemplate)

		// Restart the intrusion detection system to apply the changes
		err := exec.Command("sudo", "systemctl", "restart", "fail2ban").Run()
		if err != nil {
			logging.AppendRWZLogFile(currentUsername, logs.RWZLogFilePath, fmt.Sprintf("Failed to restart intrusion detection system: %v", err))
			pterm.Println()
			pterm.Error.Printfln("Failed to restart intrusion detection system: %v", err)
			os.Exit(1)
		}
	}

	spinner.Success("Intrusion detection system configured")
}
