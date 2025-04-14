package strfry

import (
	"fmt"
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/pterm/pterm"
	"os"
)

// Ensure the data directory exists and set permissions
func checkForDataDirectoryAndSetPermissions() {
	directories.CreateDirectory(DataDirPath, 0755)
	directories.CreateDirectory(fmt.Sprintf("%s/%s", DataDirPath, relays.DBDir), 0755)
}

// Function to set up the relay data directory
func SetUpRelayDataDir(howToHandleExistingDatabase string) {
	spinner, _ := pterm.DefaultSpinner.Start("Configuring relay data directory...")

	if howToHandleExistingDatabase == ExistingDatabaseNotFound {
		spinner.UpdateText("Creating data directory...")
		checkForDataDirectoryAndSetPermissions()
	} else if howToHandleExistingDatabase == BackupDatabaseFileOption {
		spinner.UpdateText("Checking for data directory...")
		checkForDataDirectoryAndSetPermissions()
	} else if howToHandleExistingDatabase == UseExistingDatabaseFileOption {
		spinner.UpdateText("Checking for data directory...")
		checkForDataDirectoryAndSetPermissions()
		files.SetPermissions(DatabaseFilePath, DatabaseFilePerms)
		files.SetPermissions(DatabaseLockFilePath, DatabaseLockFilePerms)
	} else if howToHandleExistingDatabase == OverwriteDatabaseFileOption {
		spinner.UpdateText("Overwriting database...")
		checkForDataDirectoryAndSetPermissions()
		files.RemoveFile(DatabaseFilePath)
		files.RemoveFile(DatabaseLockFilePath)
	} else {
		pterm.Println()
		pterm.Error.Println(("Failed to set up data directory"))
		os.Exit(1)
	}

	// Use chown command to set ownership of the data directory to the nostr user
	directories.SetOwnerAndGroup(relays.User, relays.User, DataDirPath)

	spinner.Success("Data directory set up")
}
