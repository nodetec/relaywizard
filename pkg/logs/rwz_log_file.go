package logs

import (
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/nodetec/rwz/pkg/utils/files"
)

// Function to set up rwz log file
func SetUpRWZLogFile(currentUsername string) {
	if currentUsername == relays.RootUser {
		directories.CreateAllDirectories(RWZLogDirPath, 0755)
		directories.SetPermissions(RWZLogDirPath, 0755)
		directories.SetOwnerAndGroupForAllContentUsingLinux(currentUsername, relays.RootUser, relays.RootUser, RWZLogDirPath)
	} else {
		directories.CreateAllDirectoriesUsingLinux(currentUsername, RWZLogDirPath)
		directories.SetPermissionsUsingLinux(currentUsername, RWZLogDirPath, "0755")
		directories.SetOwnerAndGroupForAllContentUsingLinux(currentUsername, relays.RootUser, relays.RootUser, RWZLogDirPath)
	}

	rwzLogFileSize := files.DetermineFileSize(RWZLogFilePath)

	if rwzLogFileSize > RWZLogFileMaxSize {
		if currentUsername == relays.RootUser {
			files.RemoveFile(RWZLogFilePath)
		} else {
			files.RemoveFileUsingLinux(currentUsername, RWZLogFilePath)
		}
	}
}
