package network

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/pterm/pterm"
)

func setDomainCertDirPerms(domainName string) {
	DomainCertificateDirPath := fmt.Sprintf("%s/%s", CertificateDirPath, domainName)

	if directories.DirExists(DomainCertificateDirPath) {
		directories.SetPermissions(DomainCertificateDirPath, 0700)
	}
}

func setDomainCertArchiveDirPerms(domainName string) {
	DomainCertificateArchiveDirPath := fmt.Sprintf("%s/%s", CertificateArchiveDirPath, domainName)

	if directories.DirExists(DomainCertificateArchiveDirPath) {
		directories.SetPermissions(DomainCertificateArchiveDirPath, 0700)
	}
}

func setDomainCertArchiveFilePerms(domainName string) {
	FullchainArchiveFilePath := fmt.Sprintf("%s/%s/%s", CertificateArchiveDirPath, domainName, FullchainArchiveFile)
	PrivkeyArchiveFilePath := fmt.Sprintf("%s/%s/%s", CertificateArchiveDirPath, domainName, PrivkeyArchiveFile)
	ChainArchiveFilePath := fmt.Sprintf("%s/%s/%s", CertificateArchiveDirPath, domainName, ChainArchiveFile)
	CertArchiveFilePath := fmt.Sprintf("%s/%s/%s", CertificateArchiveDirPath, domainName, CertArchiveFile)

	if files.FileExists(FullchainArchiveFilePath) {
		files.SetPermissions(FullchainArchiveFilePath, 0600)
	}

	if files.FileExists(PrivkeyArchiveFilePath) {
		files.SetPermissions(PrivkeyArchiveFilePath, 0600)
	}

	if files.FileExists(ChainArchiveFilePath) {
		files.SetPermissions(ChainArchiveFilePath, 0600)
	}

	if files.FileExists(CertArchiveFilePath) {
		files.SetPermissions(CertArchiveFilePath, 0600)
	}
}

