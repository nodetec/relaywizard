package network

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

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

// Function to get SSL/TLS certificates using Certbot
func GetCertificates(domainName string) bool {
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
	pterm.Println(pterm.Cyan("Do you want to obtain SSL/TLS certificates using Certbot?"))
	pterm.Println(pterm.Cyan("If you select 'yes', then this step requires that you already have a configured domain name."))
	pterm.Println(pterm.Cyan("You can always re-run this installer after you have configured your domain name."))
	pterm.Println()

	result, _ := prompt.Show()

	if result == "no" {
		pterm.Println()
		return false
	}

	pterm.Println()
	certbotSpinner, _ := pterm.DefaultSpinner.Start("Checking for Certbot email...")

	out, err := exec.Command("certbot", "show_account").Output()

	if err != nil {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Failed to retrieve Certbot account data: %v", err))
		os.Exit(1)
	}

	certbotAccountData := string(out)
	var email string

	if strings.Contains(certbotAccountData, "Email contact: none") {
		pterm.Println()
		certbotSpinner.Info("Certbot email currently set to none.")

		pterm.Println()
		pterm.Println(pterm.Cyan("Set your Certbot email to receive notifications from Let's Encrypt about your SSL/TLS certificates."))

		pterm.Println()
		pterm.Println(pterm.Yellow("Leave email empty if you don't want to receive notifications."))

		pterm.Println()
		email, _ = pterm.DefaultInteractiveTextInput.Show("Email address")

		err := exec.Command("certbot", "update_account", "--email", email, "--no-eff-email").Run()
		if err != nil {
			pterm.Println()
			pterm.Error.Println(fmt.Sprintf("Failed to set Certbot email: %v", err))
			os.Exit(1)
		}
	} else {
		_, currentEmail, _ := strings.Cut(certbotAccountData, "Email contact: ")
		pterm.Println()
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
			pterm.Println(pterm.Cyan("Set your Certbot email to receive notifications from Let's Encrypt about your SSL/TLS certificates."))

			pterm.Println()
			pterm.Println(pterm.Yellow("Leave email empty if you don't want to receive notifications."))

			pterm.Println()
			email, _ = pterm.DefaultInteractiveTextInput.Show("Email address")

			err := exec.Command("certbot", "update_account", "--email", email, "--no-eff-email").Run()
			if err != nil {
				pterm.Println()
				pterm.Error.Println(fmt.Sprintf("Failed to update Certbot email: %v", err))
				os.Exit(1)
			}
		}
	}

	pterm.Println()
	certificateSpinner, _ := pterm.DefaultSpinner.Start("Checking SSL/TLS certificates...")

	// Check if certificates already exist
	if files.FileExists(fmt.Sprintf("%s/%s/%s", CertificateDirPath, domainName, FullchainFile)) &&
		files.FileExists(fmt.Sprintf("%s/%s/%s", CertificateDirPath, domainName, PrivkeyFile)) &&
		files.FileExists(fmt.Sprintf("%s/%s/%s", CertificateDirPath, domainName, ChainFile)) {
		setDomainCertDirPerms(domainName)
		setDomainCertArchiveDirPerms(domainName)
		setDomainCertArchiveFilePerms(domainName)

		certificateSpinner.Info("SSL/TLS certificates already exist.")
		pterm.Println()
		return true
	}

	certificateSpinner.UpdateText("Obtaining SSL/TLS certificates...")
	if email == "" {
		cmd := exec.Command("certbot", "certonly", "--webroot", "-w", fmt.Sprintf("%s/%s", WWWDirPath, domainName), "-d", domainName, "--agree-tos", "--no-eff-email", "-q", "--register-unsafely-without-email")
		err := cmd.Run()
		if err != nil {
			pterm.Error.Println(fmt.Sprintf("Certbot failed to obtain the certificate for %s: %v", domainName, err))
			os.Exit(1)
		}
	} else {
		cmd := exec.Command("certbot", "certonly", "--webroot", "-w", fmt.Sprintf("%s/%s", WWWDirPath, domainName), "-d", domainName, "--email", email, "--agree-tos", "--no-eff-email", "-q")
		err := cmd.Run()
		if err != nil {
			pterm.Error.Println(fmt.Sprintf("Certbot failed to obtain the certificate for %s: %v", domainName, err))
			os.Exit(1)
		}
	}

	setDomainCertDirPerms(domainName)
	setDomainCertArchiveDirPerms(domainName)
	setDomainCertArchiveFilePerms(domainName)

	certificateSpinner.Success("SSL/TLS certificates obtained successfully.")
	return true
}
