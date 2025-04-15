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
func SetUpRelayDataDir(howToHandleExistingDatabase, dataDirPath, databaseFilePath, databaseLockFilePath string) {
	spinner, _ := pterm.DefaultSpinner.Start("Configuring relay data directory...")

	if howToHandleExistingDatabase == ExistingDatabaseNotFound {
		spinner.UpdateText("Creating data directory...")
		checkForDataDirectoryAndSetPermissions(dataDirPath)
	} else if howToHandleExistingDatabase == BackupDatabaseFileOption {
		spinner.UpdateText("Checking for data directory...")
		checkForDataDirectoryAndSetPermissions(dataDirPath)
	} else if howToHandleExistingDatabase == UseExistingDatabaseFileOption {
		spinner.UpdateText("Checking for data directory...")
		checkForDataDirectoryAndSetPermissions(dataDirPath)
		files.SetPermissions(databaseFilePath, DatabaseFilePerms)

		// TODO
		// Refactor
		if databaseLockFilePath == StrfryDatabaseLockFilePath || databaseLockFilePath == Strfry29DatabaseLockFilePath {
			files.SetPermissions(databaseLockFilePath, DatabaseLockFilePerms)
		}
	} else if howToHandleExistingDatabase == OverwriteDatabaseFileOption {
		spinner.UpdateText("Overwriting database...")
		checkForDataDirectoryAndSetPermissions(dataDirPath)
		files.RemoveFile(databaseFilePath)

		// TODO
		// Refactor
		if databaseLockFilePath == StrfryDatabaseLockFilePath || databaseLockFilePath == Strfry29DatabaseLockFilePath {
			files.RemoveFile(databaseLockFilePath)
		}
	} else {
		pterm.Println()
		pterm.Error.Println(("Failed to set up data directory"))
		os.Exit(1)
	}

	// Use chown command to set ownership of the data directory to the nostr user
	directories.SetOwnerAndGroup(relays.User, relays.User, dataDirPath)

	spinner.Success("Data directory set up")
}
