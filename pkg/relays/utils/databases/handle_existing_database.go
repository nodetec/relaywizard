package databases

import (
	"os"

	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/pterm/pterm"
)

func selectDatabaseActionOption(currentUsername, relayUser, databaseBackupsDirPath, databaseFilePath, backupFileNameBase, relayName string) string {
	ThemeDefault := pterm.ThemeDefault

	// Supported database action options
	options := []string{BackupDatabaseFileOption, UseExistingDatabaseFileOption, OverwriteDatabaseFileOption}

	databaseActionSelector := pterm.InteractiveSelectPrinter{
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
	selectedDatabaseActionOption, _ := databaseActionSelector.WithOptions(options).Show()

	// Display the selected option to the user with a green color for emphasis
	pterm.Info.Printfln("Selected option: %s", pterm.Green(selectedDatabaseActionOption))

	pterm.Println()
	var howToHandleExistingDatabase string

	if selectedDatabaseActionOption == BackupDatabaseFileOption {
		pterm.Println(pterm.LightCyan("Creating database backup..."))
		pterm.Println()
		BackupDatabase(currentUsername, relayUser, databaseBackupsDirPath, databaseFilePath, backupFileNameBase, relayName)
		howToHandleExistingDatabase = BackupDatabaseFileOption
	} else if selectedDatabaseActionOption == UseExistingDatabaseFileOption {
		pterm.Println(pterm.LightCyan("Using existing database..."))
		pterm.Println()
		howToHandleExistingDatabase = UseExistingDatabaseFileOption
	} else if selectedDatabaseActionOption == OverwriteDatabaseFileOption {
		prompt := pterm.InteractiveContinuePrinter{
			DefaultValueIndex: 0,
			DefaultText:       "Overwrite database?",
			TextStyle:         &ThemeDefault.PrimaryStyle,
			Options:           []string{"no", "yes"},
			OptionsStyle:      &ThemeDefault.SuccessMessageStyle,
			SuffixStyle:       &ThemeDefault.SecondaryStyle,
			Delimiter:         ": ",
		}

		pterm.Println(pterm.Yellow("Warning: Are you sure you want to overwrite your existing database?"))
		pterm.Printfln(pterm.Yellow("If you select 'yes', then the following database will be overwritten: %s"), databaseFilePath)
		pterm.Println()

		result, _ := prompt.Show()

		if result == "no" {
			howToHandleExistingDatabase = result
		} else if result == "yes" {
			pterm.Println()
			pterm.Println(pterm.LightCyan("Database will be overwitten..."))
			pterm.Println()
			howToHandleExistingDatabase = OverwriteDatabaseFileOption
		} else {
			pterm.Println()
			pterm.Error.Println("Failed to confirm database overwrite action")
			os.Exit(1)
		}
	} else {
		pterm.Println()
		pterm.Error.Println(("Failed to perform selected database action"))
		os.Exit(1)
	}

	return howToHandleExistingDatabase
}

// Function to handle existing database during install
func HandleExistingDatabase(currentUsername, relayUser, databaseBackupsDirPath, databaseFilePath, backupFileNameBase, relayName string) string {
	pterm.Println()
	spinner, _ := pterm.DefaultSpinner.Start("Checking for existing database...")

	if files.FileExists(databaseFilePath) {
		spinner.Info("Database found...")

		var howToHandleExistingDatabase string
		howToHandleExistingDatabase = selectDatabaseActionOption(currentUsername, relayUser, databaseBackupsDirPath, databaseFilePath, backupFileNameBase, relayName)

		for howToHandleExistingDatabase == "no" {
			howToHandleExistingDatabase = selectDatabaseActionOption(currentUsername, relayUser, databaseBackupsDirPath, databaseFilePath, backupFileNameBase, relayName)
		}

		return howToHandleExistingDatabase
	} else {
		spinner.Info("Database not found continuing with installation...")
		pterm.Println()
	}
	return ExistingDatabaseNotFound
}
