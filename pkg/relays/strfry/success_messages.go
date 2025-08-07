package strfry

import (
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/messages"
)

func SuccessMessages(domain string, httpsEnabled bool) {
	successMsgParams := messages.SuccessMsgParams{Domain: domain, HTTPSEnabled: httpsEnabled, DataDirPath: DataDirPath, ConfigFilePath: ConfigFilePath, NginxConfigFilePath: relays.StrfryNginxConfigFilePath, BinaryFilePath: relays.StrfryBinaryFilePath, ServiceFilePath: ServiceFilePath, ServiceName: ServiceName, RelayName: relays.StrfryRelayName, GitHubLink: GithubLink}
	messages.SuccessMessages(&successMsgParams)
}
