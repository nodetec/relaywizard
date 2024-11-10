package strfry29

const GitRepoBranch = "v0.4.0"
const GitRepoURL = "https://github.com/fiatjaf/relay29.git"
const GitRepoTmpDirPath = "/tmp/relay29"
const DownloadURL = "https://github.com/nodetec/relays/releases/download/v0.4.0/strfry-1.0.1-x86_64-linux-gnu.tar.gz"
const BinaryName = "strfry"
const BinaryFilePath = "/usr/local/bin/strfry"
const BinaryPluginDownloadURL = "https://github.com/nodetec/relays/releases/download/v0.4.0/relay29-0.4.0-strfry29-x86_64-linux-gnu.tar.gz"
const BinaryPluginName = "strfry29"
const BinaryPluginFilePath = "/usr/local/bin/strfry29"
const NginxConfigFilePath = "/etc/nginx/conf.d/strfry29.conf"
const DataDirPath = "/var/lib/strfry29"
const TmpConfigFilePath = "/tmp/relay29/strfry29/strfry.conf"
const ConfigDirPath = "/etc/strfry29"
const ConfigFilePath = "/etc/strfry29/strfry.conf"
const ConfigFileInfoDescription = "This is a strfry instance that only works with NIP-29 groups."

// TODO
// Currently, the strfry29 binary expects the strfry29.json file to be in the same directory
// Ideally, the location would be /etc/strfry29/strfry29.json
const PluginFilePath = "/usr/local/bin/strfry29.json"
const PluginFileTemplate = `{
  "domain": "{{.Domain}}",
  "relay_secret_key": "{{.RelaySecretKey}}",
  "strfry_config_path": "{{.ConfigFilePath}}",
  "strfry_executable_path": "{{.BinaryFilePath}}"
}
`
const ServiceName = "strfry29"
const ServiceFilePath = "/etc/systemd/system/strfry29.service"

// TODO
// Check working directory
// WorkingDirectory=/home/nostr
const ServiceFileTemplate = `[Unit]
Description=strfry29 Relay Service
After=network.target

[Service]
Type=simple
User=nostr
Group=nostr
ExecStart={{.BinaryFilePath}} --config={{.ConfigFilePath}} relay
Restart=on-failure
RestartSec=5
ProtectHome=yes
NoNewPrivileges=yes
ProtectSystem=full
LimitCORE=1000000000

[Install]
WantedBy=multi-user.target
`
const RelayName = "strfry29"
const GithubLink = "https://github.com/fiatjaf/relay29"
