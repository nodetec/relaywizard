package databases

const DatabaseFilePerms = 0644
const DatabaseBackupsDirPerms = 0755
const StrfryDatabaseLockFilePath = "/var/lib/strfry/db/lock.mdb"
const Strfry29DatabaseLockFilePath = "/var/lib/strfry29/db/lock.mdb"
const DatabaseLockFilePerms = 0644

// TODO
// Look into using an enum
const ExistingDatabaseNotFound = "Existing database not found"
const BackupDatabaseFileOption = "Backup database"
const UseExistingDatabaseFileOption = "Use existing database"
const OverwriteDatabaseFileOption = "Overwrite database"
