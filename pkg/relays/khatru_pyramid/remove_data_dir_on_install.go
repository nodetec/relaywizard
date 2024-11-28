package khatru_pyramid

import (
	"fmt"
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/pterm/pterm"
)

// Function to check if data directory should be removed on install
func RemoveDataDirOnInstall(pubKey string) {
	pubKeyCheckSpinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Checking for public key in the %s file...", UsersFilePath))
	// Check if users.json file exists
	if files.FileExists(UsersFilePath) {
		// Check if the pubKey exists in the users.json file
		lineExists := files.LineExists(fmt.Sprintf(`"%s":""`, pubKey), UsersFilePath)

		// If false remove data directory
		if !lineExists {
			pubKeyCheckSpinner.Info("Public key not found, removing data directory.")
			directories.RemoveDirectory(DataDirPath)
		} else {
			pubKeyCheckSpinner.Info("Public key found, keeping data directory.")
		}
	} else {
		pubKeyCheckSpinner.Info(fmt.Sprintf("%s file not found, keeping data directory if it exists.", UsersFilePath))
	}
}
