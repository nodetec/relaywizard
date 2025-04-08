package network

import (
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/nodetec/rwz/pkg/utils/systemd"
	"github.com/pterm/pterm"
)

// Function to configure Nginx for HTTP with HTTPS redirect
func ConfigureNginxHttpsRedirect(domainName, nginxConfigFilePath string) {
	spinner, _ := pterm.DefaultSpinner.Start("Configuring Nginx HTTPS redirect...")

	files.RemoveFile(nginxConfigFilePath)

	var configContent = DetermineNginxConfigContent(nginxConfigFilePath, HTTPSNginxRedirect, domainName)

	files.WriteFile(nginxConfigFilePath, configContent, 0644)
	files.SetOwnerAndGroup(relays.NginxUser, relays.NginxUser, nginxConfigFilePath)

	systemd.RestartService("nginx")

	spinner.Success("Nginx HTTPS redirect configured")

	pterm.Println()
	pterm.Println(pterm.Yellow("Try deleting the relay's cookies in the browser if the redirect isn't working."))
}
