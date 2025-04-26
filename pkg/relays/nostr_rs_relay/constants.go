package nostr_rs_relay

const GitRepoBranch = "0.9.0"
const GitRepoURL = "https://github.com/scsibug/nostr-rs-relay"
const GitRepoTmpDirPath = "/tmp/nostr-rs-relay"
const DownloadURL = "https://github.com/nodetec/relays/releases/download/v0.5.0/nostr-rs-relay-0.9.0-x86_64-linux-gnu.tar.gz"
const BinaryName = "nostr-rs-relay"
const BinaryFilePath = "/usr/local/bin/nostr-rs-relay"
const NginxConfigFilePath = "/etc/nginx/conf.d/nostr_rs_relay.conf"
const DataDirPath = "/var/lib/nostr-rs-relay"
const DatabaseFilePath = "/var/lib/nostr-rs-relay/db/nostr.db"
const DatabaseBackupsDirPath = "/var/lib/nostr-rs-relay/db/backups"
const BackupFileNameBase = "nostr.db-bak"
const TmpConfigFilePath = "/tmp/nostr-rs-relay/config.toml"
const ConfigDirPath = "/etc/nostr-rs-relay"
const ConfigFilePath = "/etc/nostr-rs-relay/config.toml"
const ServiceName = "nostr-rs-relay"
const ServiceFilePath = "/etc/systemd/system/nostr-rs-relay.service"
const ServiceFileTemplate = `[Unit]
Description=nostr-rs-relay Service
After=network.target

[Service]
Type=simple
User=nostr
Group=nostr
WorkingDirectory=/home/nostr
Environment=RUST_LOG=info,nostr_rs_relay=info
ExecStart={{.BinaryFilePath}} --config /etc/nostr-rs-relay/config.toml --db /var/lib/nostr-rs-relay/db
Restart=on-failure

[Install]
WantedBy=multi-user.target
`
const RelayName = "nostr-rs-relay"
const GithubLink = "https://github.com/scsibug/nostr-rs-relay"
