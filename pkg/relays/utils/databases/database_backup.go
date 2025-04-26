package databases

import (
	"fmt"
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/pterm/pterm"
	"os"
	"os/exec"
	"strconv"
)

// Function to backup sqlite3 database
func backupSQLite3Database(databaseFilePath, databaseDestPath string) {
	backupCommand := fmt.Sprintf(".backup '%s'", databaseDestPath)

	err := exec.Command("sqlite3", databaseFilePath, backupCommand).Run()
	if err != nil {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Failed to backup %s database: %v", databaseFilePath, err))
		os.Exit(1)
	}
}

// TODO
// Improve backup process by creating a unique and descriptive backup file name, e.g., <database-file-name>-<pubkey-of-relay-runner>-<utc-timestamp-of-backup>-<unique-identifier>-bak.<database-file-extension> and then check if the file exists and create the file if it doesn't or try to create a new unique file name if it already exists
func createUniqueBackupFileName(databaseBackupsDirPath, backupFileNameBase string) string {
	backupFileNumber := 0
	uniqueBackupFileName := fmt.Sprintf("%s-%s", backupFileNameBase, strconv.Itoa((backupFileNumber)))

	for files.FileExists(fmt.Sprintf("%s/%s", databaseBackupsDirPath, uniqueBackupFileName)) {
		backupFileNumber++
		uniqueBackupFileName = fmt.Sprintf("%s-%s", backupFileNameBase, strconv.Itoa(backupFileNumber))
	}

	return uniqueBackupFileName
}

func BackupDatabase(databaseBackupsDirPath, databaseFilePath, backupFileNameBase, relayName string) {
	spinner, _ := pterm.DefaultSpinner.Start("Backing up database...")

	// Ensure the backups directory exists and set permissions
	directories.CreateDirectory(databaseBackupsDirPath, DatabaseBackupsDirPerms)

	var uniqueBackupFileName string
	if relayName == NostrRsRelayName {
		spinner.UpdateText("Creating database backup in the backups directory...")
		uniqueBackupFileName = createUniqueBackupFileName(databaseBackupsDirPath, backupFileNameBase)
		backupSQLite3Database(databaseFilePath, fmt.Sprintf("%s/%s", databaseBackupsDirPath, uniqueBackupFileName))
		files.RemoveFile(databaseFilePath)
	} else if relayName == KhatruPyramidRelayName || relayName == StrfryRelayName || relayName == Khatru29RelayName || relayName == Strfry29RelayName {
		spinner.UpdateText("Moving database to the backups directory...")
		uniqueBackupFileName = createUniqueBackupFileName(databaseBackupsDirPath, backupFileNameBase)
		// TODO
		// Look into if moving the db can cause db corruption and look for a better method
		files.MoveFile(databaseFilePath, fmt.Sprintf("%s/%s", databaseBackupsDirPath, uniqueBackupFileName))
	}

	RemoveAuxiliaryDatabaseFiles(relayName)

	// Set permissions for the backup file
	files.SetPermissions(fmt.Sprintf("%s/%s", databaseBackupsDirPath, uniqueBackupFileName), DatabaseFilePerms)

	spinner.Success("Database backed up")
}
