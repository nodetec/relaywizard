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
	WWWDomainDirPath := fmt.Sprintf("%s/%s", network.WWWDirPath, domain)

	// Path to the index.html file
	IndexFilePath := fmt.Sprintf("%s/%s", WWWDomainDirPath, IndexFile)

	if currentUsername == relays.RootUser {
		// Check if the index.html file exists and remove it if it does
		files.RemoveFile(IndexFilePath)

		// Copy the index.html file to the /var/www/domain directory
		files.CopyFileUsingLinux(currentUsername, TmpIndexFilePath, WWWDomainDirPath)

		// Set permissions for the index.html file
		files.SetPermissions(IndexFilePath, 0644)

		// Use chown command to set ownership of the index.html file to the www-data user
		files.SetOwnerAndGroup(relays.NginxUser, relays.NginxUser, IndexFilePath)
	} else {
		// Check if the index.html file exists and remove it if it does
		files.RemoveFileUsingLinux(currentUsername, IndexFilePath)

		// Copy the index.html file to the /var/www/domain directory
		files.CopyFileUsingLinux(currentUsername, TmpIndexFilePath, WWWDomainDirPath)

		// Set permissions for the index.html file
		files.SetPermissionsUsingLinux(currentUsername, IndexFilePath, "0644")

		// Use chown command to set ownership of the index.html file to the www-data user
		files.SetOwnerAndGroupUsingLinux(currentUsername, relays.NginxUser, relays.NginxUser, IndexFilePath)
	}

	// Path to the static directory
	StaticDirPath := fmt.Sprintf("%s/%s", WWWDomainDirPath, StaticDir)

	if currentUsername == relays.RootUser {
		// Remove the static directory and all of its content if it exists
		spinner.UpdateText("Removing static directory...")
		directories.RemoveDirectory(StaticDirPath)

		// Copy the static directory and all of its content to the /var/www/domain directory
		directories.CopyDirectoryUsingLinux(currentUsername, TmpStaticDirPath, WWWDomainDirPath)

		// Set permissions for the static directory
		directories.SetPermissions(StaticDirPath, 0755)

		// Use chown command to set ownership of the static directory and its content to the www-data user
		directories.SetOwnerAndGroupForAllContentUsingLinux(currentUsername, relays.NginxUser, relays.NginxUser, StaticDirPath)
	} else {
		// Remove the static directory and all of its content if it exists
		spinner.UpdateText("Removing static directory...")
		directories.RemoveDirectoryUsingLinux(currentUsername, StaticDirPath)

		// Copy the static directory and all of its content to the /var/www/domain directory
		directories.CopyDirectoryUsingLinux(currentUsername, TmpStaticDirPath, WWWDomainDirPath)

		// Set permissions for the static directory
		directories.SetPermissionsUsingLinux(currentUsername, StaticDirPath, "0755")

		// Use chown command to set ownership of the static directory and its content to the www-data user
		directories.SetOwnerAndGroupForAllContentUsingLinux(currentUsername, relays.NginxUser, relays.NginxUser, StaticDirPath)
	}

	spinner.Success("Relay site set up")
}
