package khatru_pyramid

import (
	"fmt"
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/pterm/pterm"
	"os"
	"strconv"
)

// TODO
// Create utils fcns

// TODO
// Improve backup process by creating a unique and descriptive backup file name, e.g., <users-file-name>-<pubkey-of-main-user>-<utc-timestamp-of-backup>-<unique-identifier>-bak.<users-file-extension> and then check if the file exists and create the file if it doesn't or try to create a new unique file name if it already exists
func createUniqueBackupFileName() string {
	backupFileNumber := 0
	uniqueBackupFileName := fmt.Sprintf("%s-%s", UsersFileNameBase, strconv.Itoa((backupFileNumber)))

	for files.FileExists(fmt.Sprintf("%s/%s", UsersFileBackupsDirPath, uniqueBackupFileName)) {
		backupFileNumber++
		uniqueBackupFileName = fmt.Sprintf("%s-%s", UsersFileNameBase, strconv.Itoa(backupFileNumber))
	}

	return uniqueBackupFileName
}

func backupUsersFile() {
	spinner, _ := pterm.DefaultSpinner.Start("Backing up users file...")

	// Ensure the backups directory exists and set permissions
	directories.CreateDirectory(UsersFileBackupsDirPath, UsersFileBackupsDirPerms)

	var uniqueBackupFileName string
	spinner.UpdateText("Creating users file backup in the backups directory...")
	uniqueBackupFileName = createUniqueBackupFileName()
	files.MoveFile(UsersFilePath, fmt.Sprintf("%s/%s", UsersFileBackupsDirPath, uniqueBackupFileName))

	// Set permissions for the backup file
	files.SetPermissions(fmt.Sprintf("%s/%s", UsersFileBackupsDirPath, uniqueBackupFileName), UsersFilePerms)

	spinner.Success("Users file backed up")
}

// TODO
// Look more into how the public keys are handled by the relay, e.g., does the public key associated with the relay config file
// have to be in the users.json file, if yes, does it have to be the initial public key added to the file or can it just be a
// normal user in the file?

// Function to handle existing users.json file during install
func HandleExistingUsersFile(pubKey, relayUser string) {
	spinner, _ := pterm.DefaultSpinner.Start("Checking for existing users file...")

	const BackupUsersFileOption = "Backup users file (experimental)"
	const UseExistingUsersFileOption = "Use existing users file"
	const OverwriteUsersFileOption = "Overwrite users file"

	if files.FileExists(UsersFilePath) {
		spinner.UpdateText("Users file found, checking if provided public key is initial user...")

		// TODO
		// Imporove the line exits fcn and/or look into handling different possible patterns
		lineExistsWithoutSpace := files.LineExists(fmt.Sprintf(`"%s":""`, pubKey), UsersFilePath)
		lineExistsWithSpace := files.LineExists(fmt.Sprintf(`"%s": ""`, pubKey), UsersFilePath)

		ThemeDefault := pterm.ThemeDefault

		// Users file action options
		var options []string

		if !lineExistsWithoutSpace && !lineExistsWithSpace {
			spinner.Info("Public key is not initial user in users file...")
			options = []string{BackupUsersFileOption, OverwriteUsersFileOption}
		} else {
			spinner.Info("Public key found in users file...")
			options = []string{BackupUsersFileOption, UseExistingUsersFileOption, OverwriteUsersFileOption}
		}

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
		if selectedUsersFileActionOption == BackupUsersFileOption {
			pterm.Println(pterm.LightCyan("Creating users file backup..."))
			pterm.Println()
			backupUsersFile()
		} else if selectedUsersFileActionOption == UseExistingUsersFileOption {
			pterm.Println(pterm.LightCyan("Using existing users file..."))
			pterm.Println()
			// Set permissions for the users file
			files.SetPermissions(UsersFilePath, UsersFilePerms)
			// Use chown command to set ownership of the users file to the provided relay user
			files.SetOwnerAndGroup(relayUser, relayUser, UsersFilePath)
		} else if selectedUsersFileActionOption == OverwriteUsersFileOption {
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
			pterm.Printfln(pterm.Yellow("If you select 'yes', then the following users file will be overwritten: %s", UsersFilePath))
			pterm.Println()

			result, _ := prompt.Show()

			if result == "no" {
				// TODO
				// Display options again
				pterm.Println()
				pterm.Println(pterm.LightCyan("Exiting wizard..."))
				os.Exit(1)
			} else if result == "yes" {
				pterm.Println()
				pterm.Println(pterm.LightCyan("Users file overwitten..."))
				pterm.Println()
				files.RemoveFile(UsersFilePath)
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
	} else {
		spinner.Info("Users file not found continuing with installation...")
		pterm.Println()
	}
}
