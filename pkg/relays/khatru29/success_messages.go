package khatru29

import (
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/messages"
)

func SuccessMessages(domain string, httpsEnabled bool) {
	successMsgParams := messages.SuccessMsgParams{Domain: domain, HTTPSEnabled: httpsEnabled, DataDirPath: DataDirPath, NginxConfigFilePath: relays.Khatru29NginxConfigFilePath, BinaryFilePath: BinaryFilePath, EnvFilePath: EnvFilePath, ServiceFilePath: ServiceFilePath, ServiceName: ServiceName, RelayName: relays.Khatru29RelayName, GitHubLink: GithubLink}
	messages.SuccessMessages(&successMsgParams)
}
