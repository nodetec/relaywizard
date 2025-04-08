package network

import (
	"fmt"
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/nodetec/rwz/pkg/utils/systemd"
	"github.com/pterm/pterm"
)

// Function to configure Nginx for HTTP
func ConfigureNginxHttp(domainName, nginxConfigFilePath string) {
	spinner, _ := pterm.DefaultSpinner.Start("Configuring Nginx for HTTP...")

	files.RemoveFile(nginxConfigFilePath)

	directories.CreateDirectory(fmt.Sprintf("%s/%s", WWWDirPath, domainName), 0755)
	directories.CreateDirectory(fmt.Sprintf("%s/%s/%s/", WWWDirPath, domainName, AcmeChallengeDirPath), 0755)
	directories.SetOwnerAndGroup(relays.NginxUser, relays.NginxUser, fmt.Sprintf("%s/%s", WWWDirPath, domainName))

	var configContent = DetermineNginxConfigContent(nginxConfigFilePath, HTTPScheme, domainName)

	files.WriteFile(nginxConfigFilePath, configContent, 0644)
	files.SetOwnerAndGroup(relays.NginxUser, relays.NginxUser, nginxConfigFilePath)

	systemd.RestartService("nginx")

	spinner.Success("Nginx configured for HTTP")
}
