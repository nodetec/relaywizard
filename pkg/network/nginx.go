package network

import (
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/directories"
)

// Function to configure Nginx
func ConfigureNginx(currentUsername string) {
	if directories.DirExists(NginxConfDirPath) {
		if currentUsername == relays.RootUser {
			directories.SetPermissions(NginxConfDirPath, 0755)
			directories.SetOwnerAndGroup(relays.NginxUser, relays.NginxUser, NginxConfDirPath)
		} else {
			directories.SetPermissionsUsingLinux(currentUsername, NginxConfDirPath, "0755")
			directories.SetOwnerAndGroupUsingLinux(currentUsername, relays.NginxUser, relays.NginxUser, NginxConfDirPath)
		}
	}

	if directories.DirExists(WWWDirPath) {
		if currentUsername == relays.RootUser {
			directories.SetPermissions(WWWDirPath, 0755)
			directories.SetOwnerAndGroup(relays.NginxUser, relays.NginxUser, WWWDirPath)
		} else {
			directories.SetPermissionsUsingLinux(currentUsername, WWWDirPath, "0755")
			directories.SetOwnerAndGroupUsingLinux(currentUsername, relays.NginxUser, relays.NginxUser, WWWDirPath)
		}
	}
}
