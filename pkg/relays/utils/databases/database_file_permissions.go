package databases

import (
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/files"
)

// Function to set database file permissions
func SetDatabaseFilePermissions(dataDirPath, databaseFilePath, relayName string) {
	if files.FileExists(databaseFilePath) {
		files.SetPermissions(databaseFilePath, DatabaseFilePerms)
	}

	if relayName == relays.NostrRsRelayName {
		if files.FileExists(NostrRsRelayDatabaseSHMFilePath) {
			files.SetPermissions(NostrRsRelayDatabaseSHMFilePath, NostrRsRelayDatabaseSHMFilePerms)
		}
		if files.FileExists(NostrRsRelayDatabaseWALFilePath) {
			files.SetPermissions(NostrRsRelayDatabaseWALFilePath, NostrRsRelayDatabaseWALFilePerms)
		}
	} else if relayName == relays.KhatruPyramidRelayName {
		if files.FileExists(KhatruPyramidDatabaseLockFilePath) {
			files.SetPermissions(KhatruPyramidDatabaseLockFilePath, DatabaseLockFilePerms)
		}
	} else if relayName == relays.StrfryRelayName {
		if files.FileExists(StrfryDatabaseLockFilePath) {
			files.SetPermissions(StrfryDatabaseLockFilePath, DatabaseLockFilePerms)
		}
	} else if relayName == relays.Khatru29RelayName {
		if files.FileExists(Khatru29DatabaseLockFilePath) {
			files.SetPermissions(Khatru29DatabaseLockFilePath, DatabaseLockFilePerms)
		}
	} else if relayName == relays.Strfry29RelayName {
		if files.FileExists(Strfry29DatabaseLockFilePath) {
			files.SetPermissions(Strfry29DatabaseLockFilePath, DatabaseLockFilePerms)
		}
	}
}
