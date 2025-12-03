package wot_relay

import (
	"fmt"

	"github.com/nodetec/rwz/pkg/network"
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/pterm/pterm"
)

// Function to set up the relay site
func SetUpRelaySite(currentUsername, domain string) {
	spinner, _ := pterm.DefaultSpinner.Start("Configuring relay site...")

	// Path to the /var/www/domain directory
	wwwDomainDirPath := fmt.Sprintf("%s/%s", network.WWWDirPath, domain)

	// Path to the index.html file
	indexFilePath := fmt.Sprintf("%s/%s", wwwDomainDirPath, IndexFile)

	if currentUsername == relays.RootUser {
		// Check if the index.html file exists and remove it if it does
		files.RemoveFile(indexFilePath)

		// Copy the index.html file to the /var/www/domain directory
		files.CopyFileUsingLinux(currentUsername, TmpIndexFilePath, wwwDomainDirPath)

		// Set permissions for the index.html file
		files.SetPermissions(indexFilePath, 0644)

		// Use chown command to set ownership of the index.html file to the www-data user
		files.SetOwnerAndGroup(relays.NginxUser, relays.NginxUser, indexFilePath)
	} else {
		// Check if the index.html file exists and remove it if it does
		files.RemoveFileUsingLinux(currentUsername, indexFilePath)

		// Copy the index.html file to the /var/www/domain directory
		files.CopyFileUsingLinux(currentUsername, TmpIndexFilePath, wwwDomainDirPath)

		// Set permissions for the index.html file
		files.SetPermissionsUsingLinux(currentUsername, indexFilePath, "0644")

		// Use chown command to set ownership of the index.html file to the www-data user
		files.SetOwnerAndGroupUsingLinux(currentUsername, relays.NginxUser, relays.NginxUser, indexFilePath)
	}

	// Path to the static directory
	staticDirPath := fmt.Sprintf("%s/%s", wwwDomainDirPath, StaticDir)

	if currentUsername == relays.RootUser {
		// Remove the static directory and all of its content if it exists
		spinner.UpdateText("Removing static directory...")
		directories.RemoveDirectory(staticDirPath)

		// Copy the static directory and all of its content to the /var/www/domain directory
		directories.CopyDirectoryUsingLinux(currentUsername, TmpStaticDirPath, wwwDomainDirPath)

		// Set permissions for the static directory
		directories.SetPermissions(staticDirPath, 0755)

		// Use chown command to set ownership of the static directory and its content to the www-data user
		directories.SetOwnerAndGroupForAllContentUsingLinux(currentUsername, relays.NginxUser, relays.NginxUser, staticDirPath)
	} else {
		// Remove the static directory and all of its content if it exists
		spinner.UpdateText("Removing static directory...")
		directories.RemoveDirectoryUsingLinux(currentUsername, staticDirPath)

		// Copy the static directory and all of its content to the /var/www/domain directory
		directories.CopyDirectoryUsingLinux(currentUsername, TmpStaticDirPath, wwwDomainDirPath)

		// Set permissions for the static directory
		directories.SetPermissionsUsingLinux(currentUsername, staticDirPath, "0755")

		// Use chown command to set ownership of the static directory and its content to the www-data user
		directories.SetOwnerAndGroupForAllContentUsingLinux(currentUsername, relays.NginxUser, relays.NginxUser, staticDirPath)
	}

	spinner.Success("Relay site set up")
}
