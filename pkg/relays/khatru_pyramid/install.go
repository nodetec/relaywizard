package khatru_pyramid

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
func Install(relayDomain, pubKey, relayContact, relayUser string) {
	// TODO
	// Check if db writes should be allowed to finish before disabling and stopping the service
	// Re-enable the service if it exists and the user says no to overwriting the existing database

	// Check if the service file exists and disable and stop the service if it does
	systemd.DisableAndStopService(ServiceFilePath, ServiceName)

	// Determine how to handle existing database during install
	var howToHandleExistingDatabase = databases.HandleExistingDatabase(DatabaseBackupsDirPath, DatabaseFilePath, BackupFileNameBase, relays.KhatruPyramidRelayName)

	// Determine how to handle existing users file during install
	HandleExistingUsersFile(pubKey, relayUser)

	// Configure Nginx for HTTP
	network.ConfigureNginxHttp(relayDomain, relays.KhatruPyramidNginxConfigFilePath)

	// Get SSL/TLS certificates
	httpsEnabled := network.GetCertificates(relayDomain, relays.KhatruPyramidNginxConfigFilePath)
	if httpsEnabled {
		// Configure Nginx for HTTPS
		network.ConfigureNginxHttps(relayDomain, relays.KhatruPyramidNginxConfigFilePath)
	}

	// Determine the temporary file path
	tmpCompressedBinaryFilePath := files.FilePathFromFilePathBase(DownloadURL, relays.TmpDirPath)

	// Check if the temporary file exists and remove it if it does
	files.RemoveFile(tmpCompressedBinaryFilePath)

	// Download and copy the file
	downloadSpinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Downloading %s binary...", relays.KhatruPyramidRelayName))
	files.DownloadAndCopyFile(tmpCompressedBinaryFilePath, DownloadURL, 0644)
	downloadSpinner.Success(fmt.Sprintf("%s binary downloaded", relays.KhatruPyramidRelayName))

	// Verify relay binary
	verification.VerifyRelayBinary(relays.KhatruPyramidRelayName, tmpCompressedBinaryFilePath)

	// Install the compressed relay binary and make it executable
	installSpinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Installing %s binary...", relays.KhatruPyramidRelayName))
	files.InstallCompressedBinary(tmpCompressedBinaryFilePath, relays.BinaryDestDir, BinaryName, relays.BinaryFilePerms)
	installSpinner.Success(fmt.Sprintf("%s binary installed", relays.KhatruPyramidRelayName))

	// Set up the relay data directory
	databases.SetUpRelayDataDir(howToHandleExistingDatabase, DataDirPath, DatabaseFilePath, relays.KhatruPyramidRelayName)

	// Configure the relay
	ConfigureRelay(relayDomain, pubKey, relayContact)

	// Set permissions for database files
	databases.SetDatabaseFilePermissions(DataDirPath, DatabaseFilePath, relays.KhatruPyramidRelayName)

	// Use chown command to set ownership of the data directory to the provided relay user
	directories.SetOwnerAndGroup(relayUser, relayUser, DataDirPath)

	// TODO
	// Add check for database compatibility for the creating a backup case using the database backup, may have to edit the khatru pyramid env file to use the database backup to check if the version is compatible with the installed khatru pyramid binary, and then use the installed khatru pyramid binary to create potential specfic exports if compatibile

	// Set up the relay service
	SetUpRelayService(relayUser)

	// Show success messages
	SuccessMessages(relayDomain, httpsEnabled)
}
