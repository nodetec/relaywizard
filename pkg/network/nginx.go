package network

import (
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/directories"
)

// Function to configure Nginx
func ConfigureNginx() {
	if directories.DirExists(NginxConfDirPath) {
		directories.SetPermissions(NginxConfDirPath, 0755)
		directories.SetOwnerAndGroup(relays.NginxUser, relays.NginxUser, NginxConfDirPath)
	}

	if directories.DirExists(WWWDirPath) {
		directories.SetPermissions(WWWDirPath, 0755)
		directories.SetOwnerAndGroup(relays.NginxUser, relays.NginxUser, WWWDirPath)
	}
}
