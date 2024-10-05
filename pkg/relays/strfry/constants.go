package strfry

const GitRepoBranch = "1.0.1"
const GitRepoURL = "https://github.com/hoytech/strfry.git"
const GitRepoTmpDirPath = "/tmp/strfry"
const DownloadURL = "https://github.com/nodetec/relays/releases/download/v0.4.0/strfry-1.0.1-x86_64-linux-gnu.tar.gz"
const BinaryName = "strfry"
const BinaryFilePath = "/usr/local/bin/strfry"
const NginxConfigFilePath = "/etc/nginx/conf.d/strfry.conf"
const DataDirPath = "/var/lib/strfry"
const TmpConfigFilePath = "/tmp/strfry/strfry.conf"
const ConfigDirPath = "/etc/strfry"
const ConfigFilePath = "/etc/strfry/strfry.conf"
const ServiceName = "strfry"
const ServiceFilePath = "/etc/systemd/system/strfry.service"

// TODO
// Check working directory
// WorkingDirectory=/home/nostr
const ServiceFileTemplate = `[Unit]
Description=strfry Nostr Relay Service
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
const RelayName = "strfry"
const GithubLink = "https://github.com/hoytech/strfry"
