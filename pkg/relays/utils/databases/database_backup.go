package databases

import (
	"fmt"
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/pterm/pterm"
	"strconv"
)

// TODO
// Improve backup process by creating a unique and descriptive backup file name, e.g., data.mdb-<pubkey-of-relay-runner>-<backup-utc-timestamp>-<unique-identifier>.bak and then check if the file exists and create the file if it doesn't or try to create a new unique file name if it already exists
func createUniqueBackupFileName(databaseBackupsDirPath, backupFileNameBase string) string {
	backupFileNumber := 0
	uniqueBackupFileName := fmt.Sprintf("%s-%s", backupFileNameBase, strconv.Itoa((backupFileNumber)))

	for files.FileExists(fmt.Sprintf("%s/%s", databaseBackupsDirPath, uniqueBackupFileName)) {
		backupFileNumber++
		uniqueBackupFileName = fmt.Sprintf("%s-%s", backupFileNameBase, strconv.Itoa(backupFileNumber))
	}

	return uniqueBackupFileName
}

func BackupDatabase(databaseBackupsDirPath, databaseFilePath, backupFileNameBase, databaseLockFilePath string) {
	spinner, _ := pterm.DefaultSpinner.Start("Backing up database...")

	// Ensure the backups directory exists and set permissions
	directories.CreateDirectory(databaseBackupsDirPath, DatabaseBackupsDirPerms)

	spinner.UpdateText("Moving database to the backups directory...")
	uniqueBackupFileName := createUniqueBackupFileName(databaseBackupsDirPath, backupFileNameBase)
	files.MoveFile(databaseFilePath, fmt.Sprintf("%s/%s", databaseBackupsDirPath, uniqueBackupFileName))

	// TODO
	// Refactor
	if databaseLockFilePath == StrfryDatabaseLockFilePath || databaseLockFilePath == Strfry29DatabaseLockFilePath {
		files.RemoveFile(databaseLockFilePath)
	}

	// Set permissions for the backup file
	files.SetPermissions(fmt.Sprintf("%s/%s", databaseBackupsDirPath, uniqueBackupFileName), DatabaseFilePerms)

	spinner.Success("Database backed up")
}
