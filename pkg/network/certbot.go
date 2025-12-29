package network

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/nodetec/rwz/pkg/logs"
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/nodetec/rwz/pkg/utils/logging"
	"github.com/pterm/pterm"
)

func setDomainCertDirPerms(domainName string) {
	domainCertificateDirPath := fmt.Sprintf("%s/%s", CertificateDirPath, domainName)

	if directories.DirExists(domainCertificateDirPath) {
		directories.SetPermissions(domainCertificateDirPath, 0700)
	}
}

func setDomainCertArchiveDirPerms(domainName string) {
	domainCertificateArchiveDirPath := fmt.Sprintf("%s/%s", CertificateArchiveDirPath, domainName)

	if directories.DirExists(domainCertificateArchiveDirPath) {
		directories.SetPermissions(domainCertificateArchiveDirPath, 0700)
	}
}

func setDomainCertArchiveFilePerms(domainName string) {
	fullchainArchiveFilePath := fmt.Sprintf("%s/%s/%s", CertificateArchiveDirPath, domainName, FullchainArchiveFile)
	privkeyArchiveFilePath := fmt.Sprintf("%s/%s/%s", CertificateArchiveDirPath, domainName, PrivkeyArchiveFile)
	chainArchiveFilePath := fmt.Sprintf("%s/%s/%s", CertificateArchiveDirPath, domainName, ChainArchiveFile)
	certArchiveFilePath := fmt.Sprintf("%s/%s/%s", CertificateArchiveDirPath, domainName, CertArchiveFile)

	if files.FileExists(fullchainArchiveFilePath) {
		files.SetPermissions(fullchainArchiveFilePath, 0600)
	}

	if files.FileExists(privkeyArchiveFilePath) {
		files.SetPermissions(privkeyArchiveFilePath, 0600)
	}

	if files.FileExists(chainArchiveFilePath) {
		files.SetPermissions(chainArchiveFilePath, 0600)
	}

	if files.FileExists(certArchiveFilePath) {
		files.SetPermissions(certArchiveFilePath, 0600)
	}
}

// Check if certificates already exist
func checkForCertificates(currentUsername, domainName string) bool {
	// TODO
	// Check how symbolic links should be handled here
	// This can be simplified since the files in the archive directory can just be checked to see if they exist and the permissions can then be set there since the files in /etc/letsencrypt/live/domainName are just symbolic links to the archive directory
	// Also the archive files may have appropriate permissions already so this may be unecessary
	// Also there are potentially multiple archive files that are created and get appended with a number
	fullchainFilePath := fmt.Sprintf("%s/%s/%s", CertificateDirPath, domainName, FullchainFile)
	privkeyFilePath := fmt.Sprintf("%s/%s/%s", CertificateDirPath, domainName, PrivkeyFile)
	chainFilePath := fmt.Sprintf("%s/%s/%s", CertificateDirPath, domainName, ChainFile)

	if currentUsername == relays.RootUser {
		if files.FileExists(fullchainFilePath) &&
			files.FileExists(privkeyFilePath) &&
			files.FileExists(chainFilePath) {
			setDomainCertDirPerms(domainName)
			setDomainCertArchiveDirPerms(domainName)
			setDomainCertArchiveFilePerms(domainName)

			return true
		}
		return false
	} else {
		if files.CheckIfFileExistsAndSetPermissionsUsingLinux(currentUsername, fullchainFilePath, "0600") &&
			files.CheckIfFileExistsAndSetPermissionsUsingLinux(currentUsername, privkeyFilePath, "0600") &&
			files.CheckIfFileExistsAndSetPermissionsUsingLinux(currentUsername, chainFilePath, "0600") {
			domainCertificateDirPath := fmt.Sprintf("%s/%s", CertificateDirPath, domainName)
			directories.CheckIfDirectoryExistsAndSetPermissionsUsingLinux(currentUsername, domainCertificateDirPath, "0700")

			domainCertificateArchiveDirPath := fmt.Sprintf("%s/%s", CertificateArchiveDirPath, domainName)
			directories.CheckIfDirectoryExistsAndSetPermissionsUsingLinux(currentUsername, domainCertificateArchiveDirPath, "0700")

			fullchainArchiveFilePath := fmt.Sprintf("%s/%s/%s", CertificateArchiveDirPath, domainName, FullchainArchiveFile)
			privkeyArchiveFilePath := fmt.Sprintf("%s/%s/%s", CertificateArchiveDirPath, domainName, PrivkeyArchiveFile)
			chainArchiveFilePath := fmt.Sprintf("%s/%s/%s", CertificateArchiveDirPath, domainName, ChainArchiveFile)
			certArchiveFilePath := fmt.Sprintf("%s/%s/%s", CertificateArchiveDirPath, domainName, CertArchiveFile)

			files.CheckIfFileExistsAndSetPermissionsUsingLinux(currentUsername, fullchainArchiveFilePath, "0600")
			files.CheckIfFileExistsAndSetPermissionsUsingLinux(currentUsername, privkeyArchiveFilePath, "0600")
			files.CheckIfFileExistsAndSetPermissionsUsingLinux(currentUsername, chainArchiveFilePath, "0600")
			files.CheckIfFileExistsAndSetPermissionsUsingLinux(currentUsername, certArchiveFilePath, "0600")

			return true
		}
		return false
	}
}

