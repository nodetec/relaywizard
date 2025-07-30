package databases

import (
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/files"
)

// Remove auxiliary database files
func RemoveAuxiliaryDatabaseFiles(relayName string) {
	if relayName == relays.NostrRsRelayName {
		files.RemoveFile(NostrRsRelayDatabaseSHMFilePath)
		files.RemoveFile(NostrRsRelayDatabaseWALFilePath)
	} else if relayName == relays.KhatruPyramidRelayName {
		files.RemoveFile(KhatruPyramidDatabaseLockFilePath)
	} else if relayName == relays.StrfryRelayName {
		files.RemoveFile(StrfryDatabaseLockFilePath)
	} else if relayName == relays.Khatru29RelayName {
		files.RemoveFile(Khatru29DatabaseLockFilePath)
	} else if relayName == relays.Strfry29RelayName {
		files.RemoveFile(Strfry29DatabaseLockFilePath)
	}
}
