package databases

const DatabaseFilePerms = 0644
const DatabaseBackupsDirPerms = 0755

// nostr-rs-relay
const NostrRsRelayName = "nostr-rs-relay"
const NostrRsRelayDatabaseSHMFilePath = "/var/lib/nostr-rs-relay/db/nostr.db-shm"
const NostrRsRelayDatabaseSHMFilePerms = 0644
const NostrRsRelayDatabaseWALFilePath = "/var/lib/nostr-rs-relay/db/nostr.db-wal"
const NostrRsRelayDatabaseWALFilePerms = 0644

// strfry and strfry29
const StrfryRelayName = "strfry"
const StrfryDatabaseLockFilePath = "/var/lib/strfry/db/lock.mdb"
const Strfry29RelayName = "strfry29"
const Strfry29DatabaseLockFilePath = "/var/lib/strfry29/db/lock.mdb"
const DatabaseLockFilePerms = 0644

// TODO
// Look into using an enum
const ExistingDatabaseNotFound = "Existing database not found"
const BackupDatabaseFileOption = "Backup database"
const UseExistingDatabaseFileOption = "Use existing database"
const OverwriteDatabaseFileOption = "Overwrite database"
