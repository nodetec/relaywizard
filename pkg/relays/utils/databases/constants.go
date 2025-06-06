package databases

const DatabaseFilePerms = 0644
const DatabaseBackupsDirPerms = 0755
const DatabaseLockFilePerms = 0644

// nostr-rs-relay
const NostrRsRelayName = "nostr-rs-relay"
const NostrRsRelayDatabaseSHMFilePath = "/var/lib/nostr-rs-relay/db/nostr.db-shm"
const NostrRsRelayDatabaseSHMFilePerms = 0644
const NostrRsRelayDatabaseWALFilePath = "/var/lib/nostr-rs-relay/db/nostr.db-wal"
const NostrRsRelayDatabaseWALFilePerms = 0644

// Khatru Pyramid
const KhatruPyramidRelayName = "Khatru Pyramid"
const KhatruPyramidDatabaseLockFilePath = "/var/lib/khatru-pyramid/db/lock.mdb"

// strfry
const StrfryRelayName = "strfry"
const StrfryDatabaseLockFilePath = "/var/lib/strfry/db/lock.mdb"

// Khatru29
const Khatru29RelayName = "Khatru29"
const Khatru29DatabaseLockFilePath = "/var/lib/khatru29/db/lock.mdb"

// strfry29
const Strfry29RelayName = "strfry29"
const Strfry29DatabaseLockFilePath = "/var/lib/strfry29/db/lock.mdb"

// TODO
// Look into using an enum
const ExistingDatabaseNotFound = "Existing database not found"
const BackupDatabaseFileOption = "Backup database (experimental)"
const UseExistingDatabaseFileOption = "Use existing database"
const OverwriteDatabaseFileOption = "Overwrite database"
