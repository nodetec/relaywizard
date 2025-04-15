package strfry

import (
	"fmt"
	"github.com/nodetec/rwz/pkg/network"
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/relays/utils/databases"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/nodetec/rwz/pkg/utils/git"
	"github.com/nodetec/rwz/pkg/utils/systemd"
	"github.com/nodetec/rwz/pkg/verification"
	"github.com/pterm/pterm"
)

// Install the relay
func Install(relayDomain, pubKey, relayContact string) {
	// Determine how to handle existing database during install
	var howToHandleExistingDatabase = databases.HandleExistingDatabase(DatabaseBackupsDirPath, DatabaseFilePath, BackupFileNameBase, DatabaseLockFilePath)

	// Configure Nginx for HTTP
	network.ConfigureNginxHttp(relayDomain, NginxConfigFilePath)

	// Get SSL/TLS certificates
	httpsEnabled := network.GetCertificates(relayDomain, NginxConfigFilePath)
	if httpsEnabled {
		// Configure Nginx for HTTPS
		network.ConfigureNginxHttps(relayDomain, NginxConfigFilePath)
	}

	// Download the config file from the git repository
	git.RemoveThenClone(GitRepoTmpDirPath, GitRepoBranch, GitRepoURL, relays.GitRepoDirPerms)

	// Determine the temporary file path
	tmpCompressedBinaryFilePath := files.FilePathFromFilePathBase(DownloadURL, relays.TmpDirPath)

	// Check if the temporary file exists and remove it if it does
	files.RemoveFile(tmpCompressedBinaryFilePath)

	// Download and copy the file
	downloadSpinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Downloading %s binary...", RelayName))
	files.DownloadAndCopyFile(tmpCompressedBinaryFilePath, DownloadURL, 0644)
	downloadSpinner.Success(fmt.Sprintf("%s binary downloaded", RelayName))

	// Verify relay binary
	verification.VerifyRelayBinary(RelayName, tmpCompressedBinaryFilePath)

	// Check if the service file exists and disable and stop the service if it does
	systemd.DisableAndStopService(ServiceFilePath, ServiceName)

	// Install the compressed relay binary and make it executable
	installSpinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Installing %s binary...", RelayName))
	files.InstallCompressedBinary(tmpCompressedBinaryFilePath, relays.BinaryDestDir, BinaryName, relays.BinaryFilePerms)
	installSpinner.Success(fmt.Sprintf("%s binary installed", RelayName))

	// Set up the relay data directory
	databases.SetUpRelayDataDir(howToHandleExistingDatabase, DataDirPath, DatabaseFilePath, DatabaseLockFilePath)

	// Configure the relay
	ConfigureRelay(pubKey, relayContact)

	// Set up the relay service
	SetUpRelayService()

	// TODO
	// Add check for database compatibility for the creating a backup case using the database backup, may have to edit the strfry config file to use the database backup to check if the version is compatible with the installed strfry binary, and then use the installed strfry binary to create a fried export if compatibile

	// Check if installed strfry binary and existing database version are compatible
	if howToHandleExistingDatabase == databases.UseExistingDatabaseFileOption {
		databases.CheckStrfryBinaryAndDatabaseCompatibility(BinaryName, ConfigFilePath, SupportedDatabaseVersionOutput, BinaryVersion, SupportedDatabaseVersion)
	}

	// Show success messages
	SuccessMessages(relayDomain, httpsEnabled)
}
