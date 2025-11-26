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

// Function to configure Nginx for HTTPS
func ConfigureNginxHttps(currentUsername, domainName, nginxConfigFilePath string) {
	spinner, _ := pterm.DefaultSpinner.Start("Configuring Nginx for HTTPS...")

	if currentUsername == relays.RootUser {
		files.RemoveFile(nginxConfigFilePath)
	} else {
		files.RemoveFileUsingLinux(currentUsername, nginxConfigFilePath)
	}

	var configContent string

	if nginxConfigFilePath == relays.KhatruPyramidNginxConfigFilePath {
		configContent = khatru_pyramid.NginxHttpsConfigContent(domainName, WWWDirPath, AcmeChallengeDirPath, CertificateDirPath, FullchainFile, PrivkeyFile, ChainFile)
	} else if nginxConfigFilePath == relays.NostrRsRelayNginxConfigFilePath {
		configContent = nostr_rs_relay.NginxHttpsConfigContent(domainName, WWWDirPath, AcmeChallengeDirPath, CertificateDirPath, FullchainFile, PrivkeyFile, ChainFile)
	} else if nginxConfigFilePath == relays.StrfryNginxConfigFilePath {
		configContent = strfry.NginxHttpsConfigContent(domainName, WWWDirPath, AcmeChallengeDirPath, CertificateDirPath, FullchainFile, PrivkeyFile, ChainFile)
	} else if nginxConfigFilePath == relays.WotRelayNginxConfigFilePath {
		configContent = wot_relay.NginxHttpsConfigContent(domainName, WWWDirPath, AcmeChallengeDirPath, CertificateDirPath, FullchainFile, PrivkeyFile, ChainFile)
	} else if nginxConfigFilePath == relays.Khatru29NginxConfigFilePath {
		configContent = khatru29.NginxHttpsConfigContent(domainName, WWWDirPath, AcmeChallengeDirPath, CertificateDirPath, FullchainFile, PrivkeyFile, ChainFile)
	} else if nginxConfigFilePath == relays.Strfry29NginxConfigFilePath {
		configContent = strfry29.NginxHttpsConfigContent(domainName, WWWDirPath, AcmeChallengeDirPath, CertificateDirPath, FullchainFile, PrivkeyFile, ChainFile)
	} else {
		pterm.Println()
		pterm.Error.Printfln("Failed to generate Nginx config file content for %s file", nginxConfigFilePath)
		os.Exit(1)
	}

	files.WriteFile(currentUsername, nginxConfigFilePath, configContent, 0644)
	if currentUsername == relays.RootUser {
		files.SetOwnerAndGroup(relays.NginxUser, relays.NginxUser, nginxConfigFilePath)
	} else {
		files.SetOwnerAndGroupUsingLinux(currentUsername, relays.NginxUser, relays.NginxUser, nginxConfigFilePath)
	}

	systemd.RestartService(currentUsername, "nginx")

	spinner.Success("Nginx configured for HTTPS")
}