// Check if certificates already exist
func checkForCertificates(currentUsername, domainName string) bool {
	// TODO
	// Check how symbolic links should be handled here
	// This can be simplified since the files in the archive directory can just be checked to see if they exist and the permissions can then be set there since the files in /etc/letsencrypt/live/domainName are just symbolic links to the archive directory
	// Also the archive files may have appropriate permissions already so this may be unecessary
	// Also there are potentially multiple archive files that are created and get appended with a number
	if currentUsername == relays.RootUser {
		if files.FileExists(fmt.Sprintf("%s/%s/%s", CertificateDirPath, domainName, FullchainFile)) &&
			files.FileExists(fmt.Sprintf("%s/%s/%s", CertificateDirPath, domainName, PrivkeyFile)) &&
			files.FileExists(fmt.Sprintf("%s/%s/%s", CertificateDirPath, domainName, ChainFile)) {
			setDomainCertDirPerms(domainName)
			setDomainCertArchiveDirPerms(domainName)
			setDomainCertArchiveFilePerms(domainName)

			return true
		}
		return false
	} else {
		if files.CheckIfFileExistsAndSetPermissionsUsingLinux(currentUsername, fmt.Sprintf("%s/%s/%s", CertificateDirPath, domainName, FullchainFile), "0600") &&
			files.CheckIfFileExistsAndSetPermissionsUsingLinux(currentUsername, fmt.Sprintf("%s/%s/%s", CertificateDirPath, domainName, PrivkeyFile), "0600") &&
			files.CheckIfFileExistsAndSetPermissionsUsingLinux(currentUsername, fmt.Sprintf("%s/%s/%s", CertificateDirPath, domainName, ChainFile), "0600") {
			DomainCertificateDirPath := fmt.Sprintf("%s/%s", CertificateDirPath, domainName)
			directories.CheckIfDirectoryExistsAndSetPermissionsUsingLinux(currentUsername, DomainCertificateDirPath, "0700")

			DomainCertificateArchiveDirPath := fmt.Sprintf("%s/%s", CertificateArchiveDirPath, domainName)
			directories.CheckIfDirectoryExistsAndSetPermissionsUsingLinux(currentUsername, DomainCertificateArchiveDirPath, "0700")

			FullchainArchiveFilePath := fmt.Sprintf("%s/%s/%s", CertificateArchiveDirPath, domainName, FullchainArchiveFile)
			PrivkeyArchiveFilePath := fmt.Sprintf("%s/%s/%s", CertificateArchiveDirPath, domainName, PrivkeyArchiveFile)
			ChainArchiveFilePath := fmt.Sprintf("%s/%s/%s", CertificateArchiveDirPath, domainName, ChainArchiveFile)
			CertArchiveFilePath := fmt.Sprintf("%s/%s/%s", CertificateArchiveDirPath, domainName, CertArchiveFile)

			files.CheckIfFileExistsAndSetPermissionsUsingLinux(currentUsername, FullchainArchiveFilePath, "0600")
			files.CheckIfFileExistsAndSetPermissionsUsingLinux(currentUsername, PrivkeyArchiveFilePath, "0600")
			files.CheckIfFileExistsAndSetPermissionsUsingLinux(currentUsername, ChainArchiveFilePath, "0600")
			files.CheckIfFileExistsAndSetPermissionsUsingLinux(currentUsername, CertArchiveFilePath, "0600")

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
				pterm.Println()
				pterm.Error.Printfln("Failed to set Certbot email: %v", err)
				os.Exit(1)
			}
		} else {
			err := exec.Command("sudo", "certbot", "update_account", "--email", email, "--no-eff-email").Run()
			if err != nil {
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
					pterm.Println()
					pterm.Error.Printfln("Failed to update Certbot email: %v", err)
					os.Exit(1)
				}
			} else {
				err := exec.Command("sudo", "certbot", "update_account", "--email", email, "--no-eff-email").Run()
				if err != nil {
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
	if email == "" {
		if currentUsername == relays.RootUser {
			cmd := exec.Command("certbot", "certonly", "--webroot", "-w", fmt.Sprintf("%s/%s", WWWDirPath, domainName), "-d", domainName, "--agree-tos", "--no-eff-email", "-q", "--register-unsafely-without-email")
			err := cmd.Run()
			if err != nil {
				pterm.Error.Printfln("Certbot failed to obtain the certificate for %s: %v", domainName, err)
				os.Exit(1)
			}
		} else {
			cmd := exec.Command("sudo", "certbot", "certonly", "--webroot", "-w", fmt.Sprintf("%s/%s", WWWDirPath, domainName), "-d", domainName, "--agree-tos", "--no-eff-email", "-q", "--register-unsafely-without-email")
			err := cmd.Run()
			if err != nil {
				pterm.Error.Printfln("Certbot failed to obtain the certificate for %s: %v", domainName, err)
				os.Exit(1)
			}
		}
	} else {
		if currentUsername == relays.RootUser {
			cmd := exec.Command("certbot", "certonly", "--webroot", "-w", fmt.Sprintf("%s/%s", WWWDirPath, domainName), "-d", domainName, "--email", email, "--agree-tos", "--no-eff-email", "-q")
			err := cmd.Run()
			if err != nil {
				pterm.Error.Printfln("Certbot failed to obtain the certificate for %s: %v", domainName, err)
				os.Exit(1)
			}
		} else {
			cmd := exec.Command("sudo", "certbot", "certonly", "--webroot", "-w", fmt.Sprintf("%s/%s", WWWDirPath, domainName), "-d", domainName, "--email", email, "--agree-tos", "--no-eff-email", "-q")
			err := cmd.Run()
			if err != nil {
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
		DomainCertificateDirPath := fmt.Sprintf("%s/%s", CertificateDirPath, domainName)
		directories.CheckIfDirectoryExistsAndSetPermissionsUsingLinux(currentUsername, DomainCertificateDirPath, "0700")

		DomainCertificateArchiveDirPath := fmt.Sprintf("%s/%s", CertificateArchiveDirPath, domainName)
		directories.CheckIfDirectoryExistsAndSetPermissionsUsingLinux(currentUsername, DomainCertificateArchiveDirPath, "0700")

		FullchainArchiveFilePath := fmt.Sprintf("%s/%s/%s", CertificateArchiveDirPath, domainName, FullchainArchiveFile)
		PrivkeyArchiveFilePath := fmt.Sprintf("%s/%s/%s", CertificateArchiveDirPath, domainName, PrivkeyArchiveFile)
		ChainArchiveFilePath := fmt.Sprintf("%s/%s/%s", CertificateArchiveDirPath, domainName, ChainArchiveFile)
		CertArchiveFilePath := fmt.Sprintf("%s/%s/%s", CertificateArchiveDirPath, domainName, CertArchiveFile)
		files.CheckIfFileExistsAndSetPermissionsUsingLinux(currentUsername, FullchainArchiveFilePath, "0600")
		files.CheckIfFileExistsAndSetPermissionsUsingLinux(currentUsername, PrivkeyArchiveFilePath, "0600")
		files.CheckIfFileExistsAndSetPermissionsUsingLinux(currentUsername, ChainArchiveFilePath, "0600")
		files.CheckIfFileExistsAndSetPermissionsUsingLinux(currentUsername, CertArchiveFilePath, "0600")
	}

	certificateSpinner.Success("SSL/TLS certificates obtained successfully.")

	return true
}
