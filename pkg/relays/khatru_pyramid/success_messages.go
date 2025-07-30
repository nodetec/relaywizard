package khatru_pyramid

import (
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/messages"
)

func SuccessMessages(domain string, httpsEnabled bool) {
	successMsgParams := messages.SuccessMsgParams{Domain: domain, HTTPSEnabled: httpsEnabled, DataDirPath: DataDirPath, NginxConfigFilePath: relays.KhatruPyramidNginxConfigFilePath, BinaryFilePath: BinaryFilePath, EnvFilePath: EnvFilePath, ServiceFilePath: ServiceFilePath, ServiceName: ServiceName, RelayName: relays.KhatruPyramidRelayName, GitHubLink: GithubLink}
	messages.SuccessMessages(&successMsgParams)
}
