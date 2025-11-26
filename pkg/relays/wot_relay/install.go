package wot_relay

import (
	"fmt"
	"github.com/nodetec/rwz/pkg/network"
	"github.com/nodetec/rwz/pkg/relays"
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

	// Check if the service file exists and disable and stop the service if it does
	systemd.DisableAndStopService(currentUsername, ServiceFilePath, ServiceName)

	// Configure Nginx for HTTP
	network.ConfigureNginxHttp(currentUsername, relayDomain, relays.WotRelayNginxConfigFilePath)

	// Get SSL/TLS certificates
	httpsEnabled := network.GetCertificates(currentUsername, relayDomain, relays.WotRelayNginxConfigFilePath)
	if httpsEnabled {
		// Configure Nginx for HTTPS
		network.ConfigureNginxHttps(currentUsername, relayDomain, relays.WotRelayNginxConfigFilePath)
	}

	// Download the templates directory from the git repository
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
	downloadSpinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Downloading %s binary...", relays.WotRelayName))
	files.DownloadAndCopyFile(tmpCompressedBinaryFilePath, DownloadURL, 0666)
	downloadSpinner.Success(fmt.Sprintf("%s binary downloaded", relays.WotRelayName))

	// Verify relay binary
	verification.VerifyRelayBinary(currentUsername, relays.WotRelayName, tmpCompressedBinaryFilePath)

	// Install the compressed relay binary and make it executable
	installSpinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Installing %s binary...", relays.WotRelayName))
	files.InstallCompressedBinary(currentUsername, tmpCompressedBinaryFilePath, relays.BinaryDestDir, BinaryName, relays.BinaryFilePerms)
	installSpinner.Success(fmt.Sprintf("%s binary installed", relays.WotRelayName))

	// Set up the relay data directory
	SetUpRelayDataDir(currentUsername)

	// Configure the relay
	ConfigureRelay(currentUsername, relayDomain, pubKey, relayContact, httpsEnabled)

	// Set up the relay site
	SetUpRelaySite(currentUsername, relayDomain)

	// Set up the relay service
	SetUpRelayService(currentUsername, relayUser)

	// Use chown command to set ownership of the data directory to the provided relay user
	if currentUsername == relays.RootUser {
		directories.SetOwnerAndGroup(relayUser, relayUser, DataDirPath)
	} else {
		directories.SetOwnerAndGroupUsingLinux(currentUsername, relayUser, relayUser, DataDirPath)
	}

	// Show success messages
	SuccessMessages(relayDomain, httpsEnabled)
}
