package manager

import (
	"fmt"
	"github.com/pterm/pterm"
	"os"
	"os/exec"
)

// Function to check if a package is installed
func IsPackageInstalled(packageName string) bool {
	out, err := exec.Command("dpkg-query", "-W", "-f='${Status}'", packageName).Output()

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			errorCode := exitError.ExitCode()
			// Package not found
			if errorCode == 1 {
				return false
			} else {
				pterm.Error.Println(fmt.Sprintf("Failed to check if package is installed: %v", err))
				os.Exit(1)
			}
		}
	}

	status := string(out)

	if status == "'unknown ok not-installed'" {
		return false
	} else if status == "'install ok installed'" {
		return true
	}

	return false
}

// Function to install necessary packages
func AptInstallPackages() {
	spinner, _ := pterm.DefaultSpinner.Start("Updating and installing packages...")

	exec.Command("apt", "update", "-qq").Run()

	packages := []string{"nginx", "certbot", "python3-certbot-nginx", "ufw", "fail2ban"}

	// Check if package is installed, install if not
	for _, p := range packages {
		if IsPackageInstalled(p) {
			spinner.UpdateText(fmt.Sprintf("%s is already installed.", p))
		} else {
			spinner.UpdateText(fmt.Sprintf("Installing %s...", p))
			exec.Command("apt", "install", "-y", "-qq", p).Run()
		}
	}

	spinner.Success("Packages updated and installed successfully.")
}
