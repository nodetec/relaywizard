package databases

import (
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/files"
)

// Remove auxiliary database files
func RemoveAuxiliaryDatabaseFiles(currentUsername, relayName string) {
	if relayName == relays.NostrRsRelayName {
		if currentUsername == relays.RootUser {
			files.RemoveFile(NostrRsRelayDatabaseSHMFilePath)
			files.RemoveFile(NostrRsRelayDatabaseWALFilePath)
		} else {
			files.RemoveFileUsingLinux(currentUsername, NostrRsRelayDatabaseSHMFilePath)
			files.RemoveFileUsingLinux(currentUsername, NostrRsRelayDatabaseWALFilePath)
		}
	} else if relayName == relays.KhatruPyramidRelayName {
		if currentUsername == relays.RootUser {
			files.RemoveFile(KhatruPyramidDatabaseLockFilePath)
		} else {
			files.RemoveFileUsingLinux(currentUsername, KhatruPyramidDatabaseLockFilePath)
		}
	} else if relayName == relays.StrfryRelayName {
		if currentUsername == relays.RootUser {
			files.RemoveFile(StrfryDatabaseLockFilePath)
		} else {
			files.RemoveFileUsingLinux(currentUsername, StrfryDatabaseLockFilePath)
		}
	} else if relayName == relays.Khatru29RelayName {
		if currentUsername == relays.RootUser {
			files.RemoveFile(Khatru29DatabaseLockFilePath)
		} else {
			files.RemoveFileUsingLinux(currentUsername, Khatru29DatabaseLockFilePath)
		}
	} else if relayName == relays.Strfry29RelayName {
		if currentUsername == relays.RootUser {
			files.RemoveFile(Strfry29DatabaseLockFilePath)
		} else {
			files.RemoveFileUsingLinux(currentUsername, Strfry29DatabaseLockFilePath)
		}
	}
}
