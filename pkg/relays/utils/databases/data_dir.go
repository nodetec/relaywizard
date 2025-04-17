package databases

import (
	"fmt"
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/pterm/pterm"
	"os"
)

// Ensure the data directory exists and set permissions
func checkForDataDirectoryAndSetPermissions(dataDirPath string) {
	directories.CreateDirectory(dataDirPath, 0755)
	directories.CreateDirectory(fmt.Sprintf("%s/%s", dataDirPath, relays.DBDir), 0755)
}

// Function to set up the relay data directory
func SetUpRelayDataDir(howToHandleExistingDatabase, dataDirPath, databaseFilePath, relayName string) {
	spinner, _ := pterm.DefaultSpinner.Start("Configuring relay data directory...")

	if howToHandleExistingDatabase == ExistingDatabaseNotFound {
		spinner.UpdateText("Creating data directory...")
		checkForDataDirectoryAndSetPermissions(dataDirPath)
		RemoveAuxiliaryDatabaseFiles(relayName)
	} else if howToHandleExistingDatabase == BackupDatabaseFileOption {
		spinner.UpdateText("Checking for data directory...")
		checkForDataDirectoryAndSetPermissions(dataDirPath)
	} else if howToHandleExistingDatabase == UseExistingDatabaseFileOption {
		spinner.UpdateText("Checking for data directory...")
		checkForDataDirectoryAndSetPermissions(dataDirPath)
	} else if howToHandleExistingDatabase == OverwriteDatabaseFileOption {
		spinner.UpdateText("Overwriting database...")
		checkForDataDirectoryAndSetPermissions(dataDirPath)
		files.RemoveFile(databaseFilePath)
		RemoveAuxiliaryDatabaseFiles(relayName)
	} else {
		pterm.Println()
		pterm.Error.Println(("Failed to set up data directory"))
		os.Exit(1)
	}

	spinner.Success("Data directory set up")
}
