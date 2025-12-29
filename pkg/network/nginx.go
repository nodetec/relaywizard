package network

import (
	"fmt"
	"os"

	"github.com/nodetec/rwz/pkg/logs"
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/nodetec/rwz/pkg/utils/logging"
	"github.com/pterm/pterm"
)

// Function to configure Nginx
func ConfigureNginx(currentUsername string) {
	if directories.DirExists(NginxConfDirPath) {
		if currentUsername == relays.RootUser {
			directories.SetPermissions(NginxConfDirPath, 0755)
			directories.SetOwnerAndGroupForAllContentUsingLinux(currentUsername, relays.NginxUser, relays.NginxUser, NginxConfDirPath)
		} else {
			directories.SetPermissionsUsingLinux(currentUsername, NginxConfDirPath, "0755")
			directories.SetOwnerAndGroupForAllContentUsingLinux(currentUsername, relays.NginxUser, relays.NginxUser, NginxConfDirPath)
		}
	} else {
		logging.AppendRWZLogFile(currentUsername, logs.RWZLogFilePath, fmt.Sprintf("Failed to find %s directory", NginxConfDirPath))
		pterm.Println()
		pterm.Error.Printfln("Failed to find %s directory", NginxConfDirPath)
		os.Exit(1)
	}

	if directories.DirExists(WWWDirPath) {
		if currentUsername == relays.RootUser {
			directories.SetPermissions(WWWDirPath, 0755)
			directories.SetOwnerAndGroupForAllContentUsingLinux(currentUsername, relays.NginxUser, relays.NginxUser, WWWDirPath)
		} else {
			directories.SetPermissionsUsingLinux(currentUsername, WWWDirPath, "0755")
			directories.SetOwnerAndGroupForAllContentUsingLinux(currentUsername, relays.NginxUser, relays.NginxUser, WWWDirPath)
		}
	}
}
