package network

import (
	"fmt"
	"github.com/nodetec/rwz/pkg/network/relays/khatru29"
	"github.com/nodetec/rwz/pkg/network/relays/khatru_pyramid"
	"github.com/nodetec/rwz/pkg/network/relays/nostr_rs_relay"
	"github.com/nodetec/rwz/pkg/network/relays/strfry"
	"github.com/nodetec/rwz/pkg/network/relays/strfry29"
	"github.com/nodetec/rwz/pkg/network/relays/wot_relay"
	"github.com/pterm/pterm"
	"os"
)

// Function to determine config content for Nginx
func DetermineNginxConfigContent(configFilePath, configType, domainName string) string {
	var configContent string

	if configFilePath == KhatruPyramidNginxConfigFilePath {
		if configType == HTTPScheme {
			configContent = khatru_pyramid.NginxHttpConfigContent(domainName, WWWDirPath, AcmeChallengeDirPath)
		} else if configType == HTTPSScheme {
			configContent = khatru_pyramid.NginxHttpsConfigContent(domainName, WWWDirPath, AcmeChallengeDirPath, CertificateDirPath, FullchainFile, PrivkeyFile, ChainFile)
		} else if configType == HTTPSNginxRedirect {
			configContent = khatru_pyramid.NginxHttpsRedirectConfigContent(domainName, WWWDirPath, AcmeChallengeDirPath, CertificateDirPath, FullchainFile, PrivkeyFile, ChainFile)
		} else {
			pterm.Println()
			pterm.Error.Println(fmt.Sprintf("Failed to generate Nginx config file content for %s file", KhatruPyramidNginxConfigFilePath))
			os.Exit(1)
		}
	} else if configFilePath == NostrRsRelayNginxConfigFilePath {
		if configType == HTTPScheme {
			configContent = nostr_rs_relay.NginxHttpConfigContent(domainName, WWWDirPath, AcmeChallengeDirPath)
		} else if configType == HTTPSScheme {
			configContent = nostr_rs_relay.NginxHttpsConfigContent(domainName, WWWDirPath, AcmeChallengeDirPath, CertificateDirPath, FullchainFile, PrivkeyFile, ChainFile)
		} else if configType == HTTPSNginxRedirect {
			configContent = nostr_rs_relay.NginxHttpsRedirectConfigContent(domainName, WWWDirPath, AcmeChallengeDirPath, CertificateDirPath, FullchainFile, PrivkeyFile, ChainFile)
		} else {
			pterm.Println()
			pterm.Error.Println(fmt.Sprintf("Failed to generate Nginx config file content for %s file", NostrRsRelayNginxConfigFilePath))
			os.Exit(1)
		}
	} else if configFilePath == StrfryNginxConfigFilePath {
		if configType == HTTPScheme {
			configContent = strfry.NginxHttpConfigContent(domainName, WWWDirPath, AcmeChallengeDirPath)
		} else if configType == HTTPSScheme {
			configContent = strfry.NginxHttpsConfigContent(domainName, WWWDirPath, AcmeChallengeDirPath, CertificateDirPath, FullchainFile, PrivkeyFile, ChainFile)
		} else if configType == HTTPSNginxRedirect {
			configContent = strfry.NginxHttpsRedirectConfigContent(domainName, WWWDirPath, AcmeChallengeDirPath, CertificateDirPath, FullchainFile, PrivkeyFile, ChainFile)
		} else {
			pterm.Println()
			pterm.Error.Println(fmt.Sprintf("Failed to generate Nginx config file content for %s file", StrfryNginxConfigFilePath))
			os.Exit(1)
		}
	} else if configFilePath == WotRelayNginxConfigFilePath {
		if configType == HTTPScheme {
			configContent = wot_relay.NginxHttpConfigContent(domainName, WWWDirPath, AcmeChallengeDirPath)
		} else if configType == HTTPSScheme {
			configContent = wot_relay.NginxHttpsConfigContent(domainName, WWWDirPath, AcmeChallengeDirPath, CertificateDirPath, FullchainFile, PrivkeyFile, ChainFile)
		} else if configType == HTTPSNginxRedirect {
			configContent = wot_relay.NginxHttpsRedirectConfigContent(domainName, WWWDirPath, AcmeChallengeDirPath, CertificateDirPath, FullchainFile, PrivkeyFile, ChainFile)
		} else {
			pterm.Println()
			pterm.Error.Println(fmt.Sprintf("Failed to generate Nginx config file content for %s file", WotRelayNginxConfigFilePath))
			os.Exit(1)
		}
	} else if configFilePath == Khatru29NginxConfigFilePath {
		if configType == HTTPScheme {
			configContent = khatru29.NginxHttpConfigContent(domainName, WWWDirPath, AcmeChallengeDirPath)
		} else if configType == HTTPSScheme {
			configContent = khatru29.NginxHttpsConfigContent(domainName, WWWDirPath, AcmeChallengeDirPath, CertificateDirPath, FullchainFile, PrivkeyFile, ChainFile)
		} else if configType == HTTPSNginxRedirect {
			configContent = khatru29.NginxHttpsRedirectConfigContent(domainName, WWWDirPath, AcmeChallengeDirPath, CertificateDirPath, FullchainFile, PrivkeyFile, ChainFile)
		} else {
			pterm.Println()
			pterm.Error.Println(fmt.Sprintf("Failed to generate Nginx config file content for %s file", Khatru29NginxConfigFilePath))
			os.Exit(1)
		}
	} else if configFilePath == Strfry29NginxConfigFilePath {
		if configType == HTTPScheme {
			configContent = strfry29.NginxHttpConfigContent(domainName, WWWDirPath, AcmeChallengeDirPath)
		} else if configType == HTTPSScheme {
			configContent = strfry29.NginxHttpsConfigContent(domainName, WWWDirPath, AcmeChallengeDirPath, CertificateDirPath, FullchainFile, PrivkeyFile, ChainFile)
		} else if configType == HTTPSNginxRedirect {
			configContent = strfry29.NginxHttpsRedirectConfigContent(domainName, WWWDirPath, AcmeChallengeDirPath, CertificateDirPath, FullchainFile, PrivkeyFile, ChainFile)
		} else {
			pterm.Println()
			pterm.Error.Println(fmt.Sprintf("Failed to generate Nginx config file content for %s file", Strfry29NginxConfigFilePath))
			os.Exit(1)
		}
	} else {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Failed to determine Nginx config file content for %s file", configFilePath))
		os.Exit(1)
	}

	return configContent
}
