package network

import (
	"fmt"
	"github.com/nodetec/rwz/pkg/network/relays/khatru29"
	"github.com/nodetec/rwz/pkg/network/relays/khatru_pyramid"
	"github.com/nodetec/rwz/pkg/network/relays/nostr_rs_relay"
	"github.com/nodetec/rwz/pkg/network/relays/strfry"
	"github.com/nodetec/rwz/pkg/network/relays/strfry29"
	"github.com/nodetec/rwz/pkg/network/relays/wot_relay"
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/nodetec/rwz/pkg/utils/systemd"
	"github.com/pterm/pterm"
	"os"
)

// Function to configure Nginx for HTTP
func ConfigureNginxHttp(domainName, nginxConfigFilePath string) {
	spinner, _ := pterm.DefaultSpinner.Start("Configuring Nginx for HTTP...")

	files.RemoveFile(nginxConfigFilePath)

	directories.CreateDirectory(fmt.Sprintf("%s/%s", WWWDirPath, domainName), 0755)
	directories.CreateDirectory(fmt.Sprintf("%s/%s/%s/", WWWDirPath, domainName, AcmeChallengeDirPath), 0755)
	directories.SetOwnerAndGroup(relays.NginxUser, relays.NginxUser, fmt.Sprintf("%s/%s", WWWDirPath, domainName))

	var configContent string

	if nginxConfigFilePath == KhatruPyramidNginxConfigFilePath {
		configContent = khatru_pyramid.NginxHttpConfigContent(domainName, WWWDirPath, AcmeChallengeDirPath)
	} else if nginxConfigFilePath == NostrRsRelayNginxConfigFilePath {
		configContent = nostr_rs_relay.NginxHttpConfigContent(domainName, WWWDirPath, AcmeChallengeDirPath)
	} else if nginxConfigFilePath == StrfryNginxConfigFilePath {
		configContent = strfry.NginxHttpConfigContent(domainName, WWWDirPath, AcmeChallengeDirPath)
	} else if nginxConfigFilePath == WotRelayNginxConfigFilePath {
		configContent = wot_relay.NginxHttpConfigContent(domainName, WWWDirPath, AcmeChallengeDirPath)
	} else if nginxConfigFilePath == Khatru29NginxConfigFilePath {
		configContent = khatru29.NginxHttpConfigContent(domainName, WWWDirPath, AcmeChallengeDirPath)
	} else if nginxConfigFilePath == Strfry29NginxConfigFilePath {
		configContent = strfry29.NginxHttpConfigContent(domainName, WWWDirPath, AcmeChallengeDirPath)
	} else {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Failed to generate Nginx config file content for %s file", nginxConfigFilePath))
		os.Exit(1)
	}

	files.WriteFile(nginxConfigFilePath, configContent, 0644)
	files.SetOwnerAndGroup(relays.NginxUser, relays.NginxUser, nginxConfigFilePath)

	systemd.RestartService("nginx")

	spinner.Success("Nginx configured for HTTP")
}
