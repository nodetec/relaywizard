package strfry29

import (
	"github.com/nodetec/rwz/pkg/utils/messages"
)

func SuccessMessages(domain string, httpsEnabled bool) {
	successMsgParams := messages.SuccessMsgParams{Domain: domain, HTTPSEnabled: httpsEnabled, DataDirPath: DataDirPath, ConfigFilePath: ConfigFilePath, PluginFilePath: PluginFilePath, NginxConfigFilePath: NginxConfigFilePath, BinaryFilePath: BinaryFilePath, BinaryPluginFilePath: BinaryPluginFilePath, ServiceFilePath: ServiceFilePath, ServiceName: ServiceName, RelayName: RelayName, GitHubLink: GithubLink}
	messages.SuccessMessages(&successMsgParams)
}
