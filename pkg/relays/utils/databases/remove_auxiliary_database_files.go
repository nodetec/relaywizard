package databases

import (
	"github.com/nodetec/rwz/pkg/utils/files"
)

// Remove auxiliary database files
func RemoveAuxiliaryDatabaseFiles(relayName string) {
	if relayName == NostrRsRelayName {
		files.RemoveFile(NostrRsRelayDatabaseSHMFilePath)
		files.RemoveFile(NostrRsRelayDatabaseWALFilePath)
	} else if relayName == KhatruPyramidRelayName {
		files.RemoveFile(KhatruPyramidDatabaseLockFilePath)
	} else if relayName == StrfryRelayName {
		files.RemoveFile(StrfryDatabaseLockFilePath)
	} else if relayName == Khatru29RelayName {
		files.RemoveFile(Khatru29DatabaseLockFilePath)
	} else if relayName == Strfry29RelayName {
		files.RemoveFile(Strfry29DatabaseLockFilePath)
	}
}
