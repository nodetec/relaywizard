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
func backupSQLite3Database(currentUsername, relayUser, databaseFilePath, databaseDestPath string) {
	backupCommand := fmt.Sprintf(".backup '%s'", databaseDestPath)

	if currentUsername == relays.RootUser {
		err := exec.Command("sqlite3", databaseFilePath, backupCommand).Run()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to backup %s database: %v", databaseFilePath, err)
			os.Exit(1)
		}
	} else {
		err := exec.Command("sudo", "sqlite3", databaseFilePath, backupCommand).Run()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to backup %s database: %v", databaseFilePath, err)
			os.Exit(1)
		}
		directories.SetOwnerAndGroupUsingLinux(currentUsername, relayUser, relayUser, databaseDestPath)
	}
}

func BackupDatabase(currentUsername, relayUser, databaseBackupsDirPath, databaseFilePath, backupFileNameBase, relayName string) {
	spinner, _ := pterm.DefaultSpinner.Start("Backing up database...")

	// Ensure the backups directory exists and set permissions
	if currentUsername == relays.RootUser {
		directories.CreateDirectory(databaseBackupsDirPath, DatabaseBackupsDirPerms)
	} else {
		directories.CreateDirectoryUsingLinux(currentUsername, databaseBackupsDirPath)
		directories.SetPermissionsUsingLinux(currentUsername, databaseBackupsDirPath, "0755")
		directories.SetOwnerAndGroupUsingLinux(currentUsername, relayUser, relayUser, databaseBackupsDirPath)
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
		if currentUsername == relays.RootUser {
			files.MoveFile(databaseFilePath, fmt.Sprintf("%s/%s", databaseBackupsDirPath, uniqueBackupFileName))
		} else {
			files.MoveFileUsingLinux(currentUsername, databaseFilePath, fmt.Sprintf("%s/%s", databaseBackupsDirPath, uniqueBackupFileName))
		}
	}

	RemoveAuxiliaryDatabaseFiles(currentUsername, relayName)

	// Set permissions for the backup file
	if currentUsername == relays.RootUser {
		files.SetPermissions(fmt.Sprintf("%s/%s", databaseBackupsDirPath, uniqueBackupFileName), DatabaseFilePerms)
	} else {
		files.SetPermissionsUsingLinux(currentUsername, fmt.Sprintf("%s/%s", databaseBackupsDirPath, uniqueBackupFileName), "0644")
	}

	spinner.Success("Database backed up")
}
