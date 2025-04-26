package strfry

const GitRepoBranch = "1.0.4"
const GitRepoURL = "https://github.com/hoytech/strfry.git"
const GitRepoTmpDirPath = "/tmp/strfry"
const DownloadURL = "https://github.com/nodetec/relays/releases/download/v0.5.0/strfry-1.0.4-x86_64-linux-gnu.tar.gz"
const BinaryName = "strfry"
const BinaryFilePath = "/usr/local/bin/strfry"
const BinaryVersion = "1.0.4"
const NginxConfigFilePath = "/etc/nginx/conf.d/strfry.conf"
const DataDirPath = "/var/lib/strfry"
const DatabaseFilePath = "/var/lib/strfry/db/data.mdb"
const SupportedDatabaseVersionOutput = "DB version: 3"
const SupportedDatabaseVersion = "3"
const DatabaseBackupsDirPath = "/var/lib/strfry/db/backups"
const BackupFileNameBase = "data.mdb-bak"
const TmpConfigFilePath = "/tmp/strfry/strfry.conf"
const ConfigDirPath = "/etc/strfry"
const ConfigFilePath = "/etc/strfry/strfry.conf"
const ServiceName = "strfry"
const ServiceFilePath = "/etc/systemd/system/strfry.service"

// TODO
// Check working directory
// WorkingDirectory=/home/nostr
const ServiceFileTemplate = `[Unit]
Description=strfry Relay Service
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
