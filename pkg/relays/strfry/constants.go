package strfry

const GitRepoBranch = "1.0.4"
const GitRepoURL = "https://github.com/hoytech/strfry.git"
const GitRepoTmpDirPath = "/tmp/strfry"
const DownloadURL = "https://github.com/nodetec/relays/releases/download/v0.5.0/strfry-1.0.4-x86_64-linux-gnu.tar.gz"
const BinaryName = "strfry"
const BinaryVersion = "1.0.4"
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

const ServiceFileTemplate = `[Unit]
Description=strfry Relay Service
After=network.target

[Service]
Type=simple
User={{.RelayUser}}
Group={{.RelayUser}}
ExecStart={{.BinaryFilePath}} --config={{.ConfigFilePath}} relay
Restart=on-failure
RestartSec=5
NoNewPrivileges=yes
ProtectSystem=full
ProtectHome=yes
ProtectControlGroups=yes
ProtectKernelModules=yes
ProtectKernelTunables=yes
RestrictAddressFamilies=AF_UNIX AF_INET AF_INET6 AF_NETLINK
LockPersonality=yes
LimitCORE=1000000000

[Install]
WantedBy=multi-user.target
`
const GithubLink = "https://github.com/hoytech/strfry"
