package wot_relay

import (
	"fmt"
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/pterm/pterm"
)

// Function to check if data directory should be removed on install
func RemoveDataDirOnInstall(pubKey string) {
	pubKeyCheckSpinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Checking for public key in the %s file...", EnvFilePath))
	// Check if environment file exists
	if files.FileExists(EnvFilePath) {
		// Check if the pubKey exists in the environment file
		lineExists := files.LineExists(fmt.Sprintf(`RELAY_PUBKEY="%s"`, pubKey), EnvFilePath)

		// If false remove data directory
		if !lineExists {
			pubKeyCheckSpinner.Info("Public key not found, removing data directory...")
			directories.RemoveDirectory(DataDirPath)
		} else {
			pubKeyCheckSpinner.Info("Public key found, keeping data directory.")
		}
	} else {
		pubKeyCheckSpinner.Info(fmt.Sprintf("%s file not found, keeping data directory if it exists.", EnvFilePath))
	}

}
