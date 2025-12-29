package network

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/nodetec/rwz/pkg/logs"
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/nodetec/rwz/pkg/utils/logging"
	"github.com/nodetec/rwz/pkg/utils/network"
	"github.com/pterm/pterm"
)

// Function to configure remote access
func ConfigureRemoteAccess(currentUsername string) {
	if directories.DirExists(SSHDirPath) {
		if currentUsername == relays.RootUser {
			directories.SetPermissions(SSHDirPath, 0755)
		} else {
			directories.SetPermissionsUsingLinux(currentUsername, SSHDirPath, "0755")
		}
	} else {
		logging.AppendRWZLogFile(currentUsername, logs.RWZLogFilePath, fmt.Sprintf("Failed to find %s directory", SSHDirPath))
		pterm.Println()
		pterm.Error.Printfln("Failed to find %s directory", SSHDirPath)
		os.Exit(1)
	}

	if files.FileExists(SSHDConfigFilePath) {
		if currentUsername == relays.RootUser {
			files.SetPermissions(SSHDConfigFilePath, 0644)
		} else {
			files.SetPermissionsUsingLinux(currentUsername, SSHDConfigFilePath, "0644")
		}
	} else {
		logging.AppendRWZLogFile(currentUsername, logs.RWZLogFilePath, fmt.Sprintf("Failed to find %s file", SSHDConfigFilePath))
		pterm.Println()
		pterm.Error.Printfln("Failed to find %s file", SSHDConfigFilePath)
		os.Exit(1)
	}

	if !files.LineExistsUsingLinux(SSHDConfigFileIncludeAllSSHDConfigDConfFilesLinePattern, SSHDConfigFilePath) {
		logging.AppendRWZLogFile(currentUsername, logs.RWZLogFilePath, fmt.Sprintf("Failed to find %s pattern in %s", SSHDConfigFileIncludeAllSSHDConfigDConfFilesLinePattern, SSHDConfigFilePath))
		pterm.Println()
		pterm.Error.Printfln("Failed to find %s pattern in %s", SSHDConfigFileIncludeAllSSHDConfigDConfFilesLinePattern, SSHDConfigFilePath)
		os.Exit(1)
	}

	if directories.DirExists(SSHDConfigDDirPath) {
		if currentUsername == relays.RootUser {
			directories.SetPermissions(SSHDConfigDDirPath, 0755)
		} else {
			directories.SetPermissionsUsingLinux(currentUsername, SSHDConfigDDirPath, "0755")
		}
	} else {
		logging.AppendRWZLogFile(currentUsername, logs.RWZLogFilePath, fmt.Sprintf("Failed to find %s directory", SSHDConfigDDirPath))
		pterm.Println()
		pterm.Error.Printfln("Failed to find %s directory", SSHDConfigDDirPath)
		os.Exit(1)
	}

	if currentUsername == relays.RootUser {
		files.RemoveFile(SSHDConfigDRWZConfigFilePath)
	} else {
		files.RemoveFileUsingLinux(currentUsername, SSHDConfigDRWZConfigFilePath)
	}

	ThemeDefault := pterm.ThemeDefault

	prompt := pterm.InteractiveContinuePrinter{
		DefaultValueIndex: 0,
		DefaultText:       "Do you want to add an SSH public key to the server?",
		TextStyle:         &ThemeDefault.PrimaryStyle,
		Options:           []string{"yes", "no"},
		OptionsStyle:      &ThemeDefault.SuccessMessageStyle,
		SuffixStyle:       &ThemeDefault.SecondaryStyle,
		Delimiter:         ": ",
	}

	result, _ := prompt.Show()
	pterm.Println()

	var passwordAuthentication string
	var allowOnlyPubkeyAuthenticationMethod string

	if result == "yes" {
		if currentUsername == relays.RootUser {
			directories.CreateAllDirectories(RootHiddenSSHDirPath, 0700)
			directories.SetPermissions(RootHiddenSSHDirPath, 0700)
			directories.SetOwnerAndGroupForAllContentUsingLinux(currentUsername, relays.RootUser, relays.RootUser, RootHiddenSSHDirPath)
		} else {
			userHiddenSSHDirPath := fmt.Sprintf("/home/%s/.ssh", currentUsername)
			directories.CreateAllDirectories(userHiddenSSHDirPath, 0700)
			directories.SetPermissions(userHiddenSSHDirPath, 0700)
			directories.SetOwnerAndGroupForAllContentUsingLinux(currentUsername, currentUsername, currentUsername, userHiddenSSHDirPath)
		}

		var authorizedKey string
		authorizedKey, _ = pterm.DefaultInteractiveTextInput.Show("Enter the SSH public key content from the SSH public key file")
		pterm.Println()

		if currentUsername == relays.RootUser {
			if files.FileExists(RootHiddenSSHAuthorizedKeysFilePath) {
				if !files.LineExistsUsingLinux(authorizedKey, RootHiddenSSHAuthorizedKeysFilePath) {
					files.AppendContentToFile(RootHiddenSSHAuthorizedKeysFilePath, authorizedKey, 0600)
					files.SetPermissions(RootHiddenSSHAuthorizedKeysFilePath, 0600)
				} else {
					files.SetPermissions(RootHiddenSSHAuthorizedKeysFilePath, 0600)
				}
			} else {
				files.AppendContentToFile(RootHiddenSSHAuthorizedKeysFilePath, authorizedKey, 0600)
			}
		} else {
			userHiddenSSHAuthorizedKeysFilePath := fmt.Sprintf("/home/%s/.ssh/authorized_keys", currentUsername)
			if files.FileExists(userHiddenSSHAuthorizedKeysFilePath) {
				if !files.LineExistsUsingLinux(authorizedKey, userHiddenSSHAuthorizedKeysFilePath) {
					files.AppendContentToFile(userHiddenSSHAuthorizedKeysFilePath, authorizedKey, 0600)
					files.SetPermissions(userHiddenSSHAuthorizedKeysFilePath, 0600)
				} else {
					files.SetPermissions(userHiddenSSHAuthorizedKeysFilePath, 0600)
				}
			} else {
				files.AppendContentToFile(userHiddenSSHAuthorizedKeysFilePath, authorizedKey, 0600)
			}
		}

		prompt = pterm.InteractiveContinuePrinter{
			DefaultValueIndex: 0,
			DefaultText:       "Do you want to disable SSH password authentication?",
			TextStyle:         &ThemeDefault.PrimaryStyle,
			Options:           []string{"yes", "no"},
			OptionsStyle:      &ThemeDefault.SuccessMessageStyle,
			SuffixStyle:       &ThemeDefault.SecondaryStyle,
			Delimiter:         ": ",
		}

		result, _ = prompt.Show()
		pterm.Println()

		if result == "yes" {
			prompt = pterm.InteractiveContinuePrinter{
				DefaultValueIndex: 0,
				DefaultText:       "Disable SSH password authentication?",
				TextStyle:         &ThemeDefault.PrimaryStyle,
				Options:           []string{"no", "yes"},
				OptionsStyle:      &ThemeDefault.SuccessMessageStyle,
				SuffixStyle:       &ThemeDefault.SecondaryStyle,
				Delimiter:         ": ",
			}

			pterm.Println(pterm.Yellow("Warning: Are you sure you want to disable SSH password authentication?"))
			pterm.Println(pterm.Yellow("If you select 'yes', then be sure you have correctly entered your SSH public key to prevent being locked out of your server."))

			pterm.Println()
			result, _ = prompt.Show()
			pterm.Println()

			if result == "yes" {
				passwordAuthentication = "no"

				sshdConfigDConfigFileParams := network.SSHDConfigDConfigFileParams{Port: DefaultSSHPort, AllowOnlyPubkeyAuthenticationMethod: AllowOnlyPubkeyAuthenticationMethod, PasswordAuthentication: passwordAuthentication}

				network.CreateSSHDConfigDConfigFile(currentUsername, SSHDConfigDRWZConfigFilePath, SSHDConfigDRWZConfigFileTemplate, &sshdConfigDConfigFileParams)
			} else {
				passwordAuthentication = "yes"
				allowOnlyPubkeyAuthenticationMethod = ""

				sshdConfigDConfigFileParams := network.SSHDConfigDConfigFileParams{Port: DefaultSSHPort, AllowOnlyPubkeyAuthenticationMethod: allowOnlyPubkeyAuthenticationMethod, PasswordAuthentication: passwordAuthentication}

				network.CreateSSHDConfigDConfigFile(currentUsername, SSHDConfigDRWZConfigFilePath, SSHDConfigDRWZConfigFileTemplate, &sshdConfigDConfigFileParams)
			}
		} else {
			passwordAuthentication = "yes"
			allowOnlyPubkeyAuthenticationMethod = ""

			sshdConfigDConfigFileParams := network.SSHDConfigDConfigFileParams{Port: DefaultSSHPort, AllowOnlyPubkeyAuthenticationMethod: allowOnlyPubkeyAuthenticationMethod, PasswordAuthentication: passwordAuthentication}

			network.CreateSSHDConfigDConfigFile(currentUsername, SSHDConfigDRWZConfigFilePath, SSHDConfigDRWZConfigFileTemplate, &sshdConfigDConfigFileParams)
		}

		if currentUsername == relays.RootUser {
			files.SetPermissions(SSHDConfigDRWZConfigFilePath, 0600)
		} else {
			files.SetPermissionsUsingLinux(currentUsername, SSHDConfigDRWZConfigFilePath, "0600")
		}
	} else {
		passwordAuthentication = "yes"
		allowOnlyPubkeyAuthenticationMethod = ""

		sshdConfigDConfigFileParams := network.SSHDConfigDConfigFileParams{Port: DefaultSSHPort, AllowOnlyPubkeyAuthenticationMethod: allowOnlyPubkeyAuthenticationMethod, PasswordAuthentication: passwordAuthentication}

		network.CreateSSHDConfigDConfigFile(currentUsername, SSHDConfigDRWZConfigFilePath, SSHDConfigDRWZConfigFileTemplate, &sshdConfigDConfigFileParams)

		if currentUsername == relays.RootUser {
			files.SetPermissions(SSHDConfigDRWZConfigFilePath, 0600)
		} else {
			files.SetPermissionsUsingLinux(currentUsername, SSHDConfigDRWZConfigFilePath, "0600")
		}
	}

	if currentUsername == relays.RootUser {
		err := exec.Command("/usr/sbin/sshd", "-t").Run()
		if err != nil {
			files.RemoveFile(SSHDConfigDRWZConfigFilePath)
			logging.AppendRWZLogFile(currentUsername, logs.RWZLogFilePath, fmt.Sprintf("sshd configuration tests failed: %v", err))
			pterm.Println()
			pterm.Error.Printfln("sshd configuration tests failed: %v", err)
			os.Exit(1)
		}
	} else {
		err := exec.Command("sudo", "/usr/sbin/sshd", "-t").Run()
		if err != nil {
			files.RemoveFileUsingLinux(currentUsername, SSHDConfigDRWZConfigFilePath)
			logging.AppendRWZLogFile(currentUsername, logs.RWZLogFilePath, fmt.Sprintf("sshd configuration tests failed: %v", err))
			pterm.Println()
			pterm.Error.Printfln("sshd configuration tests failed: %v", err)
			os.Exit(1)
		}
	}

	pterm.Success.Println("Remote access configured")
}
