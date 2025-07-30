package network

import (
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

// Function to configure Nginx for HTTP with HTTPS redirect
func ConfigureNginxHttpsRedirect(domainName, nginxConfigFilePath string) {
	spinner, _ := pterm.DefaultSpinner.Start("Configuring Nginx HTTPS redirect...")

	files.RemoveFile(nginxConfigFilePath)

	var configContent string

	if nginxConfigFilePath == relays.KhatruPyramidNginxConfigFilePath {
		configContent = khatru_pyramid.NginxHttpsRedirectConfigContent(domainName, WWWDirPath, AcmeChallengeDirPath, CertificateDirPath, FullchainFile, PrivkeyFile, ChainFile)
	} else if nginxConfigFilePath == relays.NostrRsRelayNginxConfigFilePath {
		configContent = nostr_rs_relay.NginxHttpsRedirectConfigContent(domainName, WWWDirPath, AcmeChallengeDirPath, CertificateDirPath, FullchainFile, PrivkeyFile, ChainFile)
	} else if nginxConfigFilePath == relays.StrfryNginxConfigFilePath {
		configContent = strfry.NginxHttpsRedirectConfigContent(domainName, WWWDirPath, AcmeChallengeDirPath, CertificateDirPath, FullchainFile, PrivkeyFile, ChainFile)
	} else if nginxConfigFilePath == relays.WotRelayNginxConfigFilePath {
		configContent = wot_relay.NginxHttpsRedirectConfigContent(domainName, WWWDirPath, AcmeChallengeDirPath, CertificateDirPath, FullchainFile, PrivkeyFile, ChainFile)
	} else if nginxConfigFilePath == relays.Khatru29NginxConfigFilePath {
		configContent = khatru29.NginxHttpsRedirectConfigContent(domainName, WWWDirPath, AcmeChallengeDirPath, CertificateDirPath, FullchainFile, PrivkeyFile, ChainFile)
	} else if nginxConfigFilePath == relays.Strfry29NginxConfigFilePath {
		configContent = strfry29.NginxHttpsRedirectConfigContent(domainName, WWWDirPath, AcmeChallengeDirPath, CertificateDirPath, FullchainFile, PrivkeyFile, ChainFile)
	} else {
		pterm.Println()
		pterm.Error.Printfln("Failed to generate Nginx config file content for %s file", nginxConfigFilePath)
		os.Exit(1)
	}

	files.WriteFile(nginxConfigFilePath, configContent, 0644)
	files.SetOwnerAndGroup(relays.NginxUser, relays.NginxUser, nginxConfigFilePath)

	systemd.RestartService("nginx")

	spinner.Success("Nginx HTTPS redirect configured")

	pterm.Println()
	pterm.Println(pterm.Yellow("Try deleting the relay's cookies in the browser if the redirect isn't working."))
}
