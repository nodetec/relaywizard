package strfry

import (
	"fmt"
	"github.com/nodetec/rwz/pkg/network"
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/relays/utils/databases"
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/nodetec/rwz/pkg/utils/git"
	"github.com/nodetec/rwz/pkg/utils/systemd"
	"github.com/nodetec/rwz/pkg/verification"
	"github.com/pterm/pterm"
)

// Install the relay
func Install(currentUsername, relayDomain, pubKey, relayContact, relayUser string) {
	// TODO
	// Check if you should wait for any db writes to finish before disabling and stopping the service
	// Re-enable service if it exists and the user says no to overwriting the existing database

	// Check if the service file exists and disable and stop the service if it does
	systemd.DisableAndStopService(currentUsername, ServiceFilePath, ServiceName)

	// Determine how to handle existing database during install
	var howToHandleExistingDatabase = databases.HandleExistingDatabase(currentUsername, relayUser, DatabaseBackupsDirPath, DatabaseFilePath, BackupFileNameBase, relays.StrfryRelayName)

	// Configure Nginx for HTTP
	network.ConfigureNginxHttp(currentUsername, relayDomain, relays.StrfryNginxConfigFilePath)

	// Get SSL/TLS certificates
	httpsEnabled := network.GetCertificates(currentUsername, relayDomain, relays.StrfryNginxConfigFilePath)
	if httpsEnabled {
		// Configure Nginx for HTTPS
		network.ConfigureNginxHttps(currentUsername, relayDomain, relays.StrfryNginxConfigFilePath)
	}

	// Download the config file from the git repository
	git.RemoveThenClone(currentUsername, GitRepoTmpDirPath, GitRepoBranch, GitRepoURL, relays.GitRepoDirPerms)

	// Determine the temporary file path
	tmpCompressedBinaryFilePath := files.FilePathFromFilePathBase(DownloadURL, relays.TmpDirPath)

	// Check if the temporary file exists and remove it if it does
	if currentUsername == relays.RootUser {
		files.RemoveFile(tmpCompressedBinaryFilePath)
	} else {
		files.RemoveFileUsingLinux(currentUsername, tmpCompressedBinaryFilePath)
	}

	// Download and copy the file
	downloadSpinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Downloading %s binary...", relays.StrfryRelayName))
	files.DownloadAndCopyFile(currentUsername, tmpCompressedBinaryFilePath, DownloadURL, 0644)
	downloadSpinner.Success(fmt.Sprintf("%s binary downloaded", relays.StrfryRelayName))

	// Verify relay binary
	verification.VerifyRelayBinary(currentUsername, relays.StrfryRelayName, tmpCompressedBinaryFilePath)

	// Install the compressed relay binary and make it executable
	installSpinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Installing %s binary...", relays.StrfryRelayName))
	files.InstallCompressedBinary(currentUsername, tmpCompressedBinaryFilePath, relays.BinaryDestDir, BinaryName, relays.BinaryFilePerms)
	installSpinner.Success(fmt.Sprintf("%s binary installed", relays.StrfryRelayName))

	// Set up the relay data directory
	databases.SetUpRelayDataDir(currentUsername, relayUser, howToHandleExistingDatabase, DataDirPath, DatabaseFilePath, relays.StrfryRelayName)

	// Configure the relay
	ConfigureRelay(currentUsername, pubKey, relayContact)

	// Set permissions for database files
	databases.SetDatabaseFilePermissions(currentUsername, DataDirPath, DatabaseFilePath, relays.StrfryRelayName)

	// Use chown command to set ownership of the data directory to the provided relay user
	if currentUsername == relays.RootUser {
		directories.SetOwnerAndGroup(relayUser, relayUser, DataDirPath)
	} else {
		directories.SetOwnerAndGroupUsingLinux(currentUsername, relayUser, relayUser, DataDirPath)
	}

	// Set up the relay service
	SetUpRelayService(currentUsername, relayUser)

	// TODO
	// Add check for database compatibility for the creating a backup case using the database backup, may have to edit the strfry config file to use the database backup to check if the version is compatible with the installed strfry binary, and then use the installed strfry binary to create a fried export if compatibile

	// Check if installed strfry binary and existing database version are compatible
	if howToHandleExistingDatabase == databases.UseExistingDatabaseFileOption {
		databases.CheckStrfryBinaryAndDatabaseCompatibility(currentUsername, BinaryName, ConfigFilePath, SupportedDatabaseVersionOutput, BinaryVersion, SupportedDatabaseVersion)
	}

	// Show success messages
	SuccessMessages(relayDomain, httpsEnabled)
}
