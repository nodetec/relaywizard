package network

import (
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/nodetec/rwz/pkg/utils/systemd"
	"github.com/pterm/pterm"
)

// Function to configure Nginx for HTTPS
func ConfigureNginxHttps(domainName, nginxConfigFilePath string) {
	spinner, _ := pterm.DefaultSpinner.Start("Configuring Nginx for HTTPS...")

	files.RemoveFile(nginxConfigFilePath)

	var configContent = DetermineNginxConfigContent(nginxConfigFilePath, HTTPSScheme, domainName)

	files.WriteFile(nginxConfigFilePath, configContent, 0644)
	files.SetOwnerAndGroup(relays.NginxUser, relays.NginxUser, nginxConfigFilePath)

	systemd.RestartService("nginx")

	spinner.Success("Nginx configured for HTTPS")
}
