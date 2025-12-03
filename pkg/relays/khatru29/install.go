package khatru29

import (
	"fmt"

	"github.com/nodetec/rwz/pkg/network"
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/relays/utils/databases"
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/nodetec/rwz/pkg/utils/systemd"
	"github.com/nodetec/rwz/pkg/verification"
	"github.com/pterm/pterm"
)

// Install the relay
func Install(currentUsername, relayDomain, privKey, relayContact, relayUser string) {
	// TODO
	// Check if db writes should be allowed to finish before disabling and stopping the service
	// Re-enable the service if it exists and the user says no to overwriting the existing database

	// Check if the service file exists and disable and stop the service if it does
	systemd.DisableAndStopService(currentUsername, ServiceFilePath, ServiceName)

	// Determine how to handle existing database during install
	var howToHandleExistingDatabase = databases.HandleExistingDatabase(currentUsername, relayUser, DatabaseBackupsDirPath, DatabaseFilePath, BackupFileNameBase, relays.Khatru29RelayName)

	// Configure Nginx for HTTP
	network.ConfigureNginxHttp(currentUsername, relayDomain, relays.Khatru29NginxConfigFilePath)

	// Get SSL/TLS certificates
	httpsEnabled := network.GetCertificates(currentUsername, relayDomain, relays.Khatru29NginxConfigFilePath)
	if httpsEnabled {
		// Configure Nginx for HTTPS
		network.ConfigureNginxHttps(currentUsername, relayDomain, relays.Khatru29NginxConfigFilePath)
	}

	// Determine the temporary file path
	tmpCompressedBinaryFilePath := files.FilePathFromFilePathBase(DownloadURL, relays.TmpDirPath)

	// Check if the temporary file exists and remove it if it does
	if currentUsername == relays.RootUser {
		files.RemoveFile(tmpCompressedBinaryFilePath)
	} else {
		files.RemoveFileUsingLinux(currentUsername, tmpCompressedBinaryFilePath)
	}

	// Download and copy the file
	downloadSpinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Downloading %s binary...", relays.Khatru29RelayName))
	files.DownloadAndCopyFile(currentUsername, tmpCompressedBinaryFilePath, DownloadURL, 0644)
	downloadSpinner.Success(fmt.Sprintf("%s binary downloaded", relays.Khatru29RelayName))

	// Verify relay binary
	verification.VerifyRelayBinary(currentUsername, relays.Khatru29RelayName, tmpCompressedBinaryFilePath)

	// Install the compressed relay binary and make it executable
	installSpinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Installing %s binary...", relays.Khatru29RelayName))
	files.InstallCompressedBinary(currentUsername, tmpCompressedBinaryFilePath, relays.BinaryDestDir, BinaryName, "0755", relays.BinaryFilePerms)
	installSpinner.Success(fmt.Sprintf("%s binary installed", relays.Khatru29RelayName))

	// Set up the relay data directory
	databases.SetUpRelayDataDir(currentUsername, relayUser, howToHandleExistingDatabase, DataDirPath, DatabaseFilePath, relays.Khatru29RelayName)

	// Configure the relay
	ConfigureRelay(currentUsername, relayDomain, privKey, relayContact)

	SetUpRelaySite(currentUsername, relayDomain)

	// Set permissions for database files
	databases.SetDatabaseFilePermissions(currentUsername, DataDirPath, DatabaseFilePath, relays.Khatru29RelayName)

	// Use chown command to set ownership of the data directory to the provided relay user
	directories.SetOwnerAndGroupForAllContentUsingLinux(currentUsername, relayUser, relayUser, DataDirPath)

	// TODO
	// Add check for database compatibility for the creating a backup case using the database backup, may have to edit the khatru29 env file to use the database backup to check if the version is compatible with the installed khatru29 binary, and then use the installed khatru29 binary to create potential specfic exports if compatibile

	// Set up the relay service
	SetUpRelayService(currentUsername, relayUser)

	// Show success messages
	SuccessMessages(relayDomain, httpsEnabled)
}
