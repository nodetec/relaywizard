package network

import (
	"fmt"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/pterm/pterm"
	"os"
	"os/exec"
)

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
		return false
	}

	pterm.Println()
	pterm.Println(pterm.Yellow("Leave email empty if you don't want to receive notifications from Let's Encrypt about your SSL/TLS certificates."))

	pterm.Println()
	email, _ := pterm.DefaultInteractiveTextInput.Show("Email address")
	pterm.Println()

	spinner, _ := pterm.DefaultSpinner.Start("Checking SSL/TLS certificates...")

	// Check if certificates already exist
	if files.FileExists(fmt.Sprintf("%s/%s/%s", CertificateDirPath, domainName, FullchainFile)) &&
		files.FileExists(fmt.Sprintf("%s/%s/%s", CertificateDirPath, domainName, PrivkeyFile)) &&
		files.FileExists(fmt.Sprintf("%s/%s/%s", CertificateDirPath, domainName, ChainFile)) {
		spinner.Info("SSL/TLS certificates already exist.")
		return true
	}

	spinner.UpdateText("Obtaining SSL/TLS certificates...")
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

	spinner.Success("SSL/TLS certificates obtained successfully.")
	return true
}
