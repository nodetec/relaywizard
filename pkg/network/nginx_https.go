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
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/nodetec/rwz/pkg/utils/systemd"
	"github.com/pterm/pterm"
	"os"
)

// Function to configure Nginx for HTTPS
func ConfigureNginxHttps(domainName, nginxConfigFilePath string) {
	spinner, _ := pterm.DefaultSpinner.Start("Configuring Nginx for HTTPS...")

	files.RemoveFile(nginxConfigFilePath)

	var configContent string

	if nginxConfigFilePath == KhatruPyramidNginxConfigFilePath {
		configContent = khatru_pyramid.NginxHttpsConfigContent(domainName, WWWDirPath, AcmeChallengeDirPath, CertificateDirPath, FullchainFile, PrivkeyFile, ChainFile)
	} else if nginxConfigFilePath == NostrRsRelayNginxConfigFilePath {
		configContent = nostr_rs_relay.NginxHttpsConfigContent(domainName, WWWDirPath, AcmeChallengeDirPath, CertificateDirPath, FullchainFile, PrivkeyFile, ChainFile)
	} else if nginxConfigFilePath == StrfryNginxConfigFilePath {
		configContent = strfry.NginxHttpsConfigContent(domainName, WWWDirPath, AcmeChallengeDirPath, CertificateDirPath, FullchainFile, PrivkeyFile, ChainFile)
	} else if nginxConfigFilePath == WotRelayNginxConfigFilePath {
		configContent = wot_relay.NginxHttpsConfigContent(domainName, WWWDirPath, AcmeChallengeDirPath, CertificateDirPath, FullchainFile, PrivkeyFile, ChainFile)
	} else if nginxConfigFilePath == Khatru29NginxConfigFilePath {
		configContent = khatru29.NginxHttpsConfigContent(domainName, WWWDirPath, AcmeChallengeDirPath, CertificateDirPath, FullchainFile, PrivkeyFile, ChainFile)
	} else if nginxConfigFilePath == Strfry29NginxConfigFilePath {
		configContent = strfry29.NginxHttpsConfigContent(domainName, WWWDirPath, AcmeChallengeDirPath, CertificateDirPath, FullchainFile, PrivkeyFile, ChainFile)
	} else {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Failed to generate Nginx config file content for %s file", nginxConfigFilePath))
		os.Exit(1)
	}

	files.WriteFile(nginxConfigFilePath, configContent, 0644)
	files.SetOwnerAndGroup(relays.NginxUser, relays.NginxUser, nginxConfigFilePath)

	systemd.RestartService("nginx")

	spinner.Success("Nginx configured for HTTPS")
}
