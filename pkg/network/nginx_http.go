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
func ConfigureNginxHttp(currentUsername, domainName, nginxConfigFilePath string) {
	spinner, _ := pterm.DefaultSpinner.Start("Configuring Nginx for HTTP...")

	if currentUsername == relays.RootUser {
		files.RemoveFile(nginxConfigFilePath)
	} else {
		files.RemoveFileUsingLinux(currentUsername, nginxConfigFilePath)
	}

	if currentUsername == relays.RootUser {
		directories.CreateDirectory(fmt.Sprintf("%s/%s", WWWDirPath, domainName), 0755)
		directories.CreateDirectory(fmt.Sprintf("%s/%s/%s/", WWWDirPath, domainName, AcmeChallengeDirPath), 0755)
		directories.SetOwnerAndGroup(relays.NginxUser, relays.NginxUser, fmt.Sprintf("%s/%s", WWWDirPath, domainName))
	} else {
		directories.CreateDirectoryUsingLinux(currentUsername, fmt.Sprintf("%s/%s", WWWDirPath, domainName))
		directories.SetPermissionsUsingLinux(currentUsername, WWWDirPath, "0755")
		directories.SetOwnerAndGroupUsingLinux(currentUsername, relays.NginxUser, relays.NginxUser, WWWDirPath)
		directories.SetPermissionsUsingLinux(currentUsername, fmt.Sprintf("%s/%s", WWWDirPath, domainName), "0755")
		directories.SetOwnerAndGroupUsingLinux(currentUsername, relays.NginxUser, relays.NginxUser, fmt.Sprintf("%s/%s", WWWDirPath, domainName))

		directories.CreateDirectoryUsingLinux(currentUsername, fmt.Sprintf("%s/%s/%s/", WWWDirPath, domainName, AcmeChallengeDirPath))
		directories.SetPermissionsUsingLinux(currentUsername, fmt.Sprintf("%s/%s/%s/", WWWDirPath, domainName, ".well-known"), "0755")
		directories.SetOwnerAndGroupUsingLinux(currentUsername, relays.NginxUser, relays.NginxUser, fmt.Sprintf("%s/%s/%s/", WWWDirPath, domainName, ".well-known"))
		directories.SetPermissionsUsingLinux(currentUsername, fmt.Sprintf("%s/%s/%s/", WWWDirPath, domainName, AcmeChallengeDirPath), "0755")
		directories.SetOwnerAndGroupUsingLinux(currentUsername, relays.NginxUser, relays.NginxUser, fmt.Sprintf("%s/%s/%s/", WWWDirPath, domainName, AcmeChallengeDirPath))
	}

	var configContent string

	if nginxConfigFilePath == relays.KhatruPyramidNginxConfigFilePath {
		configContent = khatru_pyramid.NginxHttpConfigContent(domainName, WWWDirPath, AcmeChallengeDirPath)
	} else if nginxConfigFilePath == relays.NostrRsRelayNginxConfigFilePath {
		configContent = nostr_rs_relay.NginxHttpConfigContent(domainName, WWWDirPath, AcmeChallengeDirPath)
	} else if nginxConfigFilePath == relays.StrfryNginxConfigFilePath {
		configContent = strfry.NginxHttpConfigContent(domainName, WWWDirPath, AcmeChallengeDirPath)
	} else if nginxConfigFilePath == relays.WotRelayNginxConfigFilePath {
		configContent = wot_relay.NginxHttpConfigContent(domainName, WWWDirPath, AcmeChallengeDirPath)
	} else if nginxConfigFilePath == relays.Khatru29NginxConfigFilePath {
		configContent = khatru29.NginxHttpConfigContent(domainName, WWWDirPath, AcmeChallengeDirPath)
	} else if nginxConfigFilePath == relays.Strfry29NginxConfigFilePath {
		configContent = strfry29.NginxHttpConfigContent(domainName, WWWDirPath, AcmeChallengeDirPath)
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

	spinner.Success("Nginx configured for HTTP")
}
