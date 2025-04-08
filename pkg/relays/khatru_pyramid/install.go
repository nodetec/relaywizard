package khatru_pyramid

import (
	"fmt"
	"github.com/nodetec/rwz/pkg/network"
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/nodetec/rwz/pkg/utils/systemd"
	"github.com/nodetec/rwz/pkg/verification"
	"github.com/pterm/pterm"
)

// Install the relay
func Install(relayDomain, pubKey, relayContact string) {
	// Configure Nginx for HTTP
	network.ConfigureNginxHttp(relayDomain, NginxConfigFilePath)

	// Get SSL/TLS certificates
	httpsEnabled := network.GetCertificates(relayDomain, NginxConfigFilePath)
	if httpsEnabled {
		// Configure Nginx for HTTPS
		network.ConfigureNginxHttps(relayDomain, NginxConfigFilePath)
	}

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

	// Check if data directory should be removed
	pterm.Println()
	RemoveDataDirOnInstall(pubKey)
	pterm.Println()

	// Install the compressed relay binary and make it executable
	installSpinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Installing %s binary...", RelayName))
	files.InstallCompressedBinary(tmpCompressedBinaryFilePath, relays.BinaryDestDir, BinaryName, relays.BinaryFilePerms)
	installSpinner.Success(fmt.Sprintf("%s binary installed", RelayName))

	// Set up the relay data directory
	SetUpRelayDataDir()

	// Configure the relay
	ConfigureRelay(relayDomain, pubKey, relayContact)

	// Set up the relay service
	SetUpRelayService()

	// Show success messages
	SuccessMessages(relayDomain, httpsEnabled)
}
