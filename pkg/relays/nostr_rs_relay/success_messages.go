package nostr_rs_relay

import (
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/messages"
)

func SuccessMessages(domain string, httpsEnabled bool) {
	successMsgParams := messages.SuccessMsgParams{Domain: domain, HTTPSEnabled: httpsEnabled, DataDirPath: DataDirPath, ConfigFilePath: ConfigFilePath, NginxConfigFilePath: relays.NostrRsRelayNginxConfigFilePath, BinaryFilePath: BinaryFilePath, ServiceFilePath: ServiceFilePath, ServiceName: ServiceName, RelayName: relays.NostrRsRelayName, GitHubLink: GithubLink}
	messages.SuccessMessages(&successMsgParams)
}
