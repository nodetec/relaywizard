package databases

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/nodetec/rwz/pkg/logs"
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/nodetec/rwz/pkg/utils/logging"
	"github.com/pterm/pterm"
)

// Function to backup sqlite3 database
func backupSQLite3Database(currentUsername, relayUser, databaseFilePath, databaseDestFilePath string) {
	backupCommand := fmt.Sprintf(".backup '%s'", databaseDestFilePath)

	if currentUsername == relays.RootUser {
		err := exec.Command("sqlite3", databaseFilePath, backupCommand).Run()
		if err != nil {
			logging.AppendRWZLogFile(currentUsername, logs.RWZLogFilePath, fmt.Sprintf("Failed to backup %s database: %v", databaseFilePath, err))
			pterm.Println()
			pterm.Error.Printfln("Failed to backup %s database: %v", databaseFilePath, err)
			os.Exit(1)
		}
	} else {
		err := exec.Command("sudo", "sqlite3", databaseFilePath, backupCommand).Run()
		if err != nil {
			logging.AppendRWZLogFile(currentUsername, logs.RWZLogFilePath, fmt.Sprintf("Failed to backup %s database: %v", databaseFilePath, err))
			pterm.Println()
			pterm.Error.Printfln("Failed to backup %s database: %v", databaseFilePath, err)
			os.Exit(1)
		}
		files.SetOwnerAndGroupUsingLinux(currentUsername, relayUser, relayUser, databaseDestFilePath)
	}
}

func BackupDatabase(currentUsername, relayUser, databaseBackupsDirPath, databaseFilePath, backupFileNameBase, relayName string) {
	spinner, _ := pterm.DefaultSpinner.Start("Backing up database...")

	// Ensure the backups directory exists and set permissions
	if currentUsername == relays.RootUser {
		directories.CreateAllDirectories(databaseBackupsDirPath, DatabaseBackupsDirPerms)
		directories.SetPermissions(databaseBackupsDirPath, DatabaseBackupsDirPerms)
		directories.SetOwnerAndGroupForAllContentUsingLinux(currentUsername, relayUser, relayUser, databaseBackupsDirPath)
	} else {
		directories.CreateAllDirectoriesUsingLinux(currentUsername, databaseBackupsDirPath)
		directories.SetPermissionsUsingLinux(currentUsername, databaseBackupsDirPath, "0755")
		directories.SetOwnerAndGroupForAllContentUsingLinux(currentUsername, relayUser, relayUser, databaseBackupsDirPath)
	}

	var uniqueBackupFileName string
	if relayName == relays.NostrRsRelayName {
		spinner.UpdateText("Creating database backup in the backups directory...")
		uniqueBackupFileName = files.CreateUniqueBackupFileName(databaseBackupsDirPath, backupFileNameBase)
		backupSQLite3Database(currentUsername, relayUser, databaseFilePath, fmt.Sprintf("%s/%s", databaseBackupsDirPath, uniqueBackupFileName))
		if currentUsername == relays.RootUser {
			files.RemoveFile(databaseFilePath)
		} else {
			files.RemoveFileUsingLinux(currentUsername, databaseFilePath)
		}
	} else if relayName == relays.KhatruPyramidRelayName || relayName == relays.StrfryRelayName || relayName == relays.Khatru29RelayName || relayName == relays.Strfry29RelayName {
		spinner.UpdateText("Moving database to the backups directory...")
		uniqueBackupFileName = files.CreateUniqueBackupFileName(databaseBackupsDirPath, backupFileNameBase)
		// TODO
		// Look into if moving the db can cause db corruption and look for a better method
		databaseDestFilePath := fmt.Sprintf("%s/%s", databaseBackupsDirPath, uniqueBackupFileName)
		if currentUsername == relays.RootUser {
			files.MoveFileUsingLinux(currentUsername, databaseFilePath, databaseDestFilePath)
		} else {
			files.MoveFileUsingLinux(currentUsername, databaseFilePath, databaseDestFilePath)
			files.SetOwnerAndGroupUsingLinux(currentUsername, relayUser, relayUser, databaseDestFilePath)
		}
	}

	RemoveAuxiliaryDatabaseFiles(currentUsername, relayName)

	// Set permissions for the backup file
	databaseDestFilePath := fmt.Sprintf("%s/%s", databaseBackupsDirPath, uniqueBackupFileName)
	if currentUsername == relays.RootUser {
		files.SetPermissions(databaseDestFilePath, DatabaseFilePerms)
	} else {
		files.SetPermissionsUsingLinux(currentUsername, databaseDestFilePath, "0644")
	}

	spinner.Success("Database backed up")
}
