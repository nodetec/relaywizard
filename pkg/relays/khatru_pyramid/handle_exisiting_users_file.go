package khatru_pyramid

import (
	"fmt"
	"os"

	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/pterm/pterm"
)

func backupUsersFile(currentUsername, relayUser string) {
	spinner, _ := pterm.DefaultSpinner.Start("Backing up users file...")

	// Ensure the backups directory exists and set permissions
	if currentUsername == relays.RootUser {
		directories.CreateDirectory(UsersFileBackupsDirPath, UsersFileBackupsDirPerms)
	} else {
		directories.CreateDirectoryUsingLinux(currentUsername, UsersFileBackupsDirPath)
		directories.SetPermissionsUsingLinux(currentUsername, UsersFileUsersDirPath, "0755")
		directories.SetOwnerAndGroupUsingLinux(currentUsername, relayUser, relayUser, UsersFileUsersDirPath)
		directories.SetPermissionsUsingLinux(currentUsername, UsersFileBackupsDirPath, "0755")
		directories.SetOwnerAndGroupUsingLinux(currentUsername, relayUser, relayUser, UsersFileBackupsDirPath)
	}

	var uniqueBackupFileName string
	spinner.UpdateText("Creating users file backup in the backups directory...")
	uniqueBackupFileName = files.CreateUniqueBackupFileName(UsersFileBackupsDirPath, UsersFileNameBase)
	if currentUsername == relays.RootUser {
		files.MoveFile(UsersFilePath, fmt.Sprintf("%s/%s", UsersFileBackupsDirPath, uniqueBackupFileName))
	} else {
		files.MoveFileUsingLinux(currentUsername, UsersFilePath, fmt.Sprintf("%s/%s", UsersFileBackupsDirPath, uniqueBackupFileName))
	}

	// Set permissions for the backup file
	if currentUsername == relays.RootUser {
		files.SetPermissions(fmt.Sprintf("%s/%s", UsersFileBackupsDirPath, uniqueBackupFileName), UsersFilePerms)
	} else {
		files.SetPermissionsUsingLinux(currentUsername, fmt.Sprintf("%s/%s", UsersFileBackupsDirPath, uniqueBackupFileName), "0644")
	}

	spinner.Success("Users file backed up")
}

