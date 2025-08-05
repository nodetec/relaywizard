package databases

import (
	"fmt"
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/pterm/pterm"
	"os"
	"os/exec"
)

// Function to backup sqlite3 database
func backupSQLite3Database(databaseFilePath, databaseDestPath string) {
	backupCommand := fmt.Sprintf(".backup '%s'", databaseDestPath)

	err := exec.Command("sqlite3", databaseFilePath, backupCommand).Run()
	if err != nil {
		pterm.Println()
		pterm.Error.Printfln("Failed to backup %s database: %v", databaseFilePath, err)
		os.Exit(1)
	}
}

func BackupDatabase(databaseBackupsDirPath, databaseFilePath, backupFileNameBase, relayName string) {
	spinner, _ := pterm.DefaultSpinner.Start("Backing up database...")

	// Ensure the backups directory exists and set permissions
	directories.CreateDirectory(databaseBackupsDirPath, DatabaseBackupsDirPerms)

	var uniqueBackupFileName string
	if relayName == relays.NostrRsRelayName {
		spinner.UpdateText("Creating database backup in the backups directory...")
		uniqueBackupFileName = files.CreateUniqueBackupFileName(databaseBackupsDirPath, backupFileNameBase)
		backupSQLite3Database(databaseFilePath, fmt.Sprintf("%s/%s", databaseBackupsDirPath, uniqueBackupFileName))
		files.RemoveFile(databaseFilePath)
	} else if relayName == relays.KhatruPyramidRelayName || relayName == relays.StrfryRelayName || relayName == relays.Khatru29RelayName || relayName == relays.Strfry29RelayName {
		spinner.UpdateText("Moving database to the backups directory...")
		uniqueBackupFileName = files.CreateUniqueBackupFileName(databaseBackupsDirPath, backupFileNameBase)
		// TODO
		// Look into if moving the db can cause db corruption and look for a better method
		files.MoveFile(databaseFilePath, fmt.Sprintf("%s/%s", databaseBackupsDirPath, uniqueBackupFileName))
	}

	RemoveAuxiliaryDatabaseFiles(relayName)

	// Set permissions for the backup file
	files.SetPermissions(fmt.Sprintf("%s/%s", databaseBackupsDirPath, uniqueBackupFileName), DatabaseFilePerms)

	spinner.Success("Database backed up")
}
