package manager

import (
	"fmt"
	"github.com/nodetec/rwz/pkg/relays"
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
				pterm.Error.Printfln("Failed to check if package is installed: %v", err)
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
func AptInstallPackages(selectedRelayOption string) {
	spinner, _ := pterm.DefaultSpinner.Start("Updating and installing packages...")

	exec.Command("apt", "update", "-qq").Run()

	packages := []string{"curl", "gnupg", "openssh-server", "ufw", "fail2ban", "nginx", "certbot", "python3-certbot-nginx"}

	if selectedRelayOption == relays.NostrRsRelayName || selectedRelayOption == relays.StrfryRelayName || selectedRelayOption == relays.WotRelayName || selectedRelayOption == relays.Strfry29RelayName {
		packages = append(packages, "git")
	}

	if selectedRelayOption == relays.NostrRsRelayName {
		packages = append(packages, "sqlite3", "libsqlite3-dev")
	}

	if selectedRelayOption == relays.StrfryRelayName || selectedRelayOption == relays.Strfry29RelayName {
		packages = append(packages, "libssl-dev", "zlib1g-dev", "liblmdb-dev", "libflatbuffers-dev", "libsecp256k1-dev", "libzstd-dev")
	}

	// Check if package is installed, install if not
	for _, p := range packages {
		if IsPackageInstalled(p) {
			spinner.UpdateText(fmt.Sprintf("%s is already installed.", p))
		} else {
			spinner.UpdateText(fmt.Sprintf("Installing %s...", p))
			exec.Command("apt", "install", "-y", "-qq", p).Run()
		}
	}

	spinner.Success("Packages updated and installed")
}
