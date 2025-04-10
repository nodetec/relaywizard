package strfry29

import (
	"fmt"
	"github.com/nodetec/rwz/pkg/network"
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/nodetec/rwz/pkg/utils/git"
	"github.com/nodetec/rwz/pkg/utils/systemd"
	"github.com/nodetec/rwz/pkg/verification"
	"github.com/pterm/pterm"
)

// Install the relay
func Install(relayDomain, pubKey, privKey, relayContact string) {
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

	// Check if the service file exists and disable and stop the service if it does
	systemd.DisableAndStopService(ServiceFilePath, ServiceName)

	pterm.Println()
	relayBinaryCheckSpinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Checking for existing %s binary...", BinaryName))

	// Check if relay binary exists
	if !files.FileExists(BinaryFilePath) {
		relayBinaryCheckSpinner.Info(fmt.Sprintf("%s binary not found", BinaryName))
		pterm.Println()

		// Determine the temporary file path
		tmpCompressedBinaryFilePath := files.FilePathFromFilePathBase(DownloadURL, relays.TmpDirPath)

		// Check if the temporary file exists and remove it if it does
		files.RemoveFile(tmpCompressedBinaryFilePath)

		// Download and copy the file
		downloadSpinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Downloading %s binary...", BinaryName))
		files.DownloadAndCopyFile(tmpCompressedBinaryFilePath, DownloadURL, 0644)
		downloadSpinner.Success(fmt.Sprintf("%s binary downloaded", BinaryName))

		// Verify relay binary
		verification.VerifyRelayBinary(BinaryName, tmpCompressedBinaryFilePath)

		// Install the compressed relay binary and make it executable
		installSpinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Installing %s binary...", BinaryName))
		files.InstallCompressedBinary(tmpCompressedBinaryFilePath, relays.BinaryDestDir, BinaryName, relays.BinaryFilePerms)
		installSpinner.Success(fmt.Sprintf("%s binary installed", BinaryName))
	} else {
		relayBinaryCheckSpinner.Info(fmt.Sprintf("%s binary found", BinaryName))
		pterm.Println()
	}

	// Determine the temporary file path
	tmpCompressedBinaryPluginFilePath := files.FilePathFromFilePathBase(BinaryPluginDownloadURL, relays.TmpDirPath)

	// Check if the temporary file exists and remove it if it does
	files.RemoveFile(tmpCompressedBinaryPluginFilePath)

	// Download and copy the file
	binaryPluginDownloadSpinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Downloading %s plugin binary...", BinaryPluginName))
	files.DownloadAndCopyFile(tmpCompressedBinaryPluginFilePath, BinaryPluginDownloadURL, 0644)
	binaryPluginDownloadSpinner.Success(fmt.Sprintf("%s plugin binary downloaded", BinaryPluginName))

	// Verify relay binary plugin
	verification.VerifyRelayBinary(fmt.Sprintf("%s plugin", BinaryPluginName), tmpCompressedBinaryPluginFilePath)

	// Install the compressed relay binary plugin and make it executable
	binaryPluginInstallSpinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Installing %s plugin binary...", BinaryPluginName))
	files.InstallCompressedBinary(tmpCompressedBinaryPluginFilePath, relays.BinaryDestDir, BinaryPluginName, relays.BinaryFilePerms)
	binaryPluginInstallSpinner.Success(fmt.Sprintf("%s plugin binary installed", BinaryPluginName))

	// Set up the relay data directory
	SetUpRelayDataDir()

	// Configure the relay
	ConfigureRelay(relayDomain, pubKey, privKey, relayContact)

	// Set up the relay service
	SetUpRelayService()

	// Show success messages
	SuccessMessages(relayDomain, httpsEnabled)
}