func selectUsersFileActionOption(currentUsername, relayUser, backupUsersFileOption, useExistingUsersFileOption, overwriteUsersFileOption string, options []string) string {
	ThemeDefault := pterm.ThemeDefault

	usersFileActionSelector := pterm.InteractiveSelectPrinter{
		TextStyle:     &ThemeDefault.PrimaryStyle,
		DefaultText:   "Please select an option",
		Options:       []string{},
		OptionStyle:   &ThemeDefault.DefaultText,
		DefaultOption: "",
		MaxHeight:     3,
		Selector:      ">",
		SelectorStyle: &ThemeDefault.SecondaryStyle,
		Filter:        true,
	}

	pterm.Println()
	selectedUsersFileActionOption, _ := usersFileActionSelector.WithOptions(options).Show()

	// Display the selected option to the user with a green color for emphasis
	pterm.Info.Printfln("Selected option: %s", pterm.Green(selectedUsersFileActionOption))

	pterm.Println()
	var howToHandleExistingUsersFile string

	if selectedUsersFileActionOption == backupUsersFileOption {
		pterm.Println(pterm.LightCyan("Creating users file backup..."))
		pterm.Println()
		backupUsersFile(currentUsername, relayUser)
		howToHandleExistingUsersFile = backupUsersFileOption
	} else if selectedUsersFileActionOption == useExistingUsersFileOption {
		pterm.Println(pterm.LightCyan("Using existing users file..."))
		pterm.Println()
		if currentUsername == relays.RootUser {
			// Set permissions for the users file
			files.SetPermissions(UsersFilePath, UsersFilePerms)
			// Use chown command to set ownership of the users file to the provided relay user
			files.SetOwnerAndGroup(relayUser, relayUser, UsersFilePath)
		} else {
			// Set permissions for the users file
			files.SetPermissionsUsingLinux(currentUsername, UsersFilePath, "0644")
			// Use chown command to set ownership of the users file to the provided relay user
			files.SetOwnerAndGroupUsingLinux(currentUsername, relayUser, relayUser, UsersFilePath)
		}
		howToHandleExistingUsersFile = useExistingUsersFileOption
	} else if selectedUsersFileActionOption == overwriteUsersFileOption {
		prompt := pterm.InteractiveContinuePrinter{
			DefaultValueIndex: 0,
			DefaultText:       "Overwrite users file?",
			TextStyle:         &ThemeDefault.PrimaryStyle,
			Options:           []string{"no", "yes"},
			OptionsStyle:      &ThemeDefault.SuccessMessageStyle,
			SuffixStyle:       &ThemeDefault.SecondaryStyle,
			Delimiter:         ": ",
		}

		pterm.Println(pterm.Yellow("Warning: Are you sure you want to overwrite your existing users file?"))
		pterm.Printfln(pterm.Yellow("If you select 'yes', then the following users file will be overwritten: %s"), UsersFilePath)
		pterm.Println()

		result, _ := prompt.Show()

		if result == "no" {
			howToHandleExistingUsersFile = result
		} else if result == "yes" {
			if currentUsername == relays.RootUser {
				files.RemoveFile(UsersFilePath)
			} else {
				files.RemoveFileUsingLinux(currentUsername, UsersFilePath)
			}
			pterm.Println()
			pterm.Println(pterm.LightCyan("Users file overwitten..."))
			pterm.Println()
			howToHandleExistingUsersFile = overwriteUsersFileOption
		} else {
			pterm.Println()
			pterm.Error.Println(("Failed to confirm users file overwrite action"))
			os.Exit(1)
		}
	} else {
		pterm.Println()
		pterm.Error.Println(("Failed to perform selected users file action"))
		os.Exit(1)
	}

	return howToHandleExistingUsersFile
}

// TODO
// Look more into how the public keys are handled by the relay, e.g., does the public key associated with the relay config file
// have to be in the users.json file, if yes, does it have to be the initial public key added to the file or can it just be a
// normal user in the file?

// Function to handle existing users.json file during install
func HandleExistingUsersFile(currentUsername, pubKey, relayUser string) {
	spinner, _ := pterm.DefaultSpinner.Start("Checking for existing users file...")

	if files.FileExists(UsersFilePath) {
		spinner.UpdateText("Users file found, checking if provided public key is initial user...")

		const BackupUsersFileOption = "Backup users file (experimental)"
		const UseExistingUsersFileOption = "Use existing users file"
		const OverwriteUsersFileOption = "Overwrite users file"

		// TODO
		// Improve the line exits fcn and/or look into handling different possible patterns
		lineExistsWithoutSpace := files.LineExists(fmt.Sprintf(`"%s":""`, pubKey), UsersFilePath)
		lineExistsWithSpace := files.LineExists(fmt.Sprintf(`"%s": ""`, pubKey), UsersFilePath)

		// Users file action options
		var options []string

		if !lineExistsWithoutSpace && !lineExistsWithSpace {
			spinner.Info("Public key is not initial user in users file...")
			options = []string{BackupUsersFileOption, OverwriteUsersFileOption}
		} else {
			spinner.Info("Public key found in users file...")
			options = []string{BackupUsersFileOption, UseExistingUsersFileOption, OverwriteUsersFileOption}
		}

		var howToHandleExistingUsersFile string
		howToHandleExistingUsersFile = selectUsersFileActionOption(currentUsername, relayUser, BackupUsersFileOption, UseExistingUsersFileOption, OverwriteUsersFileOption, options)

		for howToHandleExistingUsersFile == "no" {
			howToHandleExistingUsersFile = selectUsersFileActionOption(currentUsername, relayUser, BackupUsersFileOption, UseExistingUsersFileOption, OverwriteUsersFileOption, options)
		}
	} else {
		spinner.Info("Users file not found continuing with installation...")
		pterm.Println()
	}
}
