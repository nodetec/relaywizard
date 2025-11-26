package databases

import (
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/files"
)

// Function to set database file permissions
func SetDatabaseFilePermissions(currentUsername, dataDirPath, databaseFilePath, relayName string) {
	if currentUsername == relays.RootUser {
		if files.FileExists(databaseFilePath) {
			files.SetPermissions(databaseFilePath, DatabaseFilePerms)
		}
	} else {
		files.CheckIfFileExistsAndSetPermissionsUsingLinux(currentUsername, databaseFilePath, "0644")
	}

	if relayName == relays.NostrRsRelayName {
		if currentUsername == relays.RootUser {
			if files.FileExists(NostrRsRelayDatabaseSHMFilePath) {
				files.SetPermissions(NostrRsRelayDatabaseSHMFilePath, NostrRsRelayDatabaseSHMFilePerms)
			}
			if files.FileExists(NostrRsRelayDatabaseWALFilePath) {
				files.SetPermissions(NostrRsRelayDatabaseWALFilePath, NostrRsRelayDatabaseWALFilePerms)
			}
		} else {
			files.CheckIfFileExistsAndSetPermissionsUsingLinux(currentUsername, NostrRsRelayDatabaseSHMFilePath, "0644")
			files.CheckIfFileExistsAndSetPermissionsUsingLinux(currentUsername, NostrRsRelayDatabaseWALFilePath, "0644")
		}
	} else if relayName == relays.KhatruPyramidRelayName {
		if currentUsername == relays.RootUser {
			if files.FileExists(KhatruPyramidDatabaseLockFilePath) {
				files.SetPermissions(KhatruPyramidDatabaseLockFilePath, DatabaseLockFilePerms)
			}
		} else {
			files.CheckIfFileExistsAndSetPermissionsUsingLinux(currentUsername, KhatruPyramidDatabaseLockFilePath, "0644")
		}
	} else if relayName == relays.StrfryRelayName {
		if currentUsername == relays.RootUser {
			if files.FileExists(StrfryDatabaseLockFilePath) {
				files.SetPermissions(StrfryDatabaseLockFilePath, DatabaseLockFilePerms)
			}
		} else {
			files.CheckIfFileExistsAndSetPermissionsUsingLinux(currentUsername, StrfryDatabaseLockFilePath, "0644")
		}
	} else if relayName == relays.Khatru29RelayName {
		if currentUsername == relays.RootUser {
			if files.FileExists(Khatru29DatabaseLockFilePath) {
				files.SetPermissions(Khatru29DatabaseLockFilePath, DatabaseLockFilePerms)
			}
		} else {
			files.CheckIfFileExistsAndSetPermissionsUsingLinux(currentUsername, Khatru29DatabaseLockFilePath, "0644")
		}
	} else if relayName == relays.Strfry29RelayName {
		if currentUsername == relays.RootUser {
			if files.FileExists(Strfry29DatabaseLockFilePath) {
				files.SetPermissions(Strfry29DatabaseLockFilePath, DatabaseLockFilePerms)
			}
		} else {
			files.CheckIfFileExistsAndSetPermissionsUsingLinux(currentUsername, Strfry29DatabaseLockFilePath, "0644")
		}
	}
}
