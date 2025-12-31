package khatru_pyramid

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

	// Path to the add-to-whitelist directory
	addToWhitelistDirPath := fmt.Sprintf("%s/%s", wwwDomainDirPath, AddToWhitelistDir)

	// Path to the remove-from-whitelist directory
	removeFromWhitelistDirPath := fmt.Sprintf("%s/%s", wwwDomainDirPath, RemoveFromWhitelistDir)

	// Path to the browse directory
	browseDirPath := fmt.Sprintf("%s/%s", wwwDomainDirPath, BrowseDir)

	// Path to the reports directory
	reportsDirPath := fmt.Sprintf("%s/%s", wwwDomainDirPath, ReportsDir)

	if currentUsername == relays.RootUser {
		// Remove the add-to-whitelist directory and all of its content if it exists
		spinner.UpdateText("Removing add-to-whitelist directory...")
		directories.RemoveDirectory(addToWhitelistDirPath)

		// Create the add-to-whitelist directory
		directories.CreateAllDirectories(addToWhitelistDirPath, 0755)

		// Set permissions for the add-to-whitelist directory
		directories.SetPermissions(addToWhitelistDirPath, 0755)

		// Use chown command to set ownership of the add-to-whitelist directory and its content to the www-data user
		directories.SetOwnerAndGroupForAllContentUsingLinux(currentUsername, relays.NginxUser, relays.NginxUser, addToWhitelistDirPath)

		// Remove the remove-from-whitelist directory and all of its content if it exists
		spinner.UpdateText("Removing remove-from-whitelist directory...")
		directories.RemoveDirectory(removeFromWhitelistDirPath)

		// Create the remove-from-whitelist directory
		directories.CreateAllDirectories(removeFromWhitelistDirPath, 0755)

		// Set permissions for the remove-from-whitelist directory
		directories.SetPermissions(removeFromWhitelistDirPath, 0755)

		// Use chown command to set ownership of the remove-from-whitelist directory and its content to the www-data user
		directories.SetOwnerAndGroupForAllContentUsingLinux(currentUsername, relays.NginxUser, relays.NginxUser, removeFromWhitelistDirPath)

		// Remove the browse directory and all of its content if it exists
		spinner.UpdateText("Removing browse directory...")
		directories.RemoveDirectory(browseDirPath)

		// Create the browse directory
		directories.CreateAllDirectories(browseDirPath, 0755)

		// Set permissions for the browse directory
		directories.SetPermissions(browseDirPath, 0755)

		// Use chown command to set ownership of the browse directory and its content to the www-data user
		directories.SetOwnerAndGroupForAllContentUsingLinux(currentUsername, relays.NginxUser, relays.NginxUser, browseDirPath)

		// Remove the reports directory and all of its content if it exists
		spinner.UpdateText("Removing reports directory...")
		directories.RemoveDirectory(reportsDirPath)

		// Create the reports directory
		directories.CreateAllDirectories(reportsDirPath, 0755)

		// Set permissions for the reports directory
		directories.SetPermissions(reportsDirPath, 0755)

		// Use chown command to set ownership of the reports directory and its content to the www-data user
		directories.SetOwnerAndGroupForAllContentUsingLinux(currentUsername, relays.NginxUser, relays.NginxUser, reportsDirPath)
	} else {
		// Remove the add-to-whitelist directory and all of its content if it exists
		spinner.UpdateText("Removing add-to-whitelist directory...")
		directories.RemoveDirectoryUsingLinux(currentUsername, addToWhitelistDirPath)

		// Create the add-to-whitelist directory
		directories.CreateAllDirectoriesUsingLinux(currentUsername, addToWhitelistDirPath)

		// Set permissions for the add-to-whitelist directory
		directories.SetPermissionsUsingLinux(currentUsername, addToWhitelistDirPath, "0755")

		// Use chown command to set ownership of the add-to-whitelist directory and its content to the www-data user
		directories.SetOwnerAndGroupForAllContentUsingLinux(currentUsername, relays.NginxUser, relays.NginxUser, addToWhitelistDirPath)

		// Remove the remove-from-whitelist directory and all of its content if it exists
		spinner.UpdateText("Removing remove-from-whitelist directory...")
		directories.RemoveDirectoryUsingLinux(currentUsername, removeFromWhitelistDirPath)

		// Create the remove-from-whitelist directory
		directories.CreateAllDirectoriesUsingLinux(currentUsername, removeFromWhitelistDirPath)

		// Set permissions for the remove-from-whitelist directory
		directories.SetPermissionsUsingLinux(currentUsername, removeFromWhitelistDirPath, "0755")

		// Use chown command to set ownership of the remove-from-whitelist directory and its content to the www-data user
		directories.SetOwnerAndGroupForAllContentUsingLinux(currentUsername, relays.NginxUser, relays.NginxUser, removeFromWhitelistDirPath)

		// Remove the browse directory and all of its content if it exists
		spinner.UpdateText("Removing browse directory...")
		directories.RemoveDirectoryUsingLinux(currentUsername, browseDirPath)

		// Create the browse directory
		directories.CreateAllDirectoriesUsingLinux(currentUsername, browseDirPath)

		// Set permissions for the browse directory
		directories.SetPermissionsUsingLinux(currentUsername, browseDirPath, "0755")

		// Use chown command to set ownership of the browse directory and its content to the www-data user
		directories.SetOwnerAndGroupForAllContentUsingLinux(currentUsername, relays.NginxUser, relays.NginxUser, browseDirPath)

		// Remove the reports directory and all of its content if it exists
		spinner.UpdateText("Removing reports directory...")
		directories.RemoveDirectoryUsingLinux(currentUsername, reportsDirPath)

		// Create the reports directory
		directories.CreateAllDirectoriesUsingLinux(currentUsername, reportsDirPath)

		// Set permissions for the reports directory
		directories.SetPermissionsUsingLinux(currentUsername, reportsDirPath, "0755")

		// Use chown command to set ownership of the reports directory and its content to the www-data user
		directories.SetOwnerAndGroupForAllContentUsingLinux(currentUsername, relays.NginxUser, relays.NginxUser, reportsDirPath)
	}

	spinner.Success("Relay site set up")
}
