package manager

import (
	"os/exec"

	"github.com/pterm/pterm"
)

// Function to check if a command exists
func commandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

// Function to install necessary packages
func AptInstallPackages() {
	spinner, _ := pterm.DefaultSpinner.Start("Updating and installing packages...")
	exec.Command("apt", "update", "-qq").Run()

	// Check if nginx is installed, install if not
	if commandExists("nginx") {
		spinner.UpdateText("nginx is already installed.")
	} else {
		spinner.UpdateText("Installing nginx...")
		exec.Command("apt", "install", "-y", "-qq", "nginx").Run()
	}

	// Check if Certbot is installed, install if not
	if commandExists("certbot") {
		spinner.UpdateText("Certbot is already installed.")
	} else {
		spinner.UpdateText("Installing Certbot...")
		exec.Command("apt", "install", "-y", "-qq", "certbot", "python3-certbot-nginx").Run()
	}

	// Check if ufw is installed, install if not
	if commandExists("ufw") {
		spinner.UpdateText("ufw is already installed.")
	} else {
		spinner.UpdateText("Installing ufw...")
		exec.Command("apt", "install", "-y", "-qq", "ufw").Run()
	}

	spinner.Success("Packages updated and installed successfully.")
}

// Function to check if a package is installed
func isPackageInstalled(packageName string) bool {
	cmd := exec.Command("dpkg", "-l", packageName)
	err := cmd.Run()
	return err == nil
}