// Function to get SSL/TLS certificates using Certbot
func GetCertificates(currentUsername, domainName, nginxConfigFilePath string) bool {
	ThemeDefault := pterm.ThemeDefault

	prompt := pterm.InteractiveContinuePrinter{
		DefaultValueIndex: 0,
		DefaultText:       "Obtain SSL/TLS certificates?",
		TextStyle:         &ThemeDefault.PrimaryStyle,
		Options:           []string{"yes", "no"},
		OptionsStyle:      &ThemeDefault.SuccessMessageStyle,
		SuffixStyle:       &ThemeDefault.SecondaryStyle,
		Delimiter:         ": ",
	}

	pterm.Println()
	pterm.Println(pterm.LightCyan("Do you want to obtain SSL/TLS certificates using Certbot?"))
	pterm.Println(pterm.LightCyan("If you select 'yes', then this step requires that you already have a configured domain name."))
	pterm.Println(pterm.LightCyan("You can always re-run this installer after you have configured your domain name."))
	pterm.Println()

	result, _ := prompt.Show()

	if result == "no" {
		var certificatesExist = checkForCertificates(currentUsername, domainName)

		if certificatesExist {
			ConfigureNginxHttpsRedirect(currentUsername, domainName, nginxConfigFilePath)
		}

		pterm.Println()
		return false
	}

	pterm.Println()
	certbotSpinner, _ := pterm.DefaultSpinner.Start("Checking for Certbot email...")

	var unableToFindExistingCertbotAccount bool
	var certbotAccountData string

	if currentUsername == relays.RootUser {
		out, err := exec.Command("certbot", "show_account").CombinedOutput()
		certbotAccountData = string(out)

		unableToFindExistingCertbotAccount = strings.Contains(certbotAccountData, "Could not find an existing account for server")

		if err != nil {
			if !unableToFindExistingCertbotAccount {
				logging.AppendRWZLogFile(currentUsername, logs.RWZLogFilePath, fmt.Sprintf("Failed to retrieve Certbot account data: %v", err))
				pterm.Println()
				pterm.Error.Printfln("Failed to retrieve Certbot account data: %v", err)
				os.Exit(1)
			}
		}
	} else {
		out, err := exec.Command("sudo", "certbot", "show_account").CombinedOutput()
		certbotAccountData = string(out)

		unableToFindExistingCertbotAccount = strings.Contains(certbotAccountData, "Could not find an existing account for server")

		if err != nil {
			if !unableToFindExistingCertbotAccount {
				logging.AppendRWZLogFile(currentUsername, logs.RWZLogFilePath, fmt.Sprintf("Failed to retrieve Certbot account data: %v", err))
				pterm.Println()
				pterm.Error.Printfln("Failed to retrieve Certbot account data: %v", err)
				os.Exit(1)
			}
		}
	}

	var email string

	if unableToFindExistingCertbotAccount {
		certbotSpinner.Info("Certbot account not found.")

		pterm.Println()
		pterm.Println(pterm.LightCyan("Set your Certbot email to receive notifications from Let's Encrypt about your SSL/TLS certificates."))

		pterm.Println()
		pterm.Println(pterm.Yellow("Leave email empty if you don't want to receive notifications."))

		pterm.Println()
		email, _ = pterm.DefaultInteractiveTextInput.Show("Email address")
	} else if strings.Contains(certbotAccountData, "Email contact: none") {
		certbotSpinner.Info("Certbot email currently set to none.")

		pterm.Println()
		pterm.Println(pterm.LightCyan("Set your Certbot email to receive notifications from Let's Encrypt about your SSL/TLS certificates."))

		pterm.Println()
		pterm.Println(pterm.Yellow("Leave email empty if you don't want to receive notifications."))

		pterm.Println()
		email, _ = pterm.DefaultInteractiveTextInput.Show("Email address")

		if currentUsername == relays.RootUser {
			err := exec.Command("certbot", "update_account", "--email", email, "--no-eff-email").Run()
			if err != nil {
				logging.AppendRWZLogFile(currentUsername, logs.RWZLogFilePath, fmt.Sprintf("Failed to set Certbot email: %v", err))
				pterm.Println()
				pterm.Error.Printfln("Failed to set Certbot email: %v", err)
				os.Exit(1)
			}
		} else {
			err := exec.Command("sudo", "certbot", "update_account", "--email", email, "--no-eff-email").Run()
			if err != nil {
				logging.AppendRWZLogFile(currentUsername, logs.RWZLogFilePath, fmt.Sprintf("Failed to set Certbot email: %v", err))
				pterm.Println()
				pterm.Error.Printfln("Failed to set Certbot email: %v", err)
				os.Exit(1)
			}
		}
	} else {
		_, currentEmail, _ := strings.Cut(certbotAccountData, "Email contact: ")
		certbotSpinner.Info(fmt.Sprintf("Email used with Certbot account: %s", currentEmail))

		prompt := pterm.InteractiveContinuePrinter{
			DefaultValueIndex: 0,
			DefaultText:       "Do you want to remove or update your Certbot email?",
			TextStyle:         &ThemeDefault.PrimaryStyle,
			Options:           []string{"yes", "no"},
			OptionsStyle:      &ThemeDefault.SuccessMessageStyle,
			SuffixStyle:       &ThemeDefault.SecondaryStyle,
			Delimiter:         ": ",
		}

		result, _ := prompt.Show()

		if result == "yes" {
			pterm.Println()
			pterm.Println(pterm.LightCyan("Set your Certbot email to receive notifications from Let's Encrypt about your SSL/TLS certificates."))

			pterm.Println()
			pterm.Println(pterm.Yellow("Leave email empty if you don't want to receive notifications."))

			pterm.Println()
			email, _ = pterm.DefaultInteractiveTextInput.Show("Email address")

			if currentUsername == relays.RootUser {
				err := exec.Command("certbot", "update_account", "--email", email, "--no-eff-email").Run()
				if err != nil {
					logging.AppendRWZLogFile(currentUsername, logs.RWZLogFilePath, fmt.Sprintf("Failed to update Certbot email: %v", err))
					pterm.Println()
					pterm.Error.Printfln("Failed to update Certbot email: %v", err)
					os.Exit(1)
				}
			} else {
				err := exec.Command("sudo", "certbot", "update_account", "--email", email, "--no-eff-email").Run()
				if err != nil {
					logging.AppendRWZLogFile(currentUsername, logs.RWZLogFilePath, fmt.Sprintf("Failed to update Certbot email: %v", err))
					pterm.Println()
					pterm.Error.Printfln("Failed to update Certbot email: %v", err)
					os.Exit(1)
				}
			}
		}
	}

	pterm.Println()
	certificateSpinner, _ := pterm.DefaultSpinner.Start("Checking SSL/TLS certificates...")

	var certificatesExist = checkForCertificates(currentUsername, domainName)

	if certificatesExist {
		certificateSpinner.Info("SSL/TLS certificates already exist.")
		pterm.Println()
		return true
	}

	certificateSpinner.UpdateText("Obtaining SSL/TLS certificates...")

	domainDirPath := fmt.Sprintf("%s/%s", WWWDirPath, domainName)

	if email == "" {
		if currentUsername == relays.RootUser {
			cmd := exec.Command("certbot", "certonly", "--webroot", "-w", domainDirPath, "-d", domainName, "--agree-tos", "--no-eff-email", "-q", "--register-unsafely-without-email")
			err := cmd.Run()
			if err != nil {
				logging.AppendRWZLogFile(currentUsername, logs.RWZLogFilePath, fmt.Sprintf("Certbot failed to obtain the certificate for %s: %v", domainName, err))
				pterm.Error.Printfln("Certbot failed to obtain the certificate for %s: %v", domainName, err)
				os.Exit(1)
			}
		} else {
			cmd := exec.Command("sudo", "certbot", "certonly", "--webroot", "-w", domainDirPath, "-d", domainName, "--agree-tos", "--no-eff-email", "-q", "--register-unsafely-without-email")
			err := cmd.Run()
			if err != nil {
				logging.AppendRWZLogFile(currentUsername, logs.RWZLogFilePath, fmt.Sprintf("Certbot failed to obtain the certificate for %s: %v", domainName, err))
				pterm.Error.Printfln("Certbot failed to obtain the certificate for %s: %v", domainName, err)
				os.Exit(1)
			}
		}
	} else {
		if currentUsername == relays.RootUser {
			cmd := exec.Command("certbot", "certonly", "--webroot", "-w", domainDirPath, "-d", domainName, "--email", email, "--agree-tos", "--no-eff-email", "-q")
			err := cmd.Run()
			if err != nil {
				logging.AppendRWZLogFile(currentUsername, logs.RWZLogFilePath, fmt.Sprintf("Certbot failed to obtain the certificate for %s: %v", domainName, err))
				pterm.Error.Printfln("Certbot failed to obtain the certificate for %s: %v", domainName, err)
				os.Exit(1)
			}
		} else {
			cmd := exec.Command("sudo", "certbot", "certonly", "--webroot", "-w", domainDirPath, "-d", domainName, "--email", email, "--agree-tos", "--no-eff-email", "-q")
			err := cmd.Run()
			if err != nil {
				logging.AppendRWZLogFile(currentUsername, logs.RWZLogFilePath, fmt.Sprintf("Certbot failed to obtain the certificate for %s: %v", domainName, err))
				pterm.Error.Printfln("Certbot failed to obtain the certificate for %s: %v", domainName, err)
				os.Exit(1)
			}
		}
	}

	if currentUsername == relays.RootUser {
		setDomainCertDirPerms(domainName)
		setDomainCertArchiveDirPerms(domainName)
		setDomainCertArchiveFilePerms(domainName)
	} else {
		domainCertificateDirPath := fmt.Sprintf("%s/%s", CertificateDirPath, domainName)
		directories.CheckIfDirectoryExistsAndSetPermissionsUsingLinux(currentUsername, domainCertificateDirPath, "0700")

		domainCertificateArchiveDirPath := fmt.Sprintf("%s/%s", CertificateArchiveDirPath, domainName)
		directories.CheckIfDirectoryExistsAndSetPermissionsUsingLinux(currentUsername, domainCertificateArchiveDirPath, "0700")

		fullchainArchiveFilePath := fmt.Sprintf("%s/%s/%s", CertificateArchiveDirPath, domainName, FullchainArchiveFile)
		privkeyArchiveFilePath := fmt.Sprintf("%s/%s/%s", CertificateArchiveDirPath, domainName, PrivkeyArchiveFile)
		chainArchiveFilePath := fmt.Sprintf("%s/%s/%s", CertificateArchiveDirPath, domainName, ChainArchiveFile)
		certArchiveFilePath := fmt.Sprintf("%s/%s/%s", CertificateArchiveDirPath, domainName, CertArchiveFile)
		files.CheckIfFileExistsAndSetPermissionsUsingLinux(currentUsername, fullchainArchiveFilePath, "0600")
		files.CheckIfFileExistsAndSetPermissionsUsingLinux(currentUsername, privkeyArchiveFilePath, "0600")
		files.CheckIfFileExistsAndSetPermissionsUsingLinux(currentUsername, chainArchiveFilePath, "0600")
		files.CheckIfFileExistsAndSetPermissionsUsingLinux(currentUsername, certArchiveFilePath, "0600")
	}

	certificateSpinner.Success("SSL/TLS certificates obtained successfully.")

	return true
}
