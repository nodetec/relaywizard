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
func checkForDataDirectoryAndSetPermissions(currentUsername, relayUser, dataDirPath string) {
	if currentUsername == relays.RootUser {
		directories.CreateDirectory(dataDirPath, 0755)
		directories.CreateDirectory(fmt.Sprintf("%s/%s", dataDirPath, relays.DBDir), 0755)
	} else {
		directories.CreateDirectoryUsingLinux(currentUsername, dataDirPath)
		directories.SetPermissionsUsingLinux(currentUsername, dataDirPath, "0755")
		directories.SetOwnerAndGroupUsingLinux(currentUsername, relayUser, relayUser, dataDirPath)

		directories.CreateDirectoryUsingLinux(currentUsername, fmt.Sprintf("%s/%s", dataDirPath, relays.DBDir))
		directories.SetPermissionsUsingLinux(currentUsername, fmt.Sprintf("%s/%s", dataDirPath, relays.DBDir), "0755")
		directories.SetOwnerAndGroupUsingLinux(currentUsername, relayUser, relayUser, fmt.Sprintf("%s/%s", dataDirPath, relays.DBDir))
	}
}

// Function to set up the relay data directory
func SetUpRelayDataDir(currentUsername, relayUser, howToHandleExistingDatabase, dataDirPath, databaseFilePath, relayName string) {
	spinner, _ := pterm.DefaultSpinner.Start("Configuring relay data directory...")

	if howToHandleExistingDatabase == ExistingDatabaseNotFound {
		spinner.UpdateText("Creating data directory...")
		checkForDataDirectoryAndSetPermissions(currentUsername, relayUser, dataDirPath)
		RemoveAuxiliaryDatabaseFiles(currentUsername, relayName)
	} else if howToHandleExistingDatabase == BackupDatabaseFileOption {
		spinner.UpdateText("Checking for data directory...")
		checkForDataDirectoryAndSetPermissions(currentUsername, relayUser, dataDirPath)
	} else if howToHandleExistingDatabase == UseExistingDatabaseFileOption {
		spinner.UpdateText("Checking for data directory...")
		checkForDataDirectoryAndSetPermissions(currentUsername, relayUser, dataDirPath)
	} else if howToHandleExistingDatabase == OverwriteDatabaseFileOption {
		spinner.UpdateText("Overwriting database...")
		checkForDataDirectoryAndSetPermissions(currentUsername, relayUser, dataDirPath)
		if currentUsername == relays.RootUser {
			files.RemoveFile(databaseFilePath)
		} else {
			files.RemoveFileUsingLinux(currentUsername, databaseFilePath)
		}
		RemoveAuxiliaryDatabaseFiles(currentUsername, relayName)
	} else {
		pterm.Println()
		pterm.Error.Println(("Failed to set up data directory"))
		os.Exit(1)
	}

	spinner.Success("Data directory set up")
}
