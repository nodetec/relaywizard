package network

import (
	"fmt"
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/pterm/pterm"
	"log"
	"os/exec"
)

// Function to get SSL certificates using Certbot
func GetCertificates(domainName string) bool {

	var ThemeDefault = pterm.ThemeDefault

	var prompt = pterm.InteractiveContinuePrinter{
		DefaultValueIndex: 0,
		DefaultText:       "Obtain SSL certificates?",
		TextStyle:         &ThemeDefault.PrimaryStyle,
		Options:           []string{"yes", "no"},
		OptionsStyle:      &ThemeDefault.SuccessMessageStyle,
		SuffixStyle:       &ThemeDefault.SecondaryStyle,
		Delimiter:         ": ",
	}

	pterm.Println()
	pterm.Println(pterm.Cyan("Do you want to obtain SSL certificates using Certbot?"))
	pterm.Println(pterm.Cyan("If you select 'yes', then this step requires that you already have a configured domain name."))
	pterm.Println(pterm.Cyan("You can always re-run this installer after you have configured your domain name."))
	pterm.Println()

	result, _ := prompt.Show()

	if result == "no" {
		return false
	}

	pterm.Println()
	pterm.Println(pterm.Yellow("Leave email empty if you don't want to receive notifications from Let's Encrypt about your SSL certificates."))

	pterm.Println()
	email, _ := pterm.DefaultInteractiveTextInput.Show("Email address")
	pterm.Println()

	spinner, _ := pterm.DefaultSpinner.Start("Checking SSL certificates...")

	var certificatePath = fmt.Sprintf("/etc/letsencrypt/live/%s", domainName)

	// Check if certificates already exist
	if files.FileExists(fmt.Sprintf("%s/fullchain.pem", certificatePath)) &&
		files.FileExists(fmt.Sprintf("%s/privkey.pem", certificatePath)) &&
		files.FileExists(fmt.Sprintf("%s/chain.pem", certificatePath)) {
		spinner.Info("SSL certificates already exist.")
		return true
	}

	directories.CreateDirectory(fmt.Sprintf("/var/www/%s/.well-known/acme-challenge/", domainName), 0755)

	spinner.UpdateText("Obtaining SSL certificates...")
	if email == "" {
		cmd := exec.Command("certbot", "certonly", "--webroot", "-w", fmt.Sprintf("/var/www/%s", domainName), "-d", domainName, "--agree-tos", "--no-eff-email", "-q", "--register-unsafely-without-email")
		err := cmd.Run()
		if err != nil {
			log.Fatalf("Certbot failed to obtain the certificate for %s: %v", domainName, err)
		}
	} else {
		cmd := exec.Command("certbot", "certonly", "--webroot", "-w", fmt.Sprintf("/var/www/%s", domainName), "-d", domainName, "--email", email, "--agree-tos", "--no-eff-email", "-q")
		err := cmd.Run()
		if err != nil {
			log.Fatalf("Certbot failed to obtain the certificate for %s: %v", domainName, err)
		}
	}

	spinner.Success("SSL certificates obtained successfully.")
	return true
}
