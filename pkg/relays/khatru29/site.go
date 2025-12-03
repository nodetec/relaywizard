package khatru29

import (
	"fmt"

	"github.com/nodetec/rwz/pkg/network"
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/pterm/pterm"
)

// Function to set up the relay site
func SetUpRelaySite(currentUsername, domain string) {
	spinner, _ := pterm.DefaultSpinner.Start("Configuring relay site...")

	// Path to the /var/www/domain directory
	wwwDomainDirPath := fmt.Sprintf("%s/%s", network.WWWDirPath, domain)

	// Path to the create directory
	createDirPath := fmt.Sprintf("%s/%s", wwwDomainDirPath, CreateDir)

	if currentUsername == relays.RootUser {
		// Remove the create directory and all of its content if it exists
		spinner.UpdateText("Removing create directory...")
		directories.RemoveDirectory(createDirPath)

		// Create the create directory
		directories.CreateAllDirectories(createDirPath, 0755)

		// Set permissions for the create directory
		directories.SetPermissions(createDirPath, 0755)

		// Use chown command to set ownership of the create directory and its content to the www-data user
		directories.SetOwnerAndGroupForAllContentUsingLinux(currentUsername, relays.NginxUser, relays.NginxUser, createDirPath)
	} else {
		// Remove the create directory and all of its content if it exists
		spinner.UpdateText("Removing create directory...")
		directories.RemoveDirectoryUsingLinux(currentUsername, createDirPath)

		// Create the create directory
		directories.CreateAllDirectoriesUsingLinux(currentUsername, createDirPath)

		// Set permissions for the create directory
		directories.SetPermissionsUsingLinux(currentUsername, createDirPath, "0755")

		// Use chown command to set ownership of the create directory and its content to the www-data user
		directories.SetOwnerAndGroupForAllContentUsingLinux(currentUsername, relays.NginxUser, relays.NginxUser, createDirPath)
	}

	spinner.Success("Relay site set up")
}
