package databases

import (
	"github.com/nodetec/rwz/pkg/utils/files"
)

// Remove auxiliary database files
func RemoveAuxiliaryDatabaseFiles(relayName string) {
	if relayName == NostrRsRelayName {
		files.RemoveFile(NostrRsRelayDatabaseSHMFilePath)
		files.RemoveFile(NostrRsRelayDatabaseWALFilePath)
	} else if relayName == StrfryRelayName {
		files.RemoveFile(StrfryDatabaseLockFilePath)
	} else if relayName == Strfry29RelayName {
		files.RemoveFile(Strfry29DatabaseLockFilePath)
	}
}

