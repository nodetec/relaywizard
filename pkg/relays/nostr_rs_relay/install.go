package nostr_rs_relay

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
func Install(relayDomain, pubKey, relayContact, relayUser string) {
	// TODO
	// Check if you should wait for any db writes to finish before disabling and stopping the service
	// Re-enable the service if it exists and the user says no to overwriting the existing database

	// Check if the service file exists and disable and stop the service if it does
	systemd.DisableAndStopService(ServiceFilePath, ServiceName)

	// Determine how to handle existing database during install
	var howToHandleExistingDatabase = databases.HandleExistingDatabase(DatabaseBackupsDirPath, DatabaseFilePath, BackupFileNameBase, relays.NostrRsRelayName)

	// Configure Nginx for HTTP
	network.ConfigureNginxHttp(relayDomain, relays.NostrRsRelayNginxConfigFilePath)

	// Get SSL/TLS certificates
	httpsEnabled := network.GetCertificates(relayDomain, relays.NostrRsRelayNginxConfigFilePath)
	if httpsEnabled {
		// Configure Nginx for HTTPS
		network.ConfigureNginxHttps(relayDomain, relays.NostrRsRelayNginxConfigFilePath)
	}

	// Download the config file from the git repository
	git.RemoveThenClone(GitRepoTmpDirPath, GitRepoBranch, GitRepoURL, relays.GitRepoDirPerms)

	// Determine the temporary file path
	tmpCompressedBinaryFilePath := files.FilePathFromFilePathBase(DownloadURL, relays.TmpDirPath)

	// Check if the temporary file exists and remove it if it does
	files.RemoveFile(tmpCompressedBinaryFilePath)

	// Download and copy the file
	downloadSpinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Downloading %s binary...", relays.NostrRsRelayName))
	files.DownloadAndCopyFile(tmpCompressedBinaryFilePath, DownloadURL, 0644)
	downloadSpinner.Success(fmt.Sprintf("%s binary downloaded", relays.NostrRsRelayName))

	// Verify relay binary
	verification.VerifyRelayBinary(relays.NostrRsRelayName, tmpCompressedBinaryFilePath)

	// Install the compressed relay binary and make it executable
	installSpinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Installing %s binary...", relays.NostrRsRelayName))
	files.InstallCompressedBinary(tmpCompressedBinaryFilePath, relays.BinaryDestDir, BinaryName, relays.BinaryFilePerms)
	installSpinner.Success(fmt.Sprintf("%s binary installed", relays.NostrRsRelayName))

	// Set up the relay data directory
	databases.SetUpRelayDataDir(howToHandleExistingDatabase, DataDirPath, DatabaseFilePath, relays.NostrRsRelayName)

	// Configure the relay
	ConfigureRelay(relayDomain, pubKey, relayContact, httpsEnabled)

	// Set permissions for database files
	databases.SetDatabaseFilePermissions(DataDirPath, DatabaseFilePath, relays.NostrRsRelayName)

	// Use chown command to set ownership of the data directory to the provided relay user
	directories.SetOwnerAndGroup(relayUser, relayUser, DataDirPath)

	// Set up the relay service
	SetUpRelayService(relayUser)

	// Show success messages
	SuccessMessages(relayDomain, httpsEnabled)
}
