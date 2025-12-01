package strfry29

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
func Install(currentUsername, relayDomain, pubKey, privKey, relayContact, relayUser string) {
	// TODO
	// Check if you should wait for any db writes to finish before disabling and stopping the service
	// Re-enable the service if it exists and the user says no to overwriting the existing database

	// Check if the service file exists and disable and stop the service if it does
	systemd.DisableAndStopService(currentUsername, ServiceFilePath, ServiceName)

	// Determine how to handle existing database during install
	var howToHandleExistingDatabase = databases.HandleExistingDatabase(currentUsername, relayUser, DatabaseBackupsDirPath, DatabaseFilePath, BackupFileNameBase, relays.Strfry29RelayName)

	// Configure Nginx for HTTP
	network.ConfigureNginxHttp(currentUsername, relayDomain, relays.Strfry29NginxConfigFilePath)

	// Get SSL/TLS certificates
	httpsEnabled := network.GetCertificates(currentUsername, relayDomain, relays.Strfry29NginxConfigFilePath)
	if httpsEnabled {
		// Configure Nginx for HTTPS
		network.ConfigureNginxHttps(currentUsername, relayDomain, relays.Strfry29NginxConfigFilePath)
	}

	// Download the config file from the git repository
	git.RemoveThenClone(currentUsername, GitRepoTmpDirPath, GitRepoBranch, GitRepoURL, "0755", relays.GitRepoDirPerms)

	// Determine the temporary file path
	tmpCompressedBinaryFilePath := files.FilePathFromFilePathBase(DownloadURL, relays.TmpDirPath)

	// Check if the temporary file exists and remove it if it does
	if currentUsername == relays.RootUser {
		files.RemoveFile(tmpCompressedBinaryFilePath)
	} else {
		files.RemoveFileUsingLinux(currentUsername, tmpCompressedBinaryFilePath)
	}

	// Download and copy the file
	downloadSpinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Downloading %s binary...", BinaryName))
	files.DownloadAndCopyFile(currentUsername, tmpCompressedBinaryFilePath, DownloadURL, 0644)
	downloadSpinner.Success(fmt.Sprintf("%s binary downloaded", BinaryName))

	// Verify relay binary
	verification.VerifyRelayBinary(currentUsername, BinaryName, tmpCompressedBinaryFilePath)

	// Install the compressed relay binary and make it executable
	installSpinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Installing %s binary...", BinaryName))
	files.InstallCompressedBinary(currentUsername, tmpCompressedBinaryFilePath, relays.BinaryDestDir, BinaryName, "0755", relays.BinaryFilePerms)
	installSpinner.Success(fmt.Sprintf("%s binary installed", BinaryName))

	// Determine the temporary file path
	tmpCompressedBinaryPluginFilePath := files.FilePathFromFilePathBase(BinaryPluginDownloadURL, relays.TmpDirPath)

	// Check if the temporary file exists and remove it if it does
	if currentUsername == relays.RootUser {
		files.RemoveFile(tmpCompressedBinaryPluginFilePath)
	} else {
		files.RemoveFileUsingLinux(currentUsername, tmpCompressedBinaryPluginFilePath)
	}

	// Download and copy the file
	binaryPluginDownloadSpinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Downloading %s plugin binary...", BinaryPluginName))
	files.DownloadAndCopyFile(currentUsername, tmpCompressedBinaryPluginFilePath, BinaryPluginDownloadURL, 0644)
	binaryPluginDownloadSpinner.Success(fmt.Sprintf("%s plugin binary downloaded", BinaryPluginName))

	// Verify relay binary plugin
	verification.VerifyRelayBinary(currentUsername, fmt.Sprintf("%s plugin", BinaryPluginName), tmpCompressedBinaryPluginFilePath)

	// Install the compressed relay binary plugin and make it executable
	binaryPluginInstallSpinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Installing %s plugin binary...", BinaryPluginName))
	files.InstallCompressedBinary(currentUsername, tmpCompressedBinaryPluginFilePath, relays.BinaryDestDir, BinaryPluginName, "0755", relays.BinaryFilePerms)
	binaryPluginInstallSpinner.Success(fmt.Sprintf("%s plugin binary installed", BinaryPluginName))

	// Set up the relay data directory
	databases.SetUpRelayDataDir(currentUsername, relayUser, howToHandleExistingDatabase, DataDirPath, DatabaseFilePath, relays.Strfry29RelayName)

	// Configure the relay
	ConfigureRelay(currentUsername, relayDomain, pubKey, privKey, relayContact)

	// Set permissions for database files
	databases.SetDatabaseFilePermissions(currentUsername, DataDirPath, DatabaseFilePath, relays.Strfry29RelayName)

	// Use chown command to set ownership of the data directory to the provided relay user
	directories.SetOwnerAndGroupForAllContentUsingLinux(currentUsername, relayUser, relayUser, DataDirPath)

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
