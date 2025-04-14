package strfry

import (
	"fmt"
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/pterm/pterm"
	"strconv"
)

// TODO
// Improve backup process by creating a unique and descriptive backup file name, e.g., data.mdb-<pubkey-of-relay-runner>-<backup-utc-timestamp>-<unique-identifier>.bak and then check if the file exists and create the file if it doesn't or try to create a new unique file name if it already exists
func createUniqueBackupFileName() string {
	const BeginningOfBackupFileName = "data.mdb-bak"
	backupFileNumber := 0
	uniqueBackupFileName := fmt.Sprintf("%s-%s", BeginningOfBackupFileName, strconv.Itoa((backupFileNumber)))

	for files.FileExists(fmt.Sprintf("%s/%s", DatabaseBackupsDirPath, uniqueBackupFileName)) {
		backupFileNumber++
		uniqueBackupFileName = fmt.Sprintf("%s-%s", BeginningOfBackupFileName, strconv.Itoa(backupFileNumber))
	}

	return uniqueBackupFileName
}

func BackupDatabase() {
	spinner, _ := pterm.DefaultSpinner.Start("Backing up database...")

	// Ensure the backups directory exists and set permissions
	directories.CreateDirectory(DatabaseBackupsDirPath, DatabaseBackupsDirPerms)

	spinner.UpdateText("Moving database to the backups directory...")
	uniqueBackupFileName := createUniqueBackupFileName()
	files.MoveFile(DatabaseFilePath, fmt.Sprintf("%s/%s", DatabaseBackupsDirPath, uniqueBackupFileName))

	files.RemoveFile(DatabaseLockFilePath)

	// Set permissions for the backup file
	files.SetPermissions(fmt.Sprintf("%s/%s", DatabaseBackupsDirPath, uniqueBackupFileName), DatabaseFilePerms)

	spinner.Success("Database backed up")
}
